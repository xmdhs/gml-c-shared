package gml

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/xmdhs/gomclauncher/download"
)

type Down struct {
	Atype         string
	Downint       int
	download      string
	run           string
	Ok            func(int, int)
	Fail          func(string)
	Minecraftpath string
}

type test struct {
	ID           string `json:"id"`
	InheritsFrom string `json:"inheritsFrom"`
}

func NewDown(Type, Minecraftpath string, downInt int, Fail func(string), Ok func(int, int)) *Down {
	if Fail == nil {
		Fail = func(s string) {}
	}
	if Ok == nil {
		Ok = func(i1, i2 int) {}
	}
	return &Down{
		Atype:         Type,
		Minecraftpath: Minecraftpath,
		Downint:       downInt,
		Fail:          Fail,
		Ok:            Ok,
	}
}

func (d *Down) Download(version string) error {
	d.download = version
	d.run = ""
	err := d.d()
	if err != nil {
		return fmt.Errorf("down.Download: %w", err)
	}
	return nil
}

func (d *Down) Check(version string) error {
	d.download = version
	d.run = version

	b, err := ioutil.ReadFile(d.Minecraftpath + "/versions/" + version + "/" + version + ".json")
	if err != nil {
		return fmt.Errorf("down.Check: %w", err)
	}

	t := test{}
	err = json.Unmarshal(b, &t)
	if err != nil {
		return fmt.Errorf("down.Check: %w", err)
	}
	if t.ID != version {
		b = bytes.ReplaceAll(b, []byte(t.ID), []byte(version))
		err := ioutil.WriteFile(d.Minecraftpath+"/versions/"+version+"/"+version+".json", b, 0777)
		if err != nil {
			return fmt.Errorf("down.Check: %w", err)
		}
	}
	if t.InheritsFrom != "" {
		d.download = t.InheritsFrom
		err = d.d()
	} else {
		d.download = version
		err = d.d()
	}
	if err != nil {
		return fmt.Errorf("down.Check: %w", err)
	}
	return nil
}

func (f Down) d() error {
	cxt := context.TODO()
	l, err := download.Getversionlist(cxt, f.Atype, f.Fail)
	if err != nil {
		return fmt.Errorf("down.d: %w", err)
	}
	err = l.Downjson(cxt, f.download, f.Fail)
	if !(f.run != "" && err != nil && errors.Is(err, download.NoSuch)) {
		if err != nil {
			return fmt.Errorf("down.d: %w", err)
		}
	}
	var b []byte
	if f.run != "" {
		b, err = ioutil.ReadFile(f.Minecraftpath + "/versions/" + f.run + "/" + f.run + ".json")
	} else {
		b, err = ioutil.ReadFile(f.Minecraftpath + "/versions/" + f.download + "/" + f.download + ".json")
	}
	if err != nil {
		return fmt.Errorf("down.d: %w", err)
	}
	dl, err := download.Newlibraries(cxt, b, f.Atype, f.Fail)
	if err != nil {
		return fmt.Errorf("down.d: %w", err)
	}
	err = dl.Downjar(f.download)
	if err != nil {
		return fmt.Errorf("down.d: %w", err)
	}
	f.Ok(1, 0)
	err = f.dd(dl.Downassets, 2)
	if err != nil {
		return fmt.Errorf("down.d: %w", err)
	}
	err = f.dd(dl.Downlibrarie, 3)
	if err != nil {
		return fmt.Errorf("down.d: %w", err)
	}
	err = dl.Unzip(f.Downint)
	if err != nil {
		return fmt.Errorf("down.d: %w", err)
	}
	return nil
}

func (f Down) dd(down func(i int, c chan int) error, a int) error {
	ch := make(chan int, 5)
	e := make(chan error, 5)
	var err error
	go func() {
		err = down(f.Downint, ch)
		if err != nil {
			e <- err
		}
	}()
b:
	for {
		select {
		case i, ok := <-ch:
			f.Ok(a, i)
			if !ok {
				break b
			}
		case err := <-e:
			if err != nil {
				return fmt.Errorf("down.dd: %w", err)
			}
		}
	}
	return nil
}
