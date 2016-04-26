package cmn

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ApEnvSet struct {
	Orgs []string `json:"orgs"`
}

var ApEnv ApEnvSet

func LoadApEnv(path string) error {
	fin, er := os.Open(path)
	if er != nil {
		fmt.Print("failure to open appenv file\n")
		return er
	}
	defer fin.Close()
	buf, er := ioutil.ReadAll(fin)
	if er != nil {
		fmt.Print("failure to read appenv file\n")
		return er
	}
	ApEnv = ApEnvSet{}
	er = json.Unmarshal(buf, &ApEnv)
	if er != nil {
		fmt.Printf("appenv file is bad format:%s\n", er.Error())
		return er
	}
	return nil
}

func SaveApEnv(path string) error {
	data, er := json.Marshal(ApEnv)
	if er != nil {
		return er
	}
	ioutil.WriteFile(path, data, 0644)
	fin, er := os.Open(path)
	if er != nil {
		fmt.Print("failure to open appenv file\n")
		return er
	}
	defer fin.Close()
	buf, er := ioutil.ReadAll(fin)
	if er != nil {
		fmt.Print("failure to read appenv file\n")
		return er
	}
	env := ApEnvSet{}
	er = json.Unmarshal(buf, &env)
	if er != nil {
		fmt.Print("appenv file is bad format\n")
		return er
	}
	return nil
}

func QueryOrgList() ([]string, error) {
	return ApEnv.Orgs, nil
}

func UpdateOrgList(val []string) error {
	ApEnv.Orgs = val
	return nil
}
