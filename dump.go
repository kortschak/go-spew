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
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	// uint8Type is a reflect.Type representing a uint8.  It is used to
	// convert cgo types to uint8 slices for hexdumping.
	uint8Type = reflect.TypeOf(uint8(0))

	// cCharRE is a regular expression that matches a cgo char.
	// It is used to detect character arrays to hexdump them.
	cCharRE = regexp.MustCompile("^.*\\._Ctype_char$")

	// cUnsignedCharRE is a regular expression that matches a cgo unsigned
	// char.  It is used to detect unsigned character arrays to hexdump
	// them.
	cUnsignedCharRE = regexp.MustCompile("^.*\\._Ctype_unsignedchar$")

	// cUint8tCharRE is a regular expression that matches a cgo uint8_t.
	// It is used to detect uint8_t arrays to hexdump them.
	cUint8tCharRE = regexp.MustCompile("^.*\\._Ctype_uint8_t$")
)

type addrType struct {
	addr uintptr
	typ  reflect.Type
}

// dumpState contains information about the state of a dump operation.
type dumpState struct {
	w                io.Writer
	depth            int
	pointers         map[uintptr]int
	nodes            map[addrType]struct{}
	ignoreNextType   bool
	ignoreNextIndent bool
	cs               *ConfigState
}

// indent performs indentation according to the depth level and cs.Indent
// option.
func (d *dumpState) indent() {
	if d.ignoreNextIndent {
		d.ignoreNextIndent = false
		return
	}
	d.w.Write(bytes.Repeat([]byte(d.cs.Indent), d.depth))
}

// unpackValue returns values inside of non-nil interfaces when possible.
// This is useful for data types like structs, arrays, slices, and maps which
// can contain varying types packed inside an interface.
func (d *dumpState) unpackValue(v reflect.Value) (val reflect.Value, wasPtr, static bool, addr uintptr) {
	if v.CanAddr() {
		addr = v.Addr().Pointer()
	}
	if v.Kind() == reflect.Interface && !v.IsNil() {
		return v.Elem(), v.Kind() == reflect.Ptr, false, addr
	}
	return v, v.Kind() == reflect.Ptr, true, addr
}

// dumpPtr handles formatting of pointers by indirecting them as necessary.
func (d *dumpState) dumpPtr(v reflect.Value) {
	// Remove pointers at or below the current depth from map used to detect
	// circular refs.
	for k, depth := range d.pointers {
		if depth >= d.depth {
			delete(d.pointers, k)
		}
	}

	// Keep list of all dereferenced pointers to show later.
	var pointerChain []uintptr

	// Figure out how many levels of indirection there are by dereferencing
	// pointers and unpacking interfaces down the chain while detecting circular
	// references.
	var nilFound, cycleFound bool
	indirects := 0
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			nilFound = true
			break
		}
		indirects++
		addr := v.Pointer()
		if d.cs.CommentPointers {
			pointerChain = append(pointerChain, addr)
		}
		if pd, ok := d.pointers[addr]; ok && pd < d.depth {
			cycleFound = true
			indirects--
			break
		}
		d.pointers[addr] = d.depth

		v = v.Elem()
		if v.Kind() == reflect.Interface {
			if v.IsNil() {
				nilFound = true
				break
			}
			v = v.Elem()
		}
	}

	// Display type information.
	d.w.Write(bytes.Repeat(ampersandBytes, indirects))
	kind := v.Kind()
	bufferedChan := kind == reflect.Chan && v.Cap() != 0
	if kind == reflect.Ptr || bufferedChan {
		d.w.Write(openParenBytes)
	}
	typeBytes := []byte(v.Type().String())
	d.w.Write(bytes.Replace(typeBytes, interfaceTypeBytes, interfaceBytes, -1))
	switch {
	case bufferedChan:
		switch len := v.Len(); len {
		case 0:
			fmt.Fprintf(d.w, ", %d", v.Cap())
		case 1:
			fmt.Fprintf(d.w, ", %d /* %d element */", v.Cap(), len)
		default:
			fmt.Fprintf(d.w, ", %d /* %d elements */", v.Cap(), len)
		}
		fallthrough
	case kind == reflect.Ptr:
		d.w.Write(closeParenBytes)
	}

	// Display pointer information.
	if len(pointerChain) > 0 {
		d.w.Write(openCommentBytes)
		for i, addr := range pointerChain {
			if i > 0 {
				d.w.Write(pointerChainBytes)
			}
			printHexPtr(d.w, addr, true)
		}
		d.w.Write(closeCommentBytes)
	}

	// Display dereferenced value.
	switch {
	case nilFound:
		d.w.Write(openParenBytes)
		d.w.Write(nilBytes)
		d.w.Write(closeParenBytes)

	case cycleFound:
		d.w.Write(circularBytes)

	default:
		d.ignoreNextType = true
		var addr uintptr
		if v.CanAddr() {
			addr = v.Addr().Pointer()
		}
		d.dump(v, true, false, addr)
	}
}

// dumpSlice handles formatting of arrays and slices.  Byte (uint8 under
// reflection) arrays and slices are dumped in hexdump -C fashion.
func (d *dumpState) dumpSlice(v reflect.Value) {
	// Determine whether this type should be hex dumped or not.  Also,
	// for types which should be hexdumped, try to use the underlying data
	// first, then fall back to trying to convert them to a uint8 slice.
	var buf []uint8
	doConvert := false
	doHexDump := false
	nPeriod := 1
	numEntries := v.Len()
	vt := v.Type().Elem()
	if numEntries > 0 {
		vts := vt.String()
		switch kind := vt.Kind(); {
		// C types that need to be converted.
		case cCharRE.MatchString(vts):
			fallthrough
		case cUnsignedCharRE.MatchString(vts):
			fallthrough
		case cUint8tCharRE.MatchString(vts):
			doConvert = true

		// Try to use existing uint8 slices and fall back to converting
		// and copying if that fails.
		case kind == reflect.Uint8:
			// We need an addressable interface to convert the type back
			// into a byte slice.  However, the reflect package won't give
			// us an interface on certain things like unexported struct
			// fields in order to enforce visibility rules.  We use unsafe
			// to bypass these restrictions since this package does not
			// mutate the values.
			vs := v
			if !vs.CanInterface() || !vs.CanAddr() {
				vs = unsafeReflectValue(vs)
			}
			vs = vs.Slice(0, numEntries)

			// Use the existing uint8 slice if it can be type
			// asserted.
			iface := vs.Interface()
			if slice, ok := iface.([]uint8); ok {
				buf = slice
				doHexDump = true
				break
			}

			// The underlying data needs to be converted if it can't
			// be type asserted to a uint8 slice.
			doConvert = true

		case isNumeric(kind):
			nPeriod = d.cs.NumericWidth

		case kind == reflect.String:
			nPeriod = d.cs.StringWidth
		}

		// Copy and convert the underlying type if needed.
		if doConvert && vt.ConvertibleTo(uint8Type) {
			// Convert and copy each element into a uint8 byte
			// slice.
			buf = make([]uint8, numEntries)
			for i := 0; i < numEntries; i++ {
				vv := v.Index(i)
				buf[i] = uint8(vv.Convert(uint8Type).Uint())
			}
			doHexDump = true
		}
	}

	// Prepare indenting for slice.
	if nPeriod == 0 {
		d.w.Write(openBraceBytes)
	} else {
		d.w.Write(openBraceNewlineBytes)
	}
	d.depth++
	defer func() {
		d.depth--
		if nPeriod != 0 {
			d.indent()
		}
		d.w.Write(closeBraceBytes)
	}()

	// Hexdump the entire slice as needed.
	if doHexDump {
		indent := strings.Repeat(d.cs.Indent, d.depth)
		hexDump(d.w, buf, indent, d.cs.BytesWidth, d.cs.CommentBytes)
		return
	}

	// Recursively call dump for each item.
	for i := 0; i < numEntries; i++ {
		vi := v.Index(i)
		if nPeriod == 0 || i%nPeriod != 0 {
			d.ignoreNextIndent = true
		}
		d.dump(d.unpackValue(vi))
		if nPeriod == 0 || (i%nPeriod != nPeriod-1 && i != numEntries-1) {
			if i < numEntries-1 {
				d.w.Write(commaSpaceBytes)
				continue
			}
			break
		}
		d.w.Write(commaNewlineBytes)
	}
}

// isNumeric returns true for all numeric and boolean kinds.
func isNumeric(k reflect.Kind) bool {
	switch k {
	case reflect.Int, reflect.Uint,
		reflect.Int8, reflect.Bool,
		reflect.Int16, reflect.Uint16,
		reflect.Int32, reflect.Uint32,
		reflect.Int64, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Complex64, reflect.Complex128,
		reflect.Uintptr, reflect.UnsafePointer:
		return true
	default:
		return false
	}
}

// dump is the main workhorse for dumping a value.  It uses the passed reflect
// value to figure out what kind of object we are dealing with and formats it
// appropriately.  It is a recursive function, however circular data structures
// are detected and annotated.
func (d *dumpState) dump(v reflect.Value, wasPtr, static bool, addr uintptr) {
	// Handle invalid reflect values immediately.
	kind := v.Kind()
	if kind == reflect.Invalid {
		d.w.Write(invalidAngleBytes)
		return
	}

	// Handle pointers specially.
	if kind == reflect.Ptr {
		d.indent()
		d.dumpPtr(v)
		return
	}

	typ := v.Type()
	wantType := true
	if d.cs.ElideType {
		defType := !wasPtr && isDefault(typ)
		wantType = (!(static || defType) || isCompound(kind)) && !(kind == reflect.Interface && v.IsNil())
	}

	// Print type information unless already handled elsewhere.
	if !d.ignoreNextType {
		d.indent()
		if wantType {
			bufferedChan := v.Kind() == reflect.Chan && v.Cap() != 0
			if bufferedChan {
				d.w.Write(openParenBytes)
			}
			typeBytes := []byte(v.Type().String())
			d.w.Write(bytes.Replace(typeBytes, interfaceTypeBytes, interfaceBytes, -1))
			if bufferedChan {
				switch len := v.Len(); len {
				case 0:
					fmt.Fprintf(d.w, ", %d)", v.Cap())
				case 1:
					fmt.Fprintf(d.w, ", %d /* %d element */)", v.Cap(), len)
				default:
					fmt.Fprintf(d.w, ", %d /* %d elements */)", v.Cap(), len)
				}
			}
		}
	}
	d.ignoreNextType = false

	if wantType {
		switch kind {
		case reflect.Invalid, reflect.Struct, reflect.Slice, reflect.Array, reflect.Map:
		default:
			d.w.Write(openParenBytes)
		}
	}

	if _, referenced := d.nodes[addrType{addr, typ}]; !wasPtr && referenced {
		d.w.Write(openCommentBytes)
		printHexPtr(d.w, addr, true)
		d.w.Write(closeCommentBytes)
	}
	switch kind {
	case reflect.Invalid:
		// We should never get here since invalid has already been handled above.
		panic("cannot reach")

	case reflect.Bool:
		printBool(d.w, v.Bool())

	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		printInt(d.w, v.Int(), 10)

	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		d.w.Write(hexZeroBytes)
		printUint(d.w, v.Uint(), 16)

	case reflect.Float32:
		printFloat(d.w, v.Float(), 32, !wantType)

	case reflect.Float64:
		printFloat(d.w, v.Float(), 64, !wantType)

	case reflect.Complex64:
		printComplex(d.w, v.Complex(), 32)

	case reflect.Complex128:
		printComplex(d.w, v.Complex(), 64)

	case reflect.Slice:
		if v.IsNil() {
			d.w.Write(openParenBytes)
			d.w.Write(nilBytes)
			d.w.Write(closeParenBytes)
			break
		}
		// Remove pointers at or below the current depth from map used to detect
		// circular refs.
		for k, depth := range d.pointers {
			if depth >= d.depth {
				delete(d.pointers, k)
			}
		}
		if pd, ok := d.pointers[addr]; ok && pd < d.depth {
			d.w.Write(circularBytes)
			break
		}
		d.pointers[addr] = d.depth

		fallthrough

	case reflect.Array:
		d.dumpSlice(v)

	case reflect.String:
		d.w.Write([]byte(strconv.Quote(v.String())))

	case reflect.Interface:
		// The only time we should get here is for nil interfaces due to
		// unpackValue calls.
		if v.IsNil() {
			d.w.Write(nilBytes)
		}

	case reflect.Ptr:
		// We should never get here since pointers have already been handled above.
		panic("cannot reach")

	case reflect.Map:
		// nil maps should be indicated as different than empty maps
		if v.IsNil() {
			d.w.Write(openParenBytes)
			d.w.Write(nilBytes)
			d.w.Write(closeParenBytes)
			break
		}

		// Remove pointers at or below the current depth from map used to detect
		// circular refs.
		for k, depth := range d.pointers {
			if depth >= d.depth {
				delete(d.pointers, k)
			}
		}
		addr := v.Pointer()
		if pd, ok := d.pointers[addr]; ok && pd < d.depth {
			d.w.Write(circularBytes)
			break
		}
		d.pointers[addr] = d.depth

		d.w.Write(openBraceNewlineBytes)
		d.depth++
		keys := v.MapKeys()
		if d.cs.SortKeys {
			sortValues(keys)
		}
		for _, key := range keys {
			d.dump(d.unpackValue(key))
			d.w.Write(colonSpaceBytes)
			d.ignoreNextIndent = true
			d.dump(d.unpackValue(v.MapIndex(key)))
			d.w.Write(commaNewlineBytes)
		}
		d.depth--
		d.indent()
		d.w.Write(closeBraceBytes)

	case reflect.Struct:
		d.w.Write(openBraceNewlineBytes)
		d.depth++
		vt := v.Type()
		numFields := v.NumField()
		for i := 0; i < numFields; i++ {
			vtf := vt.Field(i)
			if d.cs.IgnoreUnexported && vtf.PkgPath != "" {
				continue
			}
			d.indent()
			d.w.Write([]byte(vtf.Name))
			d.w.Write(colonSpaceBytes)
			d.ignoreNextIndent = true
			d.dump(d.unpackValue(v.Field(i)))
			d.w.Write(commaNewlineBytes)
		}
		d.depth--
		d.indent()
		d.w.Write(closeBraceBytes)

	case reflect.Uintptr:
		printHexPtr(d.w, uintptr(v.Uint()), false)

	case reflect.UnsafePointer, reflect.Chan, reflect.Func:
		printHexPtr(d.w, v.Pointer(), true)

	// There were not any other types at the time this code was written, but
	// fall back to letting the default fmt package handle it in case any new
	// types are added.
	default:
		if v.CanInterface() {
			fmt.Fprintf(d.w, "%v", v.Interface())
		} else {
			fmt.Fprintf(d.w, "%v", v.String())
		}
	}
	if wantType {
		switch kind {
		case reflect.Invalid, reflect.Struct, reflect.Slice, reflect.Array, reflect.Map:
		default:
			d.w.Write(closeParenBytes)
		}
	}
}

// isDefault returns whether the type is a default type absent of context.
func isDefault(typ reflect.Type) bool {
	if typ.PkgPath() != "" || typ.Name() == "" {
		return false
	}
	kind := typ.Kind()
	return kind == reflect.Int || kind == reflect.Float64 || kind == reflect.String || kind == reflect.Bool
}

// isCompound returns whether the kind is a compound data type.
func isCompound(kind reflect.Kind) bool {
	return kind == reflect.Struct || kind == reflect.Slice || kind == reflect.Array || kind == reflect.Map
}

// fdump is a helper function to consolidate the logic from the various public
// methods which take varying writers and config states.
func fdump(cs *ConfigState, w io.Writer, a interface{}) {
	if a == nil {
		w.Write(interfaceBytes)
		w.Write(openParenBytes)
		w.Write(nilBytes)
		w.Write(closeParenBytes)
		w.Write(newlineBytes)
		return
	}

	d := dumpState{w: w, cs: cs}
	d.pointers = make(map[uintptr]int)
	v := reflect.ValueOf(a)
	var addr uintptr
	if v.CanAddr() {
		addr = v.Addr().Pointer()
	}
	if cs.CommentPointers {
		d.nodes = make(map[addrType]struct{})
		d.walk(v, false, false, addr)
	}
	d.dump(v, false, false, addr)
	d.w.Write(newlineBytes)
}

// Fdump formats and displays the passed arguments to io.Writer w.  It formats
// exactly the same as Dump.
func Fdump(w io.Writer, a interface{}) {
	fdump(&Config, w, a)
}

// Sdump returns a string with the passed arguments formatted exactly the same
// as Dump.
func Sdump(a interface{}) string {
	var buf bytes.Buffer
	fdump(&Config, &buf, a)
	return buf.String()
}

/*
Dump displays the passed parameters to standard out with newlines, customizable
indentation, and additional debug information such as complete types and all
pointer addresses used to indirect to the final value.  It provides the
following features over the built-in printing facilities provided by the fmt
package:

	* Pointers are dereferenced and followed
	* Circular data structures are detected and annotated
	* Byte arrays and slices are dumped in a way similar to the hexdump -C command,
	  which includes byte values in hex, and ASCII output

The configuration options are controlled by an exported package global,
utter.Config.  See ConfigState for options documentation.

See Fdump if you would prefer dumping to an arbitrary io.Writer or Sdump to
get the formatted result as a string.
*/
func Dump(a interface{}) {
	fdump(&Config, os.Stdout, a)
}
