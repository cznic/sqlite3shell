// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crt

import (
	"fmt"
	"os"
	"time"
	"unsafe"
)

type Stm struct {
	Xtm_sec      int32
	Xtm_min      int32
	Xtm_hour     int32
	Xtm_mday     int32
	Xtm_mon      int32
	Xtm_year     int32
	Xtm_wday     int32
	Xtm_yday     int32
	Xtm_isdst    int32
	X__tm_gmtoff int64
	X__tm_zone   uintptr // *int8
}

var localtime = MustCalloc(int(unsafe.Sizeof(Stm{})))

// struct tm *localtime(const time_t *timep);
func Xlocaltime(tls *TLS, timep uintptr) uintptr { return Xlocaltime_r(tls, timep, localtime) }

// time_t time(time_t *tloc);
func Xtime(tls *TLS, tloc uintptr) int64 {
	return time.Now().Unix()
}

// struct tm *localtime_r(const time_t *timep, struct tm *result);
func Xlocaltime_r(tls *TLS, timep, tm uintptr) uintptr {
	ut := *(*int64)(unsafe.Pointer(timep))
	t := time.Unix(ut, 0)
	p := (*Stm)(unsafe.Pointer(tm))
	p.Xtm_sec = int32(t.Second())
	p.Xtm_min = int32(t.Minute())
	p.Xtm_hour = int32(t.Hour())
	p.Xtm_mday = int32(t.Day())
	p.Xtm_mon = int32(t.Month())
	p.Xtm_year = int32(t.Year())
	p.Xtm_wday = int32(t.Weekday())
	p.Xtm_yday = int32(t.YearDay())
	p.Xtm_isdst = -1 //TODO
	if strace {
		fmt.Fprintf(os.Stderr, "localtime_r(%v, %#x) %+v\n", ut, tm, p)
	}
	return tm
}
