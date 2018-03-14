// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"runtime"
	"strings"
)

type CallbackFunc func(wr io.Writer, tyMap map[string]Type, comment, ret, name, rawArgs string)

var compiledFuncs []string = []string{}

// var farProcGoTy = "func(*TLS) int64"
var size_tGoTy = "uint64"

func fileHeader(wr io.Writer, tag string, imports []string) error {
	_, err := fmt.Fprintf(wr, `// Copyright 2017 The CRT Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// Code generated by running "go generate". DO NOT EDIT.

// +build `+tag+`

package crt

import (
	"%s"
)

`, strings.Join(imports, "\"\n\t\""))
	return err
}

type Type int

const (
	TyPtr Type = iota
	TyStr
	// unicode string
	TyUstr
	TyInt32
	TyInt8
	TyUint16
	TyUint32
	TyU32Ptr
	TyI32Ptr
	TyISlicePtr
	TyI8DoublePtr
	TyU16DoublePtr
	TyHMODULEPtr
	TyFILETIMEPtr
	TyOSVERSIONINFOAPtr
	TyOSVERSIONINFOWPtr
	TyOVERLAPPEDPtr
	TySECURITY_ATTRIBUTESPtr
	TySYSTEM_INFOPtr
	TySYSTEMTIMEPtr
	TyCRITICAL_SECTIONPtr
	TyXLARGE_INTEGERPtr
	TyFarProc
	TySIZE_T
	TyGET_FILEEX_INFO_LEVELS

	TyVoid
	TyError
)

func (t Type) Syscall(val string) string {
	switch t {
	case TyPtr:
		return fmt.Sprintf("uintptr(%s)", val)
	case TyInt32, TyInt8, TyUint16, TyUint32, TySIZE_T:
		return fmt.Sprintf("uintptr(%s)", val)
	case TyStr, TyUstr, TyISlicePtr, TyU32Ptr, TyI32Ptr, TyI8DoublePtr, TyU16DoublePtr,
		TyHMODULEPtr, TyFILETIMEPtr, TyOSVERSIONINFOAPtr, TyOSVERSIONINFOWPtr, TyOVERLAPPEDPtr,
		TySECURITY_ATTRIBUTESPtr, TySYSTEM_INFOPtr, TySYSTEMTIMEPtr, TyXLARGE_INTEGERPtr,
		TyCRITICAL_SECTIONPtr, TyFarProc, TyGET_FILEEX_INFO_LEVELS:
		return fmt.Sprintf("uintptr(%s)", val)
	default:
		log.Fatal("Cannot syscall type: ", t)
		return ""
	}
}

func (t Type) Write(target, val string) string {
	switch t {
	case TyPtr,
		TyStr, TyUstr, TyISlicePtr, TyU32Ptr, TyI32Ptr, TyI8DoublePtr, TyU16DoublePtr,
		TyHMODULEPtr, TyFILETIMEPtr, TyOVERLAPPEDPtr, TyOSVERSIONINFOAPtr, TyOSVERSIONINFOWPtr,
		TySECURITY_ATTRIBUTESPtr, TySYSTEM_INFOPtr, TySYSTEMTIMEPtr, TyXLARGE_INTEGERPtr,
		TyCRITICAL_SECTIONPtr, TyFarProc:
		return fmt.Sprintf("*%s = %s", target, val)
	case TyInt32:
		return fmt.Sprintf("%s = int32(%s)", target, val)
	case TyInt8, TyGET_FILEEX_INFO_LEVELS:
		return fmt.Sprintf("%s = int8(%s)", target, val)
	case TyUint16:
		return fmt.Sprintf("%s = uint16(%s)", target, val)
	case TyUint32:
		return fmt.Sprintf("%s = uint32(%s)", target, val)
	case TySIZE_T:
		return fmt.Sprintf("%s = %s(%s)", size_tGoTy, target, val)
	case TyVoid:
		// void is usually used for a function without a return value
		// so the write is a NOP
		return ""
	default:
		log.Fatal("Cannot handle write type: ", t)
	}
	return ""
}

func (t Type) GoType() string {
	ty := ""
	switch t {
	case TyPtr, TyStr, TyUstr, TyI32Ptr, TyU32Ptr, TyI8DoublePtr, TyU16DoublePtr, TyHMODULEPtr, TyFILETIMEPtr,
		TyOSVERSIONINFOAPtr, TyISlicePtr, TyOSVERSIONINFOWPtr, TyOVERLAPPEDPtr, TySECURITY_ATTRIBUTESPtr,
		TySYSTEM_INFOPtr, TySYSTEMTIMEPtr, TyXLARGE_INTEGERPtr, TyCRITICAL_SECTIONPtr, TyFarProc:
		ty = "uintptr"
	case TyInt32:
		ty = "int32"
	case TyInt8:
		ty = "int8"
	case TyGET_FILEEX_INFO_LEVELS:
		ty = "E_GET_FILEEX_INFO_LEVELS"
	case TyUint16:
		ty = "uint16"
	case TyUint32:
		ty = "uint32"
	case TyVoid:
		ty = ""
	case TySIZE_T:
		ty = size_tGoTy
	//case TyFarProc:
	//ty = farProcGoTy
	default:
		log.Fatal("Cannot get GoType: ", t)
	}
	return ty
}

func (t Type) FormatStr(val interface{}) string {
	arg := "%s"
	switch t {
	case TyPtr, TyISlicePtr, TyU32Ptr, TyI32Ptr, TyI8DoublePtr, TyU16DoublePtr, TyHMODULEPtr, TyFILETIMEPtr,
		TyOSVERSIONINFOAPtr, TyOSVERSIONINFOWPtr, TyOVERLAPPEDPtr, TySECURITY_ATTRIBUTESPtr, TySYSTEM_INFOPtr,
		TySYSTEMTIMEPtr, TyXLARGE_INTEGERPtr, TyCRITICAL_SECTIONPtr, TyFarProc:
		arg = "%#x"
	case TyStr, TyUstr:
		arg = "%s"
	case TyInt8, TyGET_FILEEX_INFO_LEVELS, TyUint16, TyInt32, TyUint32, TySIZE_T:
		arg = "%#x"
	case TyError:
		arg = "%v"
	case TyVoid:
		arg = "%d"
	default:
		log.Fatal("Cannot format type: ", t)
	}
	return arg
}

func (t Type) Format(name string) string {
	if t == TyStr {
		return fmt.Sprintf("GoUTF16String(uintptr(%s))", name)
	}
	return name
}

func compileWinSyscall(wr io.Writer, tyMap map[string]Type, comment, ret, name, rawArgs string) {
	if _, err := fmt.Fprintf(wr, "// %s \nfunc X%s(tls *TLS", comment, name); err != nil {
		log.Fatal("cannot write function header: ", err)
	}

	// handle the arguments in the correct order (reverse, so we get the right elements of the stack)
	args := strings.Split(rawArgs, ",")
	if len(rawArgs) == 0 {
		args = []string{}
	}
	syscallArgs := make([]string, len(args))
	formatStrs := make([]string, len(args))
	printArgs := make([]string, len(args))
	for i := 0; i < len(args); i++ {
		arg := args[i]
		a := strings.Split(strings.TrimSpace(arg), " ")
		if len(a) != 2 {
			log.Fatal("expected length 2: type and argument name for ", arg)
		}

		ty, exists := tyMap[a[0]]
		if !exists {
			log.Fatal("cannot map type: ", a[0], " in ", arg)
		}

		argName := a[1]
		syscallArgs[i] = ty.Syscall(argName)
		formatStrs[i] = ty.FormatStr(argName)
		printArgs[i] = ty.Format(argName)

		_, err := fmt.Fprintf(wr, ", %s %s", argName, ty.GoType())
		if err != nil {
			log.Fatal("cannot generate stack pop: ", err)
		}
	}

	retTy, exists := tyMap[ret]
	if !exists {
		log.Fatal("cannot map return type: ", retTy)
	}

	if _, err := fmt.Fprintf(wr, ") %s {\n", retTy.GoType()); err != nil {
		log.Fatal("cannot finish func declaration: ", err)
	}

	fn := ""
	// fill with null for the respective syscall
	fillArgs := 0
	switch {
	case 0 <= len(args) && len(args) <= 3:
		fn = "Syscall"
		fillArgs = 3 - len(args)
	case 3 < len(args) && len(args) <= 6:
		fn = "Syscall6"
		fillArgs = 6 - len(args)
	case 6 < len(args) && len(args) <= 9:
		fn = "Syscall9"
		fillArgs = 9 - len(args)
	default:
		log.Fatal("Unsupported argument size: ", len(args))
	}

	for i := 0; i < fillArgs; i++ {
		syscallArgs = append(syscallArgs, "0")
	}
	argStr := strings.Join(syscallArgs, ", \n\t\t")
	if _, err := fmt.Fprintf(wr, "\n\tret, _, err := syscall.%s(proc%s.Addr(), %d, %s);\n", fn, name, len(args), argStr); err != nil {
		log.Fatal("could not write syscall: ", err)
	}

	printArgs = append(printArgs, retTy.Format("ret"))
	if _, err := fmt.Fprintf(wr, "\tif strace {\n\t\tfmt.Fprintf(os.Stderr, \"%s(%s) %s %s\\n\", %s, err)\n\t}\n\tif err != 0 {\n\t\ttls.setErrno(err)\n\t}\n",
		name, strings.Join(formatStrs, ", "), retTy.FormatStr(retTy), TyError.FormatStr("error"), strings.Join(printArgs, ", \n\t\t\t")); err != nil {
		log.Fatal("cannot generate strace: ", err)
	}

	if tyStr := retTy.GoType(); tyStr != "" {
		retVar := "ret"
		if _, err := fmt.Fprintf(wr, "\treturn (%s)(%s)\n", tyStr, retVar); err != nil {
			log.Fatal("cannot complete ret stmt: ", err)
		}
	}
	if _, err := fmt.Fprintf(wr, "}\n\n"); err != nil {
		log.Fatal("cannot end function block: ", err)
	}
}

func compileWinFile(wr io.Writer, arch string, callback CallbackFunc) {
	bytes, err := ioutil.ReadFile("windows.go")
	if err != nil {
		log.Fatal("Cannot read windows.go: ", err)
	}
	bytes2, err := ioutil.ReadFile("windows_" + arch + ".go")
	if err != nil {
		log.Fatal("Cannot read architecture specific config: ", arch)
	}
	bytes = append(bytes, bytes2...)
	content := string(bytes)

	reTy := regexp.MustCompile("//ty:(.*?): (.*)")
	reSys := regexp.MustCompile("//sys: (.*?) (.*?)\\((.*)\\);")

	// get type mappings for function signatures
	tyMatches := reTy.FindAllStringSubmatch(content, -1)
	tyMap := map[string]Type{}
	for _, match := range tyMatches {
		// the target type
		ty := match[1]
		// the aliases e.g. a list like `HANDLE, LPWXYZ`
		aliases := match[2]

		for _, alias := range strings.Split(aliases, ",") {
			alias = strings.TrimSpace(alias)
			switch ty {
			case "ptr":
				tyMap[alias] = TyPtr
			case "str":
				tyMap[alias] = TyStr
			case "ustr":
				tyMap[alias] = TyUstr
			case "int32":
				tyMap[alias] = TyInt32
			case "int8":
				tyMap[alias] = TyInt8
			case "uint16":
				tyMap[alias] = TyUint16
			case "uint32":
				tyMap[alias] = TyUint32
			case "void":
				tyMap[alias] = TyVoid
			case "isliceptr":
				tyMap[alias] = TyISlicePtr
			case "u32ptr":
				tyMap[alias] = TyU32Ptr
			case "i32ptr":
				tyMap[alias] = TyI32Ptr
			case "**i8":
				tyMap[alias] = TyI8DoublePtr
			case "**u16":
				tyMap[alias] = TyU16DoublePtr
			case "*HMODULE":
				tyMap[alias] = TyHMODULEPtr
			case "*FILETIME":
				tyMap[alias] = TyFILETIMEPtr
			case "*OSVERSIONINFOA":
				tyMap[alias] = TyOSVERSIONINFOAPtr
			case "*OSVERSIONINFOW":
				tyMap[alias] = TyOSVERSIONINFOWPtr
			case "*SECURITY_ATTRIBUTES":
				tyMap[alias] = TySECURITY_ATTRIBUTESPtr
			case "*SYSTEM_INFO":
				tyMap[alias] = TySYSTEM_INFOPtr
			case "*SYSTEMTIME":
				tyMap[alias] = TySYSTEMTIMEPtr
			case "*LARGE_INTEGER":
				tyMap[alias] = TyXLARGE_INTEGERPtr
			case "*OVERLAPPED":
				tyMap[alias] = TyOVERLAPPEDPtr
			case "*CRITICAL_SECTION":
				tyMap[alias] = TyCRITICAL_SECTIONPtr
			case "FARPROC":
				tyMap[alias] = TyFarProc
			case "size_t":
				tyMap[alias] = TySIZE_T
			case "GET_FILEEX_INFO_LEVELS":
				tyMap[alias] = TyGET_FILEEX_INFO_LEVELS
			default:
				log.Fatal("unknown target type: ", ty)
			}
		}
	}

	// compile syscalls
	sysMatches := reSys.FindAllStringSubmatch(content, -1)
	for _, match := range sysMatches {
		// the return type of the function
		ret := strings.TrimSpace(match[1])
		// the function name
		name := strings.TrimSpace(match[2])
		// the arguments
		rawArgs := strings.TrimSpace(match[3])
		compiledFuncs = append(compiledFuncs, name)
		callback(wr, tyMap, match[0], ret, name, rawArgs)
	}
}

func main() {
	var out bytes.Buffer
	var buf bytes.Buffer

	flag.Parse()

	for _, arch := range []string{"amd64", "386"} {
		compiledFuncs = []string{}
		buf.Reset()
		out.Reset()

		if arch == "386" {
			//farProcGoTy = "func(*TLS) int32"
			size_tGoTy = "uint32"
		}

		if err := fileHeader(&buf, "windows", []string{"fmt", "syscall", "os"}); err != nil {
			log.Fatal("Cannot write file header: ", err)
		}
		compileWinFile(&out, arch, compileWinSyscall)

		// procCreateFileW             = modkernel32.NewProc("CreateFileW")
		if _, err := fmt.Fprintf(&buf, "var (\n\tmodkernel32           = syscall.NewLazyDLL(\"kernel32.dll\")\n"); err != nil {
			log.Fatal("cannot begin external proc declaration: ", err)
		}

		for _, fn := range compiledFuncs {
			if _, err := fmt.Fprintf(&buf, "\tproc%-30s = modkernel32.NewProc(\"%s\")\n", fn, fn); err != nil {
				log.Fatal("cannot write sid mapping: ", err)
			}
		}

		if _, err := fmt.Fprintf(&buf, ")\n\n"); err != nil {
			log.Fatal("cannot terminate external proc declaration: ", err)
		}

		// merge headers & compiled code
		if _, err := buf.Write(out.Bytes()); err != nil {
			log.Fatal("cannot merge outputs: ", err)
		}

		if err := ioutil.WriteFile(fmt.Sprintf("windows_impl_%s.go", arch), buf.Bytes(), 0655); err != nil {
			log.Fatal("cannot write windows_impl.go: ", err)
		}
	}
}

func goArch() string {
	if s := os.Getenv("GOARCH"); s != "" {
		return s
	}

	return runtime.GOARCH
}

func goOs() string {
	if s := os.Getenv("GOOS"); s != "" {
		return s
	}

	return runtime.GOOS
}
