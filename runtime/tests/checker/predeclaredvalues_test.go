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

package checker

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/sema"
	"github.com/onflow/cadence/runtime/stdlib"
)

func TestCheckPredeclaredValues(t *testing.T) {

	t.Parallel()

	baseValueActivation := sema.NewVariableActivation(sema.BaseValueActivation)

	valueDeclaration := stdlib.NewStandardLibraryFunction(
		"foo",
		&sema.FunctionType{
			ReturnTypeAnnotation: sema.TypeAnnotation{
				Type: sema.VoidType,
			},
		},
		"",
		nil,
	)
	baseValueActivation.DeclareValue(valueDeclaration)

	_, err := ParseAndCheckWithOptions(t,
		`
            pub fun test() {
                foo()
            }
        `,
		ParseAndCheckOptions{
			Config: &sema.Config{
				BaseValueActivationHandler: func(_ common.Location) *sema.VariableActivation {
					return baseValueActivation
				},
			},
		},
	)

	require.NoError(t, err)

}
