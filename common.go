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
	"math"
	"reflect"
	"sort"
	"strconv"
	"unsafe"
)

const (
	// ptrSize is the size of a pointer on the current arch.
	ptrSize = unsafe.Sizeof((*byte)(nil))
)

var (
	// offsetPtr, offsetScalar, and offsetFlag are the offsets for the
	// internal reflect.Value fields.  These values are valid before golang
	// commit ecccf07e7f9d which changed the format.  The are also valid
	// after commit 82f48826c6c7 which changed the format again to mirror
	// the original format.  Code in the init function updates these offsets
	// as necessary.
	offsetPtr    = uintptr(ptrSize)
	offsetScalar = uintptr(0)
	offsetFlag   = uintptr(ptrSize * 2)

	// flagKindWidth and flagKindShift indicate various bits that the
	// reflect package uses internally to track kind information.
	//
	// flagRO indicates whether or not the value field of a reflect.Value is
	// read-only.
	//
	// flagIndir indicates whether the value field of a reflect.Value is
	// the actual data or a pointer to the data.
	//
	// These values are valid before golang commit 90a7c3c86944 which
	// changed their positions.  Code in the init function updates these
	// flags as necessary.
	flagKindWidth = uintptr(5)
	flagKindShift = uintptr(flagKindWidth - 1)
	flagRO        = uintptr(1 << 0)
	flagIndir     = uintptr(1 << 1)
)

func init() {
	// Older versions of reflect.Value stored small integers directly in the
	// ptr field (which is named val in the older versions).  Versions
	// between commits ecccf07e7f9d and 82f48826c6c7 added a new field named
	// scalar for this purpose which unfortunately came before the flag
	// field, so the offset of the flag field is different for those
	// versions.
	//
	// This code constructs a new reflect.Value from a known small integer
	// and checks if the size of the reflect.Value struct indicates it has
	// the scalar field. When it does, the offsets are updated accordingly.
	vv := reflect.ValueOf(0xf00)
	if unsafe.Sizeof(vv) == (ptrSize * 4) {
		offsetScalar = ptrSize * 2
		offsetFlag = ptrSize * 3
	}

	// Commit 90a7c3c86944 changed the flag positions such that the low
	// order bits are the kind.  This code extracts the kind from the flags
	// field and ensures it's the correct type.  When it's not, the flag
	// order has been changed to the newer format, so the flags are updated
	// accordingly.
	upf := unsafe.Pointer(uintptr(unsafe.Pointer(&vv)) + offsetFlag)
	upfv := *(*uintptr)(upf)
	flagKindMask := uintptr((1<<flagKindWidth - 1) << flagKindShift)
	if (upfv&flagKindMask)>>flagKindShift != uintptr(reflect.Int) {
		flagKindShift = 0
		flagRO = 1 << 5
		flagIndir = 1 << 6

		// Commit adf9b30e5594 modified the flags to separate the
		// flagRO flag into two bits which specifies whether or not the
		// field is embedded.  This causes flagIndir to move over a bit
		// and means that flagRO is the combination of either of the
		// original flagRO bit and the new bit.
		//
		// This code detects the change by extracting what used to be
		// the indirect bit to ensure it's set.  When it's not, the flag
		// order has been changed to the newer format, so the flags are
		// updated accordingly.
		if upfv&flagIndir == 0 {
			flagRO = 3 << 5
			flagIndir = 1 << 7
		}
	}
}

// unsafeReflectValue converts the passed reflect.Value into a one that bypasses
// the typical safety restrictions preventing access to unaddressable and
// unexported data.  It works by digging the raw pointer to the underlying
// value out of the protected value and generating a new unprotected (unsafe)
// reflect.Value to it.
//
// This allows us to check for implementations of the Stringer and error
// interfaces to be used for pretty printing ordinarily unaddressable and
// inaccessible values such as unexported struct fields.
func unsafeReflectValue(v reflect.Value) (rv reflect.Value) {
	indirects := 1
	vt := v.Type()
	upv := unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + offsetPtr)
	rvf := *(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&v)) + offsetFlag))
	if rvf&flagIndir != 0 {
		vt = reflect.PtrTo(v.Type())
		indirects++
	} else if offsetScalar != 0 {
		// The value is in the scalar field when it's not one of the
		// reference types.
		switch vt.Kind() {
		case reflect.Uintptr:
		case reflect.Chan:
		case reflect.Func:
		case reflect.Map:
		case reflect.Ptr:
		case reflect.UnsafePointer:
		default:
			upv = unsafe.Pointer(uintptr(unsafe.Pointer(&v)) +
				offsetScalar)
		}
	}

	pv := reflect.NewAt(vt, upv)
	rv = pv
	for i := 0; i < indirects; i++ {
		rv = rv.Elem()
	}
	return rv
}

// Some constants in the form of bytes to avoid string overhead.  This mirrors
// the technique used in the fmt package.
var (
	plusBytes             = []byte("+")
	iBytes                = []byte("i")
	trueBytes             = []byte("true")
	falseBytes            = []byte("false")
	interfaceBytes        = []byte("interface{}")
	interfaceTypeBytes    = []byte("interface {}")
	commaSpaceBytes       = []byte(", ")
	commaNewlineBytes     = []byte(",\n")
	newlineBytes          = []byte("\n")
	openBraceBytes        = []byte("{")
	openBraceNewlineBytes = []byte("{\n")
	closeBraceBytes       = []byte("}")
	ampersandBytes        = []byte("&")
	colonSpaceBytes       = []byte(": ")
	spaceBytes            = []byte(" ")
	openParenBytes        = []byte("(")
	closeParenBytes       = []byte(")")
	nilBytes              = []byte("nil")
	hexZeroBytes          = []byte("0x")
	zeroBytes             = []byte("0")
	pointZeroBytes        = []byte(".0")
	openCommentBytes      = []byte(" /*")
	closeCommentBytes     = []byte("*/ ")
	pointerChainBytes     = []byte("->")
	circularBytes         = []byte("(<already shown>)")
	invalidAngleBytes     = []byte("<invalid>")
)

// hexDigits is used to map a decimal value to a hex digit.
var hexDigits = "0123456789abcdef"

// printBool outputs a boolean value as true or false to Writer w.
func printBool(w io.Writer, val bool) {
	if val {
		w.Write(trueBytes)
	} else {
		w.Write(falseBytes)
	}
}

// printInt outputs a signed integer value to Writer w.
func printInt(w io.Writer, val int64, base int) {
	w.Write([]byte(strconv.FormatInt(val, base)))
}

// printUint outputs an unsigned integer value to Writer w.
func printUint(w io.Writer, val uint64, base int) {
	w.Write([]byte(strconv.FormatUint(val, base)))
}

// printFloat outputs a floating point value using the specified precision,
// which is expected to be 32 or 64bit, to Writer w.
func printFloat(w io.Writer, val float64, precision int, typeElided bool) {
	w.Write([]byte(strconv.FormatFloat(val, 'g', -1, precision)))
	if typeElided && !math.IsInf(val, 0) && val == math.Floor(val) {
		w.Write(pointZeroBytes)
	}
}

// printComplex outputs a complex value using the specified float precision
// for the real and imaginary parts to Writer w.
func printComplex(w io.Writer, c complex128, floatPrecision int) {
	r := real(c)
	w.Write([]byte(strconv.FormatFloat(r, 'g', -1, floatPrecision)))
	i := imag(c)
	if i >= 0 {
		w.Write(plusBytes)
	}
	w.Write([]byte(strconv.FormatFloat(i, 'g', -1, floatPrecision)))
	w.Write(iBytes)
}

// hexDump is a modified 'hexdump -C'-like that returns a commented Go syntax
// byte slice or array.
func hexDump(w io.Writer, data []byte, indent string, width int, comment bool) {
	var commentBytes []byte
	if comment {
		commentBytes = make([]byte, width)
	}

	for i, v := range data {
		if i%width == 0 {
			fmt.Fprint(w, indent)
		} else {
			w.Write(spaceBytes)
		}

		fmt.Fprintf(w, "%#02x,", v)
		if comment {
			if v < 32 || v > 126 {
				v = '.'
			}
			commentBytes[i%width] = v
		}

		if !comment {
			if i%width == width-1 || i == len(data)-1 {
				fmt.Fprintln(w)
			}
			continue
		}
		if i%width == width-1 {
			fmt.Fprintf(w, " // |%s|\n", commentBytes[:])
		} else if i == len(data)-1 {
			if len(data) > width {
				slots := width - i%width - 1
				switch slots {
				case 0:
					// Do nothing.
				case 1:
					w.Write([]byte(" /* */"))
				default:
					w.Write([]byte(" /*   "))
					w.Write(bytes.Repeat([]byte("      "), slots-2))
					w.Write([]byte("    */"))
				}
			}
			fmt.Fprintf(w, " // |%s|\n", commentBytes[:len(data)%width])
		}
	}
}

// printHexPtr outputs a uintptr formatted as hexadecimal with a leading '0x'
// prefix to Writer w.
func printHexPtr(w io.Writer, p uintptr, isPointer bool) {
	// Null pointer.
	num := uint64(p)
	if num == 0 {
		if isPointer {
			w.Write(nilBytes)
		} else {
			w.Write(zeroBytes)
		}
		return
	}

	// Max uint64 is 16 bytes in hex + 2 bytes for '0x' prefix
	buf := make([]byte, 18)

	// It's simpler to construct the hex string right to left.
	base := uint64(16)
	i := len(buf) - 1
	for num >= base {
		buf[i] = hexDigits[num%base]
		num /= base
		i--
	}
	buf[i] = hexDigits[num]

	// Add '0x' prefix.
	i--
	buf[i] = 'x'
	i--
	buf[i] = '0'

	// Strip unused leading bytes.
	buf = buf[i:]
	w.Write(buf)
}

// valuesSorter implements sort.Interface to allow a slice of reflect.Value
// elements to be sorted.
type valuesSorter struct {
	values []reflect.Value
}

// Len returns the number of values in the slice.  It is part of the
// sort.Interface implementation.
func (s *valuesSorter) Len() int {
	return len(s.values)
}

// Swap swaps the values at the passed indices.  It is part of the
// sort.Interface implementation.
func (s *valuesSorter) Swap(i, j int) {
	s.values[i], s.values[j] = s.values[j], s.values[i]
}

// valueSortLess returns whether the first value should sort before the second
// value.  It is used by valueSorter.Less as part of the sort.Interface
// implementation.
func valueSortLess(a, b reflect.Value) bool {
	switch a.Kind() {
	case reflect.Bool:
		return !a.Bool() && b.Bool()
	case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
		return a.Int() < b.Int()
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return a.Uint() < b.Uint()
	case reflect.Float32, reflect.Float64:
		return a.Float() < b.Float()
	case reflect.String:
		return a.String() < b.String()
	case reflect.Uintptr:
		return a.Uint() < b.Uint()
	case reflect.Array:
		// Compare the contents of both arrays.
		l := a.Len()
		for i := 0; i < l; i++ {
			av := a.Index(i)
			bv := b.Index(i)
			if av.Interface() == bv.Interface() {
				continue
			}
			return valueSortLess(av, bv)
		}
	}
	return a.String() < b.String()
}

// Less returns whether the value at index i should sort before the
// value at index j.  It is part of the sort.Interface implementation.
func (s *valuesSorter) Less(i, j int) bool {
	return valueSortLess(s.values[i], s.values[j])
}

// sortValues is a generic sort function for native types: int, uint, bool,
// string and uintptr.  Other inputs are sorted according to their
// Value.String() value to ensure display stability.
func sortValues(values []reflect.Value) {
	if len(values) == 0 {
		return
	}
	sort.Sort(&valuesSorter{values})
}
