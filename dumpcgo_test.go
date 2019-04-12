// Copyright (c) 2013 Dave Collins <dave@davec.name>
// Copyright (c) 2015 Dan Kortschak <dan.kortschak@adelaide.edu.au>
//
// Permission to use, copy, modify, and distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.
//
// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR
// ANY SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN
// ACTION OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF
// OR IN CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

// NOTE: Due to the following build constraints, this file will only be compiled
// when both cgo is supported and "-tags testcgo" is added to the go test
// command line.  This means the cgo tests are only added (and hence run) when
// specifially requested.  This configuration is used because utter itself
// does not require cgo to run even though it does handle certain cgo types
// specially.  Rather than forcing all clients to require cgo and an external
// C compiler just to run the tests, this scheme makes them optional.
// +build cgo,testcgo

package utter_test

import "github.com/kortschak/utter/testdata"

func addCgoDumpTests() {
	// C char pointer.
	v := testdata.GetCgoCharPointer()
	nv := testdata.GetCgoNullCharPointer()
	pv := &v
	vt := "testdata._Ctype_char"
	vs := "116"
	addDumpTest(v, "&"+vt+"("+vs+")\n")
	addDumpTest(pv, "&&"+vt+"("+vs+")\n")
	addDumpTest(&pv, "&&&"+vt+"("+vs+")\n")
	addDumpTest(nv, "(*"+vt+")(nil)\n")

	// C char array.
	v2 := testdata.GetCgoCharArray()
	v2t := "[6]testdata._Ctype_char"
	v2s := "" +
		"{\n 0x74, 0x65, 0x73, 0x74, 0x32, 0x00, // |test2.|\n}"
	addDumpTest(v2, v2t+v2s+"\n")

	// C unsigned char array.
	v3 := testdata.GetCgoUnsignedCharArray()
	v3tbefore1_6 := "[6]testdata._Ctype_unsignedchar"
	v3t1_6 := "[6]testdata._Ctype_uchar"
	v3s := "" +
		"{\n 0x74, 0x65, 0x73, 0x74, 0x33, 0x00, // |test3.|\n}"
	addDumpTest(v3, v3tbefore1_6+v3s+"\n", v3t1_6+v3s+"\n")

	// C signed char array.
	v4 := testdata.GetCgoSignedCharArray()
	v4t := "[6]testdata._Ctype_schar"
	v4t2 := "testdata._Ctype_schar"
	v4s := "{\n " +
		v4t2 + "(116),\n " +
		v4t2 + "(101),\n " +
		v4t2 + "(115),\n " +
		v4t2 + "(116),\n " +
		v4t2 + "(52),\n " +
		v4t2 + "(0),\n}"
	addDumpTest(v4, v4t+v4s+"\n")

	// C uint8_t array.
	v5 := testdata.GetCgoUint8tArray()
	v5tbefore1_10 := "[6]testdata._Ctype_uint8_t"
	v5t1_10 := "[6]testdata._Ctype_uchar"
	v5tafter1_12 := "[6]uint8"
	v5s := "" +
		"{\n 0x74, 0x65, 0x73, 0x74, 0x35, 0x00, // |test5.|\n}"
	addDumpTest(v5, v5tbefore1_10+v5s+"\n", v5t1_10+v5s+"\n", v5tafter1_12+v5s+"\n")

	// C typedefed unsigned char array.
	v6 := testdata.GetCgoTypdefedUnsignedCharArray()
	v6tbefore1_10 := "[6]testdata._Ctype_custom_uchar_t"
	v6t1_10 := "[6]testdata._Ctype_uchar"
	v6s := "" +
		"{\n 0x74, 0x65, 0x73, 0x74, 0x36, 0x00, // |test6.|\n}"
	addDumpTest(v6, v6tbefore1_10+v6s+"\n", v6t1_10+v6s+"\n")
}
