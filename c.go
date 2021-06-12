package main

import "C"

import (
	"unsafe"

	"github.com/xmdhs/gml-c-shared/c"
)

//export Getchar
func Getchar(charlist **C.char, index C.longlong) *C.char {
	c := c.Cchar{P: unsafe.Pointer(charlist)}
	char := c.Getchar(int(index))
	return (*C.char)(char)
}

//export Freechar
func Freechar(charlist **C.char, len C.longlong) {
	c := c.Cchar{P: unsafe.Pointer(charlist)}
	c.Freechar(int(len))
}

//export NewChar
func NewChar(len C.longlong) **C.char {
	c := c.NewChar(int(len))
	return (**C.char)(c.P)
}

//export SetChar
func SetChar(cc **C.char, index C.longlong, achar *C.char) {
	c := c.Cchar{P: unsafe.Pointer(cc)}
	c.SetChar(int(index), unsafe.Pointer(achar))
}

//export Malloc
func Malloc(i C.int) unsafe.Pointer {
	return c.Malloc(int(i))
}
