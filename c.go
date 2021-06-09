package main

import "C"

import (
	"unsafe"

	"github.com/xmdhs/gml-c-shared/c"
)

//export getchar
func getchar(charlist **C.char, index C.longlong) *C.char {
	c := c.Cchar{P: unsafe.Pointer(charlist)}
	char := c.Getchar(int(index))
	return (*C.char)(char)
}

//export freechar
func freechar(charlist **C.char, len C.longlong) {
	c := c.Cchar{P: unsafe.Pointer(charlist)}
	c.Freechar(int(len))
}
