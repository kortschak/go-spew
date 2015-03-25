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
Package utter implements a deep pretty printer for Go data structures to aid
data snapshotting.

A quick overview of the additional features utter provides over the built-in
printing facilities for Go data types are as follows:

	* Pointers are dereferenced and followed
	* Circular data structures are detected and annotated
	* Byte arrays and slices are dumped in a way similar to the hexdump -C command
	  which includes byte values in hex, and ASCII output

The approache utter allows for dumping Go data structures is less flexible than
its parent tool. It has just a:

	* Dump style which prints with newlines and customizable indentation

Quick Start

This section demonstrates how to quickly get started with utter.  See the
sections below for further details on formatting and configuration options.

To dump a variable with full newlines, indentation, type, and pointer
information use Dump, Fdump, or Sdump:
	utter.Dump(myVar1)
	utter.Fdump(someWriter, myVar1)
	str := utter.Sdump(myVar1)

Configuration Options

Configuration of utter is handled by fields in the ConfigState type.  For
convenience, all of the top-level functions use a global state available
via the utter.Config global.

It is also possible to create a ConfigState instance that provides methods
equivalent to the top-level functions.  This allows concurrent configuration
options.  See the ConfigState documentation for more details.

The following configuration options are available:
	* Indent
		String to use for each indentation level for Dump functions.
		It is a single space by default.  A popular alternative is "\t".

	* NumericWidth
		NumericWidth specifies the number of columns to use when dumping
		a numeric slice or array (including bool). Zero specifies all entries
		on one line.

	* StringWidth
		StringWidth specifies the number of columns to use when dumping
		a string slice or array. Zero specifies all entries on one line.

	* BytesWidth
		Number of byte columns to use when dumping byte slices and arrays.

	* CommentBytes
		Specifies whether ASCII comment annotations are attached to byte
		slice and array dumps.

	* CommentPointers
		CommentPointers specifies whether pointer information will be added
		as comments.

	* IgnoreUnexported
		Specifies that unexported fields should be ignored.

	* ElideType
		ElideType specifies that type information defined by context should
		not be printed in a dump.

	* SortKeys
		Specifies map keys should be sorted before being printed. Use
		this to have a more deterministic, diffable output.  Note that
		only native types (bool, int, uint, floats, uintptr and string)
		are supported with other types sorted according to the
		reflect.Value.String() output which guarantees display stability.
		Natural map order is used by default.

Dump Usage

Simply call utter.Dump with a list of variables you want to dump:

	utter.Dump(myVar1)

You may also call utter.Fdump if you would prefer to output to an arbitrary
io.Writer.  For example, to dump to standard error:

	utter.Fdump(os.Stderr, myVar1)

A third option is to call utter.Sdump to get the formatted output as a string:

	str := utter.Sdump(myVar1)

Sample Dump Output

See the Dump example for details on the setup of the types and variables being
shown here.

	main.Foo{
	 unexportedField: &main.Bar{
	  flag: main.Flag(1),
	  data: uintptr(0),
	 },
	 ExportedField: map[interface{}]interface{}{
	  string("one"): bool(true),
	 },
	}

Byte (and uint8) arrays and slices are displayed uniquely, similar to the hexdump -C
command as shown.

	[]uint8{
	 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, // |........|
	 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, // |....... |
	 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, // |!"#$%&'(|
	 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, // |)*+,-./0|
	 0x31, 0x32,                                     // |12|
	}
*/
package utter
