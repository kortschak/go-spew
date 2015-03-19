/*
 * Copyright (c) 2013 Dave Collins <dave@davec.name>
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
Package utter implements a deep pretty printer for Go data structures to aid in
debugging.

A quick overview of the additional features utter provides over the built-in
printing facilities for Go data types are as follows:

	* Pointers are dereferenced and followed
	* Circular data structures are detected and handled properly
	* Byte arrays and slices are dumped like the hexdump -C command which
	  includes offsets, byte values in hex, and ASCII output (only when using
	  Dump style)

There are two different approaches utter allows for dumping Go data structures:

	* Dump style which prints with newlines, customizable indentation,
	  and additional debug information such as types and all pointer addresses
	  used to indirect to the final value
	* A custom Formatter interface that integrates cleanly with the standard fmt
	  package and replaces %v, %+v, %#v, and %#+v to provide inline printing
	  similar to the default %v while providing the additional functionality
	  outlined above and passing unsupported format verbs such as %x and %q
	  along to fmt

Quick Start

This section demonstrates how to quickly get started with utter.  See the
sections below for further details on formatting and configuration options.

To dump a variable with full newlines, indentation, type, and pointer
information use Dump, Fdump, or Sdump:
	utter.Dump(myVar1, myVar2, ...)
	utter.Fdump(someWriter, myVar1, myVar2, ...)
	str := utter.Sdump(myVar1, myVar2, ...)

Alternatively, if you would prefer to use format strings with a compacted inline
printing style, use the convenience wrappers Printf, Fprintf, etc with
%v (most compact), %+v (adds pointer addresses), %#v (adds types), or
%#+v (adds types and pointer addresses):
	utter.Printf("myVar1: %v -- myVar2: %+v", myVar1, myVar2)
	utter.Printf("myVar3: %#v -- myVar4: %#+v", myVar3, myVar4)
	utter.Fprintf(someWriter, "myVar1: %v -- myVar2: %+v", myVar1, myVar2)
	utter.Fprintf(someWriter, "myVar3: %#v -- myVar4: %#+v", myVar3, myVar4)

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

	* SortKeys
		Specifies map keys should be sorted before being printed. Use
		this to have a more deterministic, diffable output.  Note that
		only native types (bool, int, uint, floats, uintptr and string)
		are supported with other types sorted according to the
		reflect.Value.String() output which guarantees display stability.
		Natural map order is used by default.

Dump Usage

Simply call utter.Dump with a list of variables you want to dump:

	utter.Dump(myVar1, myVar2, ...)

You may also call utter.Fdump if you would prefer to output to an arbitrary
io.Writer.  For example, to dump to standard error:

	utter.Fdump(os.Stderr, myVar1, myVar2, ...)

A third option is to call utter.Sdump to get the formatted output as a string:

	str := utter.Sdump(myVar1, myVar2, ...)

Sample Dump Output

See the Dump example for details on the setup of the types and variables being
shown here.

	(main.Foo) {
	 unexportedField: (*main.Bar)(0xf84002e210)({
	  flag: (main.Flag) 1,
	  data: (uintptr) <nil>
	 }),
	 ExportedField: (map[interface {}]interface {}) (len=1) {
	  (string) (len=3) "one": (bool) true
	 }
	}

Byte (and uint8) arrays and slices are displayed uniquely like the hexdump -C
command as shown.
	([]uint8) (len=32 cap=32) {
	 00000000  11 12 13 14 15 16 17 18  19 1a 1b 1c 1d 1e 1f 20  |............... |
	 00000010  21 22 23 24 25 26 27 28  29 2a 2b 2c 2d 2e 2f 30  |!"#$%&'()*+,-./0|
	 00000020  31 32                                             |12|
	}
*/
package utter
