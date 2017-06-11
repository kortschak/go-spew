utter
=====

[![Build Status](https://travis-ci.org/kortschak/utter.svg?branch=master)](https://travis-ci.org/kortschak/utter) [![Coverage Status](https://coveralls.io/repos/kortschak/utter/badge.svg?branch=master)](https://coveralls.io/r/kortschak/utter?branch=master)

utter is a fork of the outstanding [go-spew tool](https://github.com/davecgh/go-spew).
Where go-spew is an aid for debugging, providing annotation of dumped datastructures,
utter is a tool for taking snapshots of data structures to include in tests or other
code. An utter dump will not construct cyclic structure literals and a number of
pseudo-code representations of pointer-based structures will require subsequent
processing.

A comprehensive suite of tests with near 100% test coverage is provided to ensure
proper functionality.  utter is licensed under the liberal ISC license, so it may
be used in open source or commercial projects.

## Documentation

[![GoDoc](https://godoc.org/github.com/kortschak/utter?status.svg)](http://godoc.org/github.com/kortschak/utter)

Full `go doc` style documentation for the project can be viewed online without
installing this package by using the excellent GoDoc site here:
http://godoc.org/github.com/kortschak/utter

You can also view the documentation locally once the package is installed with
the `godoc` tool by running `godoc -http=":6060"` and pointing your browser to
http://localhost:6060/pkg/github.com/kortschak/utter

## Installation

```bash
$ go get -u github.com/kortschak/utter
```

## Quick Start

To dump a variable with full newlines, indentation, type, and pointer
information use Dump, Fdump, or Sdump:

```Go
utter.Dump(myVar1)
utter.Fdump(someWriter, myVar1)
str := utter.Sdump(myVar1)
```

## Sample Dump Output

```
main.Foo{
 unexportedField: &main.Bar{
  flag: main.Flag(1),
  data: uintptr(0),
 },
 ExportedField: map[interface{}]interface{}{
  string("one"): bool(true),
 },
}
```

```
[]uint8{
 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, // |........|
 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, // |....... |
 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, // |!"#$%&'(|
 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, // |)*+,-./0|
 0x31, 0x32, /*                               */ // |12|
}
```

## Configuration Options

Configuration of utter is handled by fields in the ConfigState type. For
convenience, all of the top-level functions use a global state available via the
utter.Config global.

It is also possible to create a ConfigState instance that provides methods
equivalent to the top-level functions. This allows concurrent configuration
options. See the ConfigState documentation for more details.

```
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
```

## License

utter is licensed under the liberal ISC License.
