// alg/assembly.go

package alg

import (
	"../cli"
	"../solidity"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type NotifyMine func(tx string, err error)

func Assembly_CheckMine(tx string) (string, error) {
	res, er := cli.CheckContractTransaction(tx)
	if er != nil {
		return "", er
	}
	if len(res.CA) > 0 {
		return res.CA, nil
	}
	return "", nil
}

func NewAssembly(proposal string, discuss string) (string, error) {
	// the value of param must be the hash of IPFS
	if !checkIPFSHash(proposal) {
		return "", errors.New("param for proposal is not correct format")
	}
	if !checkIPFSHash(discuss) {
		return "", errors.New("param for discuss is not correct format")
	}
	param := []string{proposal, discuss}
	abi, er := extractInputABI(solidity.Abi_Assembly.([]interface{}), "")
	if er != nil {
		return "", er
	}
	address, er := cli.NewContract(solidity.Bin_Assembly, param, abi)
	if er != nil {
		return "", er
	}
	return address, nil
}

func Assembly_GetProposal(address string) (string, string, string, error) {
	funcname := "getProposal"
	if !checkAddress(address) {
		return "", "", "", errors.New("param for address is not correct format")
	}
	param := []string{}
	abii, er := extractInputABI(solidity.Abi_Assembly, funcname)
	if er != nil {
		return "", "", "", er
	}
	abio, er := extractOutputABI(solidity.Abi_Assembly, funcname)
	if er != nil {
		return "", "", "", er
	}
	ret, er := cli.Call(address, funcname, param, abii, abio)
	if er != nil {
		return "", "", "", er
	}
	fmt.Print(strings.Join(ret, ",") + "\n")
	return ret[0], ret[1], ret[2], nil
}

func Assembly_RevisionProposal(adrs string, doc string, discuss string) (string, error) {
	funcname := "revisionProposal"
	if !checkAddress(adrs) {
		return "", errors.New("param for address is not correct format")
	}
	if !checkIPFSHash(doc) {
		return "", errors.New("param for proposal is not correct format")
	}
	if !checkIPFSHash(discuss) {
		return "", errors.New("param for discuss is not correct format")
	}
	param := []string{doc, discuss}
	abi, er := extractInputABI(solidity.Abi_Assembly, funcname)
	if er != nil {
		return "", er
	}
	tx, er := cli.Send(adrs, funcname, param, abi)
	if er != nil {
		return "", er
	}
	fmt.Print(tx + "\n")
	return tx, nil
}

func Assembly_CheckRevision(tx string) (string, uint, error) {
	funcname := "onRevisionedProposal"
	abi, er := extractInputABI(solidity.Abi_Assembly, funcname)
	if er != nil {
		return "", 0, er
	}
	res, er := cli.CheckContractTransaction(tx)
	if er != nil {
		return "", 0, er
	}
	if len(res.LOG) == 0 {
		return "", 0, er
	}
	fmt.Printf("res.LOG:%d\n", len(res.LOG))
	fmt.Printf("res.LOG[0].Data:%s\n", res.LOG[0].Data)
	data, er := binToMap(abi, res.LOG[0].Data)
	if er != nil {
		return "", 0, er
	}
	ver, _ := strconv.ParseInt(data["version"], 10, 32)
	return data["adrs"], uint(ver), nil
}
