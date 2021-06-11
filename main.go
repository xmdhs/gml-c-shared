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
	char **Flag;
	int  flag_len;
	int  independent;
} Gameinfo;

typedef struct
{
    int  code;
    char *msg;
} err;


typedef struct
{
    char **charlist;
    int len;
	err e;

} GmlReturn;

typedef void (*Fail)(char*);

typedef void (*Ok)(int,int);

typedef struct
{
    char *Username;
    char *ClientToken;
    char *UUID;
    char *AccessToken;
    char *ApiAddress;
} OldAuthDate;
*/
import "C"

import (
	"errors"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"unsafe"

	"github.com/xmdhs/gml-c-shared/gml"
	"github.com/xmdhs/gomclauncher/auth"
	"github.com/xmdhs/gomclauncher/download"
	"github.com/xmdhs/gomclauncher/launcher"

	"github.com/xmdhs/gml-c-shared/c"
)

func main() {}

//export GenCmdArgs
func GenCmdArgs(g C.Gameinfo) C.GmlReturn {
	Bool := g.independent == 1
	flag := []string{}

	for i := 0; i < int(g.flag_len); i++ {
		c := Getchar(g.Flag, C.longlong(i))
		flag = append(flag, C.GoString(c))
	}

	l := gml.Launcher{
		Gameinfo: launcher.Gameinfo{
			Minecraftpath: C.GoString(g.Minecraftpath),
			RAM:           strconv.Itoa(int(g.RAM)),
			Name:          C.GoString(g.Name),
			UUID:          C.GoString(g.UUID),
			AccessToken:   C.GoString(g.AccessToken),
			Version:       C.GoString(g.Version),
			ApiAddress:    C.GoString(g.ApiAddress),
			Flag:          flag,
		},
		Independent: Bool,
	}
	args, err := l.GenCmdArgs()
	if err != nil {
		i := errr(err)
		c := C.GmlReturn{}
		c.e = i
		return c
	}
	c := c.NewChar(len(args))

	for i, v := range args {
		c.SetChar(i, unsafe.Pointer(C.CString(v)))
	}

	var r C.GmlReturn
	r.charlist = (**C.char)(c.P)
	r.len = C.int(len(args))
	return r
}

//export SetProxy
func SetProxy(httpProxy *C.char) C.err {
	proxy, err := url.Parse(C.GoString(httpProxy))
	if err != nil {
		e := errr(err)
		return e
	}
	auth.Transport.Proxy = http.ProxyURL(proxy)
	return C.err{}
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

//export ListVersion
func ListVersion(path *C.char) C.GmlReturn {
	l, err := gml.ListVersion(C.GoString(path))
	return list(err, l)
}

//export ListDownloadType
func ListDownloadType(Type *C.char) C.GmlReturn {
	l, err := gml.ListDownloadType(C.GoString(Type))
	return list(err, l)
}

//export ListDownloadVersion
func ListDownloadVersion(VerType, Type *C.char) C.GmlReturn {
	l, err := gml.ListDownloadVersion(C.GoString(VerType), C.GoString(Type))
	return list(err, l)
}

func list(err error, l []string) C.GmlReturn {
	if err != nil {
		e := errr(err)
		c := C.GmlReturn{}
		c.e = e
		return c
	}
	c := c.NewChar(len(l))
	for i, v := range l {
		c.SetChar(i, unsafe.Pointer(C.CString(v)))
	}
	r := C.GmlReturn{}
	r.charlist = (**C.char)(c.P)
	r.len = C.int(len(l))
	return r
}

//export OldAuth
func OldAuth(ApiAddress, username, email, password, clientToken *C.char) (C.OldAuthDate, C.err) {
	apiAddress, err := auth.Getauthlibapi(C.GoString(ApiAddress))
	if err != nil {
		e := errr(err)
		return C.OldAuthDate{}, e
	}
	a, err := auth.Authenticate(apiAddress, C.GoString(username), C.GoString(email), C.GoString(password), C.GoString(clientToken))
	if err != nil {
		e := errr(err)
		return C.OldAuthDate{}, e
	}
	ca := C.OldAuthDate{}
	ca.AccessToken = C.CString(a.AccessToken)
	ca.ApiAddress = C.CString(a.ApiAddress)
	ca.ClientToken = C.CString(a.ClientToken)
	ca.UUID = C.CString(a.ID)
	ca.Username = C.CString(a.Username)
	return ca, C.err{}
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
	case errors.Is(err, launcher.JsonErr):
		c.code = 2
	case errors.Is(err, launcher.JsonNorTrue):
		c.code = 3
	case errors.Is(err, download.NoSuch):
		c.code = 4
	case errors.Is(err, download.FileDownLoadFail):
		c.code = 5
	case errors.Is(err, auth.ErrNotSelctProFile):
		c.code = 6
	case errors.Is(err, auth.ErrProFileNoExist):
		c.code = 7
	case errors.Is(err, auth.NotOk):
		c.code = 8
	case errors.Is(err, auth.NoProfiles):
		c.code = 9
	default:
		c.code = -1
	}
	c.msg = C.CString(err.Error())
	return c
}
