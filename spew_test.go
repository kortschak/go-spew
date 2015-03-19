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
	cs     *utter.ConfigState
	f      utterFunc
	format string
	in     interface{}
	want   string
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

	// depthTester is used to test max depth handling for structs, array, slices
	// and maps.
	type depthTester struct {
		ic    indirCir1
		arr   [1]string
		slice []string
		m     map[string]int
	}

	utterTests = []utterTest{
		{scsDefault, fCSFdump, "", int8(127), "int8(127)\n"},
		{scsDefault, fCSSdump, "", uint8(64), "uint8(64)\n"},
		{scsDefault, fSdump, "", complex(-10, -20), "complex128(-10-20i)\n"},
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
			t.Errorf("ConfigState #%d\n got: %s want: %s", i, s, test.want)
			continue
		}
	}
}
