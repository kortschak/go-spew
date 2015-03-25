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

import "reflect"

// walkPtr handles walking of pointers by indirecting them as necessary.
func (d *dumpState) walkPtr(v reflect.Value) {
	// Remove pointers at or below the current depth from map used to detect
	// circular refs.
	for k, depth := range d.pointers {
		if depth >= d.depth {
			delete(d.pointers, k)
		}
	}

	var nilFound, cycleFound bool
	ve := v
	for ve.Kind() == reflect.Ptr {
		if ve.IsNil() {
			nilFound = true
			break
		}
		addr := ve.Pointer()
		if pd, ok := d.pointers[addr]; ok && pd < d.depth {
			cycleFound = true
			break
		}
		d.pointers[addr] = d.depth
		d.nodes[addrType{addr, ve.Type()}] = struct{}{}

		ve = ve.Elem()
		if ve.Kind() == reflect.Interface {
			if ve.IsNil() {
				nilFound = true
				break
			}
			ve = ve.Elem()
		}
		d.nodes[addrType{addr, ve.Type()}] = struct{}{}
	}

	if !nilFound && !cycleFound {
		d.walk(ve, true, false, 0)
	}
}

// walkSlice handles walking of arrays and slices.
func (d *dumpState) walkSlice(v reflect.Value) {
	d.depth++
	// Recursively call walk for each item.
	for i := 0; i < v.Len(); i++ {
		d.walk(d.unpackValue(v.Index(i)))
	}
	d.depth--
}

// walk is the main workhorse for walking a value.  It uses the passed reflect
// value to figure out what kind of object we are dealing with and follows it
// appropriately.  It is a recursive function, however circular data structures
// are detected and escaped from.
func (d *dumpState) walk(v reflect.Value, wasPtr, static bool, addr uintptr) {
	// Handle invalid reflect values immediately.
	kind := v.Kind()
	if kind == reflect.Invalid {
		return
	}

	// Handle pointers specially.
	if kind == reflect.Ptr {
		d.walkPtr(v)
		return
	}

	switch kind {
	case reflect.Slice:
		if v.IsNil() {
			break
		}
		fallthrough

	case reflect.Array:
		d.walkSlice(v)

	case reflect.Map:
		if v.IsNil() {
			break
		}
		d.depth++
		keys := v.MapKeys()
		for _, key := range keys {
			d.walk(d.unpackValue(key))
			d.walk(d.unpackValue(v.MapIndex(key)))
		}
		d.depth--

	case reflect.Struct:
		d.depth++
		vt := v.Type()
		for i := 0; i < v.NumField(); i++ {
			vtf := vt.Field(i)
			if d.cs.IgnoreUnexported && vtf.PkgPath != "" {
				continue
			}
			d.walk(d.unpackValue(v.Field(i)))
		}
		d.depth--
	}
}
