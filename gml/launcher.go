package gml

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/xmdhs/gomclauncher/launcher"
)

type Launcher struct {
	launcher.Gameinfo
	Independent bool
}

func (f *Launcher) GenCmdArgs() ([]string, error) {
	if f.Independent {
		f.Gamedir = f.Minecraftpath + "/versions/" + f.Version
	} else {
		f.Gamedir = f.Minecraftpath
	}
	b, err := ioutil.ReadFile(f.Minecraftpath + "/versions/" + f.Version + "/" + f.Version + ".json")
	if err != nil {
		return nil, fmt.Errorf("Launcher.GenCmdArgs: %w", err)
	}
	f.Jsonbyte = b
	_, args, err := f.GenLauncherCmdArgs()
	if err != nil {
		return nil, fmt.Errorf("Launcher.GenCmdArgs: %w", err)
	}
	return args, nil
}

func ListVersion(path string) ([]string, error) {
	l, err := find(path)
	if err != nil {
		return nil, fmt.Errorf("ListVersion: %w", err)
	}
	rl := make([]string, 0, len(l))
	for _, v := range l {
		if testVersion(path + "/" + v + `/` + v + ".json") {
			rl = append(rl, v)
		}
	}
	return rl, nil
}

func find(path string) ([]string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("find: %w", err)
	}
	s := make([]string, 0)
	for _, f := range files {
		if f.IsDir() {
			s = append(s, f.Name())
		}
	}
	return s, nil
}

func testVersion(path string) bool {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}
	t := t{}
	err = json.Unmarshal(b, &t)
	if err != nil {
		return false
	}
	if len(t.Libraries) == 0 {
		return false
	}
	if t.MainClass == "" {
		return false
	}
	return true
}

type t struct {
	Libraries []interface{} `json:"Libraries"`
	MainClass string        `json:"mainClass"`
}
