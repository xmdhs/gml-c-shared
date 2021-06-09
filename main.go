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

typedef struct
{
    int  code;
    char *msg;
} err;
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
func GenCmdArgs(g C.Gameinfo) (C.GameinfoReturn, C.err) {
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
	return r, errr(err)
}

//export Download
func Download(version, Type, Minecraftpath *C.char, downInt C.int, fail C.Fail, ok C.Ok) C.err {
	d := gml.NewDown(C.GoString(Type), C.GoString(Minecraftpath), int(downInt), c.DoFail(unsafe.Pointer(fail)), c.DoOk(unsafe.Pointer(ok)))
	err := d.Download(C.GoString(version))
	return errr(err)
}

//export Check
func Check(version, Type, Minecraftpath *C.char, downInt C.int, fail C.Fail, ok C.Ok) C.err {
	d := gml.NewDown(C.GoString(Type), C.GoString(Minecraftpath), int(downInt), c.DoFail(unsafe.Pointer(fail)), c.DoOk(unsafe.Pointer(ok)))
	err := d.Check(C.GoString(version))
	return errr(err)

}

func errr(err error) C.err {
	c := C.err{}

	if err != nil {
		c.code = 0
		c.msg = nil
		return c
	}

	switch {
	case errors.Is(err, os.ErrNotExist):
		c.code = 1
		c.msg = C.CString(err.Error())
		return c
	case errors.Is(err, launcher.JsonErr):
		c.code = 2
		c.msg = C.CString(err.Error())
		return c
	case errors.Is(err, launcher.JsonNorTrue):
		c.code = 3
		c.msg = C.CString(err.Error())
		return c
	case errors.Is(err, download.NoSuch):
		c.code = 4
		c.msg = C.CString(err.Error())
		return c
	case errors.Is(err, download.FileDownLoadFail):
		c.code = 5
		c.msg = C.CString(err.Error())
		return c
	default:
		c.code = -1
		c.msg = C.CString(err.Error())
		return c
	}
}
