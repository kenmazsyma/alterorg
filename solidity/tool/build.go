// solidity/tool/build.go
// This source file provides function for building solidity source code.
//
// Usage:
//  go run build.go (file name without extension)
//
// Copyright holder is set forth in alterorg.md

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Print("invalid parameters.")
		os.Exit(1)
	}
	err := run(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run(param string) error {
	abi, err := makeABI(param)
	if err != nil {
		return err
	}
	bin, err := makeBIN(param)
	if err != nil {
		return err
	}
	err = save(param, abi, bin)
	fmt.Printf("saved!\n")
	if err != nil {
		return err
	}
	return nil
}

func makeBIN(param string) (ret map[string]string, err error) {
	ret, err = solc(param, "--bin", "Binary:")
	return ret, err
}

func makeABI(param string) (ret map[string]string, err error) {
	ret, err = solc(param, "--abi", "ContractJSONABI")
	return ret, err
}

func solc(param string, arg string, sep string) (ret map[string]string, err error) {
	out, err := exec.Command("solc", "../src/"+param+".sol", arg).CombinedOutput()
	if err != nil {
		fmt.Print(string(out))
		return nil, err
	}
	if m, _ := regexp.MatchString("^Skipping", string(out)); m {
		return nil, errors.New("no target")
	}
	exp, err2 := regexp.Compile("[ =]*")
	out2 := exp.ReplaceAllString(string(out), "")
	if err2 != nil {
		return nil, err2
	}
	out2 = strings.Replace(out2, sep, "", -1)
	out3 := strings.Split(out2, "\n")

	buf := ""
	ret = map[string]string{}
	for i := range out3 {
		if len(out3[i]) == 0 {
			continue
		}
		if buf == "" {
			buf = out3[i]
		} else {
			ret[buf] = out3[i]
			buf = ""
		}
	}
	return ret, nil
}

func save(param string, abi map[string]string, bin map[string]string) error {
	content := bytes.NewBuffer([]byte{})
	content.WriteString("package solidity\n")
	content.WriteString("import (\n")
	content.WriteString(`    "bytes"` + "\n")
	content.WriteString(`    "github.com/ethereum/go-ethereum/accounts/abi"` + "\n")
	content.WriteString(")\n")
	for k, _ := range abi {
		content.WriteString("var Abi_" + k + " abi.ABI\n")
	}
	content.WriteString("func Init_" + param + "() error{\n")
	content.WriteString("    var v *bytes.Buffer\n")
	content.WriteString("    var er error\n")
	for k, v := range abi {
		content.WriteString("    v=bytes.NewBufferString(`" + v + "`)\n")
		content.WriteString("    Abi_" + k + ", er=abi.JSON(v)\n")
		content.WriteString("    if er != nil {\n")
		content.WriteString("        return er\n")
		content.WriteString("    }\n")
	}
	content.WriteString("    return nil\n")
	content.WriteString("}\n")

	for k, v := range bin {
		content.WriteString("var Bin_" + k + `="0x` + v + `"` + "\n")
	}
	ioutil.WriteFile("../"+param+".go", content.Bytes(), 0644)
	return nil
}
