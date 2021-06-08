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
*/
import "C"

import (
	"strconv"
	"unsafe"

	"github.com/xmdhs/gml-c-shared/gml"
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
		return C.GameinfoReturn{}, 1
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
