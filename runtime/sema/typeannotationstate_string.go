// Code generated by "stringer -type=TypeAnnotationState"; DO NOT EDIT.

package sema

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[TypeAnnotationStateUnknown-0]
	_ = x[TypeAnnotationStateValid-1]
	_ = x[TypeAnnotationStateInvalidResourceAnnotation-2]
	_ = x[TypeAnnotationStateMissingResourceAnnotation-3]
	_ = x[TypeAnnotationStateDirectEntitlementTypeAnnotation-4]
	_ = x[TypeAnnotationStateDirectAttachmentTypeAnnotation-5]
}

const _TypeAnnotationState_name = "TypeAnnotationStateUnknownTypeAnnotationStateValidTypeAnnotationStateInvalidResourceAnnotationTypeAnnotationStateMissingResourceAnnotationTypeAnnotationStateDirectEntitlementTypeAnnotationTypeAnnotationStateDirectAttachmentTypeAnnotation"

var _TypeAnnotationState_index = [...]uint8{0, 26, 50, 94, 138, 188, 237}

func (i TypeAnnotationState) String() string {
	if i >= TypeAnnotationState(len(_TypeAnnotationState_index)-1) {
		return "TypeAnnotationState(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _TypeAnnotationState_name[_TypeAnnotationState_index[i]:_TypeAnnotationState_index[i+1]]
}
