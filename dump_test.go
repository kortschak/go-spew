/*
 * Copyright (c) 2013 Dave Collins <dave@davec.name>
 * Copyright (c) 2015 Dan Kortschak <dan.kortschak@adelaide.edu.au>
 *
 * Permission to use, copy, modify, and distribute this software for any
 * purpose with or without fee is hereby granted, provided that the above
 * copyright notice and this permission notice appear in all copies.
 *
 * THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
 * WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
 * MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
 * ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
 * WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
 * ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
 * OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.
 */

/*
Test Summary:
NOTE: For each test, a nil pointer, a single pointer and double pointer to the
base test element are also tested to ensure proper indirection across all types.

- Max int8, int16, int32, int64, int
- Max uint8, uint16, uint32, uint64, uint
- Boolean true and false
- Standard complex64 and complex128
- Array containing standard ints
- Array containing type with custom formatter on pointer receiver only
- Array containing interfaces
- Array containing bytes
- Slice containing standard float32 values
- Slice containing type with custom formatter on pointer receiver only
- Slice containing interfaces
- Slice containing bytes
- Nil slice
- Standard string
- Nil interface
- Sub-interface
- Map with string keys and int vals
- Map with custom formatter type on pointer receiver only keys and vals
- Map with interface keys and values
- Map with nil interface value
- Struct with primitives
- Struct that contains another struct
- Struct that contains custom type with Stringer pointer interface via both
  exported and unexported fields
- Struct that contains embedded struct and field to same struct
- Uintptr to 0 (null pointer)
- Uintptr address of real variable
- Unsafe.Pointer to 0 (null pointer)
- Unsafe.Pointer to address of real variable
- Nil channel
- Standard int channel
- Function with no params and no returns
- Function with param and no returns
- Function with multiple params and multiple returns
- Struct that is circular through self referencing
- Structs that are circular through cross referencing
- Structs that are indirectly circular
- Type that panics in its Stringer interface
*/

package utter_test

import (
	"bytes"
	"fmt"
	"testing"
	"unsafe"

	"github.com/kortschak/utter"
)

// dumpTest is used to describe a test to be perfomed against the Dump method.
type dumpTest struct {
	in    interface{}
	wants []string
}

// dumpTests houses all of the tests to be performed against the Dump method.
var dumpTests = make([]dumpTest, 0)

// addDumpTest is a helper method to append the passed input and desired result
// to dumpTests
func addDumpTest(in interface{}, wants ...string) {
	test := dumpTest{in, wants}
	dumpTests = append(dumpTests, test)
}

func addIntDumpTests() {
	// Max int8.
	v := int8(127)
	nv := (*int8)(nil)
	pv := &v
	vt := "int8"
	vs := "127"
	addDumpTest(v, vt+"("+vs+")\n")
	addDumpTest(pv, "&"+vt+"("+vs+")\n")
	addDumpTest(&pv, "&&"+vt+"("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(nil)\n")

	// Max int16.
	v2 := int16(32767)
	nv2 := (*int16)(nil)
	pv2 := &v2
	v2t := "int16"
	v2s := "32767"
	addDumpTest(v2, v2t+"("+v2s+")\n")
	addDumpTest(pv2, "&"+v2t+"("+v2s+")\n")
	addDumpTest(&pv2, "&&"+v2t+"("+v2s+")\n")
	addDumpTest(nv2, "(*"+v2t+")(nil)\n")

	// Max int32.
	v3 := int32(2147483647)
	nv3 := (*int32)(nil)
	pv3 := &v3
	v3t := "int32"
	v3s := "2147483647"
	addDumpTest(v3, v3t+"("+v3s+")\n")
	addDumpTest(pv3, "&"+v3t+"("+v3s+")\n")
	addDumpTest(&pv3, "&&"+v3t+"("+v3s+")\n")
	addDumpTest(nv3, "(*"+v3t+")(nil)\n")

	// Max int64.
	v4 := int64(9223372036854775807)
	nv4 := (*int64)(nil)
	pv4 := &v4
	v4t := "int64"
	v4s := "9223372036854775807"
	addDumpTest(v4, v4t+"("+v4s+")\n")
	addDumpTest(pv4, "&"+v4t+"("+v4s+")\n")
	addDumpTest(&pv4, "&&"+v4t+"("+v4s+")\n")
	addDumpTest(nv4, "(*"+v4t+")(nil)\n")

	// Max int.
	v5 := int(2147483647)
	nv5 := (*int)(nil)
	pv5 := &v5
	v5t := "int"
	v5s := "2147483647"
	addDumpTest(v5, v5t+"("+v5s+")\n")
	addDumpTest(pv5, "&"+v5t+"("+v5s+")\n")
	addDumpTest(&pv5, "&&"+v5t+"("+v5s+")\n")
	addDumpTest(nv5, "(*"+v5t+")(nil)\n")
}

func addUintDumpTests() {
	// Max uint8.
	v := uint8(255)
	nv := (*uint8)(nil)
	pv := &v
	vt := "uint8"
	vs := "0xff"
	addDumpTest(v, vt+"("+vs+")\n")
	addDumpTest(pv, "&"+vt+"("+vs+")\n")
	addDumpTest(&pv, "&&"+vt+"("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(nil)\n")

	// Max uint16.
	v2 := uint16(65535)
	nv2 := (*uint16)(nil)
	pv2 := &v2
	v2t := "uint16"
	v2s := "0xffff"
	addDumpTest(v2, v2t+"("+v2s+")\n")
	addDumpTest(pv2, "&"+v2t+"("+v2s+")\n")
	addDumpTest(&pv2, "&&"+v2t+"("+v2s+")\n")
	addDumpTest(nv2, "(*"+v2t+")(nil)\n")

	// Max uint32.
	v3 := uint32(4294967295)
	nv3 := (*uint32)(nil)
	pv3 := &v3
	v3t := "uint32"
	v3s := "0xffffffff"
	addDumpTest(v3, v3t+"("+v3s+")\n")
	addDumpTest(pv3, "&"+v3t+"("+v3s+")\n")
	addDumpTest(&pv3, "&&"+v3t+"("+v3s+")\n")
	addDumpTest(nv3, "(*"+v3t+")(nil)\n")

	// Max uint64.
	v4 := uint64(18446744073709551615)
	nv4 := (*uint64)(nil)
	pv4 := &v4
	v4t := "uint64"
	v4s := "0xffffffffffffffff"
	addDumpTest(v4, v4t+"("+v4s+")\n")
	addDumpTest(pv4, "&"+v4t+"("+v4s+")\n")
	addDumpTest(&pv4, "&&"+v4t+"("+v4s+")\n")
	addDumpTest(nv4, "(*"+v4t+")(nil)\n")

	// Max uint.
	v5 := uint(4294967295)
	nv5 := (*uint)(nil)
	pv5 := &v5
	v5t := "uint"
	v5s := "0xffffffff"
	addDumpTest(v5, v5t+"("+v5s+")\n")
	addDumpTest(pv5, "&"+v5t+"("+v5s+")\n")
	addDumpTest(&pv5, "&&"+v5t+"("+v5s+")\n")
	addDumpTest(nv5, "(*"+v5t+")(nil)\n")
}

func addBoolDumpTests() {
	// Boolean true.
	v := bool(true)
	nv := (*bool)(nil)
	pv := &v
	vt := "bool"
	vs := "true"
	addDumpTest(v, vt+"("+vs+")\n")
	addDumpTest(pv, "&"+vt+"("+vs+")\n")
	addDumpTest(&pv, "&&"+vt+"("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(nil)\n")

	// Boolean false.
	v2 := bool(false)
	pv2 := &v2
	v2t := "bool"
	v2s := "false"
	addDumpTest(v2, v2t+"("+v2s+")\n")
	addDumpTest(pv2, "&"+v2t+"("+v2s+")\n")
	addDumpTest(&pv2, "&&"+v2t+"("+v2s+")\n")
}

func addFloatDumpTests() {
	// Standard float32.
	v := float32(3.1415)
	nv := (*float32)(nil)
	pv := &v
	vt := "float32"
	vs := "3.1415"
	addDumpTest(v, vt+"("+vs+")\n")
	addDumpTest(pv, "&"+vt+"("+vs+")\n")
	addDumpTest(&pv, "&&"+vt+"("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(nil)\n")

	// Standard float64.
	v2 := float64(3.1415926)
	nv2 := (*float64)(nil)
	pv2 := &v2
	v2t := "float64"
	v2s := "3.1415926"
	addDumpTest(v2, v2t+"("+v2s+")\n")
	addDumpTest(pv2, "&"+v2t+"("+v2s+")\n")
	addDumpTest(&pv2, "&&"+v2t+"("+v2s+")\n")
	addDumpTest(nv2, "(*"+v2t+")(nil)\n")

	// Standard float32 - integral value.
	v3 := float32(3)
	nv3 := (*float32)(nil)
	pv3 := &v3
	v3t := "float32"
	v3s := "3"
	addDumpTest(v3, v3t+"("+v3s+")\n")
	addDumpTest(pv3, "&"+v3t+"("+v3s+")\n")
	addDumpTest(&pv3, "&&"+v3t+"("+v3s+")\n")
	addDumpTest(nv3, "(*"+v3t+")(nil)\n")

	// Standard float64 - integral value.
	v4 := float64(3)
	nv4 := (*float64)(nil)
	pv4 := &v4
	v4t := "float64"
	v4s := "3"
	addDumpTest(v4, v4t+"("+v4s+")\n")
	addDumpTest(pv4, "&"+v4t+"("+v4s+")\n")
	addDumpTest(&pv4, "&&"+v4t+"("+v4s+")\n")
	addDumpTest(nv4, "(*"+v4t+")(nil)\n")
}

func addComplexDumpTests() {
	// Standard complex64.
	v := complex(float32(6), -2)
	nv := (*complex64)(nil)
	pv := &v
	vt := "complex64"
	vs := "6-2i"
	addDumpTest(v, vt+"("+vs+")\n")
	addDumpTest(pv, "&"+vt+"("+vs+")\n")
	addDumpTest(&pv, "&&"+vt+"("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(nil)\n")

	// Standard complex128.
	v2 := complex(float64(-6), 2)
	nv2 := (*complex128)(nil)
	pv2 := &v2
	v2t := "complex128"
	v2s := "-6+2i"
	addDumpTest(v2, v2t+"("+v2s+")\n")
	addDumpTest(pv2, "&"+v2t+"("+v2s+")\n")
	addDumpTest(&pv2, "&&"+v2t+"("+v2s+")\n")
	addDumpTest(nv2, "(*"+v2t+")(nil)\n")
}

func addArrayDumpTests() {
	// Array containing standard ints.
	v := [3]int{1, 2, 3}
	nv := (*[3]int)(nil)
	pv := &v
	vt := "int"
	vs := "{\n " + vt + "(1),\n " + vt + "(2),\n " + vt + "(3),\n}"
	addDumpTest(v, "[3]"+vt+vs+"\n")
	addDumpTest(pv, "&[3]"+vt+vs+"\n")
	addDumpTest(&pv, "&&[3]"+vt+vs+"\n")
	addDumpTest(nv, "(*[3]"+vt+")(nil)\n")

	// Array containing type with custom formatter on pointer receiver only.
	v2i0 := pstringer("1")
	v2i1 := pstringer("2")
	v2i2 := pstringer("3")
	v2 := [3]pstringer{v2i0, v2i1, v2i2}
	nv2 := (*[3]pstringer)(nil)
	pv2 := &v2
	v2t := "utter_test.pstringer"
	v2s := "{\n " + v2t + "(\"1\"),\n " + v2t + "(\"2\"),\n " + v2t + "(\"3\"),\n}"
	addDumpTest(v2, "[3]"+v2t+v2s+"\n")
	addDumpTest(pv2, "&[3]"+v2t+v2s+"\n")
	addDumpTest(&pv2, "&&[3]"+v2t+v2s+"\n")
	addDumpTest(nv2, "(*[3]"+v2t+")(nil)\n")

	// Array containing interfaces.
	v3i0 := "one"
	v3 := [3]interface{}{v3i0, int(2), uint(3)}
	nv3 := (*[3]interface{})(nil)
	pv3 := &v3
	v3t := "[3]interface{}"
	v3t2 := "string"
	v3t3 := "int"
	v3t4 := "uint"
	v3s := "{\n " + v3t2 + "(\"one\"),\n " + v3t3 + "(2),\n " + v3t4 + "(0x3),\n}"
	addDumpTest(v3, v3t+v3s+"\n")
	addDumpTest(pv3, "&"+v3t+v3s+"\n")
	addDumpTest(&pv3, "&&"+v3t+v3s+"\n")
	addDumpTest(nv3, "(*"+v3t+")(nil)\n")

	// Array containing bytes.
	v4 := [34]byte{
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30,
		0x31, 0x32,
	}
	nv4 := (*[34]byte)(nil)
	pv4 := &v4
	v4t := "[34]uint8"
	v4s := "{\n" +
		" 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, // |............... |\n" +
		" 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, // |!\"#$%&'()*+,-./0|\n" +
		" 0x31, 0x32, /*                                                                               */ // |12|\n}"
	addDumpTest(v4, v4t+v4s+"\n")
	addDumpTest(pv4, "&"+v4t+v4s+"\n")
	addDumpTest(&pv4, "&&"+v4t+v4s+"\n")
	addDumpTest(nv4, "(*"+v4t+")(nil)\n")
}

func addSliceDumpTests() {
	// Slice containing standard float32 values.
	v := []float32{3.14, 6.28, 12.56}
	nv := (*[]float32)(nil)
	pv := &v
	vt := "float32"
	vs := "{\n " + vt + "(3.14),\n " + vt + "(6.28),\n " + vt + "(12.56),\n}"
	addDumpTest(v, "[]"+vt+vs+"\n")
	addDumpTest(pv, "&[]"+vt+vs+"\n")
	addDumpTest(&pv, "&&[]"+vt+vs+"\n")
	addDumpTest(nv, "(*[]"+vt+")(nil)\n")

	// Slice containing type with custom formatter on pointer receiver only.
	v2i0 := pstringer("1")
	v2i1 := pstringer("2")
	v2i2 := pstringer("3")
	v2 := []pstringer{v2i0, v2i1, v2i2}
	nv2 := (*[]pstringer)(nil)
	pv2 := &v2
	v2t := "utter_test.pstringer"
	v2s := "{\n " + v2t + "(\"1\"),\n " + v2t + "(\"2\"),\n " + v2t + "(\"3\"),\n}"
	addDumpTest(v2, "[]"+v2t+v2s+"\n")
	addDumpTest(pv2, "&[]"+v2t+v2s+"\n")
	addDumpTest(&pv2, "&&[]"+v2t+v2s+"\n")
	addDumpTest(nv2, "(*[]"+v2t+")(nil)\n")

	// Slice containing interfaces.
	v3i0 := "one"
	v3 := []interface{}{v3i0, int(2), uint(3), nil}
	nv3 := (*[]interface{})(nil)
	pv3 := &v3
	v3t := "[]interface{}"
	v3t2 := "string"
	v3t3 := "int"
	v3t4 := "uint"
	v3t5 := "interface{}"
	v3s := "{\n " + v3t2 + "(\"one\"),\n " + v3t3 + "(2),\n " + v3t4 + "(0x3),\n " + v3t5 + "(nil),\n}"
	addDumpTest(v3, v3t+v3s+"\n")
	addDumpTest(pv3, "&"+v3t+v3s+"\n")
	addDumpTest(&pv3, "&&"+v3t+v3s+"\n")
	addDumpTest(nv3, "(*"+v3t+")(nil)\n")

	// Slice containing bytes.
	v4 := []byte{
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30,
		0x31, 0x32,
	}
	nv4 := (*[]byte)(nil)
	pv4 := &v4
	v4t := "[]uint8"
	v4s := "{\n" +
		" 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, // |............... |\n" +
		" 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, // |!\"#$%&'()*+,-./0|\n" +
		" 0x31, 0x32, /*                                                                               */ // |12|\n}"
	addDumpTest(v4, v4t+v4s+"\n")
	addDumpTest(pv4, "&"+v4t+v4s+"\n")
	addDumpTest(&pv4, "&&"+v4t+v4s+"\n")
	addDumpTest(nv4, "(*"+v4t+")(nil)\n")

	// Nil slice.
	v5 := []int(nil)
	nv5 := (*[]int)(nil)
	pv5 := &v5
	v5t := "[]int"
	v5s := "nil"
	addDumpTest(v5, v5t+"("+v5s+")\n")
	addDumpTest(pv5, "&"+v5t+"("+v5s+")\n")
	addDumpTest(&pv5, "&&"+v5t+"("+v5s+")\n")
	addDumpTest(nv5, "(*"+v5t+")(nil)\n")
}

func addStringDumpTests() {
	// Standard string.
	v := "test"
	nv := (*string)(nil)
	pv := &v
	vt := "string"
	vs := "\"test\""
	addDumpTest(v, vt+"("+vs+")\n")
	addDumpTest(pv, "&"+vt+"("+vs+")\n")
	addDumpTest(&pv, "&&"+vt+"("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(nil)\n")
}

func addInterfaceDumpTests() {
	// Nil interface.
	var v interface{}
	nv := (*interface{})(nil)
	pv := &v
	vt := "interface{}"
	vs := "(nil)"
	addDumpTest(v, "interface{}"+vs+"\n")
	addDumpTest(pv, "&"+vt+vs+"\n")
	addDumpTest(&pv, "&&"+vt+vs+"\n")
	addDumpTest(nv, "(*"+vt+")"+vs+"\n")

	// Sub-interface.
	v2 := interface{}(uint16(65535))
	pv2 := &v2
	v2t := "uint16"
	v2s := "0xffff"
	addDumpTest(v2, v2t+"("+v2s+")\n")
	addDumpTest(pv2, "&"+v2t+"("+v2s+")\n")
	addDumpTest(&pv2, "&&"+v2t+"("+v2s+")\n")
}

func addMapDumpTests() {
	// Map with string keys and int vals.
	k := "one"
	kk := "two"
	m := map[string]int{k: 1, kk: 2}
	nilMap := map[string]int(nil)
	nm := (*map[string]int)(nil)
	pm := &m
	mt := "map[string]int"
	mt1 := "string"
	mt2 := "int"
	ms := "{\n " + mt1 + "(\"one\"): " + mt2 + "(1),\n " + mt1 + "(\"two\"): " + mt2 + "(2),\n}"
	ms2 := "{\n " + mt1 + "(\"two\"): " + mt2 + "(2),\n " + mt1 + "(\"one\"): " + mt2 + "(1),\n}"
	addDumpTest(m, mt+ms+"\n", mt+ms2+"\n")
	addDumpTest(pm, "&"+mt+ms+"\n", "&"+mt+ms2+"\n")
	addDumpTest(&pm, "&&"+mt+ms+"\n", "&&"+mt+ms2+"\n")
	addDumpTest(nm, "(*"+mt+")(nil)\n")
	addDumpTest(nilMap, mt+"(nil)\n")

	// Map with custom formatter type on pointer receiver only keys and vals.
	k2 := pstringer("one")
	v2 := pstringer("1")
	m2 := map[pstringer]pstringer{k2: v2}
	nilMap2 := map[pstringer]pstringer(nil)
	nm2 := (*map[pstringer]pstringer)(nil)
	pm2 := &m2
	m2t := "map[utter_test.pstringer]utter_test.pstringer"
	m2t1 := "utter_test.pstringer"
	m2t2 := "utter_test.pstringer"
	m2s := "{\n " + m2t1 + "(\"one\"): " + m2t2 + "(\"1\"),\n}"
	addDumpTest(m2, m2t+m2s+"\n")
	addDumpTest(pm2, "&"+m2t+m2s+"\n")
	addDumpTest(&pm2, "&&"+m2t+m2s+"\n")
	addDumpTest(nm2, "(*"+m2t+")(nil)\n")
	addDumpTest(nilMap2, m2t+"(nil)\n")

	// Map with interface keys and values.
	k3 := "one"
	m3 := map[interface{}]interface{}{k3: 1}
	nilMap3 := map[interface{}]interface{}(nil)
	nm3 := (*map[interface{}]interface{})(nil)
	pm3 := &m3
	m3t := "map[interface{}]interface{}"
	m3t1 := "string"
	m3t2 := "int"
	m3s := "{\n " + m3t1 + "(\"one\"): " + m3t2 + "(1),\n}"
	addDumpTest(m3, m3t+m3s+"\n")
	addDumpTest(pm3, "&"+m3t+m3s+"\n")
	addDumpTest(&pm3, "&&"+m3t+m3s+"\n")
	addDumpTest(nm3, "(*"+m3t+")(nil)\n")
	addDumpTest(nilMap3, m3t+"(nil)\n")

	// Map with nil interface value.
	k4 := "nil"
	m4 := map[string]interface{}{k4: nil}
	nilMap4 := map[string]interface{}(nil)
	nm4 := (*map[string]interface{})(nil)
	pm4 := &m4
	m4t := "map[string]interface{}"
	m4t1 := "string"
	m4t2 := "interface{}"
	m4s := "{\n " + m4t1 + "(\"nil\"): " + m4t2 + "(nil),\n}"
	addDumpTest(m4, m4t+m4s+"\n")
	addDumpTest(pm4, "&"+m4t+m4s+"\n")
	addDumpTest(&pm4, "&&"+m4t+m4s+"\n")
	addDumpTest(nm4, "(*"+m4t+")(nil)\n")
	addDumpTest(nilMap4, m4t+"(nil)\n")
}

func addStructDumpTests() {
	// Struct with primitives.
	type s1 struct {
		a int8
		b uint8
	}
	v := s1{127, 255}
	nv := (*s1)(nil)
	pv := &v
	vt := "utter_test.s1"
	vt2 := "int8"
	vt3 := "uint8"
	vs := "{\n a: " + vt2 + "(127),\n b: " + vt3 + "(0xff),\n}"
	addDumpTest(v, vt+vs+"\n")
	addDumpTest(pv, "&"+vt+vs+"\n")
	addDumpTest(&pv, "&&"+vt+vs+"\n")
	addDumpTest(nv, "(*"+vt+")(nil)\n")

	// Struct that contains another struct.
	type s2 struct {
		s1 s1
		b  bool
	}
	v2 := s2{s1{127, 255}, true}
	nv2 := (*s2)(nil)
	pv2 := &v2
	v2t := "utter_test.s2"
	v2t2 := "utter_test.s1"
	v2t3 := "int8"
	v2t4 := "uint8"
	v2t5 := "bool"
	v2s := "{\n s1: " + v2t2 + "{\n  a: " + v2t3 + "(127),\n  b: " + v2t4 + "(0xff),\n },\n b: " + v2t5 + "(true),\n}"
	addDumpTest(v2, v2t+v2s+"\n")
	addDumpTest(pv2, "&"+v2t+v2s+"\n")
	addDumpTest(&pv2, "&&"+v2t+v2s+"\n")
	addDumpTest(nv2, "(*"+v2t+")(nil)\n")

	// Struct that contains custom type with Stringer pointer interface via both
	// exported and unexported fields.
	type s3 struct {
		s pstringer
		S pstringer
	}
	v3 := s3{"test", "test2"}
	nv3 := (*s3)(nil)
	pv3 := &v3
	v3t := "utter_test.s3"
	v3t2 := "utter_test.pstringer"
	v3s := "{\n s: " + v3t2 + "(\"test\"),\n S: " + v3t2 + "(\"test2\"),\n}"
	addDumpTest(v3, v3t+v3s+"\n")
	addDumpTest(pv3, "&"+v3t+v3s+"\n")
	addDumpTest(&pv3, "&&"+v3t+v3s+"\n")
	addDumpTest(nv3, "(*"+v3t+")(nil)\n")

	// Struct that contains embedded struct and field to same struct.
	e := embed{"embedstr"}
	v4 := embedwrap{embed: &e, e: &e}
	nv4 := (*embedwrap)(nil)
	pv4 := &v4
	v4t := "utter_test.embedwrap"
	v4t2 := "utter_test.embed"
	v4t3 := "string"
	v4s := "{\n embed: &" + v4t2 + "{\n  a: " + v4t3 + "(\"embedstr\"),\n },\n e: &" + v4t2 + "{\n  a: " + v4t3 + "(\"embedstr\"),\n },\n}"
	addDumpTest(v4, v4t+v4s+"\n")
	addDumpTest(pv4, "&"+v4t+v4s+"\n")
	addDumpTest(&pv4, "&&"+v4t+v4s+"\n")
	addDumpTest(nv4, "(*"+v4t+")(nil)\n")
}

func addUintptrDumpTests() {
	// Null pointer.
	v := uintptr(0)
	pv := &v
	vt := "uintptr"
	vs := "0"
	addDumpTest(v, vt+"("+vs+")\n")
	addDumpTest(pv, "&"+vt+"("+vs+")\n")
	addDumpTest(&pv, "&&"+vt+"("+vs+")\n")

	// Address of real variable.
	i := 1
	v2 := uintptr(unsafe.Pointer(&i))
	nv2 := (*uintptr)(nil)
	pv2 := &v2
	v2t := "uintptr"
	v2s := fmt.Sprintf("(%p)", &i)
	addDumpTest(v2, v2t+v2s+"\n")
	addDumpTest(pv2, "&"+v2t+v2s+"\n")
	addDumpTest(&pv2, "&&"+v2t+v2s+"\n")
	addDumpTest(nv2, "(*"+v2t+")(nil)\n")
}

func addUnsafePointerDumpTests() {
	// Null pointer.
	v := unsafe.Pointer(uintptr(0))
	nv := (*unsafe.Pointer)(nil)
	pv := &v
	vt := "unsafe.Pointer"
	vs := "nil"
	addDumpTest(v, vt+"("+vs+")\n")
	addDumpTest(pv, "&"+vt+"("+vs+")\n")
	addDumpTest(&pv, "&&"+vt+"("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(nil)\n")

	// Address of real variable.
	i := 1
	v2 := unsafe.Pointer(&i)
	pv2 := &v2
	v2t := "unsafe.Pointer"
	v2s := fmt.Sprintf("%p", &i)
	addDumpTest(v2, v2t+"("+v2s+")\n")
	addDumpTest(pv2, "&"+v2t+"("+v2s+")\n")
	addDumpTest(&pv2, "&&"+v2t+"("+v2s+")\n")
	addDumpTest(nv, "(*"+vt+")(nil)\n")
}

func addChanDumpTests() {
	// Nil channel.
	var v chan int
	pv := &v
	nv := (*chan int)(nil)
	vt := "chan int"
	vs := "nil"
	addDumpTest(v, vt+"("+vs+")\n")
	addDumpTest(pv, "&"+vt+"("+vs+")\n")
	addDumpTest(&pv, "&&"+vt+"("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(nil)\n")

	// Real channel.
	v2 := make(chan int)
	pv2 := &v2
	v2t := "chan int"
	v2s := fmt.Sprintf("%p", v2)
	addDumpTest(v2, v2t+"("+v2s+")\n")
	addDumpTest(pv2, "&"+v2t+"("+v2s+")\n")
	addDumpTest(&pv2, "&&"+v2t+"("+v2s+")\n")

	// Real buffered channel, empty.
	v3 := make(chan int, 1)
	pv3 := &v3
	v3t := "(chan int, 1)"
	v3s := fmt.Sprintf("%p", v3)
	addDumpTest(v3, v3t+"("+v3s+")\n")
	addDumpTest(pv3, "&"+v3t+"("+v3s+")\n")
	addDumpTest(&pv3, "&&"+v3t+"("+v3s+")\n")

	// Real buffered channel with one element.
	v4 := func() chan int { c := make(chan int, 2); c <- 1; return c }()
	pv4 := &v4
	v4t := "(chan int, 2 /* 1 element */)"
	v4s := fmt.Sprintf("%p", v4)
	addDumpTest(v4, v4t+"("+v4s+")\n")
	addDumpTest(pv4, "&"+v4t+"("+v4s+")\n")
	addDumpTest(&pv4, "&&"+v4t+"("+v4s+")\n")

	// Real buffered channel with two elements.
	v5 := func() chan int { c := make(chan int, 2); c <- 1; c <- 1; return c }()
	pv5 := &v5
	v5t := "(chan int, 2 /* 2 elements */)"
	v5s := fmt.Sprintf("%p", v5)
	addDumpTest(v5, v5t+"("+v5s+")\n")
	addDumpTest(pv5, "&"+v5t+"("+v5s+")\n")
	addDumpTest(&pv5, "&&"+v5t+"("+v5s+")\n")

	// Real send only channel.
	v6 := make(chan<- int)
	pv6 := &v6
	v6t := "chan<- int"
	v6s := fmt.Sprintf("%p", v6)
	addDumpTest(v6, v6t+"("+v6s+")\n")
	addDumpTest(pv6, "&"+v6t+"("+v6s+")\n")
	addDumpTest(&pv6, "&&"+v6t+"("+v6s+")\n")

	// Real receive only channel.
	v7 := make(<-chan int)
	pv7 := &v7
	v7t := "<-chan int"
	v7s := fmt.Sprintf("%p", v7)
	addDumpTest(v7, v7t+"("+v7s+")\n")
	addDumpTest(pv7, "&"+v7t+"("+v7s+")\n")
	addDumpTest(&pv7, "&&"+v7t+"("+v7s+")\n")
}

func addFuncDumpTests() {
	// Function with no params and no returns.
	v := addIntDumpTests
	nv := (*func())(nil)
	pv := &v
	vt := "func()"
	vs := fmt.Sprintf("%p", v)
	addDumpTest(v, vt+"("+vs+")\n")
	addDumpTest(pv, "&"+vt+"("+vs+")\n")
	addDumpTest(&pv, "&&"+vt+"("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(nil)\n")

	// Function with param and no returns.
	v2 := TestDump
	nv2 := (*func(*testing.T))(nil)
	pv2 := &v2
	v2t := "func(*testing.T)"
	v2s := fmt.Sprintf("%p", v2)
	addDumpTest(v2, v2t+"("+v2s+")\n")
	addDumpTest(pv2, "&"+v2t+"("+v2s+")\n")
	addDumpTest(&pv2, "&&"+v2t+"("+v2s+")\n")
	addDumpTest(nv2, "(*"+v2t+")(nil)\n")

	// Function with multiple params and multiple returns.
	var v3 = func(i int, s string) (b bool, err error) {
		return true, nil
	}
	nv3 := (*func(int, string) (bool, error))(nil)
	pv3 := &v3
	v3t := "func(int, string) (bool, error)"
	v3s := fmt.Sprintf("%p", v3)
	addDumpTest(v3, v3t+"("+v3s+")\n")
	addDumpTest(pv3, "&"+v3t+"("+v3s+")\n")
	addDumpTest(&pv3, "&&"+v3t+"("+v3s+")\n")
	addDumpTest(nv3, "(*"+v3t+")(nil)\n")
}

func addCircularDumpTests() {
	// Struct that is circular through self referencing.
	type circular struct {
		c *circular
	}
	v := circular{nil}
	v.c = &v
	pv := &v
	vt := "utter_test.circular"
	vs := "{\n c: &" + vt + "{\n  c: (*" + vt + ")(<already shown>),\n },\n}"
	vs2 := "{\n c: (*" + vt + ")(<already shown>),\n}"
	addDumpTest(v, vt+vs+"\n")
	addDumpTest(pv, "&"+vt+vs2+"\n")
	addDumpTest(&pv, "&&"+vt+vs2+"\n")

	// Structs that are circular through cross referencing.
	v2 := xref1{nil}
	ts2 := xref2{&v2}
	v2.ps2 = &ts2
	pv2 := &v2
	v2t := "utter_test.xref1"
	v2t2 := "utter_test.xref2"
	v2s := "{\n ps2: &" + v2t2 +
		"{\n  ps1: &" + v2t +
		"{\n   ps2: (*" + v2t2 + ")(<already shown>),\n  },\n },\n}"
	v2s2 := "{\n ps2: &" + v2t2 + "{\n  ps1: (*" + v2t + ")(<already shown>),\n },\n}"
	addDumpTest(v2, v2t+v2s+"\n")
	addDumpTest(pv2, "&"+v2t+v2s2+"\n")
	addDumpTest(&pv2, "&&"+v2t+v2s2+"\n")

	// Structs that are indirectly circular.
	v3 := indirCir1{nil}
	tic2 := indirCir2{nil}
	tic3 := indirCir3{&v3}
	tic2.ps3 = &tic3
	v3.ps2 = &tic2
	pv3 := &v3
	v3t := "utter_test.indirCir1"
	v3t2 := "utter_test.indirCir2"
	v3t3 := "utter_test.indirCir3"
	v3s := "{\n ps2: &" + v3t2 +
		"{\n  ps3: &" + v3t3 +
		"{\n   ps1: &" + v3t +
		"{\n    ps2: (*" + v3t2 + ")(<already shown>),\n   },\n  },\n },\n}"
	v3s2 := "{\n ps2: &" + v3t2 +
		"{\n  ps3: &" + v3t3 +
		"{\n   ps1: (*" + v3t + ")(<already shown>),\n  },\n },\n}"
	addDumpTest(v3, v3t+v3s+"\n")
	addDumpTest(pv3, "&"+v3t+v3s2+"\n")
	addDumpTest(&pv3, "&&"+v3t+v3s2+"\n")
}

// TestDump executes all of the tests described by dumpTests.
func TestDump(t *testing.T) {
	// Setup tests.
	addIntDumpTests()
	addUintDumpTests()
	addBoolDumpTests()
	addFloatDumpTests()
	addComplexDumpTests()
	addArrayDumpTests()
	addSliceDumpTests()
	addStringDumpTests()
	addInterfaceDumpTests()
	addMapDumpTests()
	addStructDumpTests()
	addUintptrDumpTests()
	addUnsafePointerDumpTests()
	addChanDumpTests()
	addFuncDumpTests()
	addCircularDumpTests()
	addCgoDumpTests()

	t.Logf("Running %d tests", len(dumpTests))
	for i, test := range dumpTests {
		buf := new(bytes.Buffer)
		utter.Fdump(buf, test.in)
		s := buf.String()
		if testFailed(s, test.wants) {
			t.Errorf("Dump #%d\n got: %q\n %s", i, s, stringizeWants(test.wants))
			continue
		}
	}
}

func TestDumpSortedKeys(t *testing.T) {
	cfg := utter.ConfigState{SortKeys: true}
	s := cfg.Sdump(map[int]string{1: "1", 3: "3", 2: "2"})
	expected := `map[int]string{
int(1): string("1"),
int(2): string("2"),
int(3): string("3"),
}
`
	if s != expected {
		t.Errorf("Sorted keys mismatch:\n  %v %v", s, expected)
	}
}

type limitedWriter struct {
	limit int
	buf   bytes.Buffer
}

func newLimitedWriter(limit int) *limitedWriter {
	return &limitedWriter{limit: limit}
}

func (w *limitedWriter) Write(b []byte) (int, error) {
	n, err := w.buf.Write(b)
	if err != nil {
		return n, err
	}
	if len := w.buf.Len(); len > w.limit {
		panic(fmt.Sprintf("buffer longer than limit: %d > %d:\n%s",
			len, w.limit, w.buf.Bytes()))
	}
	return n, nil
}

var sliceElementCycles = []struct {
	v    interface{}
	want string
}{
	{
		v: func() interface{} {
			r := make([]interface{}, 1)
			r[0] = r
			return r
		}(),
		// We cannot detect the cycle until at least once around
		// the cycle as the initial v seen by utter.Dump was not
		// addressable.
		want: `[]interface{}{
 []interface{}(<already shown>),
}
`,
	},
	{
		v: func() interface{} {
			r := make([]interface{}, 1)
			r[0] = r
			return &r
		}(),
		want: `&[]interface{}{
 []interface{}(<already shown>),
}
`,
	},
	{
		v: func() interface{} {
			r := make([]interface{}, 1)
			r[0] = &r
			return &r
		}(),
		want: `&[]interface{}{
 &[]interface{}(<already shown>),
}
`,
	},
	{
		v: func() interface{} {
			type recurrence struct {
				v []interface{}
			}
			r := recurrence{make([]interface{}, 1)}
			r.v[0] = r
			return r
		}(),
		// We cannot detect the cycle until at least once around
		// the cycle as the initial v seen by utter.Dump was not
		// addressable.
		want: `utter_test.recurrence{
 v: []interface{}{
  utter_test.recurrence{
   v: []interface{}(<already shown>),
  },
 },
}
`,
	},
	{
		v: func() interface{} {
			type recurrence struct {
				v []interface{}
			}
			r := recurrence{make([]interface{}, 1)}
			r.v[0] = r
			return &r
		}(),
		want: `&utter_test.recurrence{
 v: []interface{}{
  utter_test.recurrence{
   v: []interface{}(<already shown>),
  },
 },
}
`,
	},
	{
		v: func() interface{} {
			type container struct {
				v []int
			}
			return &container{[]int{1}}
		}(),
		want: `&utter_test.container{
 v: []int{
  int(1),
 },
}
`,
	},
}

// https://github.com/kortschak/utter/issues/5
func TestIssue5Slices(t *testing.T) {
	for _, test := range sliceElementCycles {
		w := newLimitedWriter(512)
		func() {
			defer func() {
				r := recover()
				if r != nil {
					t.Errorf("limited writer panicked: probable cycle: %v", r)
				}
			}()
			utter.Fdump(w, test.v)
			got := w.buf.String()
			if got != test.want {
				t.Errorf("unexpected value:\ngot:\n%swant:\n%s", got, test.want)
			}
		}()
	}
}

var mapElementCycles = []struct {
	v    interface{}
	want string
}{
	{
		v: func() interface{} {
			r := make(map[int]interface{}, 1)
			r[0] = r
			return r
		}(),
		want: `map[int]interface{}{
 int(0): map[int]interface{}(<already shown>),
}
`,
	},
	{
		v: func() interface{} {
			r := make(map[int]interface{}, 1)
			r[0] = r
			return &r
		}(),
		want: `&map[int]interface{}{
 int(0): map[int]interface{}(<already shown>),
}
`,
	},
	{
		v: func() interface{} {
			r := make(map[int]interface{}, 1)
			r[0] = &r
			return &r
		}(),
		want: `&map[int]interface{}{
 int(0): &map[int]interface{}(<already shown>),
}
`,
	},
	{
		v: func() interface{} {
			type recurrence struct {
				v map[int]interface{}
			}
			r := recurrence{make(map[int]interface{}, 1)}
			r.v[0] = r
			return r
		}(),
		want: `utter_test.recurrence{
 v: map[int]interface{}{
  int(0): utter_test.recurrence{
   v: map[int]interface{}(<already shown>),
  },
 },
}
`,
	},
	{
		v: func() interface{} {
			type recurrence struct {
				v map[int]interface{}
			}
			r := recurrence{make(map[int]interface{}, 1)}
			r.v[0] = r
			return &r
		}(),
		want: `&utter_test.recurrence{
 v: map[int]interface{}{
  int(0): utter_test.recurrence{
   v: map[int]interface{}(<already shown>),
  },
 },
}
`,
	},
	// The following test is to confirm that the recursion detection
	// is not overly zealous by missing identifying the address of slices.
	// This is https://github.com/kortschak/utter/issues/12.
	{
		v: map[interface{}][]interface{}{
			"outer": []interface{}{
				map[interface{}]interface{}{
					"inner": []interface{}{"value"},
				},
			},
		},
		want: `map[interface{}][]interface{}{
 string("outer"): []interface{}{
  map[interface{}]interface{}{
   string("inner"): []interface{}{
    string("value"),
   },
  },
 },
}
`,
	},
}

// https://github.com/kortschak/utter/issues/5
// https://github.com/kortschak/utter/issues/12
func TestIssue5Maps(t *testing.T) {
	for _, test := range mapElementCycles {
		w := newLimitedWriter(512)
		func() {
			defer func() {
				r := recover()
				if r != nil {
					t.Errorf("limited writer panicked: probable cycle: %v", r)
				}
			}()
			utter.Fdump(w, test.v)
			got := w.buf.String()
			if got != test.want {
				t.Errorf("unexpected value:\ngot:\n%swant:\n%s", got, test.want)
			}
		}()
	}
}
