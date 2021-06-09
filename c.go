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
