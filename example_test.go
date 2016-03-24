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
	"fmt"

	"github.com/kortschak/utter"
)

type Flag int

const (
	flagOne Flag = iota
	flagTwo
)

var flagStrings = map[Flag]string{
	flagOne: "flagOne",
	flagTwo: "flagTwo",
}

func (f Flag) String() string {
	if s, ok := flagStrings[f]; ok {
		return s
	}
	return fmt.Sprintf("Unknown flag (%d)", int(f))
}

type Bar struct {
	flag Flag
	data uintptr
}

type Foo struct {
	unexportedField Bar
	ExportedField   map[interface{}]interface{}
}

// This example demonstrates how to use Dump to dump variables to stdout.
func ExampleDump() {
	// The following package level declarations are assumed for this example:
	/*
		type Flag int

		const (
			flagOne Flag = iota
			flagTwo
		)

		var flagStrings = map[Flag]string{
			flagOne: "flagOne",
			flagTwo: "flagTwo",
		}

		func (f Flag) String() string {
			if s, ok := flagStrings[f]; ok {
				return s
			}
			return fmt.Sprintf("Unknown flag (%d)", int(f))
		}

		type Bar struct {
			flag Flag
			data uintptr
		}

		type Foo struct {
			unexportedField Bar
			ExportedField   map[interface{}]interface{}
		}
	*/

	// Setup some sample data structures for the example.
	bar := Bar{Flag(flagTwo), uintptr(0)}
	s1 := Foo{bar, map[interface{}]interface{}{"one": true}}
	f := Flag(5)
	b := []byte{
		0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18,
		0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20,
		0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28,
		0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30,
		0x31, 0x32,
	}

	// Dump!
	utter.Dump([]interface{}{s1, f, b})

	// Output:
	//
	// []interface{}{
	//  utter_test.Foo{
	//   unexportedField: utter_test.Bar{
	//    flag: utter_test.Flag(1),
	//    data: uintptr(0),
	//   },
	//   ExportedField: map[interface{}]interface{}{
	//    string("one"): bool(true),
	//   },
	//  },
	//  utter_test.Flag(5),
	//  []uint8{
	//   0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f, 0x20, // |............... |
	//   0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27, 0x28, 0x29, 0x2a, 0x2b, 0x2c, 0x2d, 0x2e, 0x2f, 0x30, // |!"#$%&'()*+,-./0|
	//   0x31, 0x32, /*                                                                               */ // |12|
	//  },
	// }
}

// This example demonstrates how to use a ConfigState.
func ExampleConfigState() {
	// Modify the indent level of the ConfigState only.  The global
	// configuration is not modified.
	scs := utter.ConfigState{Indent: "\t"}

	// Output using the ConfigState instance.
	v := map[string]int{"one": 1}
	scs.Dump(v)

	// Output:
	//
	// map[string]int{
	// 	string("one"): int(1),
	// }
}

// This example demonstrates how to use ConfigState.Dump to dump variables to
// stdout
func ExampleConfigState_Dump() {
	// See the top-level Dump example for details on the types used in this
	// example.

	// Create two ConfigState instances with different indentation.
	scs := utter.ConfigState{Indent: "\t"}
	scs2 := utter.ConfigState{Indent: " "}

	// Setup some sample data structures for the example.
	bar := Bar{Flag(flagTwo), uintptr(0)}
	s1 := Foo{bar, map[interface{}]interface{}{"one": true}}

	// Dump using the ConfigState instances.
	scs.Dump(s1)
	scs2.Dump(s1)

	// Output:
	//
	// utter_test.Foo{
	// 	unexportedField: utter_test.Bar{
	// 		flag: utter_test.Flag(1),
	// 		data: uintptr(0),
	// 	},
	// 	ExportedField: map[interface{}]interface{}{
	// 		string("one"): bool(true),
	// 	},
	// }
	// utter_test.Foo{
	//  unexportedField: utter_test.Bar{
	//   flag: utter_test.Flag(1),
	//   data: uintptr(0),
	//  },
	//  ExportedField: map[interface{}]interface{}{
	//   string("one"): bool(true),
	//  },
	// }
}
