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

package utter_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/kortschak/utter"
)

// utterFunc is used to identify which public function of the utter package or
// ConfigState a test applies to.
type utterFunc int

const (
	fCSFdump utterFunc = iota
	fCSSdump
	fSdump
)

// Map of utterFunc values to names for pretty printing.
var utterFuncStrings = map[utterFunc]string{
	fCSFdump: "ConfigState.Fdump",
	fCSSdump: "ConfigState.Sdump",
	fSdump:   "utter.Sdump",
}

func (f utterFunc) String() string {
	if s, ok := utterFuncStrings[f]; ok {
		return s
	}
	return fmt.Sprintf("Unknown utterFunc (%d)", int(f))
}

// utterTest is used to describe a test to be performed against the public
// functions of the utter package or ConfigState.
type utterTest struct {
	cs   *utter.ConfigState
	f    utterFunc
	in   interface{}
	want string
}

// utterTests houses the tests to be performed against the public functions of
// the utter package and ConfigState.
//
// These tests are only intended to ensure the public functions are exercised
// and are intentionally not exhaustive of types.  The exhaustive type
// tests are handled in the dump and format tests.
var utterTests []utterTest

// redirStdout is a helper function to return the standard output from f as a
// byte slice.
func redirStdout(f func()) ([]byte, error) {
	tempFile, err := ioutil.TempFile("", "ss-test")
	if err != nil {
		return nil, err
	}
	fileName := tempFile.Name()
	defer os.Remove(fileName) // Ignore error

	origStdout := os.Stdout
	os.Stdout = tempFile
	f()
	os.Stdout = origStdout
	tempFile.Close()

	return ioutil.ReadFile(fileName)
}

func initSpewTests() {
	// Config states with various settings.
	scsDefault := utter.NewDefaultConfig()

	// Byte slice without comments.
	noComDefault := utter.NewDefaultConfig()
	noComDefault.CommentBytes = false

	// Byte slice with 8 columns.
	bs8Default := utter.NewDefaultConfig()
	bs8Default.BytesWidth = 8

	// Ignore unexported fields.
	ignUnexDefault := utter.NewDefaultConfig()
	ignUnexDefault.IgnoreUnexported = true

	// Elide implicit types.
	elideTypeDefault := utter.NewDefaultConfig()
	elideTypeDefault.ElideType = true

	// depthTester is used to test max depth handling for structs, array, slices
	// and maps.
	type depthTester struct {
		ic    indirCir1
		arr   [1]string
		slice []string
		m     map[string]int
	}

	utterTests = []utterTest{
		{scsDefault, fCSFdump, int8(127), "int8(127)\n"},
		{scsDefault, fCSSdump, uint8(64), "uint8(0x40)\n"},
		{scsDefault, fSdump, complex(-10, -20), "complex128(-10-20i)\n"},
		{noComDefault, fCSFdump, []byte{1, 2, 3, 4, 5, 0},
			"[]uint8{\n 0x01, 0x02, 0x03, 0x04, 0x05, 0x00,\n}\n",
		},
		{bs8Default, fCSFdump, []byte{1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0}, "[]uint8{\n" +
			" 0x01, 0x02, 0x03, 0x04, 0x05, 0x00, 0x01, 0x02, // |........|\n" +
			" 0x03, 0x04, 0x05, 0x00,                         // |....|\n}\n",
		},
		{ignUnexDefault, fCSFdump, Foo{Bar{flag: 1}, map[interface{}]interface{}{"one": true}},
			"utter_test.Foo{\n ExportedField: map[interface{}]interface{}{\n  string(\"one\"): bool(true),\n },\n}\n",
		},
		{elideTypeDefault, fCSFdump, float64(1), "1.0\n"},
		{elideTypeDefault, fCSFdump, float32(1), "float32(1)\n"},
		{elideTypeDefault, fCSFdump, int(1), "1\n"},
		{elideTypeDefault, fCSFdump, []interface{}{true, 1.0, float32(1), "one", 1, 'a'},
			"[]interface{}{\n true,\n 1.0,\n float32(1),\n \"one\",\n 1,\n int32(97),\n}\n",
		},
		{elideTypeDefault, fCSFdump, Foo{Bar{flag: 1}, map[interface{}]interface{}{"one": true}}, "utter_test.Foo{\n" +
			" unexportedField: utter_test.Bar{\n  flag: 1,\n  data: 0,\n },\n" +
			" ExportedField: map[interface{}]interface{}{\n  \"one\": true,\n },\n}\n",
		},
		{elideTypeDefault, fCSFdump, map[interface{}]interface{}{"one": nil}, "map[interface{}]interface{}{\n \"one\": nil,\n}\n"},
		{elideTypeDefault, fCSFdump, float32(1), "float32(1)\n"},
		{elideTypeDefault, fCSFdump, float64(1), "1.0\n"},
		{elideTypeDefault, fCSFdump, func() *float64 { f := 1.0; return &f }(), "&float64(1)\n"},
		{elideTypeDefault, fCSFdump, []float32{1, 2, 3, 4, 5}, "[]float32{\n 1.0,\n 2.0,\n 3.0,\n 4.0,\n 5.0,\n}\n"},
	}
}

// TestSpew executes all of the tests described by utterTests.
func TestSpew(t *testing.T) {
	initSpewTests()

	t.Logf("Running %d tests", len(utterTests))
	for i, test := range utterTests {
		buf := new(bytes.Buffer)
		switch test.f {
		case fCSFdump:
			test.cs.Fdump(buf, test.in)

		case fCSSdump:
			str := test.cs.Sdump(test.in)
			buf.WriteString(str)

		case fSdump:
			str := utter.Sdump(test.in)
			buf.WriteString(str)

		default:
			t.Errorf("%v #%d unrecognized function", test.f, i)
			continue
		}
		s := buf.String()
		if test.want != s {
			t.Errorf("ConfigState #%d\n got: %q\nwant: %q", i, s, test.want)
			continue
		}
	}
}
