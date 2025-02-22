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
	_ = x[PrimitiveStaticTypeHashableStruct-12]
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
	_ = x[PrimitiveStaticType_Count-105]
}

const _PrimitiveStaticType_name = "UnknownVoidAnyNeverAnyStructAnyResourceBoolAddressStringCharacterMetaTypeBlockHashableStructNumberSignedNumberIntegerSignedIntegerFixedPointSignedFixedPointIntInt8Int16Int32Int64Int128Int256UIntUInt8UInt16UInt32UInt64UInt128UInt256Word8Word16Word32Word64Word128Word256Fix64UFix64PathCapabilityStoragePathCapabilityPathPublicPathPrivatePathAuthAccountPublicAccountDeployedContractAuthAccountContractsPublicAccountContractsAuthAccountKeysPublicAccountKeysAccountKeyAuthAccountInboxStorageCapabilityControllerAccountCapabilityControllerAuthAccountStorageCapabilitiesAuthAccountAccountCapabilitiesAuthAccountCapabilitiesPublicAccountCapabilities_Count"

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
	12:  _PrimitiveStaticType_name[78:92],
	18:  _PrimitiveStaticType_name[92:98],
	19:  _PrimitiveStaticType_name[98:110],
	24:  _PrimitiveStaticType_name[110:117],
	25:  _PrimitiveStaticType_name[117:130],
	30:  _PrimitiveStaticType_name[130:140],
	31:  _PrimitiveStaticType_name[140:156],
	36:  _PrimitiveStaticType_name[156:159],
	37:  _PrimitiveStaticType_name[159:163],
	38:  _PrimitiveStaticType_name[163:168],
	39:  _PrimitiveStaticType_name[168:173],
	40:  _PrimitiveStaticType_name[173:178],
	41:  _PrimitiveStaticType_name[178:184],
	42:  _PrimitiveStaticType_name[184:190],
	44:  _PrimitiveStaticType_name[190:194],
	45:  _PrimitiveStaticType_name[194:199],
	46:  _PrimitiveStaticType_name[199:205],
	47:  _PrimitiveStaticType_name[205:211],
	48:  _PrimitiveStaticType_name[211:217],
	49:  _PrimitiveStaticType_name[217:224],
	50:  _PrimitiveStaticType_name[224:231],
	53:  _PrimitiveStaticType_name[231:236],
	54:  _PrimitiveStaticType_name[236:242],
	55:  _PrimitiveStaticType_name[242:248],
	56:  _PrimitiveStaticType_name[248:254],
	57:  _PrimitiveStaticType_name[254:261],
	58:  _PrimitiveStaticType_name[261:268],
	64:  _PrimitiveStaticType_name[268:273],
	72:  _PrimitiveStaticType_name[273:279],
	76:  _PrimitiveStaticType_name[279:283],
	77:  _PrimitiveStaticType_name[283:293],
	78:  _PrimitiveStaticType_name[293:304],
	79:  _PrimitiveStaticType_name[304:318],
	80:  _PrimitiveStaticType_name[318:328],
	81:  _PrimitiveStaticType_name[328:339],
	90:  _PrimitiveStaticType_name[339:350],
	91:  _PrimitiveStaticType_name[350:363],
	92:  _PrimitiveStaticType_name[363:379],
	93:  _PrimitiveStaticType_name[379:399],
	94:  _PrimitiveStaticType_name[399:421],
	95:  _PrimitiveStaticType_name[421:436],
	96:  _PrimitiveStaticType_name[436:453],
	97:  _PrimitiveStaticType_name[453:463],
	98:  _PrimitiveStaticType_name[463:479],
	99:  _PrimitiveStaticType_name[479:506],
	100: _PrimitiveStaticType_name[506:533],
	101: _PrimitiveStaticType_name[533:563],
	102: _PrimitiveStaticType_name[563:593],
	103: _PrimitiveStaticType_name[593:616],
	104: _PrimitiveStaticType_name[616:641],
	105: _PrimitiveStaticType_name[641:647],
}

func (i PrimitiveStaticType) String() string {
	if str, ok := _PrimitiveStaticType_map[i]; ok {
		return str
	}
	return "PrimitiveStaticType(" + strconv.FormatInt(int64(i), 10) + ")"
}
