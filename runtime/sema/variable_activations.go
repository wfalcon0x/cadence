/*
 * Cadence - The resource-oriented smart contract programming language
 *
 * Copyright Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package sema

import (
	"sync"

	"github.com/onflow/cadence/runtime/ast"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/errors"
)

// An VariableActivation is a map of strings to variables.
// It is used to represent an active scope in a program,
// i.e. it is used as a symbol table during semantic analysis.
type VariableActivation struct {
	entries        *StringVariableOrderedMap
	Parent         *VariableActivation
	LeaveCallbacks []func(EndPositionGetter)
	Depth          int
}

type EndPositionGetter func(common.MemoryGauge) ast.Position

// NewVariableActivation returns as new activation with the given parent.
// The parent may be nil.
func NewVariableActivation(parent *VariableActivation) *VariableActivation {
	activation := &VariableActivation{}
	activation.SetParent(parent)
	return activation
}

// SetParent sets the parent of this activation to the given parent
// and updates the depth.
func (a *VariableActivation) SetParent(parent *VariableActivation) {
	a.Parent = parent

	var depth int
	if parent != nil {
		depth = parent.Depth + 1
	}
	a.Depth = depth
}

// Find returns the variable for a given name in the activation.
// It returns nil if no variable is found.
func (a *VariableActivation) Find(name string) *Variable {

	current := a

	for current != nil {
		if current.entries != nil {
			result, ok := current.entries.Get(name)
			if ok {
				return result
			}
		}

		current = current.Parent
	}

	return nil
}

// Set sets the given variable.
func (a *VariableActivation) Set(name string, variable *Variable) {
	if a.entries == nil {
		a.entries = &StringVariableOrderedMap{}
	}

	a.entries.Set(name, variable)
}

// Clear removes all variables from this activation.
func (a *VariableActivation) Clear() {
	a.LeaveCallbacks = nil

	if a.entries == nil {
		return
	}

	a.entries.Clear()
}

// ForEach calls the given function for each name-variable pair in the activation.
// It can be used to iterate over all entries of the activation.
func (a *VariableActivation) ForEach(cb func(string, *Variable) error) error {

	activation := a

	for activation != nil {

		if activation.entries != nil {
			for pair := activation.entries.Oldest(); pair != nil; pair = pair.Next() {
				err := cb(pair.Key, pair.Value)
				if err != nil {
					return err
				}
			}
		}

		activation = activation.Parent
	}

	return nil
}

func (a *VariableActivation) DeclareValue(declaration ValueDeclaration) {
	name := declaration.ValueDeclarationName()

	a.Set(name, &Variable{
		Identifier:      name,
		DeclarationKind: declaration.ValueDeclarationKind(),
		Type:            declaration.ValueDeclarationType(),
		// TODO: add access to ValueDeclaration and use declaration's access instead here
		Access:          PrimitiveAccess(ast.AccessAll),
		IsConstant:      declaration.ValueDeclarationIsConstant(),
		ArgumentLabels:  declaration.ValueDeclarationArgumentLabels(),
		Pos:             declaration.ValueDeclarationPosition(),
		DocString:       declaration.ValueDeclarationDocString(),
		ActivationDepth: 0,
	})
}

func (a *VariableActivation) DeclareType(declaration TypeDeclaration) {
	name := declaration.TypeDeclarationName()

	a.Set(name, &Variable{
		Identifier:      name,
		DeclarationKind: declaration.TypeDeclarationKind(),
		Type:            declaration.TypeDeclarationType(),
		// TODO: add access to TypeDeclaration and use declaration's access instead here
		Access:         PrimitiveAccess(ast.AccessAll),
		IsConstant:     true,
		ArgumentLabels: nil,
		Pos:            declaration.TypeDeclarationPosition(),
		// TODO: DocString
		ActivationDepth: 0,
	})
}

var variableActivationPool = sync.Pool{
	New: func() any {
		return &VariableActivation{}
	},
}

func getVariableActivation() *VariableActivation {
	activation, ok := variableActivationPool.Get().(*VariableActivation)
	if !ok {
		panic(errors.NewUnreachableError())
	}
	activation.Clear()
	return activation
}

// VariableActivations is a stack of activation records.
// Each entry represents a new activation record.
//
// The current / most nested activation record can be found
// at the top of the stack (see function `Current`).
type VariableActivations struct {
	activations []*VariableActivation
}

func NewVariableActivations(parent *VariableActivation) *VariableActivations {
	activations := &VariableActivations{}
	activations.pushNewWithParent(parent)
	return activations
}

// pushNewWithParent pushes a new empty activation
// to the top of the activation stack.
// The new activation has the given parent as its parent.
func (a *VariableActivations) pushNewWithParent(parent *VariableActivation) *VariableActivation {
	activation := getVariableActivation()
	activation.SetParent(parent)
	a.push(activation)
	return activation
}

// push pushes the given activation
// onto the top of the activation stack.
func (a *VariableActivations) push(activation *VariableActivation) {
	a.activations = append(
		a.activations,
		activation,
	)
}

// Enter pushes a new empty activation
// to the top of the activation stack.
// The new activation has the current activation as its parent.
func (a *VariableActivations) Enter() {
	a.pushNewWithParent(a.Current())
}

// Leave pops the top-most (current) activation
// from the top of the activation stack.
func (a *VariableActivations) Leave(getEndPosition func(common.MemoryGauge) ast.Position) {
	count := len(a.activations)
	if count < 1 {
		return
	}
	lastIndex := count - 1
	activation := a.activations[lastIndex]
	a.activations[lastIndex] = nil
	a.activations = a.activations[:lastIndex]
	for _, callback := range activation.LeaveCallbacks {
		callback(getEndPosition)
	}
	variableActivationPool.Put(activation)
}

// Set sets the variable in the current activation.
func (a *VariableActivations) Set(name string, variable *Variable) {
	current := a.Current()
	// create the first scope if there is no scope
	if current == nil {
		current = a.pushNewWithParent(nil)
	}
	current.Set(name, variable)
}

// Find returns the variable for a given name in the current activation.
// It returns nil if no variable is found
// or if there is no current activation.
func (a *VariableActivations) Find(name string) *Variable {

	current := a.Current()
	if current == nil {
		return nil
	}

	return current.Find(name)
}

// Depth returns the depth (size) of the activation stack.
func (a *VariableActivations) Depth() int {
	return len(a.activations)
}

type variableDeclaration struct {
	ty                       Type
	identifier               string
	docString                string
	argumentLabels           []string
	pos                      ast.Position
	access                   Access
	kind                     common.DeclarationKind
	isConstant               bool
	allowOuterScopeShadowing bool
}

func (a *VariableActivations) declare(declaration variableDeclaration) (*Variable, error) {

	depth := a.Depth()

	// Check if a variable with this name is already declared.
	// Report an error if shadowing variables of outer scopes is not allowed,
	// or the existing variable is declared in the current scope,
	// or the existing variable is a built-in.

	existingVariable := a.Find(declaration.identifier)
	if existingVariable != nil &&
		(!declaration.allowOuterScopeShadowing ||
			existingVariable.ActivationDepth == depth ||
			existingVariable.ActivationDepth == 0) {

		return nil, &RedeclarationError{
			Kind:        declaration.kind,
			Name:        declaration.identifier,
			Pos:         declaration.pos,
			PreviousPos: existingVariable.Pos,
		}
	}

	// A variable with this name is not yet declared in the current scope,
	// declare it.

	variable := &Variable{
		Identifier:      declaration.identifier,
		Access:          declaration.access,
		DeclarationKind: declaration.kind,
		IsConstant:      declaration.isConstant,
		ActivationDepth: depth,
		Type:            declaration.ty,
		Pos:             &declaration.pos,
		ArgumentLabels:  declaration.argumentLabels,
		DocString:       declaration.docString,
	}
	a.Set(declaration.identifier, variable)
	return variable, nil
}

func (a *VariableActivations) DeclareValue(declaration ValueDeclaration) (*Variable, error) {
	var variablePos ast.Position
	declarationPos := declaration.ValueDeclarationPosition()
	if declarationPos != nil {
		variablePos = *declarationPos
	}

	return a.declare(variableDeclaration{
		identifier: declaration.ValueDeclarationName(),
		kind:       declaration.ValueDeclarationKind(),
		ty:         declaration.ValueDeclarationType(),
		// TODO: add access to ValueDeclaration and use declaration's access instead here
		access:         PrimitiveAccess(ast.AccessAll),
		isConstant:     declaration.ValueDeclarationIsConstant(),
		argumentLabels: declaration.ValueDeclarationArgumentLabels(),
		pos:            variablePos,
		docString:      declaration.ValueDeclarationDocString(),
	})
}

type typeDeclaration struct {
	ty                       Type
	docString                string
	identifier               ast.Identifier
	declarationKind          common.DeclarationKind
	access                   Access
	allowOuterScopeShadowing bool
}

func (a *VariableActivations) declareType(declaration typeDeclaration) (*Variable, error) {
	return a.declare(
		variableDeclaration{
			identifier:               declaration.identifier.Identifier,
			ty:                       declaration.ty,
			access:                   declaration.access,
			kind:                     declaration.declarationKind,
			pos:                      declaration.identifier.Pos,
			isConstant:               true,
			argumentLabels:           nil,
			allowOuterScopeShadowing: declaration.allowOuterScopeShadowing,
			docString:                declaration.docString,
		},
	)
}

func (a *VariableActivations) declareImplicitConstant(
	identifier string,
	ty Type,
	kind common.DeclarationKind,
) (*Variable, error) {
	return a.declare(
		variableDeclaration{
			identifier:               identifier,
			ty:                       ty,
			access:                   PrimitiveAccess(ast.AccessAll),
			kind:                     kind,
			isConstant:               true,
			allowOuterScopeShadowing: false,
		},
	)
}

func (a *VariableActivations) ForEachVariableDeclaredInAndBelow(depth int, f func(name string, value *Variable)) {

	activation := a.Current()

	_ = activation.ForEach(func(name string, variable *Variable) error {

		if variable.ActivationDepth >= depth {
			f(name, variable)
		}

		return nil
	})
}

// Current returns the current / most nested activation,
// which can be found at the top of the stack.
// It returns nil if there is no active activation.
func (a *VariableActivations) Current() *VariableActivation {
	count := len(a.activations)
	if count < 1 {
		return nil
	}
	return a.activations[count-1]
}
