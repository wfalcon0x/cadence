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

package ccf_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/encoding/ccf"
	"github.com/onflow/cadence/runtime"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/interpreter"
	"github.com/onflow/cadence/runtime/sema"
	"github.com/onflow/cadence/runtime/tests/checker"
	"github.com/onflow/cadence/runtime/tests/utils"
)

type encodeTest struct {
	name     string
	val      cadence.Value
	expected []byte
}

func TestEncodeVoid(t *testing.T) {

	t.Parallel()

	testEncodeAndDecode(
		t,
		cadence.NewVoid(),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Void"}
			//
			// language=edn, format=ccf
			// 130([137(50), null])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// void type ID (50)
			0x18, 0x32,
			// nil
			0xf6,
		},
	)
}

func TestEncodeOptional(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Optional(nil)",
			cadence.NewOptional(nil),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Optional","value":null}
				//
				// language=edn, format=ccf
				// 130([138(137(42)), null])
				//
				// language=cbor, format=ccf, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagOptionalType,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// never type ID (42)
				0x18, 0x2a,
				// nil
				0xf6,
			},
		},
		{
			"Optional(non-nil)",
			cadence.NewOptional(cadence.NewInt(42)),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Optional","value":{"type":"Int","value":"42"}}
				//
				// language=edn, format=ccf
				// 130([138(137(4)), 42])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagOptionalType,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int type ID (4)
				0x04,
				// tag (big number)
				0xc2,
				// bytes, 1 byte follow
				0x41,
				// 42
				0x2a,
			},
		},
		{
			"Optional(Optional(nil))",
			cadence.NewOptional(cadence.NewOptional(nil)),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"value":{"value":null,"type":"Optional"},"type":"Optional"}
				//
				// language=edn, format=ccf
				// 130([138(138(137(42))), null])
				//
				// language=cbor, format=ccf, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagOptionalType,
				// tag
				0xd8, ccf.CBORTagOptionalType,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// never type ID (42)
				0x18, 0x2a,
				// nil
				0xf6,
			},
		},
		{
			"Optional(Optional(non-nil))",
			cadence.NewOptional(cadence.NewOptional(cadence.NewInt(42))),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"value":{"value":{"value":"42","type":"Int"},"type":"Optional"},"type":"Optional"}
				//
				// language=edn, format=ccf
				// 130([138(138(137(4))), 42])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagOptionalType,
				// tag
				0xd8, ccf.CBORTagOptionalType,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int type ID (4)
				0x04,
				// tag (big number)
				0xc2,
				// bytes, 1 byte follow
				0x41,
				// 42
				0x2a,
			},
		},
		{
			"Optional(Optional(Optional(nil)))",
			cadence.NewOptional(cadence.NewOptional(cadence.NewOptional(nil))),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"value":{"value":{"value":null,"type":"Optional"},"type":"Optional"},"type":"Optional"}
				//
				// language=edn, format=ccf
				// 130([138(138(138(137(42)))), null])
				//
				// language=cbor, format=ccf, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagOptionalType,
				// tag
				0xd8, ccf.CBORTagOptionalType,
				// tag
				0xd8, ccf.CBORTagOptionalType,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// never type ID (42)
				0x18, 0x2a,
				// nil
				0xf6,
			},
		},
		{
			"Optional(Optional(Optional(non-nil)))",
			cadence.NewOptional(cadence.NewOptional(cadence.NewOptional(cadence.NewInt(42)))),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"value":{"value":{"value":{"value":"42","type":"Int"},"type":"Optional"},"type":"Optional"},"type":"Optional"}
				//
				// language=edn, format=ccf
				// 130([138(138(138(137(4)))), 42])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagOptionalType,
				// tag
				0xd8, ccf.CBORTagOptionalType,
				// tag
				0xd8, ccf.CBORTagOptionalType,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int type ID (4)
				0x04,
				// tag (big number)
				0xc2,
				// bytes, 1 byte follow
				0x41,
				// 42
				0x2a,
			},
		},
	}...)
}

func TestEncodeBool(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"True",
			cadence.NewBool(true),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Bool","value":true}
				//
				// language=edn, format=ccf
				// 130([137(0), true])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Bool type ID (0)
				0x00,
				// true
				0xf5,
			},
		},
		{
			"False",
			cadence.NewBool(false),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Bool","value":false}
				//
				// language=edn, format=ccf
				// 130([137(0), false])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Bool type ID (0)
				0x00,
				// false
				0xf4,
			},
		},
	}...)
}

func TestEncodeCharacter(t *testing.T) {

	t.Parallel()

	a, _ := cadence.NewCharacter("a")
	b, _ := cadence.NewCharacter("b")

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"a",
			a,
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Character","value":"a"}
				//
				// language=edn, format=ccf
				// 130([137(2), "a"])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Character type ID (2)
				0x02,
				// UTF-8 string, 1 bytes follow
				0x61,
				// a
				0x61,
			},
		},
		{
			"b",
			b,
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Character","value":"b"}
				//
				// language=edn, format=ccf
				// 130([137(2), "b"])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Character type ID (2)
				0x02,
				// UTF-8 string, 1 bytes follow
				0x61,
				// b
				0x62,
			},
		},
	}...)
}

func TestEncodeString(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Empty",
			cadence.String(""),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"String","value":""}
				//
				// language=edn, format=ccf
				// 130([137(1), ""])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// String type ID (1)
				0x01,
				// UTF-8 string, 0 bytes follow
				0x60,
			},
		},
		{
			"Non-empty",
			cadence.String("foo"),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"String","value":"foo"}
				//
				// language=edn, format=ccf
				// 130([137(1), "foo"])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// String type ID (1)
				0x01,
				// UTF-8 string, 3 bytes follow
				0x63,
				// f, o, o
				0x66, 0x6f, 0x6f,
			},
		},
	}...)
}

func TestEncodeAddress(t *testing.T) {

	t.Parallel()

	testEncodeAndDecode(
		t,
		cadence.BytesToAddress([]byte{1, 2, 3, 4, 5}),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Address","value":"0x0000000102030405"}
			//
			// language=edn, format=ccf
			// 130([137(3), h'0000000102030405'])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Address type ID (3)
			0x03,
			// bytes, 8 bytes follow
			0x48,
			// 0, 0, 0, 1, 2, 3, 4, 5
			0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x5,
		},
	)
}

func TestEncodeInt(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Negative",
			cadence.NewInt(-42),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int","value":"-42"}
				//
				// language=edn, format=ccf
				// 130([137(4), -42])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int type ID (4)
				0x04,
				// tag (negative big number)
				0xc3,
				// bytes, 1 byte follow
				0x41,
				// -42
				0x29,
			},
		},
		{
			"Zero",
			cadence.NewInt(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(4), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int type ID (4)
				0x04,
				// tag (positive big number)
				0xc2,
				// bytes, 0 byte follow
				0x40,
			},
		},
		{
			"Positive",
			cadence.NewInt(42),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int","value":"42"}
				//
				// language=edn, format=ccf
				// 130([137(4), 42])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int type ID (4)
				0x04,
				// tag (positive big number)
				0xc2,
				// bytes, 1 byte follow
				0x41,
				// 42
				0x2a,
			},
		},
		{
			"SmallerThanMinInt256",
			cadence.NewIntFromBig(new(big.Int).Sub(sema.Int256TypeMinIntBig, big.NewInt(10))),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int","value":"-57896044618658097711785492504343953926634992332820282019728792003956564819978"}
				//
				// language=edn, format=ccf
				// 130([137(4), -57896044618658097711785492504343953926634992332820282019728792003956564819978])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int type ID (4)
				0x04,
				// tag (negative big number)
				0xc3,
				// bytes, 32 bytes follow
				0x58, 0x20,
				// -57896044618658097711785492504343953926634992332820282019728792003956564819978
				0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09,
			},
		},
		{
			"LargerThanMaxUInt256",
			cadence.NewIntFromBig(new(big.Int).Add(sema.UInt256TypeMaxIntBig, big.NewInt(10))),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int","value":"115792089237316195423570985008687907853269984665640564039457584007913129639945"}
				//
				// language=edn, format=ccf
				// 130([137(4), 115792089237316195423570985008687907853269984665640564039457584007913129639945])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int type ID (4)
				0x04,
				// tag (positive big number)
				0xc2,
				// bytes, 33 bytes follow
				0x58, 0x21,
				// 115792089237316195423570985008687907853269984665640564039457584007913129639945
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x09,
			},
		},
	}...)
}

func TestEncodeInt8(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Min",
			cadence.NewInt8(math.MinInt8),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int8","value":"-128"}
				//
				// language=edn, format=ccf
				// 130([137(5), -128])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int8 type ID (5)
				0x05,
				// -128
				0x38, 0x7f,
			},
		},
		{
			"Zero",
			cadence.NewInt8(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int8","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(5), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int8 type ID (5)
				0x05,
				// 0
				0x00,
			},
		},
		{
			"Max",
			cadence.NewInt8(math.MaxInt8),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int8","value":"127"}
				//
				// language=edn, format=ccf
				// 130([137(5), 127])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int8 type ID (5)
				0x05,
				// 127
				0x18, 0x7f,
			},
		},
	}...)
}

func TestEncodeInt16(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Min",
			cadence.NewInt16(math.MinInt16),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int16","value":"-32768"}
				//
				// language=edn, format=ccf
				// 130([137(6), -32768])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int16 type ID (6)
				0x06,
				// -32768
				0x39, 0x7F, 0xFF,
			},
		},
		{
			"Zero",
			cadence.NewInt16(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int16","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(6), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int16 type ID (6)
				0x06,
				// 0
				0x00,
			},
		},
		{
			"Max",
			cadence.NewInt16(math.MaxInt16),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int16","value":"32767"}
				//
				// language=edn, format=ccf
				// 130([137(6), 32767])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int16 type ID (6)
				0x06,
				// 32767
				0x19, 0x7F, 0xFF,
			},
		},
	}...)
}

func TestEncodeInt32(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Min",
			cadence.NewInt32(math.MinInt32),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int32","value":"-2147483648"}
				//
				// language=edn, format=ccf
				// 130([137(7), -2147483648])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int32 type ID (7)
				0x07,
				// -2147483648
				0x3a, 0x7f, 0xff, 0xff, 0xff,
			},
		},
		{
			"Zero",
			cadence.NewInt32(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int32","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(7), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int32 type ID (7)
				0x07,
				// 0
				0x00,
			},
		},
		{
			"Max",
			cadence.NewInt32(math.MaxInt32),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int32","value":"2147483647"}
				//
				// language=edn, format=ccf
				// 130([137(7), 2147483647])
				//
				// language=cbor, format=ccf, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int32 type ID (7)
				0x07,
				// 2147483647
				0x1a, 0x7f, 0xff, 0xff, 0xff,
			},
		},
	}...)
}

func TestEncodeInt64(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Min",
			cadence.NewInt64(math.MinInt64),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int64","value":"-9223372036854775808"}
				//
				// language=edn, format=ccf
				// 130([137(8), -9223372036854775808])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int64 type ID (8)
				0x08,
				// -9223372036854775808
				0x3b, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			},
		},
		{
			"Zero",
			cadence.NewInt64(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int64","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(8), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int64 type ID (8)
				0x08,
				// 0
				0x00,
			},
		},
		{
			"Max",
			cadence.NewInt64(math.MaxInt64),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int64","value":"9223372036854775807"}
				//
				// language=edn, format=ccf
				// 130([137(8), 9223372036854775807])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int64 type ID (8)
				0x08,
				// 9223372036854775807
				0x1b, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			},
		},
	}...)
}

func TestEncodeInt128(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Min",
			cadence.Int128{Value: sema.Int128TypeMinIntBig},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int128","value":"-170141183460469231731687303715884105728"}
				//
				// language=edn, format=ccf
				// 130([137(9), -170141183460469231731687303715884105728])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int128 type ID (9)
				0x09,
				// tag big num
				0xc3,
				// bytes, 16 bytes follow
				0x50,
				// -170141183460469231731687303715884105728
				0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			},
		},
		{
			"Zero",
			cadence.NewInt128(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int128","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(9), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int128 type ID (9)
				0x09,
				// tag big num
				0xc2,
				// bytes, 0 bytes follow
				0x40,
			},
		},
		{
			"Max",
			cadence.Int128{Value: sema.Int128TypeMaxIntBig},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int128","value":"170141183460469231731687303715884105727"}
				//
				// language=edn, format=ccf
				// 130([137(9), 170141183460469231731687303715884105727])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int128 type ID (9)
				0x09,
				// tag big num
				0xc2,
				// bytes, 16 bytes follow
				0x50,
				// 170141183460469231731687303715884105727
				0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			},
		},
	}...)
}

func TestEncodeInt256(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Min",
			cadence.Int256{Value: sema.Int256TypeMinIntBig},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int256","value":"-57896044618658097711785492504343953926634992332820282019728792003956564819968"}
				//
				// language=edn, format=ccf
				// 130([137(10), -57896044618658097711785492504343953926634992332820282019728792003956564819968])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int256 type ID (10)
				0x0a,
				// tag big num
				0xc3,
				// bytes, 32 bytes follow
				0x58, 0x20,
				// -57896044618658097711785492504343953926634992332820282019728792003956564819968
				0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			},
		},
		{
			"Zero",
			cadence.NewInt256(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int256","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(10), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int256 type ID (10)
				0x0a,
				// tag big num
				0xc2,
				// bytes, 0 bytes follow
				0x40,
			},
		},
		{
			"Max",
			cadence.Int256{Value: sema.Int256TypeMaxIntBig},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Int256","value":"57896044618658097711785492504343953926634992332820282019728792003956564819967"}
				//
				// language=edn, format=ccf
				// 130([137(10), 57896044618658097711785492504343953926634992332820282019728792003956564819967])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Int256 type ID (10)
				0x0a,
				// tag big num
				0xc2,
				// bytes, 32 bytes follow
				0x58, 0x20,
				// 57896044618658097711785492504343953926634992332820282019728792003956564819967
				0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			},
		},
	}...)
}

func TestEncodeUInt(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Zero",
			cadence.NewUInt(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UInt","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(11), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UInt type ID (11)
				0x0b,
				// tag big num
				0xc2,
				// bytes, 0 bytes follow
				0x40,
			},
		},
		{
			"Positive",
			cadence.NewUInt(42),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UInt","value":"42"}
				//
				// language=edn, format=ccf
				// 130([137(11), 42])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UInt type ID (11)
				0x0b,
				// tag big num
				0xc2,
				// bytes, 1 bytes follow
				0x41,
				// 42
				0x2a,
			},
		},
		{
			"LargerThanMaxUInt256",
			cadence.UInt{Value: new(big.Int).Add(sema.UInt256TypeMaxIntBig, big.NewInt(10))},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UInt","value":"115792089237316195423570985008687907853269984665640564039457584007913129639945"}
				//
				// language=edn, format=ccf
				// 130([137(11), 115792089237316195423570985008687907853269984665640564039457584007913129639945])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UInt type ID (11)
				0x0b,
				// tag big num
				0xc2,
				// bytes, 32 bytes follow
				0x58, 0x21,
				// 115792089237316195423570985008687907853269984665640564039457584007913129639945
				0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x09,
			},
		},
	}...)
}

func TestEncodeUInt8(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Zero",
			cadence.NewUInt8(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UInt8","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(12), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UInt8 type ID (12)
				0x0c,
				// 0
				0x00,
			},
		},
		{
			"Max",
			cadence.NewUInt8(math.MaxUint8),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UInt8","value":"255"}
				//
				// language=edn, format=ccf
				// 130([137(12), 255])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UInt8 type ID (12)
				0x0c,
				// 255
				0x18, 0xff,
			},
		},
	}...)
}

func TestEncodeUInt16(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Zero",
			cadence.NewUInt16(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UInt16","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(13), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UInt16 type ID (13)
				0x0d,
				// 0
				0x00,
			},
		},
		{
			"Max",
			cadence.NewUInt16(math.MaxUint16),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UInt16","value":"65535"}
				//
				// language=edn, format=ccf
				// 130([137(13), 65535])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UInt16 type ID (13)
				0x0d,
				// 65535
				0x19, 0xff, 0xff,
			},
		},
	}...)
}

func TestEncodeUInt32(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Zero",
			cadence.NewUInt32(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UInt32","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(14), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UInt32 type ID (14)
				0x0e,
				// 0
				0x00,
			},
		},
		{
			"Max",
			cadence.NewUInt32(math.MaxUint32),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UInt32","value":"4294967295"}
				//
				// language=edn, format=ccf
				// 130([137(14), 4294967295])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UInt32 type ID (14)
				0x0e,
				// 4294967295
				0x1a, 0xff, 0xff, 0xff, 0xff,
			},
		},
	}...)
}

func TestEncodeUInt64(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Zero",
			cadence.NewUInt64(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UInt64","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(15), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UInt64 type ID (15)
				0x0f,
				// 0
				0x00,
			},
		},
		{
			"Max",
			cadence.NewUInt64(uint64(math.MaxUint64)),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UInt64","value":"18446744073709551615"}
				//
				// language=edn, format=ccf
				// 130([137(15), 18446744073709551615])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UInt64 type ID (15)
				0x0f,
				// 18446744073709551615
				0x1b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			},
		},
	}...)
}

func TestEncodeUInt128(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Zero",
			cadence.NewUInt128(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UInt128","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(16), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UInt128 type ID (16)
				0x10,
				// tag (big num)
				0xc2,
				// bytes, 0 bytes follow
				0x40,
			},
		},
		{
			"Max",
			cadence.UInt128{Value: sema.UInt128TypeMaxIntBig},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UInt128","value":"340282366920938463463374607431768211455"}
				//
				// language=edn, format=ccf
				// 130([137(16), 340282366920938463463374607431768211455])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UInt128 type ID (16)
				0x10,
				// tag (big num)
				0xc2,
				// bytes, 16 bytes follow
				0x50,
				// 340282366920938463463374607431768211455
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			},
		},
	}...)
}

func TestEncodeUInt256(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Zero",
			cadence.NewUInt256(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UInt256","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(17), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UInt256 type ID (17)
				0x11,
				// tag (big num)
				0xc2,
				// bytes, 0 bytes follow
				0x40,
			},
		},
		{
			"Max",
			cadence.UInt256{Value: sema.UInt256TypeMaxIntBig},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UInt256","value":"115792089237316195423570985008687907853269984665640564039457584007913129639935"}
				//
				// language=edn, format=ccf
				// 130([137(17), 115792089237316195423570985008687907853269984665640564039457584007913129639935])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UInt256 type ID (17)
				0x11,
				// tag (big num)
				0xc2,
				// bytes, 32 bytes follow
				0x58, 0x20,
				// 115792089237316195423570985008687907853269984665640564039457584007913129639935
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			},
		},
	}...)
}

func TestEncodeWord8(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Zero",
			cadence.NewWord8(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Word8","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(18), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Word8 type ID (18)
				0x12,
				// 0
				0x00,
			},
		},
		{
			"Max",
			cadence.NewWord8(math.MaxUint8),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Word8","value":"255"}
				//
				// language=edn, format=ccf
				// 130([137(18), 255])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Word8 type ID (18)
				0x12,
				// 255
				0x18, 0xff,
			},
		},
	}...)
}

func TestEncodeWord16(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Zero",
			cadence.NewWord16(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Word16","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(19), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Word16 type ID (19)
				0x13,
				// 0
				0x00,
			},
		},
		{
			"Max",
			cadence.NewWord16(math.MaxUint16),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Word16","value":"65535"}
				//
				// language=edn, format=ccf
				// 130([137(19), 65535])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Word16 type ID (19)
				0x13,
				// 65535
				0x19, 0xff, 0xff,
			},
		},
	}...)
}

func TestEncodeWord32(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Zero",
			cadence.NewWord32(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Word32","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(20), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Word32 type ID (20)
				0x14,
				// 0
				0x00,
			},
		},
		{
			"Max",
			cadence.NewWord32(math.MaxUint32),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Word32","value":"4294967295"}
				//
				// language=edn, format=ccf
				// 130([137(20), 4294967295])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Word32 type ID (20)
				0x14,
				// 4294967295
				0x1a, 0xff, 0xff, 0xff, 0xff,
			},
		},
	}...)
}

func TestEncodeWord64(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Zero",
			cadence.NewWord64(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Word64","value":"0"}
				//
				// language=edn, format=ccf
				// 130([137(21), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Word64 type ID (21)
				0x15,
				// 0
				0x00,
			},
		},
		{
			"Max",
			cadence.NewWord64(math.MaxUint64),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Word64","value":"18446744073709551615"}
				//
				// language=edn, format=ccf
				// 130([137(21), 18446744073709551615])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Word64 type ID (21)
				0x15,
				// 18446744073709551615
				0x1b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			},
		},
	}...)
}

func TestEncodeFix64(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Zero",
			cadence.Fix64(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Fix64","value":"0.00000000"}
				//
				// language=edn, format=ccf
				// 130([137(22), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Fix64 type ID (22)
				0x16,
				// 0
				0x00,
			},
		},
		{
			"789.00123010",
			cadence.Fix64(78_900_123_010),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Fix64","value":"789.00123010"}
				//
				// language=edn, format=ccf
				// 130([137(22), 78900123010])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Fix64 type ID (22)
				0x16,
				// 78900123010
				0x1b, 0x00, 0x00, 0x00, 0x12, 0x5e, 0xd0, 0x55, 0x82,
			},
		},
		{
			"1234.056",
			cadence.Fix64(123_405_600_000),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Fix64","value":"1234.05600000"}
				//
				// language=edn, format=ccf
				// 130([137(22), 123405600000])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Fix64 type ID (22)
				0x16,
				// 123405600000
				0x1b, 0x00, 0x00, 0x00, 0x1c, 0xbb, 0x8c, 0x05, 0x00,
			},
		},
		{
			"-12345.006789",
			cadence.Fix64(-1_234_500_678_900),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Fix64","value":"-12345.00678900"}
				//
				// language=edn, format=ccf
				// 130([137(22), -1234500678900])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Fix64 type ID (22)
				0x16,
				// -1234500678900
				0x3b, 0x00, 0x00, 0x01, 0x1f, 0x6d, 0xf9, 0x74, 0xf3,
			},
		},
	}...)
}

func TestEncodeUFix64(t *testing.T) {

	t.Parallel()

	testAllEncodeAndDecode(t, []encodeTest{
		{
			"Zero",
			cadence.UFix64(0),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UFix64","value":"0.00000000"}
				//
				// language=edn, format=ccf
				// 130([137(23), 0])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,
				// 0
				0x00,
			},
		},
		{
			"789.00123010",
			cadence.UFix64(78_900_123_010),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UFix64","value":"789.00123010"}
				//
				// language=edn, format=ccf
				// 130([137(23), 78900123010])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,
				// 78900123010
				0x1b, 0x00, 0x00, 0x00, 0x12, 0x5e, 0xd0, 0x55, 0x82,
			},
		},
		{
			"1234.056",
			cadence.UFix64(123_405_600_000),
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"UFix64","value":"1234.05600000"}
				//
				// language=edn, format=ccf
				// 130([137(23), 123405600000])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,
				// 123405600000
				0x1b, 0x00, 0x00, 0x00, 0x1c, 0xbb, 0x8c, 0x05, 0x00,
			},
		},
	}...)
}

func TestEncodeArray(t *testing.T) {

	t.Parallel()

	// []
	emptyArray := encodeTest{
		"Empty",
		cadence.NewArray(
			[]cadence.Value{},
		).WithType(cadence.NewVariableSizedArrayType(cadence.NewIntType())),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Array","value":[]}
			//
			// language=edn, format=ccf
			// 130([139(137(4)), []])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// type []int
			// tag
			0xd8, ccf.CBORTagVarsizedArrayType,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// array data without inlined type
			// array, 0 items follow
			0x80,
		},
	}

	// [1, 2, 3]
	intArray := encodeTest{
		"Integers",
		cadence.NewArray([]cadence.Value{
			cadence.NewInt(1),
			cadence.NewInt(2),
			cadence.NewInt(3),
		}).WithType(cadence.NewVariableSizedArrayType(cadence.NewIntType())),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Array","value":[{"type":"Int","value":"1"},{"type":"Int","value":"2"},{"type":"Int","value":"3"}]}
			//
			// language=edn, format=ccf
			// 130([139(137(4)), [1, 2, 3]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// type []int
			// tag
			0xd8, ccf.CBORTagVarsizedArrayType,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// array data without inlined type definition
			// array, 3 items follow
			0x83,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 2
			0x02,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 3
			0x03,
		},
	}

	// [S.test.Foo{1}, S.test.Foo{2}, S.test.Foo{3}]
	resourceArray := encodeTest{
		"Resources",
		cadence.NewArray([]cadence.Value{
			cadence.NewResource([]cadence.Value{
				cadence.NewInt(1),
			}).WithType(fooResourceType),
			cadence.NewResource([]cadence.Value{
				cadence.NewInt(2),
			}).WithType(fooResourceType),
			cadence.NewResource([]cadence.Value{
				cadence.NewInt(3),
			}).WithType(fooResourceType),
		}).WithType(cadence.NewVariableSizedArrayType(fooResourceType)),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Array","value":[{"type":"Resource","value":{"id":"S.test.Foo","fields":[{"name":"bar","value":{"type":"Int","value":"1"}}]}},{"type":"Resource","value":{"id":"S.test.Foo","fields":[{"name":"bar","value":{"type":"Int","value":"2"}}]}},{"type":"Resource","value":{"id":"S.test.Foo","fields":[{"name":"bar","value":{"type":"Int","value":"3"}}]}}]}
			//
			// language=edn, format=ccf
			// 129([[161([h'', "S.test.Foo", [["bar", 137(4)]]])], [139(136(h'')), [[1], [2], [3]]]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definitions
			// array, 1 items follow
			0x81,
			// resource type:
			// id: []byte{}
			// cadence-type-id: "S.test.Foo"
			// fields: [["bar", int type]]
			// tag
			0xd8, ccf.CBORTagResourceType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 10 bytes follow
			0x6a,
			// S.test.Foo
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f,
			// fields
			// array, 1 items follow
			0x81,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 3 bytes follow
			0x63,
			// bar
			0x62, 0x61, 0x72,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagVarsizedArrayType,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 0 bytes follow
			0x40,
			// [S.test.Foo{1}, S.test.Foo{2}, S.test.Foo{3}]
			// array, 3 items follow
			0x83,
			// array, 1 items follow
			0x81,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// array, 1 items follow
			0x81,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 2
			0x02,
			// array, 1 items follow
			0x81,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 3
			0x03,
		},
	}

	s, err := cadence.NewString("a")
	require.NoError(t, err)

	resourceWithAbstractFieldArray := encodeTest{
		"Resources with abstract field",
		cadence.NewArray([]cadence.Value{
			cadence.NewResource([]cadence.Value{
				cadence.NewInt(1),
				cadence.NewInt(1),
			}).WithType(foooResourceTypeWithAbstractField),
			cadence.NewResource([]cadence.Value{
				cadence.NewInt(2),
				s,
			}).WithType(foooResourceTypeWithAbstractField),
			cadence.NewResource([]cadence.Value{
				cadence.NewInt(3),
				cadence.NewBool(true),
			}).WithType(foooResourceTypeWithAbstractField),
		}).WithType(cadence.NewVariableSizedArrayType(foooResourceTypeWithAbstractField)),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Array","value":[{"type":"Resource","value":{"id":"S.test.Fooo","fields":[{"name":"bar","value":{"type":"Int","value":"1"}},{"name":"baz","value":{"type":"Int","value":"1"}}]}},{"type":"Resource","value":{"id":"S.test.Fooo","fields":[{"name":"bar","value":{"type":"Int","value":"2"}},{"name":"baz","value":{"type":"String","value":"a"}}]}},{"type":"Resource","value":{"id":"S.test.Fooo","fields":[{"name":"bar","value":{"type":"Int","value":"3"}},{"name":"baz","value":{"type":"Bool","value":true}}]}}]}
			//
			// language=edn, format=ccf
			// 129([[161([h'', "S.test.Fooo", [["bar", 137(4)], ["baz", 137(39)]]])], [139(136(h'')), [[1, 130([137(4), 1])], [2, 130([137(1), "a"])], [3, 130([137(0), true])]]]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definitions
			// array, 1 items follow
			0x81,
			// resource type:
			// id: []byte{}
			// cadence-type-id: "S.test.Fooo"
			// fields: [["bar", int type], ["baz", any type]]
			// tag
			0xd8, ccf.CBORTagResourceType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 11 bytes follow
			0x6b,
			// S.test.Fooo
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f, 0x6f,
			// fields
			// array, 2 items follow
			0x82,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 3 bytes follow
			0x63,
			// bar
			0x62, 0x61, 0x72,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// field 1
			// array, 2 items follow
			0x82,
			// text, 3 bytes follow
			0x63,
			// baz
			0x62, 0x61, 0x7a,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// AnyStruct type ID (39)
			0x18, 0x27,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagVarsizedArrayType,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 0 bytes follow
			0x40,
			// [S.test.Foo{1, 1}, S.test.Foo{2, "a"}, S.test.Foo{3, true}]
			// array, 3 items follow
			0x83,
			// element 0
			// array, 2 items follow
			0x82,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// element 1
			// array, 2 items follow
			0x82,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 2
			0x02,
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// String type ID (1)
			0x01,
			// text string, 1 byte
			0x61,
			// "a"
			0x61,
			// element 2
			// array, 2 items follow
			0x82,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 3
			0x03,
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Bool type ID (0)
			0x00,
			// true
			0xf5,
		},
	}

	// [1, "a", true]
	heterogeneousSimpleTypeArray := encodeTest{
		"Heterogenous AnyStruct Array with Simple Values",
		cadence.NewArray([]cadence.Value{
			cadence.NewInt(1),
			s,
			cadence.NewBool(true),
		}).WithType(cadence.NewVariableSizedArrayType(cadence.NewAnyStructType())),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Array","value":[{"type":"Int","value":"1"},{"type":"String","value":"a"},{"type":"Bool","value":true}]}
			//
			// language=edn, format=ccf
			// 130([139(137(39)), [130([137(4), 1]), 130([137(1), "a"]), 130([137(0), true])]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// type ([]AnyStruct)
			// tag
			0xd8, ccf.CBORTagVarsizedArrayType,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// AnyStruct type ID (39)
			0x18, 0x27,
			// array data with inlined type because static array element type is abstract (AnyStruct)
			// array, 3 items follow
			0x83,
			// element 0 (inline type and value)
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// element 1 (inline type and value)
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// String type ID (1)
			0x01,
			// text string, length 1
			0x61,
			// "a"
			0x61,
			// element 2 (inline type and value)
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Bool type ID (0)
			0x00,
			// true
			0xf5,
		},
	}

	// [Int8(1), Int16(2), Int32(3)]
	heterogeneousNumberTypeArray := encodeTest{
		"Heterogeous Number Array",
		cadence.NewArray([]cadence.Value{
			cadence.NewInt8(1),
			cadence.NewInt16(2),
			cadence.NewInt32(3),
		}).WithType(cadence.NewVariableSizedArrayType(cadence.NewNumberType())),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"value":[{"value":"1","type":"Int8"},{"value":"2","type":"Int16"},{"value":"3","type":"Int32"}],"type":"Array"}
			//
			// language=edn, format=ccf
			// 130([139(137(43)), [130([137(5), 1]), 130([137(6), 2]), 130([137(7), 3])]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// type ([]Integer)
			// tag
			0xd8, ccf.CBORTagVarsizedArrayType,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Number type ID (43)
			0x18, 0x2b,
			// array data with inlined type because static array element type is abstract (Number)
			// array, 3 items follow
			0x83,
			// element 0 (inline type and value)
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int8 type ID (5)
			0x05,
			// 1
			0x01,
			// element 1 (inline type and value)
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int16 type ID (6)
			0x06,
			// 2
			0x02,
			// element 2 (inline type and value)
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int32 type ID (7)
			0x07,
			// 3
			0x03,
		},
	}

	// [1, S.test.Foo{1}]
	heterogeneousCompositeTypeArray := encodeTest{
		"Heterogenous AnyStruct Array with Composite Value",
		cadence.NewArray([]cadence.Value{
			cadence.NewInt(1),
			cadence.NewResource([]cadence.Value{
				cadence.NewInt(1),
			}).WithType(fooResourceType),
		}).WithType(cadence.NewVariableSizedArrayType(cadence.NewAnyStructType())),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"value":[{"value":"1","type":"Int"},{"value":{"id":"S.test.Foo","fields":[{"value":{"value":"1","type":"Int"},"name":"bar"}]},"type":"Resource"}],"type":"Array"}
			//
			// language=edn, format=ccf
			// 129([[161([h'', "S.test.Foo", [["bar", 137(4)]]])], [139(137(39)), [130([137(4), 1]), 130([136(h''), [1]])]]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definition
			// array, 1 items follow
			0x81,
			// resource type:
			// id: []byte{}
			// cadence-type-id: "S.test.Foo"
			// fields: [["bar", int type]]
			// tag
			0xd8, ccf.CBORTagResourceType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 10 bytes follow
			0x6a,
			// S.test.Foo
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f,
			// fields
			// array, 1 items follow
			0x81,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 3 bytes follow
			0x63,
			// bar
			0x62, 0x61, 0x72,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// type ([]AnyStruct)
			// tag
			0xd8, ccf.CBORTagVarsizedArrayType,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// AnyStruct type ID (39)
			0x18, 0x27,
			// array data with inlined type because static array element type is abstract (AnyStruct)
			// array, 2 items follow
			0x82,
			// element 0 (inline type and value)
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// element 1 (inline type and value)
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 0 bytes follow
			0x40,
			// S.test.Foo{1}
			// array, 1 items follow
			0x81,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
		},
	}

	resourceInterfaceType := &cadence.ResourceInterfaceType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "Bar",
	}

	// [S.test.Foo{1}, S.test.Fooo{2, "a"}]
	heterogeneousInterfaceTypeArray := encodeTest{
		"Heterogenous Interface Array",
		cadence.NewArray([]cadence.Value{
			cadence.NewResource([]cadence.Value{
				cadence.NewInt(1),
			}).WithType(fooResourceType),
			cadence.NewResource([]cadence.Value{
				cadence.NewInt(2),
				s,
			}).WithType(foooResourceTypeWithAbstractField),
		}).WithType(cadence.NewVariableSizedArrayType(resourceInterfaceType)),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"value":[{"value":{"id":"S.test.Foo","fields":[{"value":{"value":"1","type":"Int"},"name":"bar"}]},"type":"Resource"},{"value":{"id":"S.test.Fooo","fields":[{"value":{"value":"2","type":"Int"},"name":"bar"},{"value":{"value":"a","type":"String"},"name":"baz"}]},"type":"Resource"}],"type":"Array"}
			//
			// language=edn, format=ccf
			// 129([[177([h'', "S.test.Bar"]), 161([h'01', "S.test.Foo", [["bar", 137(4)]]]), 161([h'02', "S.test.Fooo", [["bar", 137(4)], ["baz", 137(39)]]])], [139(136(h'')), [130([136(h'01'), [1]]), 130([136(h'02'), [2, 130([137(1), "a"])]])]]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definitions
			// array, 3 items follow
			0x83,
			// type definition 0
			// resource interface type:
			// id: []byte{}
			// cadence-type-id: "S.test.Bar"
			// tag
			0xd8, ccf.CBORTagResourceInterfaceType,
			// array, 2 items follow
			0x82,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 10 bytes follow
			0x6a,
			// S.test.Boo
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x42, 0x61, 0x72,
			// type definition 1
			// resource type:
			// id: []byte{1}
			// cadence-type-id: "S.test.Foo"
			// fields: [["bar", int type]]
			// tag
			0xd8, ccf.CBORTagResourceType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// cadence-type-id
			// string, 10 bytes follow
			0x6a,
			// S.test.Foo
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f,
			// fields
			// array, 1 items follow
			0x81,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 3 bytes follow
			0x63,
			// bar
			0x62, 0x61, 0x72,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// type definition 2:
			// resource type:
			// id: []byte{}
			// cadence-type-id: "S.test.Fooo"
			// fields: [["bar", int type], ["baz", any type]]
			// tag
			0xd8, ccf.CBORTagResourceType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 1 bytes follow
			0x41,
			// 2
			0x02,
			// cadence-type-id
			// string, 11 bytes follow
			0x6b,
			// S.test.Fooo
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f, 0x6f,
			// fields
			// array, 2 items follow
			0x82,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 3 bytes follow
			0x63,
			// bar
			0x62, 0x61, 0x72,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// field 1
			// array, 2 items follow
			0x82,
			// text, 3 bytes follow
			0x63,
			// baz
			0x62, 0x61, 0x7a,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// AnyStruct type ID (39)
			0x18, 0x27,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagVarsizedArrayType,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 0 bytes follow
			0x40,
			// array, 2 item follow
			0x82,
			// element 0 (inline type and value)
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// S.test.Foo{1}
			// array, 1 items follow
			0x81,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// element 1 (inline type and value)
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 1 bytes follow
			0x41,
			// 2
			0x02,
			// S.test.Fooo{2, "a"}
			// array, 2 items follow
			0x82,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 2
			0x02,
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// String type ID (1)
			0x01,
			// text string, 1 byte
			0x61,
			// "a"
			0x61,
		},
	}

	testAllEncodeAndDecode(t,
		emptyArray,
		intArray,
		resourceArray,
		resourceWithAbstractFieldArray,
		heterogeneousSimpleTypeArray,
		heterogeneousNumberTypeArray,
		heterogeneousCompositeTypeArray,
		heterogeneousInterfaceTypeArray,
	)
}

func TestEncodeDictionary(t *testing.T) {

	t.Parallel()

	emptyDict := encodeTest{
		"empty",
		cadence.NewDictionary(
			[]cadence.KeyValuePair{},
		).WithType(cadence.NewDictionaryType(cadence.NewStringType(), cadence.NewIntType())),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"value":[],"type":"Dictionary"}
			//
			// language=edn, format=ccf
			// 130([141([137(1), 137(4)]), []])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// type (map[string]int)
			// tag
			0xd8, ccf.CBORTagDictType,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// String type ID (1)
			0x01,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// array, 6 items follow
			0x80,
		},
	}

	// {"a":1, "b":2, "c":3}
	simpleDict := encodeTest{
		"Simple",
		cadence.NewDictionary([]cadence.KeyValuePair{
			{
				Key:   cadence.String("a"),
				Value: cadence.NewInt(1),
			},
			{
				Key:   cadence.String("b"),
				Value: cadence.NewInt(2),
			},
			{
				Key:   cadence.String("c"),
				Value: cadence.NewInt(3),
			},
		}).WithType(cadence.NewDictionaryType(cadence.NewStringType(), cadence.NewIntType())),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Dictionary","value":[{"key":{"type":"String","value":"a"},"value":{"type":"Int","value":"1"}},{"key":{"type":"String","value":"b"},"value":{"type":"Int","value":"2"}},{"key":{"type":"String","value":"c"},"value":{"type":"Int","value":"3"}}]}
			//
			// language=edn, format=ccf
			// 130([141([137(1), 137(4)]), ["a", 1, "b", 2, "c", 3]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// type (map[string]int)
			// tag
			0xd8, ccf.CBORTagDictType,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// String type ID (1)
			0x01,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// array data without inlined type definition
			// array, 6 items follow
			0x86,
			// string, 1 bytes follow
			0x61,
			// a
			0x61,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// string, 1 bytes follow
			0x61,
			// b
			0x62,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 2
			0x02,
			// string, 1 bytes follow
			0x61,
			// c
			0x63,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 3
			0x03,
		},
	}

	// {"a":{"1":1}, "b":{"2":2}, "c":{"3:3"}}
	nestedDict := encodeTest{
		"Nested",
		cadence.NewDictionary([]cadence.KeyValuePair{
			{
				Key: cadence.String("a"),
				Value: cadence.NewDictionary([]cadence.KeyValuePair{
					{
						Key:   cadence.String("1"),
						Value: cadence.NewInt(1),
					},
				}).WithType(cadence.NewDictionaryType(cadence.NewStringType(), cadence.NewIntType())),
			},
			{
				Key: cadence.String("b"),
				Value: cadence.NewDictionary([]cadence.KeyValuePair{
					{
						Key:   cadence.String("2"),
						Value: cadence.NewInt(2),
					},
				}).WithType(cadence.NewDictionaryType(cadence.NewStringType(), cadence.NewIntType())),
			},
			{
				Key: cadence.String("c"),
				Value: cadence.NewDictionary([]cadence.KeyValuePair{
					{
						Key:   cadence.String("3"),
						Value: cadence.NewInt(3),
					},
				}).WithType(cadence.NewDictionaryType(cadence.NewStringType(), cadence.NewIntType())),
			},
		}).WithType(cadence.NewDictionaryType(
			cadence.NewStringType(),
			cadence.NewDictionaryType(cadence.NewStringType(), cadence.NewIntType())),
		),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Dictionary","value":[{"key":{"type":"String","value":"a"},"value":{"type":"Dictionary","value":[{"key":{"type":"String","value":"1"},"value":{"type":"Int","value":"1"}}]}},{"key":{"type":"String","value":"b"},"value":{"type":"Dictionary","value":[{"key":{"type":"String","value":"2"},"value":{"type":"Int","value":"2"}}]}},{"key":{"type":"String","value":"c"},"value":{"type":"Dictionary","value":[{"key":{"type":"String","value":"3"},"value":{"type":"Int","value":"3"}}]}}]}
			//
			// language=edn, format=ccf
			// 130([141([137(1), 141([137(1), 137(4)])]), ["a", ["1", 1], "b", ["2", 2], "c", ["3", 3]]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// type (map[string]map[string, int])
			// tag
			0xd8, ccf.CBORTagDictType,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// String type ID (1)
			0x01,
			// tag
			0xd8, ccf.CBORTagDictType,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// String type ID (1)
			0x01,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// array data without inlined type definition
			// array, 6 items follow
			0x86,
			// string, 1 bytes follow
			0x61,
			// a
			0x61,
			// nested dictionary
			// array, 2 items follow
			0x82,
			// string, 1 bytes follow
			0x61,
			// "1"
			0x31,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// string, 1 bytes follow
			0x61,
			// b
			0x62,
			// nested dictionary
			// array, 2 items follow
			0x82,
			// string, 1 bytes follow
			0x61,
			// "2"
			0x32,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 2
			0x02,
			// string, 1 bytes follow
			0x61,
			// c
			0x63,
			// nested dictionary
			// array, 2 items follow
			0x82,
			// string, 1 bytes follow
			0x61,
			// "3"
			0x33,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 3
			0x03,
		},
	}

	// {"a":foo{1}, "b":foo{2}, "c":foo{3}}
	resourceDict := encodeTest{
		"Resources",
		cadence.NewDictionary([]cadence.KeyValuePair{
			{
				Key: cadence.String("a"),
				Value: cadence.NewResource([]cadence.Value{
					cadence.NewInt(1),
				}).WithType(fooResourceType),
			},
			{
				Key: cadence.String("b"),
				Value: cadence.NewResource([]cadence.Value{
					cadence.NewInt(2),
				}).WithType(fooResourceType),
			},
			{
				Key: cadence.String("c"),
				Value: cadence.NewResource([]cadence.Value{
					cadence.NewInt(3),
				}).WithType(fooResourceType),
			},
		}).WithType(cadence.NewDictionaryType(
			cadence.NewStringType(),
			fooResourceType,
		)),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Dictionary","value":[{"key":{"type":"String","value":"a"},"value":{"type":"Resource","value":{"id":"S.test.Foo","fields":[{"name":"bar","value":{"type":"Int","value":"1"}}]}}},{"key":{"type":"String","value":"b"},"value":{"type":"Resource","value":{"id":"S.test.Foo","fields":[{"name":"bar","value":{"type":"Int","value":"2"}}]}}},{"key":{"type":"String","value":"c"},"value":{"type":"Resource","value":{"id":"S.test.Foo","fields":[{"name":"bar","value":{"type":"Int","value":"3"}}]}}}]}
			//
			// language=edn, format=ccf
			// 129([[161([h'', "S.test.Foo", [["bar", 137(4)]]])], [141([137(1), 136(h'')]), ["a", [1], "b", [2], "c", [3]]]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definition
			// array, 1 items follow
			0x81,
			// resource type:
			// id: []byte{}
			// cadence-type-id: "S.test.Foo"
			// fields: [["bar", int type]]
			// tag
			0xd8, ccf.CBORTagResourceType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 10 bytes follow
			0x6a,
			// S.test.Foo
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f,
			// fields
			// array, 1 items follow
			0x81,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 3 bytes follow
			0x63,
			// bar
			0x62, 0x61, 0x72,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagDictType,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// String type ID (1)
			0x01,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 0 bytes follow
			0x40,
			// array, 6 items follow
			0x86,
			// string, 1 bytes follow
			0x61,
			// a
			0x61,
			// array, 1 items follow
			0x81,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// string, 1 bytes follow
			0x61,
			// b
			0x62,
			// array, 1 items follow
			0x81,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 2
			0x02,
			// string, 1 bytes follow
			0x61,
			// c
			0x63,
			// array, 1 items follow
			0x81,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 3
			0x03,
		},
	}

	testAllEncodeAndDecode(t,
		emptyDict,
		simpleDict,
		nestedDict,
		resourceDict,
	)
}

func TestEncodeSortedDictionary(t *testing.T) {
	type testcase struct {
		name         string
		val          cadence.Value
		expectedVal  cadence.Value
		expectedCBOR []byte
	}

	dict := cadence.NewDictionary([]cadence.KeyValuePair{
		{
			Key:   cadence.String("c"),
			Value: cadence.NewInt(3),
		},
		{
			Key:   cadence.String("a"),
			Value: cadence.NewInt(1),
		},
		{
			Key:   cadence.String("b"),
			Value: cadence.NewInt(2),
		},
	}).WithType(cadence.NewDictionaryType(cadence.NewStringType(), cadence.NewIntType()))

	expectedDict := cadence.NewDictionary([]cadence.KeyValuePair{
		{
			Key:   cadence.String("a"),
			Value: cadence.NewInt(1),
		},
		{
			Key:   cadence.String("b"),
			Value: cadence.NewInt(2),
		},
		{
			Key:   cadence.String("c"),
			Value: cadence.NewInt(3),
		},
	}).WithType(cadence.NewDictionaryType(cadence.NewStringType(), cadence.NewIntType()))

	simpleDict := testcase{
		"Simple",
		dict,
		expectedDict,
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Dictionary","value":[{"key":{"type":"String","value":"a"},"value":{"type":"Int","value":"1"}},{"key":{"type":"String","value":"b"},"value":{"type":"Int","value":"2"}},{"key":{"type":"String","value":"c"},"value":{"type":"Int","value":"3"}}]}
			//
			// language=edn, format=ccf
			// 130([141([137(1), 137(4)]), ["a", 1, "b", 2, "c", 3]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 items follow
			0x82,
			// type (map[string]int)
			// tag
			0xd8, ccf.CBORTagDictType,
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// String type ID (1)
			0x01,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// array data without inlined type definition
			// array, 6 items follow
			0x86,
			// string, 1 bytes follow
			0x61,
			// a
			0x61,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// string, 1 bytes follow
			0x61,
			// b
			0x62,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 2
			0x02,
			// string, 1 bytes follow
			0x61,
			// c
			0x63,
			// tag (big num)
			0xc2,
			// bytes, 1 bytes follow
			0x41,
			// 3
			0x03,
		},
	}

	testcases := []testcase{
		simpleDict,
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			actualCBOR := testEncode(t, tc.val, tc.expectedCBOR)
			testDecode(t, actualCBOR, tc.expectedVal)
		})
	}
}

func exportFromScript(t *testing.T, code string) cadence.Value {
	checker, err := checker.ParseAndCheck(t, code)
	require.NoError(t, err)

	var uuid uint64

	inter, err := interpreter.NewInterpreter(
		interpreter.ProgramFromChecker(checker),
		checker.Location,
		&interpreter.Config{
			UUIDHandler: func() (uint64, error) {
				uuid++
				return uuid, nil
			},
			AtreeStorageValidationEnabled: true,
			AtreeValueValidationEnabled:   true,
			Storage:                       interpreter.NewInMemoryStorage(nil),
		},
	)
	require.NoError(t, err)

	err = inter.Interpret()
	require.NoError(t, err)

	result, err := inter.Invoke("main")
	require.NoError(t, err)

	exported, err := runtime.ExportValue(result, inter, interpreter.EmptyLocationRange)
	require.NoError(t, err)

	return exported
}

func TestEncodeResource(t *testing.T) {

	t.Parallel()

	t.Run("Simple", func(t *testing.T) {

		t.Parallel()

		actual := exportFromScript(t, `
			resource Foo {
				let bar: Int

				init(bar: Int) {
					self.bar = bar
				}
			}

			fun main(): @Foo {
				return <- create Foo(bar: 42)
			}
		`)

		// expectedVal is different from actual because "bar" field is
		// encoded before "uuid" field due to deterministic encoding.
		expectedVal := cadence.NewResource([]cadence.Value{
			cadence.NewInt(42),
			cadence.NewUInt64(1),
		}).WithType(cadence.NewResourceType(
			common.NewStringLocation(nil, "test"),
			"Foo",
			[]cadence.Field{
				{Type: cadence.NewIntType(), Identifier: "bar"},
				{Type: cadence.NewUInt64Type(), Identifier: "uuid"},
			},
			nil,
		))

		expectedCBOR := []byte{
			// language=json, format=json-cadence data interchange format
			// {"type":"Resource","value":{"id":"S.test.Foo","fields":[{"name":"uuid","value":{"type":"UInt64","value":"1"}},{"name":"bar","value":{"type":"Int","value":"42"}}]}}
			//
			// language=edn, format=ccf
			// 129([[161([h'', "S.test.Foo", [["bar", 137(4)], ["uuid", 137(15)]]])], [136(h''), [42, 1]]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definitions
			// array, 1 items follow
			0x81,
			// resource type:
			// id: []byte{}
			// cadence-type-id: "S.test.Foo"
			// 2 fields: [["bar", type(int)], ["uuid", type(uint64)]]
			// tag
			0xd8, ccf.CBORTagResourceType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 10 bytes follow
			0x6a,
			// S.test.Foo
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f,
			// fields
			// array, 2 items follow
			0x82,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 3 bytes follow
			0x63,
			// bar
			0x62, 0x61, 0x72,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// field 1
			// array, 2 items follow
			0x82,
			// text, 4 bytes follow
			0x64,
			// uuid
			0x75, 0x75, 0x69, 0x64,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Uint type ID (15)
			0x0f,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 0 bytes follow
			0x40,
			// array, 2 items follow
			0x82,
			// tag (big number)
			0xc2,
			// bytes, 1 byte follow
			0x41,
			// 42
			0x2a,
			// 1
			0x01,
		}

		testEncodeAndDecodeEx(t, actual, expectedCBOR, expectedVal)
	})

	t.Run("With function member", func(t *testing.T) {

		t.Parallel()

		actual := exportFromScript(t, `
				resource Foo {
					let bar: Int

					fun foo(): String {
						return "foo"
					}

					init(bar: Int) {
						self.bar = bar
					}
				}

				fun main(): @Foo {
					return <- create Foo(bar: 42)
				}
			`)

		// expectedVal is different from actual because "bar" field is
		// encoded before "uuid" field due to deterministic encoding.
		expectedVal := cadence.NewResource([]cadence.Value{
			cadence.NewInt(42),
			cadence.NewUInt64(1),
		}).WithType(cadence.NewResourceType(
			common.NewStringLocation(nil, "test"),
			"Foo",
			[]cadence.Field{
				{Type: cadence.NewIntType(), Identifier: "bar"},
				{Type: cadence.NewUInt64Type(), Identifier: "uuid"},
			},
			nil,
		))

		// function "foo" should be omitted from resulting CBOR
		expectedCBOR := []byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Resource","value":{"id":"S.test.Foo","fields":[{"name":"uuid","value":{"type":"UInt64","value":"1"}},{"name":"bar","value":{"type":"Int","value":"42"}}]}}
			//
			// language=edn, format=ccf
			// 129([[161([h'', "S.test.Foo", [["bar", 137(4)], ["uuid", 137(15)]]])], [136(h''), [42, 1]]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definitions
			// array, 1 items follow
			0x81,
			// resource type:
			// id: []byte{}
			// cadence-type-id: "S.test.Foo"
			// 2 fields: [["bar", type(int)], ["uuid", type(uint64)]]
			// tag
			0xd8, ccf.CBORTagResourceType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 10 bytes follow
			0x6a,
			// S.test.Foo
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f,
			// fields
			// array, 2 items follow
			0x82,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 3 bytes follow
			0x63,
			// bar
			0x62, 0x61, 0x72,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// field 1
			// array, 2 items follow
			0x82,
			// text, 4 bytes follow
			0x64,
			// uuid
			0x75, 0x75, 0x69, 0x64,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Uint type ID (15)
			0x0f,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 0 bytes follow
			0x40,
			// array, 2 items follow
			0x82,
			// tag (big number)
			0xc2,
			// bytes, 1 byte follow
			0x41,
			// 42
			0x2a,
			// 1
			0x01,
		}

		testEncodeAndDecodeEx(t, actual, expectedCBOR, expectedVal)
	})

	t.Run("Nested resource", func(t *testing.T) {

		t.Parallel()

		actual := exportFromScript(t, `
				resource Bar {
					let x: Int

					init(x: Int) {
						self.x = x
					}
				}

				resource Foo {
					let bar: @Bar

					init(bar: @Bar) {
						self.bar <- bar
					}

					destroy() {
						destroy self.bar
					}
				}

				fun main(): @Foo {
					return <- create Foo(bar: <- create Bar(x: 42))
				}
			`)

		// S.test.Foo(uuid: 2, bar: S.test.Bar(uuid: 1, x: 42)) (cadence.Resource)

		// expectedVal is different from actual because "bar" field is
		// encoded before "uuid" field due to deterministic encoding.
		expectedBarResourceType := cadence.NewResourceType(
			common.NewStringLocation(nil, "test"),
			"Bar",
			[]cadence.Field{
				{Type: cadence.NewIntType(), Identifier: "x"},
				{Type: cadence.NewUInt64Type(), Identifier: "uuid"},
			},
			nil,
		)
		expectedBarResource := cadence.NewResource(
			[]cadence.Value{
				cadence.NewInt(42),
				cadence.NewUInt64(1),
			},
		).WithType(expectedBarResourceType)

		expectedVal := cadence.NewResource(
			[]cadence.Value{
				expectedBarResource,
				cadence.NewUInt64(2),
			}).WithType(cadence.NewResourceType(
			common.NewStringLocation(nil, "test"),
			"Foo",
			[]cadence.Field{
				{Type: expectedBarResourceType, Identifier: "bar"},
				{Type: cadence.NewUInt64Type(), Identifier: "uuid"},
			},
			nil,
		))

		expectedCBOR := []byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Resource","value":{"id":"S.test.Foo","fields":[{"name":"uuid","value":{"type":"UInt64","value":"2"}},{"name":"bar","value":{"type":"Resource","value":{"id":"S.test.Bar","fields":[{"name":"uuid","value":{"type":"UInt64","value":"1"}},{"name":"x","value":{"type":"Int","value":"42"}}]}}}]}}
			//
			// language=edn, format=ccf
			// 129([[161([h'', "S.test.Bar", [["x", 137(4)], ["uuid", 137(15)]]]), 161([h'01', "S.test.Foo", [["bar", 136(h'')], ["uuid", 137(15)]]])], [136(h'01'), [[42, 1], 2]]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definitions
			// array, 2 items follow
			0x82,

			// resource type:
			// id: []byte{01}
			// cadence-type-id: "S.test.Bar"
			// 2 fields: [["x", type(int)], ["uuid", type(uint64)], ]
			// tag
			0xd8, ccf.CBORTagResourceType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 10 bytes follow
			0x6a,
			// S.test.Bar
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x42, 0x61, 0x72,
			// fields
			// array, 2 items follow
			0x82,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 3 bytes follow
			0x61,
			// x
			0x78,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// field 1
			// array, 2 items follow
			0x82,
			// text, 4 bytes follow
			0x64,
			// uuid
			0x75, 0x75, 0x69, 0x64,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Uint64 type ID (15)
			0x0f,

			// resource type:
			// id: []byte{}
			// cadence-type-id: "S.test.Foo"
			// 2 fields: [["bar", type ref(1)], ["uuid", type(uint64)]]
			// tag
			0xd8, ccf.CBORTagResourceType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// cadence-type-id
			// string, 10 bytes follow
			0x6a,
			// S.test.Foo
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f,
			// fields
			// array, 2 items follow
			0x82,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 3 bytes follow
			0x63,
			// bar
			0x62, 0x61, 0x72,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// type definition ID (0)
			// bytes, 0 bytes follow
			0x40,
			// field 1
			// array, 2 items follow
			0x82,
			// text, 4 bytes follow
			0x64,
			// uuid
			0x75, 0x75, 0x69, 0x64,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Uint64 type ID (15)
			0x0f,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// array, 2 items follow
			0x82,
			// array, 2 items follow
			0x82,
			// tag (big number)
			0xc2,
			// bytes, 1 byte follow
			0x41,
			// 42
			0x2a,
			// 1
			0x01,
			// 2
			0x02,
		}

		testEncodeAndDecodeEx(t, actual, expectedCBOR, expectedVal)
	})
}

func TestEncodeStruct(t *testing.T) {

	t.Parallel()

	simpleStructType := &cadence.StructType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "FooStruct",
		Fields: []cadence.Field{
			{
				Identifier: "a",
				Type:       cadence.IntType{},
			},
			{
				Identifier: "b",
				Type:       cadence.StringType{},
			},
		},
	}

	simpleStruct := encodeTest{
		"Simple",
		cadence.NewStruct(
			[]cadence.Value{
				cadence.NewInt(1),
				cadence.String("foo"),
			},
		).WithType(simpleStructType),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Struct","value":{"id":"S.test.FooStruct","fields":[{"name":"a","value":{"type":"Int","value":"1"}},{"name":"b","value":{"type":"String","value":"foo"}}]}}
			//
			// language=edn, format=ccf
			// 129([[160([h'', "S.test.FooStruct", [["a", 137(4)], ["b", 137(1)]]])], [136(h''), [1, "foo"]]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definitions
			// array, 1 items follow
			0x81,
			// struct type:
			// id: []byte{}
			// cadence-type-id: "S.test.FooStruct"
			// 2 fields: [["a", type(int)], ["b", type(string)]]
			// tag
			0xd8, ccf.CBORTagStructType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 16 bytes follow
			0x70,
			// S.test.FooStruct
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74,
			// fields
			// array, 2 items follow
			0x82,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 1 bytes follow
			0x61,
			// a
			0x61,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// field 1
			// array, 2 items follow
			0x82,
			// text, 1 bytes follow
			0x61,
			// b
			0x62,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// String type ID (1)
			0x01,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 0 bytes follow
			0x40,
			// array, 2 items follow
			0x82,
			// tag (big number)
			0xc2,
			// bytes, 1 byte follow
			0x41,
			// 1
			0x01,
			// String, 3 bytes follow
			0x63,
			// foo
			0x66, 0x6f, 0x6f,
		},
	}

	resourceStructType := &cadence.StructType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "FooStruct",
		Fields: []cadence.Field{
			{
				Identifier: "a",
				Type:       cadence.StringType{},
			},
			{
				Identifier: "b",
				Type:       fooResourceType,
			},
		},
	}

	resourceStruct := encodeTest{
		"Resources",
		cadence.NewStruct(
			[]cadence.Value{
				cadence.String("foo"),
				cadence.NewResource(
					[]cadence.Value{
						cadence.NewInt(42),
					},
				).WithType(fooResourceType),
			},
		).WithType(resourceStructType),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Struct","value":{"id":"S.test.FooStruct","fields":[{"name":"a","value":{"type":"String","value":"foo"}},{"name":"b","value":{"type":"Resource","value":{"id":"S.test.Foo","fields":[{"name":"bar","value":{"type":"Int","value":"42"}}]}}}]}}
			//
			// language=edn, format=ccf
			// 129([[161([h'', "S.test.Foo", [["bar", 137(4)]]]), 160([h'01', "S.test.FooStruct", [["a", 137(1)], ["b", 136(h'')]]])], [136(h'01'), ["foo", [42]]]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definitions
			// array, 2 items follow
			0x82,

			// resource type:
			// id: []byte{}
			// cadence-type-id: "S.test.Foo"
			// fields: [["bar", int type]]
			// tag
			0xd8, ccf.CBORTagResourceType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 10 bytes follow
			0x6a,
			// S.test.Foo
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f,
			// fields
			// array, 1 items follow
			0x81,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 3 bytes follow
			0x63,
			// bar
			0x62, 0x61, 0x72,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,

			// struct type:
			// id: []byte{}
			// cadence-type-id: "S.test.FooStruct"
			// 2 fields: [["a", type(int)], ["b", type(string)]]
			// tag
			0xd8, ccf.CBORTagStructType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// cadence-type-id
			// string, 16 bytes follow
			0x70,
			// S.test.FooStruct
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74,
			// fields
			// array, 2 items follow
			0x82,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 1 bytes follow
			0x61,
			// a
			0x61,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// String type ID (1)
			0x01,
			// field 1
			// array, 2 items follow
			0x82,
			// text, 1 bytes follow
			0x61,
			// b
			0x62,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// type reference ID (1)
			// bytes, 0 bytes follow
			0x40,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// array, 2 items follow
			0x82,
			// String, 3 bytes follow
			0x63,
			// foo
			0x66, 0x6f, 0x6f,
			// array, 1 items follow
			0x81,
			// tag (big number)
			0xc2,
			// bytes, 1 byte follow
			0x41,
			// 42
			0x2a,
		},
	}

	testAllEncodeAndDecode(t, simpleStruct, resourceStruct)
}

func TestEncodeEvent(t *testing.T) {

	t.Parallel()

	simpleEventType := &cadence.EventType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "FooEvent",
		Fields: []cadence.Field{
			{
				Identifier: "a",
				Type:       cadence.IntType{},
			},
			{
				Identifier: "b",
				Type:       cadence.StringType{},
			},
		},
	}

	simpleEvent := encodeTest{
		"Simple",
		cadence.NewEvent(
			[]cadence.Value{
				cadence.NewInt(1),
				cadence.String("foo"),
			},
		).WithType(simpleEventType),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Event","value":{"id":"S.test.FooEvent","fields":[{"name":"a","value":{"type":"Int","value":"1"}},{"name":"b","value":{"type":"String","value":"foo"}}]}}
			//
			// language=edn, format=ccf
			// 129([[162([h'', "S.test.FooEvent", [["a", 137(4)], ["b", 137(1)]]])], [136(h''), [1, "foo"]]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definitions
			// array, 1 items follow
			0x81,
			// event type:
			// id: []byte{}
			// cadence-type-id: "S.test.FooEvent"
			// 2 fields: [["a", type(int)], ["b", type(string)]]
			// tag
			0xd8, ccf.CBORTagEventType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 15 bytes follow
			0x6f,
			// S.test.FooEvent
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f, 0x45, 0x76, 0x65, 0x6e, 0x74,
			// fields
			// array, 2 items follow
			0x82,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 1 bytes follow
			0x61,
			// a
			0x61,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// field 1
			// array, 2 items follow
			0x82,
			// text, 1 bytes follow
			0x61,
			// b
			0x62,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// String type ID (1)
			0x01,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 0 bytes follow
			0x40,
			// array, 2 items follow
			0x82,
			// tag (big number)
			0xc2,
			// bytes, 1 byte follow
			0x41,
			// 1
			0x01,
			// String, 3 bytes follow
			0x63,
			// foo
			0x66, 0x6f, 0x6f,
		},
	}

	resourceEventType := &cadence.EventType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "FooEvent",
		Fields: []cadence.Field{
			{
				Identifier: "a",
				Type:       cadence.StringType{},
			},
			{
				Identifier: "b",
				Type:       fooResourceType,
			},
		},
	}

	resourceEvent := encodeTest{
		"Resources",
		cadence.NewEvent(
			[]cadence.Value{
				cadence.String("foo"),
				cadence.NewResource(
					[]cadence.Value{
						cadence.NewInt(42),
					},
				).WithType(fooResourceType),
			},
		).WithType(resourceEventType),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Event","value":{"id":"S.test.FooEvent","fields":[{"name":"a","value":{"type":"String","value":"foo"}},{"name":"b","value":{"type":"Resource","value":{"id":"S.test.Foo","fields":[{"name":"bar","value":{"type":"Int","value":"42"}}]}}}]}}
			//
			// language=edn, format=ccf
			// 129([[161([h'', "S.test.Foo", [["bar", 137(4)]]]), 162([h'01', "S.test.FooEvent", [["a", 137(1)], ["b", 136(h'')]]])], [136(h'01'), ["foo", [42]]]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definitions
			// array, 2 items follow
			0x82,

			// resource type:
			// id: []byte{}
			// cadence-type-id: "S.test.Foo"
			// fields: [["bar", int type]]
			// tag
			0xd8, ccf.CBORTagResourceType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 10 bytes follow
			0x6a,
			// S.test.Foo
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f,
			// fields
			// array, 1 items follow
			0x81,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 3 bytes follow
			0x63,
			// bar
			0x62, 0x61, 0x72,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,

			// event type:
			// id: []byte{0x01}
			// cadence-type-id: "S.test.FooEvent"
			// 2 fields: [["a", type(int)], ["b", type(string)]]
			// tag
			0xd8, ccf.CBORTagEventType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// cadence-type-id
			// string, 15 bytes follow
			0x6f,
			// S.test.FooEvent
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f, 0x45, 0x76, 0x65, 0x6e, 0x74,
			// fields
			// array, 2 items follow
			0x82,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 1 bytes follow
			0x61,
			// a
			0x61,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// String type ID (1)
			0x01,
			// field 1
			// array, 2 items follow
			0x82,
			// text, 1 bytes follow
			0x61,
			// b
			0x62,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 0 bytes follow
			0x40,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// array, 2 items follow
			0x82,
			// String, 3 bytes follow
			0x63,
			// foo
			0x66, 0x6f, 0x6f,
			// array, 1 items follow
			0x81,
			// tag (big number)
			0xc2,
			// bytes, 1 byte follow
			0x41,
			// 42
			0x2a,
		},
	}

	testAllEncodeAndDecode(t, simpleEvent, resourceEvent)
}

func TestEncodeContract(t *testing.T) {

	t.Parallel()

	simpleContractType := &cadence.ContractType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "FooContract",
		Fields: []cadence.Field{
			{
				Identifier: "a",
				Type:       cadence.IntType{},
			},
			{
				Identifier: "b",
				Type:       cadence.StringType{},
			},
		},
	}

	simpleContract := encodeTest{
		"Simple",
		cadence.NewContract(
			[]cadence.Value{
				cadence.NewInt(1),
				cadence.String("foo"),
			},
		).WithType(simpleContractType),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Contract","value":{"id":"S.test.FooContract","fields":[{"name":"a","value":{"type":"Int","value":"1"}},{"name":"b","value":{"type":"String","value":"foo"}}]}}
			//
			// language=edn, format=ccf
			// 129([[163([h'', "S.test.FooContract", [["a", 137(4)], ["b", 137(1)]]])], [136(h''), [1, "foo"]]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definitions
			// array, 1 items follow
			0x81,
			// contract type:
			// id: []byte{}
			// cadence-type-id: "S.test.FooContract"
			// 2 fields: [["a", type(int)], ["b", type(string)]]
			// tag
			0xd8, ccf.CBORTagContractType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 18 bytes follow
			0x72,
			// S.test.FooContract
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74,
			// fields
			// array, 2 items follow
			0x82,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 1 bytes follow
			0x61,
			// a
			0x61,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// field 1
			// array, 2 items follow
			0x82,
			// text, 1 bytes follow
			0x61,
			// b
			0x62,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// String type ID (1)
			0x01,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 0 bytes follow
			0x40,
			// array, 2 items follow
			0x82,
			// tag (big number)
			0xc2,
			// bytes, 1 byte follow
			0x41,
			// 1
			0x01,
			// String, 3 bytes follow
			0x63,
			// foo
			0x66, 0x6f, 0x6f,
		},
	}

	resourceContractType := &cadence.ContractType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "FooContract",
		Fields: []cadence.Field{
			{
				Identifier: "a",
				Type:       cadence.StringType{},
			},
			{
				Identifier: "b",
				Type:       fooResourceType,
			},
		},
	}

	resourceContract := encodeTest{
		"Resources",
		cadence.NewContract(
			[]cadence.Value{
				cadence.String("foo"),
				cadence.NewResource(
					[]cadence.Value{
						cadence.NewInt(42),
					},
				).WithType(fooResourceType),
			},
		).WithType(resourceContractType),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Contract","value":{"id":"S.test.FooContract","fields":[{"name":"a","value":{"type":"String","value":"foo"}},{"name":"b","value":{"type":"Resource","value":{"id":"S.test.Foo","fields":[{"name":"bar","value":{"type":"Int","value":"42"}}]}}}]}}
			//
			// language=edn, format=ccf
			// 129([[161([h'', "S.test.Foo", [["bar", 137(4)]]]), 163([h'01', "S.test.FooContract", [["a", 137(1)], ["b", 136(h'')]]])], [136(h'01'), ["foo", [42]]]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definitions
			// array, 2 items follow
			0x82,

			// resource type:
			// id: []byte{}
			// cadence-type-id: "S.test.Foo"
			// fields: [["bar", int type]]
			// tag
			0xd8, ccf.CBORTagResourceType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 10 bytes follow
			0x6a,
			// S.test.Foo
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f,
			// fields
			// array, 1 items follow
			0x81,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 3 bytes follow
			0x63,
			// bar
			0x62, 0x61, 0x72,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,

			// contract type:
			// id: []byte{}
			// cadence-type-id: "S.test.FooContract"
			// 2 fields: [["a", type(int)], ["b", type(string)]]
			// tag
			0xd8, ccf.CBORTagContractType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// cadence-type-id
			// string, 18 bytes follow
			0x72,
			// S.test.FooContract
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x61, 0x63, 0x74,
			// fields
			// array, 2 items follow
			0x82,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 1 bytes follow
			0x61,
			// a
			0x61,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// String type ID (1)
			0x01,
			// field 1
			// array, 2 items follow
			0x82,
			// text, 1 bytes follow
			0x61,
			// b
			0x62,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 0 bytes follow
			0x40,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 1 bytes follow
			0x41,
			// 1
			0x01,
			// array, 2 items follow
			0x82,
			// String, 3 bytes follow
			0x63,
			// foo
			0x66, 0x6f, 0x6f,
			// array, 1 items follow
			0x81,
			// tag (big number)
			0xc2,
			// bytes, 1 byte follow
			0x41,
			// 42
			0x2a,
		},
	}

	testAllEncodeAndDecode(t, simpleContract, resourceContract)
}

func TestEncodeSimpleTypes(t *testing.T) {

	t.Parallel()

	type simpleTypes struct {
		typ              cadence.Type
		cborSimpleTypeID int
	}

	var tests []encodeTest

	for _, ty := range []simpleTypes{
		{cadence.AnyType{}, ccf.TypeAny},
		{cadence.AnyResourceType{}, ccf.TypeAnyResource},
		{cadence.MetaType{}, ccf.TypeMetaType},
		{cadence.VoidType{}, ccf.TypeVoid},
		{cadence.NeverType{}, ccf.TypeNever},
		{cadence.BoolType{}, ccf.TypeBool},
		{cadence.StringType{}, ccf.TypeString},
		{cadence.CharacterType{}, ccf.TypeCharacter},
		{cadence.BytesType{}, ccf.TypeBytes},
		{cadence.AddressType{}, ccf.TypeAddress},
		{cadence.SignedNumberType{}, ccf.TypeSignedNumber},
		{cadence.IntegerType{}, ccf.TypeInteger},
		{cadence.SignedIntegerType{}, ccf.TypeSignedInteger},
		{cadence.FixedPointType{}, ccf.TypeFixedPoint},
		{cadence.IntType{}, ccf.TypeInt},
		{cadence.Int8Type{}, ccf.TypeInt8},
		{cadence.Int16Type{}, ccf.TypeInt16},
		{cadence.Int32Type{}, ccf.TypeInt32},
		{cadence.Int64Type{}, ccf.TypeInt64},
		{cadence.Int128Type{}, ccf.TypeInt128},
		{cadence.Int256Type{}, ccf.TypeInt256},
		{cadence.UIntType{}, ccf.TypeUInt},
		{cadence.UInt8Type{}, ccf.TypeUInt8},
		{cadence.UInt16Type{}, ccf.TypeUInt16},
		{cadence.UInt32Type{}, ccf.TypeUInt32},
		{cadence.UInt64Type{}, ccf.TypeUInt64},
		{cadence.UInt128Type{}, ccf.TypeUInt128},
		{cadence.UInt256Type{}, ccf.TypeUInt256},
		{cadence.Word8Type{}, ccf.TypeWord8},
		{cadence.Word16Type{}, ccf.TypeWord16},
		{cadence.Word32Type{}, ccf.TypeWord32},
		{cadence.Word64Type{}, ccf.TypeWord64},
		{cadence.Fix64Type{}, ccf.TypeFix64},
		{cadence.UFix64Type{}, ccf.TypeUFix64},
		{cadence.BlockType{}, ccf.TypeBlock},
		{cadence.PathType{}, ccf.TypePath},
		{cadence.CapabilityPathType{}, ccf.TypeCapabilityPath},
		{cadence.StoragePathType{}, ccf.TypeStoragePath},
		{cadence.PublicPathType{}, ccf.TypePublicPath},
		{cadence.PrivatePathType{}, ccf.TypePrivatePath},
		{cadence.AccountKeyType{}, ccf.TypeAccountKey},
		{cadence.AuthAccountContractsType{}, ccf.TypeAuthAccountContracts},
		{cadence.AuthAccountKeysType{}, ccf.TypeAuthAccountKeys},
		{cadence.AuthAccountType{}, ccf.TypeAuthAccount},
		{cadence.PublicAccountContractsType{}, ccf.TypePublicAccountContracts},
		{cadence.PublicAccountKeysType{}, ccf.TypePublicAccountKeys},
		{cadence.PublicAccountType{}, ccf.TypePublicAccount},
		{cadence.DeployedContractType{}, ccf.TypeDeployedContract},
	} {
		var w bytes.Buffer

		encoder := ccf.CBOREncMode.NewStreamEncoder(&w)

		err := encoder.EncodeRawBytes([]byte{
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 elements follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Meta type ID (41)
			0x18, 0x29,
			// tag
			0xd8, ccf.CBORTagSimpleTypeValue,
		})
		require.NoError(t, err)

		err = encoder.EncodeInt(ty.cborSimpleTypeID)
		require.NoError(t, err)

		encoder.Flush()

		tests = append(tests, encodeTest{
			name: fmt.Sprintf("with static %s", ty.typ.ID()),
			val: cadence.TypeValue{
				StaticType: ty.typ,
			},
			expected: w.Bytes(),
			// language=json, format=json-cadence data interchange format
			// {"type":"Type","value":{"staticType":{"kind":"[ty.ID()]"}}}
			//
			// language=edn, format=ccf
			// 130([137(41), 185(simple_type_id)])
		})
	}

	testAllEncodeAndDecode(t, tests...)
}

func TestEncodeType(t *testing.T) {

	t.Parallel()

	t.Run("with static int?", func(t *testing.T) {

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: &cadence.OptionalType{Type: cadence.IntType{}},
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":{"kind":"Optional", "type" : {"kind" : "Int"}}}}
				//
				// language=edn, format=ccf
				// 130([137(41), 186(185(4))])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagOptionalTypeValue,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
			},
		)

	})

	t.Run("with static int??", func(t *testing.T) {

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: &cadence.OptionalType{Type: &cadence.OptionalType{Type: cadence.IntType{}}},
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":{"kind":"Optional", "type" : {"kind" : "Int"}}}}
				//
				// language=edn, format=ccf
				// 130([137(41), 186(186(185(4)))])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagOptionalTypeValue,
				// tag
				0xd8, ccf.CBORTagOptionalTypeValue,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
			},
		)

	})
	t.Run("with static [int]", func(t *testing.T) {

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: &cadence.VariableSizedArrayType{ElementType: cadence.IntType{}},
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":{"kind":"VariableSizedArray", "type" : {"kind" : "Int"}}}}
				//
				// language=edn, format=ccf
				// 130([137(41), 187(185(4))])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagVarsizedArrayTypeValue,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
			},
		)

	})

	t.Run("with static [int; 3]", func(t *testing.T) {

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: &cadence.ConstantSizedArrayType{
					ElementType: cadence.IntType{},
					Size:        3,
				},
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":{"kind":"ConstantSizedArray", "type" : {"kind" : "Int"}, "size" : 3}}}
				//
				// language=edn, format=ccf
				// 130([137(41), 188([3, 185(4)])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagConstsizedArrayTypeValue,
				// array, 2 elements follow
				0x82,
				// 3
				0x03,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
			},
		)

	})

	t.Run("with static {int:string}", func(t *testing.T) {

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: &cadence.DictionaryType{
					ElementType: cadence.StringType{},
					KeyType:     cadence.IntType{},
				},
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":{"kind":"Dictionary", "key" : {"kind" : "Int"}, "value" : {"kind" : "String"}}}}
				//
				// language=edn, format=ccf
				// 130([137(41), 189([185(4), 185(1)])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagDictTypeValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// String type (1)
				0x01,
			},
		)

	})

	t.Run("with static struct", func(t *testing.T) {

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: &cadence.StructType{
					Location:            utils.TestLocation,
					QualifiedIdentifier: "S",
					Fields: []cadence.Field{
						{Identifier: "foo", Type: cadence.IntType{}},
					},
					Initializers: [][]cadence.Parameter{
						{{Label: "foo", Identifier: "bar", Type: cadence.IntType{}}},
						{{Label: "qux", Identifier: "baz", Type: cadence.StringType{}}},
					},
				},
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type", "value": {"staticType": {"kind": "Struct", "type" : "", "typeID" : "S.test.S", "fields" : [ {"id" : "foo", "type": {"kind" : "Int"} } ], "initializers" : [ [{"label" : "foo", "id" : "bar", "type": {"kind" : "Int"}}], [{"label" : "qux", "id" : "baz", "type": {"kind" : "String"}}] ] } } }
				//
				// language=edn, format=ccf
				// 130([137(41), 208([h'', "S.test.S", null, [["foo", 185(4)]], [[["foo", "bar", 185(4)]], [["qux", "baz", 185(1)]]]])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagStructTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 8 bytes follow
				0x68,
				// S.test.So
				0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x53,
				// type (nil for struct)
				0xf6,
				// fields
				// array, 1 element follows
				0x81,
				// array, 2 elements follow
				0x82,
				// string, 3 bytes follow
				0x63,
				// foo
				0x66, 0x6f, 0x6f,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
				// initializers
				// array, 2 elements follow
				0x82,
				// array, 1 element follows
				0x81,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// foo
				0x66, 0x6f, 0x6f,
				// string, 3 bytes follow
				0x63,
				// bar
				0x62, 0x61, 0x72,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
				// array, 1 element follows
				0x81,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// qux
				0x71, 0x75, 0x78,
				// string, 3 bytes follow
				0x63,
				// bax
				0x62, 0x61, 0x7a,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// String type (1)
				0x01,
			},
		)
	})
	t.Run("with static resource", func(t *testing.T) {

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: &cadence.ResourceType{
					Location:            utils.TestLocation,
					QualifiedIdentifier: "R",
					Fields: []cadence.Field{
						{Identifier: "foo", Type: cadence.IntType{}},
					},
					Initializers: [][]cadence.Parameter{
						{{Label: "foo", Identifier: "bar", Type: cadence.IntType{}}},
						{{Label: "qux", Identifier: "baz", Type: cadence.StringType{}}},
					},
				},
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type", "value": {"staticType": {"kind": "Resource", "type" : "", "typeID" : "S.test.R", "fields" : [ {"id" : "foo", "type": {"kind" : "Int"} } ], "initializers" : [ [{"label" : "foo", "id" : "bar", "type": {"kind" : "Int"}}], [{"label" : "qux", "id" : "baz", "type": {"kind" : "String"}}] ] } } }
				//
				// language=edn, format=ccf
				// 130([137(41), 209([h'', "S.test.R", null, [["foo", 185(4)]], [[["foo", "bar", 185(4)]], [["qux", "baz", 185(1)]]]])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagResourceTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 8 bytes follow
				0x68,
				// S.test.R
				0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x52,
				// type (nil for struct)
				0xf6,
				// fields
				// array, 1 element follows
				0x81,
				// array, 2 elements follow
				0x82,
				// string, 3 bytes follow
				0x63,
				// foo
				0x66, 0x6f, 0x6f,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
				// initializers
				// array, 2 elements follow
				0x82,
				// array, 1 element follows
				0x81,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// foo
				0x66, 0x6f, 0x6f,
				// string, 3 bytes follow
				0x63,
				// bar
				0x62, 0x61, 0x72,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
				// array, 1 element follows
				0x81,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// qux
				0x71, 0x75, 0x78,
				// string, 3 bytes follow
				0x63,
				// bax
				0x62, 0x61, 0x7a,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// String type (1)
				0x01,
			},
		)
	})

	t.Run("with static contract", func(t *testing.T) {

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: &cadence.ContractType{
					Location:            utils.TestLocation,
					QualifiedIdentifier: "C",
					Fields: []cadence.Field{
						{Identifier: "foo", Type: cadence.IntType{}},
					},
					Initializers: [][]cadence.Parameter{
						{{Label: "foo", Identifier: "bar", Type: cadence.IntType{}}},
						{{Label: "qux", Identifier: "baz", Type: cadence.StringType{}}},
					},
				},
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type", "value": {"staticType": {"kind": "Contract", "type" : "", "typeID" : "S.test.C", "fields" : [ {"id" : "foo", "type": {"kind" : "Int"} } ], "initializers" : [ [{"label" : "foo", "id" : "bar", "type": {"kind" : "Int"}}], [{"label" : "qux", "id" : "baz", "type": {"kind" : "String"}}] ] } } }
				//
				// language=edn, format=ccf
				// 130([137(41), 211([h'', "S.test.C", null, [["foo", 185(4)]], [[["foo", "bar", 185(4)]], [["qux", "baz", 185(1)]]]])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (42)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagContractTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 8 bytes follow
				0x68,
				// S.test.C
				0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x43,
				// type (nil for struct)
				0xf6,
				// fields
				// array, 1 element follows
				0x81,
				// array, 2 elements follow
				0x82,
				// string, 3 bytes follow
				0x63,
				// foo
				0x66, 0x6f, 0x6f,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
				// initializers
				// array, 2 elements follow
				0x82,
				// array, 1 element follows
				0x81,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// foo
				0x66, 0x6f, 0x6f,
				// string, 3 bytes follow
				0x63,
				// bar
				0x62, 0x61, 0x72,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
				// array, 1 element follows
				0x81,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// qux
				0x71, 0x75, 0x78,
				// string, 3 bytes follow
				0x63,
				// bax
				0x62, 0x61, 0x7a,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// String type (1)
				0x01,
			},
		)
	})

	t.Run("with static struct interface", func(t *testing.T) {

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: &cadence.StructInterfaceType{
					Location:            utils.TestLocation,
					QualifiedIdentifier: "S",
					Fields: []cadence.Field{
						{Identifier: "foo", Type: cadence.IntType{}},
					},
					Initializers: [][]cadence.Parameter{
						{{Label: "foo", Identifier: "bar", Type: cadence.IntType{}}},
						{{Label: "qux", Identifier: "baz", Type: cadence.StringType{}}},
					},
				},
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type", "value": {"staticType": {"kind": "StructInterface", "type" : "", "typeID" : "S.test.S", "fields" : [ {"id" : "foo", "type": {"kind" : "Int"} } ], "initializers" : [ [{"label" : "foo", "id" : "bar", "type": {"kind" : "Int"}}], [{"label" : "qux", "id" : "baz", "type": {"kind" : "String"}}] ] } } }
				//
				// language=edn, format=ccf
				// 130([137(41), 224([h'', "S.test.S", null, [["foo", 185(4)]], [[["foo", "bar", 185(4)]], [["qux", "baz", 185(1)]]]])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagStructInterfaceTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 8 bytes follow
				0x68,
				// S.test.S
				0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x53,
				// type (nil for struct)
				0xf6,
				// fields
				// array, 1 element follows
				0x81,
				// array, 2 elements follow
				0x82,
				// string, 3 bytes follow
				0x63,
				// foo
				0x66, 0x6f, 0x6f,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
				// initializers
				// array, 2 elements follow
				0x82,
				// array, 1 element follows
				0x81,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// foo
				0x66, 0x6f, 0x6f,
				// string, 3 bytes follow
				0x63,
				// bar
				0x62, 0x61, 0x72,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
				// array, 1 element follows
				0x81,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// qux
				0x71, 0x75, 0x78,
				// string, 3 bytes follow
				0x63,
				// bax
				0x62, 0x61, 0x7a,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// String type (1)
				0x01,
			},
		)
	})

	t.Run("with static resource interface", func(t *testing.T) {

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: &cadence.ResourceInterfaceType{
					Location:            utils.TestLocation,
					QualifiedIdentifier: "R",
					Fields: []cadence.Field{
						{Identifier: "foo", Type: cadence.IntType{}},
					},
					Initializers: [][]cadence.Parameter{
						{{Label: "foo", Identifier: "bar", Type: cadence.IntType{}}},
						{{Label: "qux", Identifier: "baz", Type: cadence.StringType{}}},
					},
				},
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type", "value": {"staticType": {"kind": "ResourceInterface", "type" : "", "typeID" : "S.test.R", "fields" : [ {"id" : "foo", "type": {"kind" : "Int"} } ], "initializers" : [ [{"label" : "foo", "id" : "bar", "type": {"kind" : "Int"}}], [{"label" : "qux", "id" : "baz", "type": {"kind" : "String"}}] ] } } }
				//
				// language=edn, format=ccf
				// 130([137(41), 225([h'', "S.test.R", null, [["foo", 185(4)]], [[["foo", "bar", 185(4)]], [["qux", "baz", 185(1)]]]])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagResourceInterfaceTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 8 bytes follow
				0x68,
				// S.test.R
				0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x52,
				// type (nil for struct)
				0xf6,
				// fields
				// array, 1 element follows
				0x81,
				// array, 2 elements follow
				0x82,
				// string, 3 bytes follow
				0x63,
				// foo
				0x66, 0x6f, 0x6f,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
				// initializers
				// array, 2 elements follow
				0x82,
				// array, 1 element follows
				0x81,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// foo
				0x66, 0x6f, 0x6f,
				// string, 3 bytes follow
				0x63,
				// bar
				0x62, 0x61, 0x72,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
				// array, 1 element follows
				0x81,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// qux
				0x71, 0x75, 0x78,
				// string, 3 bytes follow
				0x63,
				// bax
				0x62, 0x61, 0x7a,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// String type (1)
				0x01,
			},
		)
	})

	t.Run("with static contract interface", func(t *testing.T) {

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: &cadence.ContractInterfaceType{
					Location:            utils.TestLocation,
					QualifiedIdentifier: "C",
					Fields: []cadence.Field{
						{Identifier: "foo", Type: cadence.IntType{}},
					},
					Initializers: [][]cadence.Parameter{
						{{Label: "foo", Identifier: "bar", Type: cadence.IntType{}}},
						{{Label: "qux", Identifier: "baz", Type: cadence.StringType{}}},
					},
				},
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type", "value": {"staticType": {"kind": "ContractInterface", "type" : "", "typeID" : "S.test.C", "fields" : [ {"id" : "foo", "type": {"kind" : "Int"} } ], "initializers" : [ [{"label" : "foo", "id" : "bar", "type": {"kind" : "Int"}}], [{"label" : "qux", "id" : "baz", "type": {"kind" : "String"}}] ] } } }
				//
				// language=edn, format=ccf
				// 130([137(41), 226([h'', "S.test.C", null, [["foo", 185(4)]], [[["foo", "bar", 185(4)]], [["qux", "baz", 185(1)]]]])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagContractInterfaceTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 8 bytes follow
				0x68,
				// S.test.C
				0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x43,
				// type (nil for contract interface)
				0xf6,
				// fields
				// array, 1 element follows
				0x81,
				// array, 2 elements follow
				0x82,
				// string, 3 bytes follow
				0x63,
				// foo
				0x66, 0x6f, 0x6f,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
				// initializers
				// array, 2 elements follow
				0x82,
				// array, 1 element follows
				0x81,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// foo
				0x66, 0x6f, 0x6f,
				// string, 3 bytes follow
				0x63,
				// bar
				0x62, 0x61, 0x72,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
				// array, 1 element follows
				0x81,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// qux
				0x71, 0x75, 0x78,
				// string, 3 bytes follow
				0x63,
				// bax
				0x62, 0x61, 0x7a,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// String type (1)
				0x01,
			},
		)
	})

	t.Run("with static event", func(t *testing.T) {

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: &cadence.EventType{
					Location:            utils.TestLocation,
					QualifiedIdentifier: "E",
					Fields: []cadence.Field{
						{Identifier: "foo", Type: cadence.IntType{}},
					},
					Initializer: []cadence.Parameter{
						{Label: "foo", Identifier: "bar", Type: cadence.IntType{}},
						{Label: "qux", Identifier: "baz", Type: cadence.StringType{}},
					},
				},
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type", "value": {"staticType": {"kind": "Event", "type" : "", "typeID" : "S.test.E", "fields" : [ {"id" : "foo", "type": {"kind" : "Int"} } ], "initializers" : [[{"label" : "foo", "id" : "bar", "type": {"kind" : "Int"}}, {"label" : "qux", "id" : "baz", "type": {"kind" : "String"}}]] } } }
				//
				// language=edn, format=ccf
				// 130([137(41), 210([h'', "S.test.E", null, [["foo", 185(4)]], [[["foo", "bar", 185(4)], ["qux", "baz", 185(1)]]]])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagEventTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 8 bytes follow
				0x68,
				// S.test.E
				0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x45,
				// type (nil for event)
				0xf6,
				// fields
				// array, 1 element follows
				0x81,
				// array, 2 elements follow
				0x82,
				// string, 3 bytes follow
				0x63,
				// foo
				0x66, 0x6f, 0x6f,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
				// initializers
				// array, 1 elements follow
				0x81,
				// array, 2 element follows
				0x82,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// foo
				0x66, 0x6f, 0x6f,
				// string, 3 bytes follow
				0x63,
				// bar
				0x62, 0x61, 0x72,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// qux
				0x71, 0x75, 0x78,
				// string, 3 bytes follow
				0x63,
				// baz
				0x62, 0x61, 0x7a,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// String type (1)
				0x01,
			},
		)
	})

	t.Run("with static enum", func(t *testing.T) {

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: &cadence.EnumType{
					Location:            utils.TestLocation,
					QualifiedIdentifier: "E",
					RawType:             cadence.StringType{},
					Fields: []cadence.Field{
						{Identifier: "foo", Type: cadence.IntType{}},
					},
					Initializers: [][]cadence.Parameter{
						{{Label: "foo", Identifier: "bar", Type: cadence.IntType{}}},
						{{Label: "qux", Identifier: "baz", Type: cadence.StringType{}}},
					},
				},
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type", "value": {"staticType": {"kind": "Enum", "type" : {"kind" : "String"}, "typeID" : "S.test.E", "fields" : [ {"id" : "foo", "type": {"kind" : "Int"} } ], "initializers" : [ [{"label" : "foo", "id" : "bar", "type": {"kind" : "Int"}}], [{"label" : "qux", "id" : "baz", "type": {"kind" : "String"}}] ] } } }
				//
				// language=edn, format=ccf
				// 130([137(41), 212([h'', "S.test.E", 185(1), [["foo", 185(4)]], [[["foo", "bar", 185(4)]], [["qux", "baz", 185(1)]]]])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagEnumTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 8 bytes follow
				0x68,
				// S.test.E
				0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x45,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// String type ID (1)
				0x01,
				// fields
				// array, 1 element follows
				0x81,
				// array, 2 elements follow
				0x82,
				// string, 3 bytes follow
				0x63,
				// foo
				0x66, 0x6f, 0x6f,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
				// initializers
				// array, 2 elements follow
				0x82,
				// array, 1 element follows
				0x81,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// foo
				0x66, 0x6f, 0x6f,
				// string, 3 bytes follow
				0x63,
				// bar
				0x62, 0x61, 0x72,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type (4)
				0x04,
				// array, 1 element follows
				0x81,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// qux
				0x71, 0x75, 0x78,
				// string, 3 bytes follow
				0x63,
				// bax
				0x62, 0x61, 0x7a,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// String type (1)
				0x01,
			},
		)
	})

	t.Run("with static &int", func(t *testing.T) {

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: &cadence.ReferenceType{
					Authorized: false,
					Type:       cadence.IntType{},
				},
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":{"kind":"Reference", "type" : {"kind" : "Int"}, "authorized" : false}}}`
				//
				// language=edn, format=ccf
				// 130([137(41), 190([false, 185(4)])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagReferenceTypeValue,
				// array, 2 elements follow
				0x82,
				// authorized
				// bool
				0xf4,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type ID (4)
				0x04,
			},
		)

	})

	t.Run("with static function", func(t *testing.T) {

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: (&cadence.FunctionType{
					Parameters: []cadence.Parameter{
						{Label: "qux", Identifier: "baz", Type: cadence.StringType{}},
					},
					ReturnType: cadence.IntType{},
				}).WithID("Foo"),
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType": { "kind" : "Function", "typeID":"Foo", "return" : {"kind" : "Int"}, "parameters" : [ {"label" : "qux", "id" : "baz", "type": {"kind" : "String"}} ]} } }
				//
				// language=edn, format=ccf
				// 130([137(41), 193(["Foo", [["qux", "baz", 185(1)]], 185(4)])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagFunctionTypeValue,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// Foo
				0x46, 0x6f, 0x6f,
				// array, 1 elements follow
				0x81,
				// array, 3 elements follow
				0x83,
				// string, 3 bytes follow
				0x63,
				// qux
				0x71, 0x75, 0x78,
				// string, 3 bytes follow
				0x63,
				// bax
				0x62, 0x61, 0x7a,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// String type (1)
				0x01,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type ID (4)
				0x04,
			},
		)

	})

	t.Run("with static Capability<Int>", func(t *testing.T) {

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: &cadence.CapabilityType{
					BorrowType: cadence.IntType{},
				},
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":{"kind":"Capability", "type" : {"kind" : "Int"}}}}
				//
				// language=edn, format=ccf
				// 130([137(41), 192([185(4)])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagCapabilityTypeValue,
				// array, 1 element follows
				0x81,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type ID (4)
				0x04,
			},
		)
	})

	t.Run("with static restricted type", func(t *testing.T) {

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: (&cadence.RestrictedType{
					Restrictions: []cadence.Type{
						cadence.StringType{},
					},
					Type: cadence.IntType{},
				}).WithID("Int{String}"),
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType": { "kind": "Restriction", "typeID":"Int{String}", "type" : {"kind" : "Int"}, "restrictions" : [ {"kind" : "String"} ]} } }
				//
				// language=edn, format=ccf
				// 130([137(41), 191(["Int{String}", 185(4), [185(1)]])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagRestrictedTypeValue,
				// array, 3 elements follow
				0x83,
				// ID
				// string, 11 bytes follow
				0x6b,
				// Int{String}
				0x49, 0x6e, 0x74, 0x7b, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x7d,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// Int type ID (4)
				0x04,
				// array, 1 element follows
				0x81,
				// tag
				0xd8, ccf.CBORTagSimpleTypeValue,
				// String type ID (1)
				0x01,
			},
		)

	})

	t.Run("without static type", func(t *testing.T) {

		t.Parallel()

		testEncodeAndDecode(
			t,
			cadence.TypeValue{},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":""}}
				//
				// language=edn, format=ccf
				// 130([137(41), null])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// nil
				0xf6,
			},
		)
	})
}

func TestEncodeCapability(t *testing.T) {

	t.Parallel()

	testEncodeAndDecode(
		t,
		cadence.StorageCapability{
			Path:       cadence.NewPath("storage", "foo"),
			Address:    cadence.BytesToAddress([]byte{1, 2, 3, 4, 5}),
			BorrowType: cadence.IntType{},
		},
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Capability","value":{"path":{"type":"Path","value":{"domain":"storage","identifier":"foo"}},"borrowType":{"kind":"Int"},"address":"0x0000000102030405"}}
			//
			// language=edn, format=ccf
			// 130([144([137(4)]), [h'0000000102030405', [1, "foo"]]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 elements follow
			0x82,
			// tag
			0xd8, ccf.CBORTagCapabilityType,
			// array, 1 element follows
			0x81,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Int type ID (4)
			0x04,
			// array, 2 elements follow
			0x82,
			// address
			// bytes, 8 bytes folloow
			0x48,
			// {1,2,3,4,5}
			0x00, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05,
			// array, 2 elements follow
			0x82,
			// 1
			0x01,
			// string, 3 bytes follow
			0x63,
			// foo
			0x66, 0x6f, 0x6f,
		},
	)
}

func TestDecodeFix64(t *testing.T) {

	t.Parallel()

	var maxInt int64 = sema.Fix64TypeMaxInt
	var minInt int64 = sema.Fix64TypeMinInt
	var maxFrac int64 = sema.Fix64TypeMaxFractional
	var minFrac int64 = sema.Fix64TypeMinFractional
	var factor int64 = sema.Fix64Factor

	type test struct {
		name        string
		expected    cadence.Fix64
		encodedData []byte
		check       func(t *testing.T, actual cadence.Value, err error)
	}

	var tests = []test{
		{
			name:     "12.3",
			expected: cadence.Fix64(12_30000000),
			encodedData: []byte{
				// language=json, format=json-cadence data interchange format
				// {"type": "Fix64", "value": "12.3"}
				//
				// language=edn, format=ccf
				// 130([137(22), 1230000000])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// fix64 type ID (22)
				0x16,
				// 1230000000
				0x1a, 0x49, 0x50, 0x4f, 0x80,
			},
		},
		{
			name:     "12.03",
			expected: cadence.Fix64(12_03000000),
			encodedData: []byte{
				// language=json, format=json-cadence data interchange format
				// {"type": "Fix64", "value": "12.03"}
				//
				// language=edn, format=ccf
				// 130([137(22), 1203000000])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// fix64 type ID (22)
				0x16,
				// 1203000000
				0x1a, 0x47, 0xb4, 0x52, 0xc0,
			},
		},
		{
			name:     "12.003",
			expected: cadence.Fix64(12_00300000),
			encodedData: []byte{
				// language=json, format=json-cadence data interchange format
				// {"type": "Fix64", "value": "12.003"}
				//
				// language=edn, format=ccf
				// 130([137(22), 1200300000])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// fix64 type ID (22)
				0x16,
				// 1200300000
				0x1a, 0x47, 0x8b, 0x1f, 0xe0,
			},
		},
		{
			name:     "12.0003",
			expected: cadence.Fix64(12_00030000),
			encodedData: []byte{
				// language=json, format=json-cadence data interchange format
				// {"type": "Fix64", "value": "12.0003"}
				//
				// language=edn, format=ccf
				// 130([137(22), 1200030000])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// fix64 type ID (22)
				0x16,
				// 1200030000
				0x1a, 0x47, 0x87, 0x01, 0x30,
			},
		},
		{
			name:     "12.00003",
			expected: cadence.Fix64(12_00003000),
			encodedData: []byte{
				// language=json, format=json-cadence data interchange format
				// {"type": "Fix64", "value": "12.00003"}
				//
				// language=edn, format=ccf
				// 130([137(22), 1200003000])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// fix64 type ID (22)
				0x16,
				// 1200003000
				0x1a, 0x47, 0x86, 0x97, 0xb8,
			},
		},
		{
			name:     "12.000003",
			expected: cadence.Fix64(12_00000300),
			encodedData: []byte{
				// language=json, format=json-cadence data interchange format
				// {"type": "Fix64", "value": "12.000003"}
				//
				// language=edn, format=ccf
				// 130([137(22), 1200000300])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// fix64 type ID (22)
				0x16,
				// 1200000300
				0x1a, 0x47, 0x86, 0x8d, 0x2c,
			},
		},
		{
			name:     "12.0000003",
			expected: cadence.Fix64(12_00000030),
			encodedData: []byte{
				// language=json, format=json-cadence data interchange format
				// {"type": "Fix64", "value": "12.0000003"}
				//
				// language=edn, format=ccf
				// 130([137(22), 1200000030])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// fix64 type ID (22)
				0x16,
				// 1200000030
				0x1a, 0x47, 0x86, 0x8c, 0x1e,
			},
		},
		{
			name:     "120.3",
			expected: cadence.Fix64(120_30000000),
			encodedData: []byte{
				// language=json, format=json-cadence data interchange format
				// {"type": "Fix64", "value": "120.3"}
				//
				// language=edn, format=ccf
				// 130([137(22), 12030000000])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// fix64 type ID (22)
				0x16,
				// 12030000000
				0x1b, 0x00, 0x00, 0x00, 0x02, 0xcd, 0x0b, 0x3b, 0x80,
			},
		},
		{
			// 92233720368.1
			name:     fmt.Sprintf("%d.1", maxInt),
			expected: cadence.Fix64(9223372036810000000),
			encodedData: []byte{
				// language=json, format=json-cadence data interchange format
				// {"type": "Fix64", "value": "92233720368.1"}
				//
				// language=edn, format=ccf
				// 130([137(22), 9223372036810000000])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// fix64 type ID (22)
				0x16,
				// 9223372036810000000
				0x1b, 0x7f, 0xff, 0xff, 0xff, 0xfd, 0x54, 0xc6, 0x80,
			},
		},
		{
			// 92233720369.1
			name: fmt.Sprintf("%d.1", maxInt+1),
			encodedData: []byte{
				// language=json, format=json-cadence data interchange format
				// {"type": "Fix64", "value": "92233720369.1"}
				//
				// language=edn, format=ccf
				// 130([137(22), 9223372036910000000])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// fix64 type ID (22)
				0x16,
				// 9223372036910000000
				0x1b, 0x80, 0x00, 0x00, 0x00, 0x03, 0x4a, 0xa7, 0x80,
			},
			check: func(t *testing.T, actual cadence.Value, err error) {
				assert.Error(t, err)
			},
		},
		{
			// -92233720368.1
			name:     fmt.Sprintf("%d.1", minInt),
			expected: cadence.Fix64(-9223372036810000000),
			encodedData: []byte{
				// language=json, format=json-cadence data interchange format
				// {"type": "Fix64", "value": "-92233720368.1"}
				//
				// language=edn, format=ccf
				// 130([137(22), -9223372036810000000])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// fix64 type ID (22)
				0x16,
				// -9223372036810000000
				0x3b, 0x7f, 0xff, 0xff, 0xff, 0xfd, 0x54, 0xc6, 0x7f,
			},
		},
		{
			// -92233720369.1
			name: fmt.Sprintf("%d.1", minInt-1),
			encodedData: []byte{
				// language=json, format=json-cadence data interchange format
				// {"type": "Fix64", "value": "-92233720369.1"}
				//
				// language=edn, format=ccf
				// 130([137(22), -9223372036910000000])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// fix64 type ID (22)
				0x16,
				// -9223372036910000000
				0x3b, 0x80, 0x00, 0x00, 0x00, 0x03, 0x4a, 0xa7, 0x7f,
			},
			check: func(t *testing.T, actual cadence.Value, err error) {
				assert.Error(t, err)
			},
		},
		{
			// 92233720368.54775807
			name:     fmt.Sprintf("%d.%d", maxInt, maxFrac),
			expected: cadence.Fix64(maxInt*factor + maxFrac),
			encodedData: []byte{
				// language=json, format=json-cadence data interchange format
				// {"type": "Fix64", "value": "92233720368.54775807"}
				//
				// language=edn, format=ccf
				// 130([137(22), 9223372036854775807])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// fix64 type ID (22)
				0x16,
				// 9223372036854775807
				0x1b, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			},
		},
		{
			// 92233720368.54775808
			name: fmt.Sprintf("%d.%d", maxInt, maxFrac+1),
			encodedData: []byte{
				// language=json, format=json-cadence data interchange format
				// {"type": "Fix64", "value": "92233720368.54775808"}
				//
				// language=edn, format=ccf
				// 130([137(22), 9223372036854775808])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// fix64 type ID (22)
				0x16,
				// 9223372036854775808
				0x1b, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			check: func(t *testing.T, actual cadence.Value, err error) {
				assert.Error(t, err)
			},
		},
		{
			// -92233720368.54775808
			name:     fmt.Sprintf("%d.%d", minInt, -(minFrac)),
			expected: cadence.Fix64(-9223372036854775808),
			encodedData: []byte{
				// language=json, format=json-cadence data interchange format
				// {"type": "Fix64", "value": "-92233720368.54775808"}
				//
				// language=edn, format=ccf
				// 130([137(22), -9223372036854775808])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// fix64 type ID (22)
				0x16,
				// -9223372036854775808
				0x3b, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
			},
		},
		{
			// -92233720368.54775809
			name: fmt.Sprintf("%d.%d", minInt, -(minFrac - 1)),
			encodedData: []byte{
				// language=json, format=json-cadence data interchange format
				// {"type": "Fix64", "value": "-92233720368.54775809"}
				//
				// language=edn, format=ccf
				// 130([137(22), -9223372036854775809])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// fix64 type ID (22)
				0x16,
				// -9223372036854775809
				0x3b, 0x80, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			check: func(t *testing.T, actual cadence.Value, err error) {
				assert.Error(t, err)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ccf.Decode(nil, tt.encodedData)
			if tt.check != nil {
				tt.check(t, actual, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

func TestExportRecursiveType(t *testing.T) {

	t.Parallel()

	ty := &cadence.ResourceType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "Foo",
		Fields: []cadence.Field{
			{
				Identifier: "foo",
			},
		},
	}

	ty.Fields[0].Type = &cadence.OptionalType{
		Type: ty,
	}

	testEncode(
		t,
		cadence.Resource{
			Fields: []cadence.Value{
				cadence.Optional{},
			},
		}.WithType(ty),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Resource","value":{"id":"S.test.Foo","fields":[{"name":"foo","value":{"type": "Optional","value":null}}]}}
			//
			// language=edn, format=ccf
			// 129([[161([h'', "S.test.Foo", [["foo", 138(136(h''))]]])], [136(h''), [null]]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definition
			// array, 1 items follow
			0x81,
			// resource type:
			// id: []byte{}
			// cadence-type-id: "S.test.Foo"
			// 1 fields: [["foo", optional(type ref id(0))]]
			// tag
			0xd8, ccf.CBORTagResourceType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 10 bytes follow
			0x6a,
			// S.test.Foo
			0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f,
			// fields
			// array, 1 items follow
			0x81,
			// field 0
			// array, 2 items follow
			0x82,
			// text, 3 bytes follow
			0x63,
			// foo
			0x66, 0x6f, 0x6f,
			// tag
			0xd8, ccf.CBORTagOptionalType,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 0 bytes follow
			0x40,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 0 bytes follow
			0x40,
			// array, 1 items follow
			0x81,
			// nil
			0xf6,
		},
	)

}

func TestExportTypeValueRecursiveType(t *testing.T) {

	t.Parallel()

	t.Run("recursive", func(t *testing.T) {

		t.Parallel()

		ty := &cadence.ResourceType{
			Location:            utils.TestLocation,
			QualifiedIdentifier: "Foo",
			Fields: []cadence.Field{
				{
					Identifier: "foo",
				},
			},
			Initializers: [][]cadence.Parameter{},
		}

		ty.Fields[0].Type = &cadence.OptionalType{
			Type: ty,
		}

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: ty,
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":{"kind":"Resource","typeID":"S.test.Foo","fields":[{"id":"foo","type":{"kind":"Optional","type":"S.test.Foo"}}],"initializers":[],"type":""}}}
				//
				// language=edn, format=ccf
				// 130([137(41), 209([h'', "S.test.Foo", null, [["foo", 186(184(h''))]], []])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagResourceTypeValue,
				// array, 4 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 10 bytes follow
				0x6a,
				// S.test.Foo
				0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6F, 0x6F,
				// type (nil for struct)
				0xf6,
				// fields
				// array, 1 element follows
				0x81,
				// array, 2 elements follow
				0x82,
				// string, 3 bytes follow
				0x63,
				// foo
				0x66, 0x6f, 0x6f,
				// tag
				0xd8, ccf.CBORTagOptionalTypeValue,
				// tag
				0xd8, ccf.CBORTagTypeValueRef,
				// bytes, 0 bytes follow
				0x40,
				// initializers
				// array, 0 elements follow
				0x80,
			},
		)

	})

	t.Run("non-recursive, repeated", func(t *testing.T) {

		t.Parallel()

		fooTy := &cadence.ResourceType{
			Location:            utils.TestLocation,
			QualifiedIdentifier: "Foo",
			Fields:              []cadence.Field{},
			Initializers:        [][]cadence.Parameter{},
		}

		barTy := &cadence.ResourceType{
			Location:            utils.TestLocation,
			QualifiedIdentifier: "Bar",
			Fields: []cadence.Field{
				{
					Identifier: "foo1",
					Type:       fooTy,
				},
				{
					Identifier: "foo2",
					Type:       fooTy,
				},
			},
			Initializers: [][]cadence.Parameter{},
		}

		testEncodeAndDecode(
			t,
			cadence.TypeValue{
				StaticType: barTy,
			},
			[]byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":{"kind":"Resource","typeID":"S.test.Bar","fields":[{"id":"foo1","type":{"kind":"Resource","typeID":"S.test.Foo","fields":[],"initializers":[],"type":""}},{"id":"foo2","type":"S.test.Foo"}],"initializers":[],"type":""}}}
				//
				// language=edn, format=ccf
				// 130([137(41), 209([h'', "S.test.Bar", null, [["foo1", 209([h'01', "S.test.Foo", null, [], []])], ["foo2", 184(h'01')]], []])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagResourceTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 10 bytes follow
				0x6a,
				// S.test.Bar
				0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x42, 0x61, 0x72,
				// type (nil for struct)
				0xf6,
				// fields
				// array, 2 element follows
				0x82,
				// array, 2 elements follow
				0x82,
				// string, 4 bytes follow
				0x64,
				// foo1
				0x66, 0x6f, 0x6f, 0x31,
				// tag
				0xd8, ccf.CBORTagResourceTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 1 bytes follow
				0x41,
				// 1
				0x01,
				// string, 10 bytes follow
				0x6a,
				// S.test.Foo
				0x53, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x2e, 0x46, 0x6f, 0x6f,
				// type (nil for struct)
				0xf6,
				// fields
				// array, 0 elements follow
				0x80,
				// initializer
				// array, 0 elements follow
				0x80,
				// array, 2 elements follow
				0x82,
				// string, 4 bytes follow
				0x64,
				// foo2
				0x66, 0x6f, 0x6f, 0x32,
				// tag
				0xd8, ccf.CBORTagTypeValueRef,
				// bytes, 1 bytes follow
				0x41,
				// 1
				0x01,
				// initializers
				// array, 0 elements follow
				0x80,
			},
		)
	})
}

func TestEncodePath(t *testing.T) {

	t.Parallel()

	testEncodeAndDecode(
		t,
		cadence.NewPath("storage", "foo"),
		[]byte{ // language=json, format=json-cadence data interchange format
			// {"type":"Path","value":{"domain":"storage","identifier":"foo"}}
			//
			// language=edn, format=ccf
			// 130([137(24), [1, "foo"]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 elements follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Path type ID (24)
			0x18, 0x18,
			// array, 2 elements follow
			0x82,
			// 1
			0x01,
			// string, 3 bytes follow
			0x63,
			// foo
			0x66, 0x6f, 0x6f,
		},
	)
}

func testAllEncodeAndDecode(t *testing.T, tests ...encodeTest) {

	test := func(testCase encodeTest) {

		t.Run(testCase.name, func(t *testing.T) {

			t.Parallel()

			testEncodeAndDecode(t, testCase.val, testCase.expected)
		})
	}

	for _, testCase := range tests {
		test(testCase)
	}
}

func TestDecodeInvalidType(t *testing.T) {

	t.Parallel()

	t.Run("empty type", func(t *testing.T) {
		t.Parallel()

		encodedData := []byte{
			// language=json, format=json-cadence data interchange format
			// { "type":"Struct", "value":{ "id":"", "fields":[] } }
			//
			// language=edn, format=ccf
			// 129([[160([h'', "", []])], [136(h''), []]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definition
			// array, 1 items follow
			0x81,
			// struct type:
			// id: []byte{}
			// cadence-type-id: ""
			// 0 field
			// tag
			0xd8, ccf.CBORTagStructType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 0 bytes follow
			0x60,
			// fields
			// array, 0 items follow
			0x80,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 0 bytes follow
			0x40,
			// array, 0 items follow
			0x80,
		}
		_, err := ccf.Decode(nil, encodedData)
		require.Error(t, err)
		assert.Equal(t, "ccf: failed to decode: invalid type ID for built-in: ``", err.Error())
	})

	t.Run("undefined type", func(t *testing.T) {
		t.Parallel()

		encodedData := []byte{
			// language=json, format=json-cadence data interchange format
			// { "type":"Struct", "value":{ "id":"I.Foo", "fields":[] } }
			//
			// language=edn, format=ccf
			// 129([[160([h'', "I.Foo", []])], [136(h''), []]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definition
			// array, 1 items follow
			0x81,
			// struct type:
			// id: []byte{}
			// cadence-type-id: "I.Foo"
			// 0 field
			// tag
			0xd8, ccf.CBORTagStructType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 5 bytes follow
			0x65,
			// I.Foo
			0x49, 0x2e, 0x46, 0x6f, 0x6f,
			// fields
			// array, 0 items follow
			0x80,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 0 bytes follow
			0x40,
			// array, 0 items follow
			0x80,
		}
		_, err := ccf.Decode(nil, encodedData)
		require.Error(t, err)
		assert.Equal(t, "ccf: failed to decode: invalid type ID `I.Foo`: invalid identifier location type ID: missing qualified identifier", err.Error())
	})

	t.Run("unknown location prefix", func(t *testing.T) {
		t.Parallel()

		encodedData := []byte{
			// language=json, format=json-cadence data interchange format
			// { "type":"Struct", "value":{ "id":"N.PublicKey", "fields":[] } }
			//
			// language=edn, format=ccf
			// 129([[160([h'', "N.PublicKey", []])], [136(h''), []]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeDefAndValue,
			// array, 2 items follow
			0x82,
			// element 0: type definitions
			// array, 1 items follow
			0x81,
			// struct type:
			// id: []byte{}
			// cadence-type-id: "N.PublicKey"
			// 0 field
			// tag
			0xd8, ccf.CBORTagStructType,
			// array, 3 items follow
			0x83,
			// id
			// bytes, 0 bytes follow
			0x40,
			// cadence-type-id
			// string, 11 bytes follow
			0x6b,
			// N.PublicKey
			0x4e, 0x2e, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x4b, 0x65, 0x79,
			// fields
			// array, 0 items follow
			0x80,

			// element 1: type and value
			// array, 2 items follow
			0x82,
			// tag
			0xd8, ccf.CBORTagTypeRef,
			// bytes, 0 bytes follow
			0x40,
			// array, 0 items follow
			0x80,
		}
		_, err := ccf.Decode(nil, encodedData)
		require.Error(t, err)
		assert.Equal(t, "ccf: failed to decode: invalid type ID for built-in: `N.PublicKey`", err.Error())
	})
}

func testEncodeAndDecode(t *testing.T, val cadence.Value, expectedCBOR []byte) {
	actualCBOR := testEncode(t, val, expectedCBOR)
	testDecode(t, actualCBOR, val)
}

// testEncodeAndDecodeEx is used when val != expectedVal because of deterministic encoding.
func testEncodeAndDecodeEx(t *testing.T, val cadence.Value, expectedCBOR []byte, expectedVal cadence.Value) {
	actualCBOR := testEncode(t, val, expectedCBOR)
	testDecode(t, actualCBOR, expectedVal)
}

func testEncode(t *testing.T, val cadence.Value, expectedCBOR []byte) (actualCBOR []byte) {
	actualCBOR, err := ccf.Encode(val)
	require.NoError(t, err)
	assert.True(t, bytes.Equal(expectedCBOR, actualCBOR), fmt.Sprintf("actual: 0x%x", actualCBOR))
	return actualCBOR
}

func testDecode(t *testing.T, actualCBOR []byte, expectedVal cadence.Value) {
	decodedVal, err := ccf.Decode(nil, actualCBOR)
	require.NoError(t, err)
	assert.Equal(
		t,
		cadence.ValueWithCachedTypeID(expectedVal),
		cadence.ValueWithCachedTypeID(decodedVal),
	)
}

var fooResourceType = &cadence.ResourceType{
	Location:            utils.TestLocation,
	QualifiedIdentifier: "Foo",
	Fields: []cadence.Field{
		{
			Identifier: "bar",
			Type:       cadence.IntType{},
		},
	},
}

var foooResourceTypeWithAbstractField = &cadence.ResourceType{
	Location:            utils.TestLocation,
	QualifiedIdentifier: "Fooo",
	Fields: []cadence.Field{
		{
			Identifier: "bar",
			Type:       cadence.IntType{},
		},
		{
			Identifier: "baz",
			Type:       cadence.AnyStructType{},
		},
	},
}

func TestEncodeBuiltinComposites(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		typ     cadence.Type
		encoded []byte
	}{
		{
			name: "Struct",
			typ: &cadence.StructType{
				Location:            nil,
				QualifiedIdentifier: "Foo",
			},
			encoded: []byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":{"kind":"Struct","typeID":"Foo","fields":[],"initializers":[],"type":""}}}
				//
				// language=edn, format=ccf
				// 130([137(41), 208([h'', "Foo", null, [], []])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagStructTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 3 bytes follow
				0x63,
				// Foo
				0x46, 0x6f, 0x6f,
				// type (nil for event)
				0xf6,
				// fields
				// array, 0 element follows
				0x80,
				// initializers
				// array, 0 elements follow
				0x80,
			},
		},
		{
			name: "StructInterface",
			typ: &cadence.StructInterfaceType{
				Location:            nil,
				QualifiedIdentifier: "Foo",
			},
			encoded: []byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":{"kind":"StructInterface","typeID":"Foo","fields":[],"initializers":[],"type":""}}}
				//
				// language=edn, format=ccf
				// 130([137(41), 224([h'', "Foo", null, [], []])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagStructInterfaceTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 3 bytes follow
				0x63,
				// Foo
				0x46, 0x6f, 0x6f,
				// type (nil for event)
				0xf6,
				// fields
				// array, 0 element follows
				0x80,
				// initializers
				// array, 0 elements follow
				0x80,
			},
		},
		{
			name: "Resource",
			typ: &cadence.ResourceType{
				Location:            nil,
				QualifiedIdentifier: "Foo",
			},
			encoded: []byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":{"kind":"Resource","typeID":"Foo","fields":[],"initializers":[],"type":""}}}
				//
				// language=edn, format=ccf
				// 130([137(41), 209([h'', "Foo", null, [], []])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagResourceTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 3 bytes follow
				0x63,
				// Foo
				0x46, 0x6f, 0x6f,
				// type (nil for event)
				0xf6,
				// fields
				// array, 0 element follows
				0x80,
				// initializers
				// array, 0 elements follow
				0x80,
			},
		},
		{
			name: "ResourceInterface",
			typ: &cadence.ResourceInterfaceType{
				Location:            nil,
				QualifiedIdentifier: "Foo",
			},
			encoded: []byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":{"kind":"ResourceInterface","typeID":"Foo","fields":[],"initializers":[],"type":""}}}
				//
				// language=edn, format=ccf
				// 130([137(41), 225([h'', "Foo", null, [], []])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagResourceInterfaceTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 3 bytes follow
				0x63,
				// Foo
				0x46, 0x6f, 0x6f,
				// type (nil for event)
				0xf6,
				// fields
				// array, 0 element follows
				0x80,
				// initializers
				// array, 0 elements follow
				0x80,
			},
		},
		{
			name: "Contract",
			typ: &cadence.ContractType{
				Location:            nil,
				QualifiedIdentifier: "Foo",
			},
			encoded: []byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":{"kind":"Contract","typeID":"Foo","fields":[],"initializers":[],"type":""}}}
				//
				// language=edn, format=ccf
				// 130([137(41), 211([h'', "Foo", null, [], []])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagContractTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 3 bytes follow
				0x63,
				// Foo
				0x46, 0x6f, 0x6f,
				// type (nil for event)
				0xf6,
				// fields
				// array, 0 element follows
				0x80,
				// initializers
				// array, 0 elements follow
				0x80,
			},
		},
		{
			name: "ContractInterface",
			typ: &cadence.ContractInterfaceType{
				Location:            nil,
				QualifiedIdentifier: "Foo",
			},
			encoded: []byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":{"kind":"ContractInterface","typeID":"Foo","fields":[],"initializers":[],"type":""}}}
				//
				// language=edn, format=ccf
				// 130([137(41), 226([h'', "Foo", null, [], []])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagContractInterfaceTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 3 bytes follow
				0x63,
				// Foo
				0x46, 0x6f, 0x6f,
				// type (nil for event)
				0xf6,
				// fields
				// array, 0 element follows
				0x80,
				// initializers
				// array, 0 elements follow
				0x80,
			},
		},
		{
			name: "Enum",
			typ: &cadence.EnumType{
				Location:            nil,
				QualifiedIdentifier: "Foo",
			},
			encoded: []byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":{"kind":"Enum","typeID":"Foo","fields":[],"initializers":[],"type":""}}}
				//
				// language=edn, format=ccf
				// 130([137(41), 212([h'', "Foo", null, [], []])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagEnumTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 3 bytes follow
				0x63,
				// Foo
				0x46, 0x6f, 0x6f,
				// type (nil for event)
				0xf6,
				// fields
				// array, 0 element follows
				0x80,
				// initializers
				// array, 0 elements follow
				0x80,
			},
		},
		{
			name: "Event",
			typ: &cadence.EventType{
				Location:            nil,
				QualifiedIdentifier: "Foo",
			},
			encoded: []byte{ // language=json, format=json-cadence data interchange format
				// {"type":"Type","value":{"staticType":{"kind":"Event","typeID":"Foo","fields":[],"initializers":[],"type":""}}}
				//
				// language=edn, format=ccf
				// 130([137(41), 210([h'', "Foo", null, [], [[]]])])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeAndValue,
				// array, 2 elements follow
				0x82,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Meta type ID (41)
				0x18, 0x29,
				// tag
				0xd8, ccf.CBORTagEventTypeValue,
				// array, 5 elements follow
				0x85,
				// bytes, 0 bytes follow
				0x40,
				// string, 3 bytes follow
				0x63,
				// Foo
				0x46, 0x6f, 0x6f,
				// type (nil for event)
				0xf6,
				// fields
				// array, 0 element follows
				0x80,
				// initializers
				// array, 1 elements follow
				0x81,
				// array, 0 element follow
				0x80,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			typeValue := cadence.NewTypeValue(test.typ)
			testEncode(t, typeValue, test.encoded)
		})
	}
}

func TestExportFunctionValue(t *testing.T) {

	t.Parallel()

	ty := &cadence.ResourceType{
		Location:            utils.TestLocation,
		QualifiedIdentifier: "Foo",
		Fields: []cadence.Field{
			{
				Identifier: "foo",
			},
		},
	}

	ty.Fields[0].Type = &cadence.OptionalType{
		Type: ty,
	}

	testEncode(
		t,
		cadence.Function{
			FunctionType: (&cadence.FunctionType{
				Parameters: []cadence.Parameter{},
				ReturnType: cadence.VoidType{},
			}).WithID("(():Void)"),
		},
		[]byte{ // language=json, format=json-cadence data interchange format
			// { "type": "Function", "value": { "functionType": { "kind": "Function", "typeID": "(():Void)", "parameters": [], "return": { "kind": "Void" } } } }
			//
			// language=edn, format=ccf
			// 130([137(51), ["(():Void)", [], 185(50)]])
			//
			// language=cbor, format=ccf
			// tag
			0xd8, ccf.CBORTagTypeAndValue,
			// array, 2 elements follow
			0x82,
			// tag
			0xd8, ccf.CBORTagSimpleType,
			// Function type ID (51)
			0x18, 0x33,
			// array, 3 elements follow
			0x83,
			// element 0: cadence-type-id
			// string, 9 bytes follow
			0x69,
			// (():Void)
			0x28, 0x28, 0x29, 0x3a, 0x56, 0x6f, 0x69, 0x64, 0x29,
			// element 1: parameters
			// array, 0 element
			0x80,
			// element 2: return type
			// tag
			0xd8, ccf.CBORTagSimpleTypeValue,
			// Void type ID (50)
			0x18, 0x32,
		},
	)
}

func TestDeployedEvents(t *testing.T) {
	var tests = []struct {
		name         string
		event        cadence.Event
		expectedCBOR []byte
	}{
		{
			name:  "FlowFees.FeesDeducted",
			event: createFlowFeesFeesDeductedEvent(),
			expectedCBOR: []byte{
				// language=json, format=json-cadence data interchange format
				// {"value":{"id":"A.f919ee77447b7497.FlowFees.FeesDeducted","fields":[{"value":{"value":"0.01797293","type":"UFix64"},"name":"amount"},{"value":{"value":"1.00000000","type":"UFix64"},"name":"inclusionEffort"},{"value":{"value":"0.00360123","type":"UFix64"},"name":"executionEffort"}]},"type":"Event"}
				//
				// language=edn, format=ccf
				// 129([[162([h'', "A.f919ee77447b7497.FlowFees.FeesDeducted", [["amount", 137(23)], ["executionEffort", 137(23)], ["inclusionEffort", 137(23)]]])], [136(h''), [1797293, 360123, 100000000]]])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeDefAndValue,
				// array, 2 element follows
				0x82,
				// element 0: type definitions
				// array, 1 element follows
				0x81,
				// event type:
				// id: []byte{}
				// cadence-type-id: "A.f919ee77447b7497.FlowFees.FeesDeducted"
				// 3 fields: [["amount", type(ufix64)], ["executionEffort", type(ufix64)], ["inclusionEffort", type(ufix64)]]
				// tag
				0xd8, ccf.CBORTagEventType,
				// array, 3 elements follow
				0x83,
				// id
				// bytes, 0 bytes follow
				0x40,
				// cadence-type-id
				// string, 40 bytes follow
				0x78, 0x28,
				// A.f919ee77447b7497.FlowFees.FeesDeducted
				0x41, 0x2e, 0x66, 0x39, 0x31, 0x39, 0x65, 0x65, 0x37, 0x37, 0x34, 0x34, 0x37, 0x62, 0x37, 0x34, 0x39, 0x37, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x46, 0x65, 0x65, 0x73, 0x2e, 0x46, 0x65, 0x65, 0x73, 0x44, 0x65, 0x64, 0x75, 0x63, 0x74, 0x65, 0x64,
				// fields
				// array, 3 items follow
				0x83,
				// field 0
				// array, 2 items follow
				0x82,
				// text, 6 bytes follow
				0x66,
				// amount
				0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,
				// field 1
				// array, 2 items follow
				0x82,
				// text, 15 bytes follow
				0x6f,
				// executionEffort
				0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x66, 0x66, 0x6f, 0x72, 0x74,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,
				// field 2
				// array, 2 items follow
				0x82,
				// text, 15 bytes follow
				0x6f,
				// inclusionEffort
				0x69, 0x6e, 0x63, 0x6c, 0x75, 0x73, 0x69, 0x6f, 0x6e, 0x45, 0x66, 0x66, 0x6f, 0x72, 0x74,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,

				// element 1: type and value
				// array, 2 items follow
				0x82,
				// tag
				0xd8, ccf.CBORTagTypeRef,
				// bytes, 0 bytes follow
				0x40,
				// array, 3 items follow
				0x83,
				// 1797293
				0x1a, 0x00, 0x1b, 0x6c, 0xad,
				// 360123
				0x1a, 0x00, 0x05, 0x7e, 0xbb,
				// 100000000
				0x1a, 0x05, 0xf5, 0xe1, 0x00,
			},
		},
		{
			name:  "FlowFees.TokensWithdrawn",
			event: createFlowFeesTokensWithdrawnEvent(),
			expectedCBOR: []byte{
				// language=json, format=json-cadence data interchange format
				// {"value":{"id":"A.f919ee77447b7497.FlowFees.TokensWithdrawn","fields":[{"value":{"value":"53.04112895","type":"UFix64"},"name":"amount"}]},"type":"Event"}
				//
				// language=edn, format=ccf
				// 129([[162([h'', "A.f919ee77447b7497.FlowFees.TokensWithdrawn", [["amount", 137(23)]]])], [136(h''), [5304112895]]])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeDefAndValue,
				// array, 2 element follows
				0x82,
				// element 0: type definitions
				// array, 1 element follows
				0x81,
				// event type:
				// id: []byte{}
				// cadence-type-id: "A.f919ee77447b7497.FlowFees.TokensWithdrawn"
				// 1 field: [["amount", type(ufix64)]]
				// tag
				0xd8, ccf.CBORTagEventType,
				// array, 3 element follows
				0x83,
				// id
				// bytes, 0 bytes follow
				0x40,
				// cadence-type-id
				// string, 43 bytes follow
				0x78, 0x2b,
				// "A.f919ee77447b7497.FlowFees.TokensWithdrawn"
				0x41, 0x2e, 0x66, 0x39, 0x31, 0x39, 0x65, 0x65, 0x37, 0x37, 0x34, 0x34, 0x37, 0x62, 0x37, 0x34, 0x39, 0x37, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x46, 0x65, 0x65, 0x73, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x57, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x6e,
				// fields
				// array, 1 items follow
				0x81,
				// field 0
				// array, 2 element follows
				0x82,
				// text, 6 bytes follow
				0x66,
				// "amount"
				0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,

				// element 1: type and value
				// array, 2 element follows
				0x82,
				// tag
				0xd8, ccf.CBORTagTypeRef,
				// bytes, 0 bytes follow
				0x40,
				// array, 1 items follow
				0x81,
				// 5304112895
				0x1b, 0x00, 0x00, 0x00, 0x01, 0x3c, 0x26, 0x56, 0xff,
			},
		},
		{
			name:  "FlowIDTableStaking.DelegatorRewardsPaid",
			event: createFlowIDTableStakingDelegatorRewardsPaidEvent(),
			expectedCBOR: []byte{
				// language=json, format=json-cadence data interchange format
				// {"value":{"id":"A.8624b52f9ddcd04a.FlowIDTableStaking.DelegatorRewardsPaid","fields":[{"value":{"value":"e52cbcd825e328acac8db6bcbdcbb6e7724862c8b89b09d85edccf41ff9981eb","type":"String"},"name":"nodeID"},{"value":{"value":"92","type":"UInt32"},"name":"delegatorID"},{"value":{"value":"4.38760261","type":"UFix64"},"name":"amount"}]},"type":"Event"}
				//
				// language=edn, format=ccf
				// 129([[162([h'', "A.8624b52f9ddcd04a.FlowIDTableStaking.DelegatorRewardsPaid", [["amount", 137(23)], ["nodeID", 137(1)], ["delegatorID", 137(14)]]])], [136(h''), [438760261, "e52cbcd825e328acac8db6bcbdcbb6e7724862c8b89b09d85edccf41ff9981eb", 92]]])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeDefAndValue,
				// array, 2 element follows
				0x82,
				// element 0: type definitions
				// array, 1 element follows
				0x81,

				// event type:
				// id: []byte{}
				// cadence-type-id: "A.8624b52f9ddcd04a.FlowIDTableStaking.DelegatorRewardsPaid"
				// 3 field: [["amount", type(ufix64)], ["nodeID", type(string)], ["delegatorID", type(uint32)]]
				// tag
				0xd8, ccf.CBORTagEventType,
				// array, 3 element follows
				0x83,
				// id
				// bytes, 0 bytes follow
				0x40,
				// cadence-type-id
				// string, 58 bytes follow
				0x78, 0x3a,
				// "A.8624b52f9ddcd04a.FlowIDTableStaking.DelegatorRewardsPaid"
				0x41, 0x2e, 0x38, 0x36, 0x32, 0x34, 0x62, 0x35, 0x32, 0x66, 0x39, 0x64, 0x64, 0x63, 0x64, 0x30, 0x34, 0x61, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x49, 0x44, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x74, 0x61, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x73, 0x50, 0x61, 0x69, 0x64,
				// fields
				// array, 3 items follow
				0x83,
				// field 0
				// array, 2 element follows
				0x82,
				// text, 6 bytes follow
				0x66,
				// "amount"
				0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,
				// field 1
				// array, 2 element follows
				0x82,
				// text, 6 bytes follow
				0x66,
				// "nodeID"
				0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// String type ID (1)
				0x01,
				// field 2
				// array, 2 element follows
				0x82,
				// text, 11 bytes follow
				0x6b,
				// "delegatorID"
				0x64, 0x65, 0x6c, 0x65, 0x67, 0x61, 0x74, 0x6f, 0x72, 0x49, 0x44,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UInt32 type ID (14)
				0x0e,

				// element 1: type and value
				// array, 2 element follows
				0x82,
				// tag
				0xd8, ccf.CBORTagTypeRef,
				// bytes, 0 bytes follow
				0x40,
				// array, 3 items follow
				0x83,
				// 438760261
				0x1a, 0x1a, 0x26, 0xf3, 0x45,
				// text, 64 bytes follow
				0x78, 0x40,
				// "e52cbcd825e328acac8db6bcbdcbb6e7724862c8b89b09d85edccf41ff9981eb"
				0x65, 0x35, 0x32, 0x63, 0x62, 0x63, 0x64, 0x38, 0x32, 0x35, 0x65, 0x33, 0x32, 0x38, 0x61, 0x63, 0x61, 0x63, 0x38, 0x64, 0x62, 0x36, 0x62, 0x63, 0x62, 0x64, 0x63, 0x62, 0x62, 0x36, 0x65, 0x37, 0x37, 0x32, 0x34, 0x38, 0x36, 0x32, 0x63, 0x38, 0x62, 0x38, 0x39, 0x62, 0x30, 0x39, 0x64, 0x38, 0x35, 0x65, 0x64, 0x63, 0x63, 0x66, 0x34, 0x31, 0x66, 0x66, 0x39, 0x39, 0x38, 0x31, 0x65, 0x62,
				// 92
				0x18, 0x5c,
			},
		},
		{
			name:  "FlowIDTableStaking.EpochTotalRewardsPaid",
			event: createFlowIDTableStakingEpochTotalRewardsPaidEvent(),
			expectedCBOR: []byte{
				// language=json, format=json-cadence data interchange format
				// {"value":{"id":"A.8624b52f9ddcd04a.FlowIDTableStaking.EpochTotalRewardsPaid","fields":[{"value":{"value":"1316543.00000000","type":"UFix64"},"name":"total"},{"value":{"value":"53.04112895","type":"UFix64"},"name":"fromFees"},{"value":{"value":"1316489.95887105","type":"UFix64"},"name":"minted"},{"value":{"value":"6.04080767","type":"UFix64"},"name":"feesBurned"}]},"type":"Event"}
				//
				// language=edn, format=ccf
				// 129([[162([h'', "A.8624b52f9ddcd04a.FlowIDTableStaking.EpochTotalRewardsPaid", [["total", 137(23)], ["minted", 137(23)], ["fromFees", 137(23)], ["feesBurned", 137(23)]]])], [136(h''), [131654300000000, 131648995887105, 5304112895, 604080767]]])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeDefAndValue,
				// array, 2 element follows
				0x82,
				// element 0: type definitions
				// array, 1 element follows
				0x81,
				// event type:
				// id: []byte{}
				// cadence-type-id: "A.8624b52f9ddcd04a.FlowIDTableStaking.EpochTotalRewardsPaid"
				// 4 field: [["total", type(ufix64)], ["minted", type(ufix64)], ["fromFees", type(ufix64)], ["feesBurned", type(ufix64)]]
				// tag
				0xd8, ccf.CBORTagEventType,
				// array, 3 element follows
				0x83,
				// id
				// bytes, 0 bytes follow
				0x40,
				// cadence-type-id
				// string, 59 bytes follow
				0x78, 0x3b,
				// "A.8624b52f9ddcd04a.FlowIDTableStaking.EpochTotalRewardsPaid"
				0x41, 0x2e, 0x38, 0x36, 0x32, 0x34, 0x62, 0x35, 0x32, 0x66, 0x39, 0x64, 0x64, 0x63, 0x64, 0x30, 0x34, 0x61, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x49, 0x44, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x74, 0x61, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x73, 0x50, 0x61, 0x69, 0x64,
				// fields
				// array, 4 items follow
				0x84,
				// field 0
				// array, 2 element follows
				0x82,
				// text, 5 bytes follow
				0x65,
				// "total"
				0x74, 0x6f, 0x74, 0x61, 0x6c,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,
				// field 1
				// array, 2 element follows
				0x82,
				// text, 6 bytes follow
				0x66,
				// "minted"
				0x6d, 0x69, 0x6e, 0x74, 0x65, 0x64,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,
				// field 2
				// array, 2 element follows
				0x82,
				// text, 8 bytes follow
				0x68,
				// "fromFees"
				0x66, 0x72, 0x6f, 0x6d, 0x46, 0x65, 0x65, 0x73,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,
				// field 3
				// array, 2 element follows
				0x82,
				// text, 10 bytes follow
				0x6a,
				// "feesBurned"
				0x66, 0x65, 0x65, 0x73, 0x42, 0x75, 0x72, 0x6e, 0x65, 0x64,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,

				// element 1: type and value
				// array, 2 element follows
				0x82,
				// tag
				0xd8, ccf.CBORTagTypeRef,
				// bytes, 0 bytes follow
				0x40,
				// array, 4 items follow
				0x84,
				// 131654300000000
				0x1b, 0x00, 0x00, 0x77, 0xbd, 0x27, 0xc8, 0xdf, 0x00,
				// 131648995887105
				0x1b, 0x00, 0x00, 0x77, 0xbb, 0xeb, 0xa2, 0x88, 0x01,
				// 5304112895
				0x1b, 0x00, 0x00, 0x00, 0x01, 0x3c, 0x26, 0x56, 0xff,
				// 604080767
				0x1a, 0x24, 0x01, 0x8a, 0x7f,
			},
		},
		{
			name:  "FlowIDTableStaking.NewWeeklyPayout",
			event: createFlowIDTableStakingNewWeeklyPayoutEvent(),
			expectedCBOR: []byte{
				// language=json, format=json-cadence data interchange format
				// {"value":{"id":"A.8624b52f9ddcd04a.FlowIDTableStaking.NewWeeklyPayout","fields":[{"value":{"value":"1317778.00000000","type":"UFix64"},"name":"newPayout"}]},"type":"Event"}
				//
				// language=edn, format=ccf
				// 129([[162([h'', "A.8624b52f9ddcd04a.FlowIDTableStaking.NewWeeklyPayout", [["newPayout", 137(23)]]])], [136(h''), [131777800000000]]])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeDefAndValue,
				// array, 2 element follows
				0x82,
				// element 0: type definitions
				// array, 1 element follows
				0x81,
				// event type:
				// id: []byte{}
				// cadence-type-id: "A.8624b52f9ddcd04a.FlowIDTableStaking.NewWeeklyPayout"
				// 1 field: [["newPayout", type(ufix64)]]
				// tag
				0xd8, ccf.CBORTagEventType,
				// array, 3 element follows
				0x83,
				// id
				// bytes, 0 bytes follow
				0x40,
				// cadence-type-id
				// string, 53 bytes follow
				0x78, 0x35,
				// "A.8624b52f9ddcd04a.FlowIDTableStaking.NewWeeklyPayout"
				0x41, 0x2e, 0x38, 0x36, 0x32, 0x34, 0x62, 0x35, 0x32, 0x66, 0x39, 0x64, 0x64, 0x63, 0x64, 0x30, 0x34, 0x61, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x49, 0x44, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x74, 0x61, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x4e, 0x65, 0x77, 0x57, 0x65, 0x65, 0x6b, 0x6c, 0x79, 0x50, 0x61, 0x79, 0x6f, 0x75, 0x74,
				// fields
				// array, 1 items follow
				0x81,
				// field 0
				// array, 2 element follows
				0x82,
				// text, 9 bytes follow
				0x69,
				// "newPayout"
				0x6e, 0x65, 0x77, 0x50, 0x61, 0x79, 0x6f, 0x75, 0x74,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,

				// element 1: type and value
				// array, 2 element follows
				0x82,
				// tag
				0xd8, ccf.CBORTagTypeRef,
				// bytes, 0 bytes follow
				0x40,
				// array, 1 items follow
				0x81,
				// 131777800000000
				0x1b, 0x00, 0x00, 0x77, 0xd9, 0xe8, 0xf5, 0x52, 0x00,
			},
		},
		{
			name:  "FlowIDTableStaking.RewardsPaid",
			event: createFlowIDTableStakingRewardsPaidEvent(),
			expectedCBOR: []byte{
				// language=json, format=json-cadence data interchange format
				// {"value":{"id":"A.8624b52f9ddcd04a.FlowIDTableStaking.RewardsPaid","fields":[{"value":{"value":"e52cbcd825e328acac8db6bcbdcbb6e7724862c8b89b09d85edccf41ff9981eb","type":"String"},"name":"nodeID"},{"value":{"value":"1745.49955740","type":"UFix64"},"name":"amount"}]},"type":"Event"}
				//
				// language=edn, format=ccf
				// 129([[162([h'', "A.8624b52f9ddcd04a.FlowIDTableStaking.RewardsPaid", [["amount", 137(23)], ["nodeID", 137(1)]]])], [136(h''), [174549955740, "e52cbcd825e328acac8db6bcbdcbb6e7724862c8b89b09d85edccf41ff9981eb"]]])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeDefAndValue,
				// array, 2 element follows
				0x82,
				// element 0: type definitions
				// array, 1 element follows
				0x81,
				// event type:
				// id: []byte{}
				// cadence-type-id: "A.8624b52f9ddcd04a.FlowIDTableStaking.RewardsPaid"
				// 2 field: [["amount", type(ufix64)], ["nodeID", type(string)]]
				// tag
				0xd8, ccf.CBORTagEventType,
				// array, 3 element follows
				0x83,
				// id
				// bytes, 0 bytes follow
				0x40,
				// cadence-type-id
				// string, 49 bytes follow
				0x78, 0x31,
				// "A.8624b52f9ddcd04a.FlowIDTableStaking.RewardsPaid"
				0x41, 0x2e, 0x38, 0x36, 0x32, 0x34, 0x62, 0x35, 0x32, 0x66, 0x39, 0x64, 0x64, 0x63, 0x64, 0x30, 0x34, 0x61, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x49, 0x44, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x53, 0x74, 0x61, 0x6b, 0x69, 0x6e, 0x67, 0x2e, 0x52, 0x65, 0x77, 0x61, 0x72, 0x64, 0x73, 0x50, 0x61, 0x69, 0x64,
				// fields
				// array, 2 items follow
				0x82,
				// field 0
				// array, 2 element follows
				0x82,
				// text, 6 bytes follow
				0x66,
				// "amount"
				0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,
				// field 1
				// array, 2 element follows
				0x82,
				// text, 6 bytes follow
				0x66,
				// "nodeID"
				0x6e, 0x6f, 0x64, 0x65, 0x49, 0x44,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// String type ID (1)
				0x01,

				// element 1: type and value
				// array, 2 element follows
				0x82,
				// tag
				0xd8, ccf.CBORTagTypeRef,
				// bytes, 0 bytes follow
				0x40,
				// array, 2 items follow
				0x82,
				// 174549955740
				0x1b, 0x00, 0x00, 0x00, 0x28, 0xa3, 0xfc, 0xf4, 0x9c,
				// string, 64 bytes follow
				0x78, 0x40,
				// "e52cbcd825e328acac8db6bcbdcbb6e7724862c8b89b09d85edccf41ff9981eb"
				0x65, 0x35, 0x32, 0x63, 0x62, 0x63, 0x64, 0x38, 0x32, 0x35, 0x65, 0x33, 0x32, 0x38, 0x61, 0x63, 0x61, 0x63, 0x38, 0x64, 0x62, 0x36, 0x62, 0x63, 0x62, 0x64, 0x63, 0x62, 0x62, 0x36, 0x65, 0x37, 0x37, 0x32, 0x34, 0x38, 0x36, 0x32, 0x63, 0x38, 0x62, 0x38, 0x39, 0x62, 0x30, 0x39, 0x64, 0x38, 0x35, 0x65, 0x64, 0x63, 0x63, 0x66, 0x34, 0x31, 0x66, 0x66, 0x39, 0x39, 0x38, 0x31, 0x65, 0x62,
			},
		},
		{
			name:  "FlowToken.TokensDeposited with nil receiver",
			event: createFlowTokenTokensDepositedEventNoReceiver(),
			expectedCBOR: []byte{
				// language=json, format=json-cadence data interchange format
				// {"value":{"id":"A.1654653399040a61.FlowToken.TokensDeposited","fields":[{"value":{"value":"1316489.95887105","type":"UFix64"},"name":"amount"},{"value":{"value":null,"type":"Optional"},"name":"to"}]},"type":"Event"}
				//
				// language=edn, format=ccf
				// 129([[162([h'', "A.1654653399040a61.FlowToken.TokensDeposited", [["to", 138(137(3))], ["amount", 137(23)]]])], [136(h''), [null, 131648995887105]]])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeDefAndValue,
				// array, 2 element follows
				0x82,
				// element 0: type definitions
				// array, 1 element follows
				0x81,
				// event type:
				// id: []byte{}
				// cadence-type-id: "A.1654653399040a61.FlowToken.TokensDeposited"
				// 2 field: [["to", type(optional(address))], ["amount", type(ufix64)]]
				// tag
				0xd8, ccf.CBORTagEventType,
				// array, 3 element follows
				0x83,
				// id
				// bytes, 0 bytes follow
				0x40,
				// cadence-type-id
				// string, 44 bytes follow
				0x78, 0x2c,
				// "A.1654653399040a61.FlowToken.TokensDeposited"
				0x41, 0x2e, 0x31, 0x36, 0x35, 0x34, 0x36, 0x35, 0x33, 0x33, 0x39, 0x39, 0x30, 0x34, 0x30, 0x61, 0x36, 0x31, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x44, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x65, 0x64,
				// fields
				// array, 2 items follow
				0x82,
				// field 0
				// array, 2 element follows
				0x82,
				// text, 2 bytes follow
				0x62,
				// "to"
				0x74, 0x6f,
				// tag
				0xd8, ccf.CBORTagOptionalType,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Address type ID (3)
				0x03,
				// field 1
				// array, 2 element follows
				0x82,
				// text, 6 bytes follow
				0x66,
				// "amount"
				0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,

				// element 1: type and value
				// array, 2 element follows
				0x82,
				// tag
				0xd8, ccf.CBORTagTypeRef,
				// bytes, 0 bytes follow
				0x40,
				// array, 2 items follow
				0x82,
				// null
				0xf6,
				// 131648995887105
				0x1b, 0x00, 0x00, 0x77, 0xbb, 0xeb, 0xa2, 0x88, 0x01,
			},
		},
		{
			name:  "FlowToken.TokensDeposited",
			event: createFlowTokenTokensDepositedEvent(),
			expectedCBOR: []byte{
				// language=json, format=json-cadence data interchange format
				// {"value":{"id":"A.1654653399040a61.FlowToken.TokensDeposited","fields":[{"value":{"value":"1745.49955740","type":"UFix64"},"name":"amount"},{"value":{"value":{"value":"0x8624b52f9ddcd04a","type":"Address"},"type":"Optional"},"name":"to"}]},"type":"Event"}
				//
				// language=edn, format=ccf
				// 129([[162([h'', "A.1654653399040a61.FlowToken.TokensDeposited", [["to", 138(137(3))], ["amount", 137(23)]]])], [136(h''), [h'8624B52F9DDCD04A', 174549955740]]])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeDefAndValue,
				// array, 2 element follows
				0x82,
				// element 0: type definitions
				// array, 1 element follows
				0x81,
				// event type:
				// id: []byte{}
				// cadence-type-id: "A.1654653399040a61.FlowToken.TokensDeposited"
				// 2 field: [["to", type(optional(address))], ["amount", type(ufix64)]]
				// tag
				0xd8, ccf.CBORTagEventType,
				// array, 3 element follows
				0x83,
				// id
				// bytes, 0 bytes follow
				0x40,
				// cadence-type-id
				// string, 44 bytes follow
				0x78, 0x2c,
				// "A.1654653399040a61.FlowToken.TokensDeposited"
				0x41, 0x2e, 0x31, 0x36, 0x35, 0x34, 0x36, 0x35, 0x33, 0x33, 0x39, 0x39, 0x30, 0x34, 0x30, 0x61, 0x36, 0x31, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x44, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x65, 0x64,
				// fields
				// array, 2 items follow
				0x82,
				// field 0
				// array, 2 element follows
				0x82,
				// text, 2 bytes follow
				0x62,
				// "to"
				0x74, 0x6f,
				// tag
				0xd8, ccf.CBORTagOptionalType,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Address type ID (3)
				0x03,
				// field 1
				// array, 2 element follows
				0x82,
				// text, 6 bytes follow
				0x66,
				// "amount"
				0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,

				// element 1: type and value
				// array, 2 element follows
				0x82,
				// tag
				0xd8, ccf.CBORTagTypeRef,
				// bytes, 0 bytes follow
				0x40,
				// array, 2 items follow
				0x82,
				// bytes, 8 bytes follow
				0x48,
				// 0x8624b52f9ddcd04a
				0x86, 0x24, 0xb5, 0x2f, 0x9d, 0xdc, 0xd0, 0x4a,
				// 174549955740
				0x1b, 0x00, 0x00, 0x00, 0x28, 0xa3, 0xfc, 0xf4, 0x9c,
			},
		},
		{
			name:  "FlowToken.TokensMinted",
			event: createFlowTokenTokensMintedEvent(),
			expectedCBOR: []byte{
				// language=json, format=json-cadence data interchange format
				// {"value":{"id":"A.1654653399040a61.FlowToken.TokensMinted","fields":[{"value":{"value":"1316489.95887105","type":"UFix64"},"name":"amount"}]},"type":"Event"}
				//
				// language=edn, format=ccf
				// 129([[162([h'', "A.1654653399040a61.FlowToken.TokensMinted", [["amount", 137(23)]]])], [136(h''), [131648995887105]]])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeDefAndValue,
				// array, 2 element follows
				0x82,
				// element 0: type definitions
				// array, 1 element follows
				0x81,
				// event type:
				// id: []byte{}
				// cadence-type-id: "A.1654653399040a61.FlowToken.TokensMinted"
				// 1 field: [["amount", type(ufix64)]]
				// tag
				0xd8, ccf.CBORTagEventType,
				// array, 3 element follows
				0x83,
				// id
				// bytes, 0 bytes follow
				0x40,
				// cadence-type-id
				// string, 41 bytes follow
				0x78, 0x29,
				// "A.1654653399040a61.FlowToken.TokensMinted"
				0x41, 0x2e, 0x31, 0x36, 0x35, 0x34, 0x36, 0x35, 0x33, 0x33, 0x39, 0x39, 0x30, 0x34, 0x30, 0x61, 0x36, 0x31, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x4d, 0x69, 0x6e, 0x74, 0x65, 0x64,
				// fields
				// array, 1 items follow
				0x81,
				// field 0
				// array, 2 element follows
				0x82,
				// text, 6 bytes follow
				0x66,
				// "amount"
				0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,

				// element 1: type and value
				// array, 2 element follows
				0x82,
				// tag
				0xd8, ccf.CBORTagTypeRef,
				// bytes, 0 bytes follow
				0x40,
				// array, 1 items follow
				0x81,
				// 131648995887105
				0x1b, 0x00, 0x00, 0x77, 0xbb, 0xeb, 0xa2, 0x88, 0x01,
			},
		},
		{
			name:  "FlowToken.TokensWithdrawn",
			event: createFlowTokenTokensWithdrawnEvent(),
			expectedCBOR: []byte{
				// language=json, format=json-cadence data interchange format
				// {"value":{"id":"A.1654653399040a61.FlowToken.TokensWithdrawn","fields":[{"value":{"value":"53.04112895","type":"UFix64"},"name":"amount"},{"value":{"value":{"value":"0xf919ee77447b7497","type":"Address"},"type":"Optional"},"name":"from"}]},"type":"Event"}
				//
				// language=edn, format=ccf
				// 129([[162([h'', "A.1654653399040a61.FlowToken.TokensWithdrawn", [["from", 138(137(3))], ["amount", 137(23)]]])], [136(h''), [h'F919EE77447B7497', 5304112895]]])
				//
				// language=cbor, format=ccf
				// tag
				0xd8, ccf.CBORTagTypeDefAndValue,
				// array, 2 element follows
				0x82,
				// element 0: type definitions
				// array, 1 element follows
				0x81,
				// event type:
				// id: []byte{}
				// cadence-type-id: "A.1654653399040a61.FlowToken.TokensWithdrawn"
				// 2 field: [["from", type(optional(address))], ["amount", type(ufix64)]]
				// tag
				0xd8, ccf.CBORTagEventType,
				// array, 3 element follows
				0x83,
				// id
				// bytes, 0 bytes follow
				0x40,
				// cadence-type-id
				// string, 44 bytes follow
				0x78, 0x2c,
				// "A.1654653399040a61.FlowToken.TokensWithdrawn"
				0x41, 0x2e, 0x31, 0x36, 0x35, 0x34, 0x36, 0x35, 0x33, 0x33, 0x39, 0x39, 0x30, 0x34, 0x30, 0x61, 0x36, 0x31, 0x2e, 0x46, 0x6c, 0x6f, 0x77, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x73, 0x57, 0x69, 0x74, 0x68, 0x64, 0x72, 0x61, 0x77, 0x6e,
				// fields
				// array, 2 items follow
				0x82,
				// field 0
				// array, 2 element follows
				0x82,
				// text, 4 bytes follow
				0x64,
				// "from"
				0x66, 0x72, 0x6f, 0x6d,
				// tag
				0xd8, ccf.CBORTagOptionalType,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// Address type ID (3)
				0x03,
				// field 1
				// array, 2 element follows
				0x82,
				// text, 6 bytes follow
				0x66,
				// "amount"
				0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74,
				// tag
				0xd8, ccf.CBORTagSimpleType,
				// UFix64 type ID (23)
				0x17,

				// element 1: type and value
				// array, 2 element follows
				0x82,
				// tag
				0xd8, ccf.CBORTagTypeRef,
				// bytes, 0 bytes follow
				0x40,
				// array, 2 items follow
				0x82,
				// bytes, 8 bytes follow
				0x48,
				// 0xf919ee77447b7497
				0xf9, 0x19, 0xee, 0x77, 0x44, 0x7b, 0x74, 0x97,
				// 5304112895
				0x1b, 0x00, 0x00, 0x00, 0x01, 0x3c, 0x26, 0x56, 0xff,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encode Cadence value to CCF
			actualCBOR, err := ccf.Encode(tt.event)
			require.NoError(t, err)
			require.Equal(t, tt.expectedCBOR, actualCBOR)

			// Decode CCF to Cadence value
			decodedEvent, err := ccf.Decode(nil, actualCBOR)
			require.NoError(t, err)

			// Test original event and decoded events are equal even if
			// fields are ordered differently due to deterministic encoding.
			testEventEquality(t, tt.event, decodedEvent.(cadence.Event))
		})
	}
}

func newFlowFeesFeesDeductedEventType() *cadence.EventType {
	// pub event FeesDeducted(amount: UFix64, inclusionEffort: UFix64, executionEffort: UFix64)

	address, _ := common.HexToAddress("f919ee77447b7497")
	location := common.NewAddressLocation(nil, address, "FlowFees")

	return &cadence.EventType{
		Location:            location,
		QualifiedIdentifier: "FlowFees.FeesDeducted",
		Fields: []cadence.Field{
			{
				Identifier: "amount",
				Type:       cadence.UFix64Type{},
			},
			{
				Identifier: "inclusionEffort",
				Type:       cadence.UFix64Type{},
			},
			{
				Identifier: "executionEffort",
				Type:       cadence.UFix64Type{},
			},
		},
	}
}

func createFlowFeesFeesDeductedEvent() cadence.Event {
	/*
		A.f919ee77447b7497.FlowFees.FeesDeducted
		{
			"amount": "0.01797293",
			"inclusionEffort": "1.00000000",
			"executionEffort": "0.00360123"
		}
	*/
	amount, _ := cadence.NewUFix64("0.01797293")
	inclusionEffort, _ := cadence.NewUFix64("1.00000000")
	executionEffort, _ := cadence.NewUFix64("0.00360123")

	return cadence.NewEvent(
		[]cadence.Value{amount, inclusionEffort, executionEffort},
	).WithType(newFlowFeesFeesDeductedEventType())
}

func newFlowFeesTokensWithdrawnEventType() *cadence.EventType {
	// pub event TokensWithdrawn(amount: UFix64)

	address, _ := common.HexToAddress("f919ee77447b7497")
	location := common.NewAddressLocation(nil, address, "FlowFees")

	return &cadence.EventType{
		Location:            location,
		QualifiedIdentifier: "FlowFees.TokensWithdrawn",
		Fields: []cadence.Field{
			{
				Identifier: "amount",
				Type:       cadence.UFix64Type{},
			},
		},
	}
}

func createFlowFeesTokensWithdrawnEvent() cadence.Event {
	/*
		A.f919ee77447b7497.FlowFees.TokensWithdrawn
		{
			"amount": "53.04112895"
		}
	*/
	amount, _ := cadence.NewUFix64("53.04112895")

	return cadence.NewEvent(
		[]cadence.Value{amount},
	).WithType(newFlowFeesTokensWithdrawnEventType())
}

func newFlowTokenTokensDepositedEventType() *cadence.EventType {
	// pub event TokensDeposited(amount: UFix64, to: Address?)

	address, _ := common.HexToAddress("1654653399040a61")
	location := common.NewAddressLocation(nil, address, "FlowToken")

	return &cadence.EventType{
		Location:            location,
		QualifiedIdentifier: "FlowToken.TokensDeposited",
		Fields: []cadence.Field{
			{
				Identifier: "amount",
				Type:       cadence.UFix64Type{},
			},
			{
				Identifier: "to",
				Type: &cadence.OptionalType{
					Type: cadence.NewAddressType(),
				},
			},
		},
	}
}

func createFlowTokenTokensDepositedEventNoReceiver() cadence.Event {
	/*
		A.1654653399040a61.FlowToken.TokensDeposited
		{
			"amount": "1316489.95887105",
			"to": null
		}
	*/
	amount, _ := cadence.NewUFix64("1316489.95887105")
	to := cadence.NewOptional(nil)

	return cadence.NewEvent(
		[]cadence.Value{amount, to},
	).WithType(newFlowTokenTokensDepositedEventType())
}

func createFlowTokenTokensDepositedEvent() cadence.Event {
	/*
		A.1654653399040a61.FlowToken.TokensDeposited
		{
			"amount": "1745.49955740",
			"to": "0x8624b52f9ddcd04a"
		}
	*/
	addressBytes, _ := hex.DecodeString("8624b52f9ddcd04a")

	amount, _ := cadence.NewUFix64("1745.49955740")
	to := cadence.NewOptional(cadence.BytesToAddress(addressBytes))

	return cadence.NewEvent(
		[]cadence.Value{amount, to},
	).WithType(newFlowTokenTokensDepositedEventType())
}

func newFlowTokenTokensMintedEventType() *cadence.EventType {
	// pub event TokensMinted(amount: UFix64)

	address, _ := common.HexToAddress("1654653399040a61")
	location := common.NewAddressLocation(nil, address, "FlowToken")

	return &cadence.EventType{
		Location:            location,
		QualifiedIdentifier: "FlowToken.TokensMinted",
		Fields: []cadence.Field{
			{
				Identifier: "amount",
				Type:       cadence.UFix64Type{},
			},
		},
	}
}

func createFlowTokenTokensMintedEvent() cadence.Event {
	/*
		A.1654653399040a61.FlowToken.TokensMinted
		{
			"amount": "1316489.95887105"
		}
	*/
	amount, _ := cadence.NewUFix64("1316489.95887105")

	return cadence.NewEvent(
		[]cadence.Value{amount},
	).WithType(newFlowTokenTokensMintedEventType())
}

func newFlowTokenTokensWithdrawnEventType() *cadence.EventType {
	// pub event TokensWithdrawn(amount: UFix64, from: Address?)

	address, _ := common.HexToAddress("1654653399040a61")
	location := common.NewAddressLocation(nil, address, "FlowToken")

	return &cadence.EventType{
		Location:            location,
		QualifiedIdentifier: "FlowToken.TokensWithdrawn",
		Fields: []cadence.Field{
			{
				Identifier: "amount",
				Type:       cadence.UFix64Type{},
			},
			{
				Identifier: "from",
				Type: &cadence.OptionalType{
					Type: cadence.NewAddressType(),
				},
			},
		},
	}
}

func createFlowTokenTokensWithdrawnEvent() cadence.Event {
	/*
		A.1654653399040a61.FlowToken.TokensWithdrawn
		{
			"amount": "53.04112895",
			"from": "0xf919ee77447b7497"
		}
	*/
	addressBytes, _ := hex.DecodeString("f919ee77447b7497")

	amount, _ := cadence.NewUFix64("53.04112895")
	to := cadence.NewOptional(cadence.BytesToAddress(addressBytes))

	return cadence.NewEvent(
		[]cadence.Value{amount, to},
	).WithType(newFlowTokenTokensWithdrawnEventType())
}

func newFlowIDTableStakingDelegatorRewardsPaidEventType() *cadence.EventType {
	// pub event DelegatorRewardsPaid(nodeID: String, delegatorID: UInt32, amount: UFix64)

	address, _ := common.HexToAddress("8624b52f9ddcd04a")
	location := common.NewAddressLocation(nil, address, "FlowIDTableStaking")

	return &cadence.EventType{
		Location:            location,
		QualifiedIdentifier: "FlowIDTableStaking.DelegatorRewardsPaid",
		Fields: []cadence.Field{
			{
				Identifier: "nodeID",
				Type:       cadence.StringType{},
			},
			{
				Identifier: "delegatorID",
				Type:       cadence.UInt32Type{},
			},
			{
				Identifier: "amount",
				Type:       cadence.UFix64Type{},
			},
		},
	}
}

func createFlowIDTableStakingDelegatorRewardsPaidEvent() cadence.Event {
	/*
		A.8624b52f9ddcd04a.FlowIDTableStaking.DelegatorRewardsPaid
		{
			"nodeID": "e52cbcd825e328acac8db6bcbdcbb6e7724862c8b89b09d85edccf41ff9981eb",
			"delegatorID": 92,
			"amount": "4.38760261"
		}
	*/
	nodeID := cadence.String("e52cbcd825e328acac8db6bcbdcbb6e7724862c8b89b09d85edccf41ff9981eb")
	delegatorID := cadence.UInt32(92)
	amount, _ := cadence.NewUFix64("4.38760261")

	return cadence.NewEvent(
		[]cadence.Value{nodeID, delegatorID, amount},
	).WithType(newFlowIDTableStakingDelegatorRewardsPaidEventType())
}

func newFlowIDTableStakingEpochTotalRewardsPaidEventType() *cadence.EventType {
	// pub event EpochTotalRewardsPaid(total: UFix64, fromFees: UFix64, minted: UFix64, feesBurned: UFix64)

	address, _ := common.HexToAddress("8624b52f9ddcd04a")
	location := common.NewAddressLocation(nil, address, "FlowIDTableStaking")

	return &cadence.EventType{
		Location:            location,
		QualifiedIdentifier: "FlowIDTableStaking.EpochTotalRewardsPaid",
		Fields: []cadence.Field{
			{
				Identifier: "total",
				Type:       cadence.UFix64Type{},
			},
			{
				Identifier: "fromFees",
				Type:       cadence.UFix64Type{},
			},
			{
				Identifier: "minted",
				Type:       cadence.UFix64Type{},
			},
			{
				Identifier: "feesBurned",
				Type:       cadence.UFix64Type{},
			},
		},
	}
}

func createFlowIDTableStakingEpochTotalRewardsPaidEvent() cadence.Event {
	/*
		A.8624b52f9ddcd04a.FlowIDTableStaking.EpochTotalRewardsPaid
		{
			"total": "1316543.00000000",
			"fromFees": "53.04112895",
			"minted": "1316489.95887105",
			"feesBurned": "6.04080767"
		}
	*/
	total, _ := cadence.NewUFix64("1316543.00000000")
	fromFees, _ := cadence.NewUFix64("53.04112895")
	minted, _ := cadence.NewUFix64("1316489.95887105")
	feesBurned, _ := cadence.NewUFix64("6.04080767")

	return cadence.NewEvent(
		[]cadence.Value{total, fromFees, minted, feesBurned},
	).WithType(newFlowIDTableStakingEpochTotalRewardsPaidEventType())
}

func newFlowIDTableStakingNewWeeklyPayoutEventType() *cadence.EventType {
	// pub event NewWeeklyPayout(newPayout: UFix64)

	address, _ := common.HexToAddress("8624b52f9ddcd04a")
	location := common.NewAddressLocation(nil, address, "FlowIDTableStaking")

	return &cadence.EventType{
		Location:            location,
		QualifiedIdentifier: "FlowIDTableStaking.NewWeeklyPayout",
		Fields: []cadence.Field{
			{
				Identifier: "newPayout",
				Type:       cadence.UFix64Type{},
			},
		},
	}
}

func createFlowIDTableStakingNewWeeklyPayoutEvent() cadence.Event {
	/*
		A.8624b52f9ddcd04a.FlowIDTableStaking.NewWeeklyPayout
		{
			"newPayout": "1317778.00000000"
		}
	*/
	newPayout, _ := cadence.NewUFix64("1317778.00000000")

	return cadence.NewEvent(
		[]cadence.Value{newPayout},
	).WithType(newFlowIDTableStakingNewWeeklyPayoutEventType())
}

func newFlowIDTableStakingRewardsPaidEventType() *cadence.EventType {
	// pub event RewardsPaid(nodeID: String, amount: UFix64)

	address, _ := common.HexToAddress("8624b52f9ddcd04a")
	location := common.NewAddressLocation(nil, address, "FlowIDTableStaking")

	return &cadence.EventType{
		Location:            location,
		QualifiedIdentifier: "FlowIDTableStaking.RewardsPaid",
		Fields: []cadence.Field{
			{
				Identifier: "nodeID",
				Type:       cadence.StringType{},
			},
			{
				Identifier: "amount",
				Type:       cadence.UFix64Type{},
			},
		},
	}
}

func createFlowIDTableStakingRewardsPaidEvent() cadence.Event {
	nodeID, _ := cadence.NewString("e52cbcd825e328acac8db6bcbdcbb6e7724862c8b89b09d85edccf41ff9981eb")
	amount, _ := cadence.NewUFix64("1745.49955740")

	return cadence.NewEvent(
		[]cadence.Value{nodeID, amount},
	).WithType(newFlowIDTableStakingRewardsPaidEventType())
}

func testEventEquality(t *testing.T, event1 cadence.Event, event2 cadence.Event) {
	require.True(t, event1.Type().Equal(event2.Type()))
	require.Equal(t, len(event1.Fields), len(event2.Fields))
	require.Equal(t, len(event1.EventType.Fields), len(event2.EventType.Fields))

	for i, event1FieldType := range event1.EventType.Fields {

		foundField := false

		for j, event2FieldType := range event2.EventType.Fields {
			if event1FieldType.Identifier == event2FieldType.Identifier {
				require.Equal(t, event1.Fields[i], event2.Fields[j])
				foundField = true
				break
			}
		}

		require.True(t, foundField)
	}
}
