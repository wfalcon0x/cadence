// Code generated by "stringer -type=PrimitiveStaticType -trimprefix=PrimitiveStaticType"; DO NOT EDIT.

package interpreter

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[PrimitiveStaticTypeUnknown-0]
	_ = x[PrimitiveStaticTypeVoid-1]
	_ = x[PrimitiveStaticTypeAny-2]
	_ = x[PrimitiveStaticTypeNever-3]
	_ = x[PrimitiveStaticTypeAnyStruct-4]
	_ = x[PrimitiveStaticTypeAnyResource-5]
	_ = x[PrimitiveStaticTypeBool-6]
	_ = x[PrimitiveStaticTypeAddress-7]
	_ = x[PrimitiveStaticTypeString-8]
	_ = x[PrimitiveStaticTypeCharacter-9]
	_ = x[PrimitiveStaticTypeMetaType-10]
	_ = x[PrimitiveStaticTypeBlock-11]
	_ = x[PrimitiveStaticTypeAnyResourceAttachment-12]
	_ = x[PrimitiveStaticTypeAnyStructAttachment-13]
	_ = x[PrimitiveStaticTypeNumber-18]
	_ = x[PrimitiveStaticTypeSignedNumber-19]
	_ = x[PrimitiveStaticTypeInteger-24]
	_ = x[PrimitiveStaticTypeSignedInteger-25]
	_ = x[PrimitiveStaticTypeFixedPoint-30]
	_ = x[PrimitiveStaticTypeSignedFixedPoint-31]
	_ = x[PrimitiveStaticTypeInt-36]
	_ = x[PrimitiveStaticTypeInt8-37]
	_ = x[PrimitiveStaticTypeInt16-38]
	_ = x[PrimitiveStaticTypeInt32-39]
	_ = x[PrimitiveStaticTypeInt64-40]
	_ = x[PrimitiveStaticTypeInt128-41]
	_ = x[PrimitiveStaticTypeInt256-42]
	_ = x[PrimitiveStaticTypeUInt-44]
	_ = x[PrimitiveStaticTypeUInt8-45]
	_ = x[PrimitiveStaticTypeUInt16-46]
	_ = x[PrimitiveStaticTypeUInt32-47]
	_ = x[PrimitiveStaticTypeUInt64-48]
	_ = x[PrimitiveStaticTypeUInt128-49]
	_ = x[PrimitiveStaticTypeUInt256-50]
	_ = x[PrimitiveStaticTypeWord8-53]
	_ = x[PrimitiveStaticTypeWord16-54]
	_ = x[PrimitiveStaticTypeWord32-55]
	_ = x[PrimitiveStaticTypeWord64-56]
	_ = x[PrimitiveStaticTypeWord128-57]
	_ = x[PrimitiveStaticTypeWord256-58]
	_ = x[PrimitiveStaticTypeFix64-64]
	_ = x[PrimitiveStaticTypeUFix64-72]
	_ = x[PrimitiveStaticTypePath-76]
	_ = x[PrimitiveStaticTypeCapability-77]
	_ = x[PrimitiveStaticTypeStoragePath-78]
	_ = x[PrimitiveStaticTypeCapabilityPath-79]
	_ = x[PrimitiveStaticTypePublicPath-80]
	_ = x[PrimitiveStaticTypePrivatePath-81]
	_ = x[PrimitiveStaticTypeAuthAccount-90]
	_ = x[PrimitiveStaticTypePublicAccount-91]
	_ = x[PrimitiveStaticTypeDeployedContract-92]
	_ = x[PrimitiveStaticTypeAuthAccountContracts-93]
	_ = x[PrimitiveStaticTypePublicAccountContracts-94]
	_ = x[PrimitiveStaticTypeAuthAccountKeys-95]
	_ = x[PrimitiveStaticTypePublicAccountKeys-96]
	_ = x[PrimitiveStaticTypeAccountKey-97]
	_ = x[PrimitiveStaticTypeAuthAccountInbox-98]
	_ = x[PrimitiveStaticTypeStorageCapabilityController-99]
	_ = x[PrimitiveStaticTypeAccountCapabilityController-100]
	_ = x[PrimitiveStaticTypeAuthAccountStorageCapabilities-101]
	_ = x[PrimitiveStaticTypeAuthAccountAccountCapabilities-102]
	_ = x[PrimitiveStaticTypeAuthAccountCapabilities-103]
	_ = x[PrimitiveStaticTypePublicAccountCapabilities-104]
	_ = x[PrimitiveStaticTypeAccount-105]
	_ = x[PrimitiveStaticTypeAccount_Contracts-106]
	_ = x[PrimitiveStaticTypeAccount_Keys-107]
	_ = x[PrimitiveStaticTypeAccount_Inbox-108]
	_ = x[PrimitiveStaticTypeAccount_StorageCapabilities-109]
	_ = x[PrimitiveStaticTypeAccount_AccountCapabilities-110]
	_ = x[PrimitiveStaticTypeAccount_Capabilities-111]
	_ = x[PrimitiveStaticTypeAccount_Storage-112]
	_ = x[PrimitiveStaticTypeMutate-118]
	_ = x[PrimitiveStaticTypeInsert-119]
	_ = x[PrimitiveStaticTypeRemove-120]
	_ = x[PrimitiveStaticTypeIdentity-121]
	_ = x[PrimitiveStaticTypeStorage-125]
	_ = x[PrimitiveStaticTypeSaveValue-126]
	_ = x[PrimitiveStaticTypeLoadValue-127]
	_ = x[PrimitiveStaticTypeBorrowValue-128]
	_ = x[PrimitiveStaticTypeContracts-129]
	_ = x[PrimitiveStaticTypeAddContract-130]
	_ = x[PrimitiveStaticTypeUpdateContract-131]
	_ = x[PrimitiveStaticTypeRemoveContract-132]
	_ = x[PrimitiveStaticTypeKeys-133]
	_ = x[PrimitiveStaticTypeAddKey-134]
	_ = x[PrimitiveStaticTypeRevokeKey-135]
	_ = x[PrimitiveStaticTypeInbox-136]
	_ = x[PrimitiveStaticTypePublishInboxCapability-137]
	_ = x[PrimitiveStaticTypeUnpublishInboxCapability-138]
	_ = x[PrimitiveStaticTypeClaimInboxCapability-139]
	_ = x[PrimitiveStaticTypeCapabilities-140]
	_ = x[PrimitiveStaticTypeStorageCapabilities-141]
	_ = x[PrimitiveStaticTypeAccountCapabilities-142]
	_ = x[PrimitiveStaticTypePublishCapability-143]
	_ = x[PrimitiveStaticTypeUnpublishCapability-144]
	_ = x[PrimitiveStaticTypeGetStorageCapabilityController-145]
	_ = x[PrimitiveStaticTypeIssueStorageCapabilityController-146]
	_ = x[PrimitiveStaticTypeGetAccountCapabilityController-147]
	_ = x[PrimitiveStaticTypeIssueAccountCapabilityController-148]
	_ = x[PrimitiveStaticTypeCapabilitiesMapping-149]
	_ = x[PrimitiveStaticTypeAccountMapping-150]
	_ = x[PrimitiveStaticType_Count-151]
}

const _PrimitiveStaticType_name = "UnknownVoidAnyNeverAnyStructAnyResourceBoolAddressStringCharacterMetaTypeBlockAnyResourceAttachmentAnyStructAttachmentNumberSignedNumberIntegerSignedIntegerFixedPointSignedFixedPointIntInt8Int16Int32Int64Int128Int256UIntUInt8UInt16UInt32UInt64UInt128UInt256Word8Word16Word32Word64Word128Word256Fix64UFix64PathCapabilityStoragePathCapabilityPathPublicPathPrivatePathAuthAccountPublicAccountDeployedContractAuthAccountContractsPublicAccountContractsAuthAccountKeysPublicAccountKeysAccountKeyAuthAccountInboxStorageCapabilityControllerAccountCapabilityControllerAuthAccountStorageCapabilitiesAuthAccountAccountCapabilitiesAuthAccountCapabilitiesPublicAccountCapabilitiesAccountAccount_ContractsAccount_KeysAccount_InboxAccount_StorageCapabilitiesAccount_AccountCapabilitiesAccount_CapabilitiesAccount_StorageMutateInsertRemoveIdentityStorageSaveValueLoadValueBorrowValueContractsAddContractUpdateContractRemoveContractKeysAddKeyRevokeKeyInboxPublishInboxCapabilityUnpublishInboxCapabilityClaimInboxCapabilityCapabilitiesStorageCapabilitiesAccountCapabilitiesPublishCapabilityUnpublishCapabilityGetStorageCapabilityControllerIssueStorageCapabilityControllerGetAccountCapabilityControllerIssueAccountCapabilityControllerCapabilitiesMappingAccountMapping_Count"

var _PrimitiveStaticType_map = map[PrimitiveStaticType]string{
	0:   _PrimitiveStaticType_name[0:7],
	1:   _PrimitiveStaticType_name[7:11],
	2:   _PrimitiveStaticType_name[11:14],
	3:   _PrimitiveStaticType_name[14:19],
	4:   _PrimitiveStaticType_name[19:28],
	5:   _PrimitiveStaticType_name[28:39],
	6:   _PrimitiveStaticType_name[39:43],
	7:   _PrimitiveStaticType_name[43:50],
	8:   _PrimitiveStaticType_name[50:56],
	9:   _PrimitiveStaticType_name[56:65],
	10:  _PrimitiveStaticType_name[65:73],
	11:  _PrimitiveStaticType_name[73:78],
	12:  _PrimitiveStaticType_name[78:99],
	13:  _PrimitiveStaticType_name[99:118],
	18:  _PrimitiveStaticType_name[118:124],
	19:  _PrimitiveStaticType_name[124:136],
	24:  _PrimitiveStaticType_name[136:143],
	25:  _PrimitiveStaticType_name[143:156],
	30:  _PrimitiveStaticType_name[156:166],
	31:  _PrimitiveStaticType_name[166:182],
	36:  _PrimitiveStaticType_name[182:185],
	37:  _PrimitiveStaticType_name[185:189],
	38:  _PrimitiveStaticType_name[189:194],
	39:  _PrimitiveStaticType_name[194:199],
	40:  _PrimitiveStaticType_name[199:204],
	41:  _PrimitiveStaticType_name[204:210],
	42:  _PrimitiveStaticType_name[210:216],
	44:  _PrimitiveStaticType_name[216:220],
	45:  _PrimitiveStaticType_name[220:225],
	46:  _PrimitiveStaticType_name[225:231],
	47:  _PrimitiveStaticType_name[231:237],
	48:  _PrimitiveStaticType_name[237:243],
	49:  _PrimitiveStaticType_name[243:250],
	50:  _PrimitiveStaticType_name[250:257],
	53:  _PrimitiveStaticType_name[257:262],
	54:  _PrimitiveStaticType_name[262:268],
	55:  _PrimitiveStaticType_name[268:274],
	56:  _PrimitiveStaticType_name[274:280],
	57:  _PrimitiveStaticType_name[280:287],
	58:  _PrimitiveStaticType_name[287:294],
	64:  _PrimitiveStaticType_name[294:299],
	72:  _PrimitiveStaticType_name[299:305],
	76:  _PrimitiveStaticType_name[305:309],
	77:  _PrimitiveStaticType_name[309:319],
	78:  _PrimitiveStaticType_name[319:330],
	79:  _PrimitiveStaticType_name[330:344],
	80:  _PrimitiveStaticType_name[344:354],
	81:  _PrimitiveStaticType_name[354:365],
	90:  _PrimitiveStaticType_name[365:376],
	91:  _PrimitiveStaticType_name[376:389],
	92:  _PrimitiveStaticType_name[389:405],
	93:  _PrimitiveStaticType_name[405:425],
	94:  _PrimitiveStaticType_name[425:447],
	95:  _PrimitiveStaticType_name[447:462],
	96:  _PrimitiveStaticType_name[462:479],
	97:  _PrimitiveStaticType_name[479:489],
	98:  _PrimitiveStaticType_name[489:505],
	99:  _PrimitiveStaticType_name[505:532],
	100: _PrimitiveStaticType_name[532:559],
	101: _PrimitiveStaticType_name[559:589],
	102: _PrimitiveStaticType_name[589:619],
	103: _PrimitiveStaticType_name[619:642],
	104: _PrimitiveStaticType_name[642:667],
	105: _PrimitiveStaticType_name[667:674],
	106: _PrimitiveStaticType_name[674:691],
	107: _PrimitiveStaticType_name[691:703],
	108: _PrimitiveStaticType_name[703:716],
	109: _PrimitiveStaticType_name[716:743],
	110: _PrimitiveStaticType_name[743:770],
	111: _PrimitiveStaticType_name[770:790],
	112: _PrimitiveStaticType_name[790:805],
	118: _PrimitiveStaticType_name[805:811],
	119: _PrimitiveStaticType_name[811:817],
	120: _PrimitiveStaticType_name[817:823],
	121: _PrimitiveStaticType_name[823:831],
	125: _PrimitiveStaticType_name[831:838],
	126: _PrimitiveStaticType_name[838:847],
	127: _PrimitiveStaticType_name[847:856],
	128: _PrimitiveStaticType_name[856:867],
	129: _PrimitiveStaticType_name[867:876],
	130: _PrimitiveStaticType_name[876:887],
	131: _PrimitiveStaticType_name[887:901],
	132: _PrimitiveStaticType_name[901:915],
	133: _PrimitiveStaticType_name[915:919],
	134: _PrimitiveStaticType_name[919:925],
	135: _PrimitiveStaticType_name[925:934],
	136: _PrimitiveStaticType_name[934:939],
	137: _PrimitiveStaticType_name[939:961],
	138: _PrimitiveStaticType_name[961:985],
	139: _PrimitiveStaticType_name[985:1005],
	140: _PrimitiveStaticType_name[1005:1017],
	141: _PrimitiveStaticType_name[1017:1036],
	142: _PrimitiveStaticType_name[1036:1055],
	143: _PrimitiveStaticType_name[1055:1072],
	144: _PrimitiveStaticType_name[1072:1091],
	145: _PrimitiveStaticType_name[1091:1121],
	146: _PrimitiveStaticType_name[1121:1153],
	147: _PrimitiveStaticType_name[1153:1183],
	148: _PrimitiveStaticType_name[1183:1215],
	149: _PrimitiveStaticType_name[1215:1234],
	150: _PrimitiveStaticType_name[1234:1248],
	151: _PrimitiveStaticType_name[1248:1254],
}

func (i PrimitiveStaticType) String() string {
	if str, ok := _PrimitiveStaticType_map[i]; ok {
		return str
	}
	return "PrimitiveStaticType(" + strconv.FormatInt(int64(i), 10) + ")"
}
