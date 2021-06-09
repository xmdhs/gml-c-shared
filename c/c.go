package c

/*
#include "char.h"
*/
import "C"

import "unsafe"

type Cchar struct {
	P unsafe.Pointer
}

func NewChar(i int) Cchar {
	c := C.newchar(C.longlong(int64(i)))
	return Cchar{unsafe.Pointer(c)}
}

func (c Cchar) SetChar(index int, char unsafe.Pointer) {
	C.setchar((**C.char)(c.P), C.longlong(int64(index)), (*C.char)(char))
}

func (c Cchar) Getchar(index int) unsafe.Pointer {
	ch := C.getchar((**C.char)(c.P), C.longlong(int64(index)))
	return unsafe.Pointer(ch)
}

func (c Cchar) Freechar(len int) {
	C.freechar((**C.char)(c.P), C.longlong(int64(len)))
}
