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

package utter

import (
	"bytes"
	"io"
	"os"
)

// ConfigState houses the configuration options used by utter to format and
// display values.  There is a global instance, Config, that is used to control
// all top-level Formatter and Dump functionality.  Each ConfigState instance
// provides methods equivalent to the top-level functions.
//
// The zero value for ConfigState provides no indentation.  You would typically
// want to set it to a space or a tab.
//
// Alternatively, you can use NewDefaultConfig to get a ConfigState instance
// with default settings.  See the documentation of NewDefaultConfig for default
// values.
type ConfigState struct {
	// Indent specifies the string to use for each indentation level.  The
	// global config instance that all top-level functions use set this to a
	// single space by default.  If you would like more indentation, you might
	// set this to a tab with "\t" or perhaps two spaces with "  ".
	Indent string

	// NumericWidth specifies the number of columns to use when dumping
	// a numeric slice or array (including bool). Zero specifies all entries
	// on one line.
	NumericWidth int

	// StringWidth specifies the number of columns to use when dumping
	// a string slice or array. Zero specifies all entries on one line.
	StringWidth int

	// BytesWidth specifies the number of byte columns to use when dumping a
	// byte slice or array.
	BytesWidth int

	// CommentBytes specifies whether byte slice or array dumps have ASCII
	// comment annotations.
	CommentBytes bool

	// CommentPointers specifies whether pointer information will be added
	// as comments.
	CommentPointers bool

	// IgnoreUnexported specifies that unexported struct fields should be
	// ignored during a dump.
	IgnoreUnexported bool

	// ElideType specifies that type information defined by context should
	// not be printed in a dump.
	ElideType bool

	// SortKeys specifies map keys should be sorted before being printed. Use
	// this to have a more deterministic, diffable output.  Note that only
	// native types (bool, int, uint, floats, uintptr and string) are supported
	// with other types sorted according to the reflect.Value.String() output
	// which guarantees display stability.
	SortKeys bool
}

// Config is the active configuration of the top-level functions.
// The configuration can be changed by modifying the contents of utter.Config.
var Config = ConfigState{
	Indent:       " ",
	NumericWidth: 1,
	StringWidth:  1,
	BytesWidth:   16,
	CommentBytes: true,
}

// Fdump formats and displays the passed arguments to io.Writer w.  It formats
// exactly the same as Dump.
func (c *ConfigState) Fdump(w io.Writer, a interface{}) {
	fdump(c, w, a)
}

/*
Dump displays the passed parameters to standard out with newlines, customizable
indentation, and additional debug information such as complete types and all
pointer addresses used to indirect to the final value.  It provides the
following features over the built-in printing facilities provided by the fmt
package:

	* Pointers are dereferenced and followed
	* Circular data structures are detected and handled properly
	* Byte arrays and slices are dumped in a way similar to the hexdump -C command
	  which includes byte values in hex, and ASCII output

The configuration options are controlled by modifying the public members
of c.  See ConfigState for options documentation.

See Fdump if you would prefer dumping to an arbitrary io.Writer or Sdump to
get the formatted result as a string.
*/
func (c *ConfigState) Dump(a interface{}) {
	fdump(c, os.Stdout, a)
}

// Sdump returns a string with the passed arguments formatted exactly the same
// as Dump.
func (c *ConfigState) Sdump(a interface{}) string {
	var buf bytes.Buffer
	fdump(c, &buf, a)
	return buf.String()
}

// NewDefaultConfig returns a ConfigState with the following default settings.
//
// 	Indent: " "
//  NumericWidth: 1,
//  StringWidth: 1,
// 	BytesWidth: 16
// 	CommentBytes: true
// 	CommentPointers: false
//  IgnoreUnexported: false
//  ElideType: false
// 	SortKeys: false
func NewDefaultConfig() *ConfigState {
	return &ConfigState{
		Indent:       " ",
		NumericWidth: 1,
		StringWidth:  1,
		BytesWidth:   16,
		CommentBytes: true,
	}
}
