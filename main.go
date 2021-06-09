package main

/*
typedef struct
{
    char *Minecraftpath;
    long long RAM;
    char *Name;
    char *UUID;
    char *AccessToken;
    char *Gamedir;
    char *Version;
    char *ApiAddress;
} Gameinfo;

typedef struct
{
    char **args;
    long long len;
} GameinfoReturn;

typedef void (*Fail)(char*);

typedef void (*Ok)(int,int);
*/
import "C"

import (
	"errors"
	"os"
	"strconv"
	"unsafe"

	"github.com/xmdhs/gml-c-shared/gml"
	"github.com/xmdhs/gomclauncher/download"
	"github.com/xmdhs/gomclauncher/launcher"

	"github.com/xmdhs/gml-c-shared/c"
)

func main() {}

//export GenCmdArgs
func GenCmdArgs(g C.Gameinfo) (C.GameinfoReturn, int) {
	l := gml.Launcher{
		Gameinfo: launcher.Gameinfo{
			Minecraftpath: C.GoString(g.Minecraftpath),
			RAM:           strconv.Itoa(int(g.RAM)),
			Name:          C.GoString(g.Name),
			UUID:          C.GoString(g.UUID),
			AccessToken:   C.GoString(g.AccessToken),
			Version:       C.GoString(g.Version),
			ApiAddress:    C.GoString(g.ApiAddress),
		},
		Independent: false,
	}
	args, err := l.GenCmdArgs()
	if err != nil {
		i := errr(err)
		return C.GameinfoReturn{}, i
	}
	c := c.NewChar(len(args))

	for i, v := range args {
		c.SetChar(i, unsafe.Pointer(C.CString(v)))
	}

	var r C.GameinfoReturn
	r.args = (**C.char)(c.P)
	r.len = C.longlong(int64(len(args)))
	return r, 0
}

//export Download
func Download(version, Type, Minecraftpath *C.char, downInt C.int, fail C.Fail, ok C.Ok) int {
	d := gml.NewDown(C.GoString(Type), C.GoString(Minecraftpath), int(downInt), c.DoFail(unsafe.Pointer(fail)), c.DoOk(unsafe.Pointer(ok)))
	err := d.Download(C.GoString(version))
	if err != nil {
		return errr(err)
	}
	return 0
}

//export Check
func Check(version, Type, Minecraftpath *C.char, downInt C.int, fail C.Fail, ok C.Ok) int {
	d := gml.NewDown(C.GoString(Type), C.GoString(Minecraftpath), int(downInt), c.DoFail(unsafe.Pointer(fail)), c.DoOk(unsafe.Pointer(ok)))
	err := d.Check(C.GoString(version))
	if err != nil {
		return errr(err)
	}
	return 0

}

func errr(err error) int {
	switch {
	case errors.Is(err, os.ErrNotExist):
		return 1
	case errors.Is(err, launcher.JsonErr):
		return 2
	case errors.Is(err, launcher.JsonNorTrue):
		return 3
	case errors.Is(err, download.NoSuch):
		return 4
	case errors.Is(err, download.FileDownLoadFail):
		return 5
	default:
		return -1
	}
}
