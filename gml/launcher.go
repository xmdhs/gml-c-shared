package gml

import (
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
