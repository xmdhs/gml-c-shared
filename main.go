package main

/*
#cgo LDFLAGS: -L${SRCDIR}/c -lchar

char **newchar(long long i);

void setchar(char **c, long long index, char *s);

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

	"github.com/xmdhs/gml-c-shared/gml"
	"github.com/xmdhs/gomclauncher/launcher"
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

	c := C.newchar(C.longlong(int64(len(args))))

	for i, v := range args {
		C.setchar(c, C.longlong(int64(i)), C.CString(v))
	}

	var r C.GameinfoReturn
	r.args = c
	r.len = C.longlong(int64(len(args)))
	return r, 0
}
