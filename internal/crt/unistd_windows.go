package crt

// int access(const char *path, int amode);
func X_access(tls *TLS, path uintptr, amode int32) int32 {
	panic("TODO: NOT IMPLEMENTED")
}

func X_setmode(tls *TLS, fd uintptr, amode int32) int32 {
	panic("TODO: NOT IMPLEMENTED")
}

func X_fileno(tls *TLS, fd uintptr) uintptr {
	panic("TODO: NOT IMPLEMENTED")
}

func Xfputs(tls *TLS, fd uintptr, data uintptr) uintptr {
	panic("TODO: NOT IMPLEMENTED")
}

func Xisdigit(tls *TLS, chr int32) uintptr {
	panic("TODO: NOT IMPLEMENTED")
}

func Xisspace(tls *TLS, chr int32) uintptr {
	panic("TODO: NOT IMPLEMENTED")
}

func Xisalpha(tls *TLS, chr int32) uintptr {
	panic("TODO: NOT IMPLEMENTED")
}

func Xisalnum(tls *TLS, chr int32) uintptr {
	panic("TODO: NOT IMPLEMENTED")
}
