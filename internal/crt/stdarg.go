// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
)

func VAuintptr(ap *[]interface{}) (r uintptr) {
	s := *ap
	switch x := s[0].(type) {
	case uintptr:
		r = x
	default:
		panic(fmt.Errorf("crt.VAuintptr %T", x))
	}
	*ap = s[1:]
	return r
}

func VAfloat64(ap *[]interface{}) float64 {
	s := *ap
	v := s[0].(float64)
	*ap = s[1:]
	return v
}

func VAint32(ap *[]interface{}) (v int32) {
	s := *ap
	switch x := s[0].(type) {
	case int32:
		v = x
	case uint32:
		v = int32(x)
	case int64:
		v = int32(x)
	case uint64:
		v = int32(x)
	case uintptr:
		v = int32(x)
	default:
		panic(fmt.Errorf("crt.VAint32 %T", x))
	}
	*ap = s[1:]
	return v
}

func VAuint32(ap *[]interface{}) (v uint32) {
	s := *ap
	switch x := s[0].(type) {
	case int32:
		v = uint32(x)
	case uint32:
		v = x
	case int64:
		v = uint32(x)
	case uint64:
		v = uint32(x)
	case uintptr:
		v = uint32(x)
	default:
		panic(fmt.Errorf("crt.VAuint32 %T", x))
	}
	*ap = s[1:]
	return v
}

func VAint64(ap *[]interface{}) (v int64) {
	s := *ap
	switch x := s[0].(type) {
	case int32:
		v = int64(x)
	case uint32:
		v = int64(x)
	case int64:
		v = x
	case uint64:
		v = int64(x)
	case uintptr:
		v = int64(x)
	default:
		panic(fmt.Errorf("crt.VAint64 %T", x))
	}
	*ap = s[1:]
	return v
}

func VAuint64(ap *[]interface{}) (v uint64) {
	s := *ap
	switch x := s[0].(type) {
	case int32:
		v = uint64(x)
	case uint32:
		v = uint64(x)
	case int64:
		v = uint64(x)
	case uint64:
		v = x
	case uintptr:
		v = uint64(x)
	default:
		panic(fmt.Errorf("crt.VAuint64 %T", x))
	}
	*ap = s[1:]
	return v
}

func vaLong(ap *[]interface{}) int64   { return VAint64(ap) }
func vaULong(ap *[]interface{}) uint64 { return VAuint64(ap) }
