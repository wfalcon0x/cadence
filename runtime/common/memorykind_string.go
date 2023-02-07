// Code generated by "stringer -type=MemoryKind -trimprefix=MemoryKind"; DO NOT EDIT.

package common

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[MemoryKindUnknown-0]
	_ = x[MemoryKindAddressValue-1]
	_ = x[MemoryKindStringValue-2]
	_ = x[MemoryKindCharacterValue-3]
	_ = x[MemoryKindNumberValue-4]
	_ = x[MemoryKindArrayValueBase-5]
	_ = x[MemoryKindDictionaryValueBase-6]
	_ = x[MemoryKindCompositeValueBase-7]
	_ = x[MemoryKindSimpleCompositeValueBase-8]
	_ = x[MemoryKindOptionalValue-9]
	_ = x[MemoryKindTypeValue-10]
	_ = x[MemoryKindPathValue-11]
	_ = x[MemoryKindStorageCapabilityValue-12]
	_ = x[MemoryKindPathLinkValue-13]
	_ = x[MemoryKindAccountLinkValue-14]
	_ = x[MemoryKindStorageReferenceValue-15]
	_ = x[MemoryKindAccountReferenceValue-16]
	_ = x[MemoryKindEphemeralReferenceValue-17]
	_ = x[MemoryKindInterpretedFunctionValue-18]
	_ = x[MemoryKindHostFunctionValue-19]
	_ = x[MemoryKindBoundFunctionValue-20]
	_ = x[MemoryKindBigInt-21]
	_ = x[MemoryKindSimpleCompositeValue-22]
	_ = x[MemoryKindPublishedValue-23]
	_ = x[MemoryKindCapabilityControllerValue-24]
	_ = x[MemoryKindAtreeArrayDataSlab-25]
	_ = x[MemoryKindAtreeArrayMetaDataSlab-26]
	_ = x[MemoryKindAtreeArrayElementOverhead-27]
	_ = x[MemoryKindAtreeMapDataSlab-28]
	_ = x[MemoryKindAtreeMapMetaDataSlab-29]
	_ = x[MemoryKindAtreeMapElementOverhead-30]
	_ = x[MemoryKindAtreeMapPreAllocatedElement-31]
	_ = x[MemoryKindAtreeEncodedSlab-32]
	_ = x[MemoryKindPrimitiveStaticType-33]
	_ = x[MemoryKindCompositeStaticType-34]
	_ = x[MemoryKindInterfaceStaticType-35]
	_ = x[MemoryKindVariableSizedStaticType-36]
	_ = x[MemoryKindConstantSizedStaticType-37]
	_ = x[MemoryKindDictionaryStaticType-38]
	_ = x[MemoryKindOptionalStaticType-39]
	_ = x[MemoryKindRestrictedStaticType-40]
	_ = x[MemoryKindReferenceStaticType-41]
	_ = x[MemoryKindCapabilityStaticType-42]
	_ = x[MemoryKindFunctionStaticType-43]
	_ = x[MemoryKindCapabilityControllerStaticType-44]
	_ = x[MemoryKindCadenceVoidValue-45]
	_ = x[MemoryKindCadenceOptionalValue-46]
	_ = x[MemoryKindCadenceBoolValue-47]
	_ = x[MemoryKindCadenceStringValue-48]
	_ = x[MemoryKindCadenceCharacterValue-49]
	_ = x[MemoryKindCadenceAddressValue-50]
	_ = x[MemoryKindCadenceIntValue-51]
	_ = x[MemoryKindCadenceNumberValue-52]
	_ = x[MemoryKindCadenceArrayValueBase-53]
	_ = x[MemoryKindCadenceArrayValueLength-54]
	_ = x[MemoryKindCadenceDictionaryValue-55]
	_ = x[MemoryKindCadenceKeyValuePair-56]
	_ = x[MemoryKindCadenceStructValueBase-57]
	_ = x[MemoryKindCadenceStructValueSize-58]
	_ = x[MemoryKindCadenceResourceValueBase-59]
	_ = x[MemoryKindCadenceAttachmentValueBase-60]
	_ = x[MemoryKindCadenceResourceValueSize-61]
	_ = x[MemoryKindCadenceAttachmentValueSize-62]
	_ = x[MemoryKindCadenceEventValueBase-63]
	_ = x[MemoryKindCadenceEventValueSize-64]
	_ = x[MemoryKindCadenceContractValueBase-65]
	_ = x[MemoryKindCadenceContractValueSize-66]
	_ = x[MemoryKindCadenceEnumValueBase-67]
	_ = x[MemoryKindCadenceEnumValueSize-68]
	_ = x[MemoryKindCadencePathLinkValue-69]
	_ = x[MemoryKindCadenceAccountLinkValue-70]
	_ = x[MemoryKindCadencePathValue-71]
	_ = x[MemoryKindCadenceTypeValue-72]
	_ = x[MemoryKindCadenceStorageCapabilityValue-73]
	_ = x[MemoryKindCadenceFunctionValue-74]
	_ = x[MemoryKindCadenceOptionalType-75]
	_ = x[MemoryKindCadenceVariableSizedArrayType-76]
	_ = x[MemoryKindCadenceConstantSizedArrayType-77]
	_ = x[MemoryKindCadenceDictionaryType-78]
	_ = x[MemoryKindCadenceField-79]
	_ = x[MemoryKindCadenceParameter-80]
	_ = x[MemoryKindCadenceTypeParameter-81]
	_ = x[MemoryKindCadenceStructType-82]
	_ = x[MemoryKindCadenceResourceType-83]
	_ = x[MemoryKindCadenceAttachmentType-84]
	_ = x[MemoryKindCadenceEventType-85]
	_ = x[MemoryKindCadenceContractType-86]
	_ = x[MemoryKindCadenceStructInterfaceType-87]
	_ = x[MemoryKindCadenceResourceInterfaceType-88]
	_ = x[MemoryKindCadenceContractInterfaceType-89]
	_ = x[MemoryKindCadenceFunctionType-90]
	_ = x[MemoryKindCadenceReferenceType-91]
	_ = x[MemoryKindCadenceRestrictedType-92]
	_ = x[MemoryKindCadenceCapabilityType-93]
	_ = x[MemoryKindCadenceEnumType-94]
	_ = x[MemoryKindRawString-95]
	_ = x[MemoryKindAddressLocation-96]
	_ = x[MemoryKindBytes-97]
	_ = x[MemoryKindVariable-98]
	_ = x[MemoryKindCompositeTypeInfo-99]
	_ = x[MemoryKindCompositeField-100]
	_ = x[MemoryKindInvocation-101]
	_ = x[MemoryKindStorageMap-102]
	_ = x[MemoryKindStorageKey-103]
	_ = x[MemoryKindTypeToken-104]
	_ = x[MemoryKindErrorToken-105]
	_ = x[MemoryKindSpaceToken-106]
	_ = x[MemoryKindProgram-107]
	_ = x[MemoryKindIdentifier-108]
	_ = x[MemoryKindArgument-109]
	_ = x[MemoryKindBlock-110]
	_ = x[MemoryKindFunctionBlock-111]
	_ = x[MemoryKindParameter-112]
	_ = x[MemoryKindParameterList-113]
	_ = x[MemoryKindTypeParameter-114]
	_ = x[MemoryKindTypeParameterList-115]
	_ = x[MemoryKindTransfer-116]
	_ = x[MemoryKindMembers-117]
	_ = x[MemoryKindTypeAnnotation-118]
	_ = x[MemoryKindDictionaryEntry-119]
	_ = x[MemoryKindFunctionDeclaration-120]
	_ = x[MemoryKindCompositeDeclaration-121]
	_ = x[MemoryKindAttachmentDeclaration-122]
	_ = x[MemoryKindInterfaceDeclaration-123]
	_ = x[MemoryKindEnumCaseDeclaration-124]
	_ = x[MemoryKindFieldDeclaration-125]
	_ = x[MemoryKindTransactionDeclaration-126]
	_ = x[MemoryKindImportDeclaration-127]
	_ = x[MemoryKindVariableDeclaration-128]
	_ = x[MemoryKindSpecialFunctionDeclaration-129]
	_ = x[MemoryKindPragmaDeclaration-130]
	_ = x[MemoryKindAssignmentStatement-131]
	_ = x[MemoryKindBreakStatement-132]
	_ = x[MemoryKindContinueStatement-133]
	_ = x[MemoryKindEmitStatement-134]
	_ = x[MemoryKindExpressionStatement-135]
	_ = x[MemoryKindForStatement-136]
	_ = x[MemoryKindIfStatement-137]
	_ = x[MemoryKindReturnStatement-138]
	_ = x[MemoryKindSwapStatement-139]
	_ = x[MemoryKindSwitchStatement-140]
	_ = x[MemoryKindWhileStatement-141]
	_ = x[MemoryKindRemoveStatement-142]
	_ = x[MemoryKindBooleanExpression-143]
	_ = x[MemoryKindVoidExpression-144]
	_ = x[MemoryKindNilExpression-145]
	_ = x[MemoryKindStringExpression-146]
	_ = x[MemoryKindIntegerExpression-147]
	_ = x[MemoryKindFixedPointExpression-148]
	_ = x[MemoryKindArrayExpression-149]
	_ = x[MemoryKindDictionaryExpression-150]
	_ = x[MemoryKindIdentifierExpression-151]
	_ = x[MemoryKindInvocationExpression-152]
	_ = x[MemoryKindMemberExpression-153]
	_ = x[MemoryKindIndexExpression-154]
	_ = x[MemoryKindConditionalExpression-155]
	_ = x[MemoryKindUnaryExpression-156]
	_ = x[MemoryKindBinaryExpression-157]
	_ = x[MemoryKindFunctionExpression-158]
	_ = x[MemoryKindCastingExpression-159]
	_ = x[MemoryKindCreateExpression-160]
	_ = x[MemoryKindDestroyExpression-161]
	_ = x[MemoryKindReferenceExpression-162]
	_ = x[MemoryKindForceExpression-163]
	_ = x[MemoryKindPathExpression-164]
	_ = x[MemoryKindAttachExpression-165]
	_ = x[MemoryKindConstantSizedType-166]
	_ = x[MemoryKindDictionaryType-167]
	_ = x[MemoryKindFunctionType-168]
	_ = x[MemoryKindInstantiationType-169]
	_ = x[MemoryKindNominalType-170]
	_ = x[MemoryKindOptionalType-171]
	_ = x[MemoryKindReferenceType-172]
	_ = x[MemoryKindRestrictedType-173]
	_ = x[MemoryKindVariableSizedType-174]
	_ = x[MemoryKindPosition-175]
	_ = x[MemoryKindRange-176]
	_ = x[MemoryKindElaboration-177]
	_ = x[MemoryKindActivation-178]
	_ = x[MemoryKindActivationEntries-179]
	_ = x[MemoryKindVariableSizedSemaType-180]
	_ = x[MemoryKindConstantSizedSemaType-181]
	_ = x[MemoryKindDictionarySemaType-182]
	_ = x[MemoryKindOptionalSemaType-183]
	_ = x[MemoryKindRestrictedSemaType-184]
	_ = x[MemoryKindReferenceSemaType-185]
	_ = x[MemoryKindCapabilitySemaType-186]
	_ = x[MemoryKindCapabilityControllerSemaType-187]
	_ = x[MemoryKindOrderedMap-188]
	_ = x[MemoryKindOrderedMapEntryList-189]
	_ = x[MemoryKindOrderedMapEntry-190]
	_ = x[MemoryKindLast-191]
}

const _MemoryKind_name = "UnknownAddressValueStringValueCharacterValueNumberValueArrayValueBaseDictionaryValueBaseCompositeValueBaseSimpleCompositeValueBaseOptionalValueTypeValuePathValueStorageCapabilityValuePathLinkValueAccountLinkValueStorageReferenceValueAccountReferenceValueEphemeralReferenceValueInterpretedFunctionValueHostFunctionValueBoundFunctionValueBigIntSimpleCompositeValuePublishedValueCapabilityControllerValueAtreeArrayDataSlabAtreeArrayMetaDataSlabAtreeArrayElementOverheadAtreeMapDataSlabAtreeMapMetaDataSlabAtreeMapElementOverheadAtreeMapPreAllocatedElementAtreeEncodedSlabPrimitiveStaticTypeCompositeStaticTypeInterfaceStaticTypeVariableSizedStaticTypeConstantSizedStaticTypeDictionaryStaticTypeOptionalStaticTypeRestrictedStaticTypeReferenceStaticTypeCapabilityStaticTypeFunctionStaticTypeCapabilityControllerStaticTypeCadenceVoidValueCadenceOptionalValueCadenceBoolValueCadenceStringValueCadenceCharacterValueCadenceAddressValueCadenceIntValueCadenceNumberValueCadenceArrayValueBaseCadenceArrayValueLengthCadenceDictionaryValueCadenceKeyValuePairCadenceStructValueBaseCadenceStructValueSizeCadenceResourceValueBaseCadenceAttachmentValueBaseCadenceResourceValueSizeCadenceAttachmentValueSizeCadenceEventValueBaseCadenceEventValueSizeCadenceContractValueBaseCadenceContractValueSizeCadenceEnumValueBaseCadenceEnumValueSizeCadencePathLinkValueCadenceAccountLinkValueCadencePathValueCadenceTypeValueCadenceStorageCapabilityValueCadenceFunctionValueCadenceOptionalTypeCadenceVariableSizedArrayTypeCadenceConstantSizedArrayTypeCadenceDictionaryTypeCadenceFieldCadenceParameterCadenceTypeParameterCadenceStructTypeCadenceResourceTypeCadenceAttachmentTypeCadenceEventTypeCadenceContractTypeCadenceStructInterfaceTypeCadenceResourceInterfaceTypeCadenceContractInterfaceTypeCadenceFunctionTypeCadenceReferenceTypeCadenceRestrictedTypeCadenceCapabilityTypeCadenceEnumTypeRawStringAddressLocationBytesVariableCompositeTypeInfoCompositeFieldInvocationStorageMapStorageKeyTypeTokenErrorTokenSpaceTokenProgramIdentifierArgumentBlockFunctionBlockParameterParameterListTypeParameterTypeParameterListTransferMembersTypeAnnotationDictionaryEntryFunctionDeclarationCompositeDeclarationAttachmentDeclarationInterfaceDeclarationEnumCaseDeclarationFieldDeclarationTransactionDeclarationImportDeclarationVariableDeclarationSpecialFunctionDeclarationPragmaDeclarationAssignmentStatementBreakStatementContinueStatementEmitStatementExpressionStatementForStatementIfStatementReturnStatementSwapStatementSwitchStatementWhileStatementRemoveStatementBooleanExpressionVoidExpressionNilExpressionStringExpressionIntegerExpressionFixedPointExpressionArrayExpressionDictionaryExpressionIdentifierExpressionInvocationExpressionMemberExpressionIndexExpressionConditionalExpressionUnaryExpressionBinaryExpressionFunctionExpressionCastingExpressionCreateExpressionDestroyExpressionReferenceExpressionForceExpressionPathExpressionAttachExpressionConstantSizedTypeDictionaryTypeFunctionTypeInstantiationTypeNominalTypeOptionalTypeReferenceTypeRestrictedTypeVariableSizedTypePositionRangeElaborationActivationActivationEntriesVariableSizedSemaTypeConstantSizedSemaTypeDictionarySemaTypeOptionalSemaTypeRestrictedSemaTypeReferenceSemaTypeCapabilitySemaTypeCapabilityControllerSemaTypeOrderedMapOrderedMapEntryListOrderedMapEntryLast"

var _MemoryKind_index = [...]uint16{0, 7, 19, 30, 44, 55, 69, 88, 106, 130, 143, 152, 161, 183, 196, 212, 233, 254, 277, 301, 318, 336, 342, 362, 376, 401, 419, 441, 466, 482, 502, 525, 552, 568, 587, 606, 625, 648, 671, 691, 709, 729, 748, 768, 786, 816, 832, 852, 868, 886, 907, 926, 941, 959, 980, 1003, 1025, 1044, 1066, 1088, 1112, 1138, 1162, 1188, 1209, 1230, 1254, 1278, 1298, 1318, 1338, 1361, 1377, 1393, 1422, 1442, 1461, 1490, 1519, 1540, 1552, 1568, 1588, 1605, 1624, 1645, 1661, 1680, 1706, 1734, 1762, 1781, 1801, 1822, 1843, 1858, 1867, 1882, 1887, 1895, 1912, 1926, 1936, 1946, 1956, 1965, 1975, 1985, 1992, 2002, 2010, 2015, 2028, 2037, 2050, 2063, 2080, 2088, 2095, 2109, 2124, 2143, 2163, 2184, 2204, 2223, 2239, 2261, 2278, 2297, 2323, 2340, 2359, 2373, 2390, 2403, 2422, 2434, 2445, 2460, 2473, 2488, 2502, 2517, 2534, 2548, 2561, 2577, 2594, 2614, 2629, 2649, 2669, 2689, 2705, 2720, 2741, 2756, 2772, 2790, 2807, 2823, 2840, 2859, 2874, 2888, 2904, 2921, 2935, 2947, 2964, 2975, 2987, 3000, 3014, 3031, 3039, 3044, 3055, 3065, 3082, 3103, 3124, 3142, 3158, 3176, 3193, 3211, 3239, 3249, 3268, 3283, 3287}

func (i MemoryKind) String() string {
	if i >= MemoryKind(len(_MemoryKind_index)-1) {
		return "MemoryKind(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _MemoryKind_name[_MemoryKind_index[i]:_MemoryKind_index[i+1]]
}
