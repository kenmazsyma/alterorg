package api

import (
	"../alg"
	"fmt"
)

type Assembly struct {
}

type ArgPpslParam struct {
	Proposal   string `json:"proposal"`
	Discussion string `json:"discussion"`
}

type ArgRevPpslParam struct {
	Address    string `json:"address"`
	Proposal   string `json:"proposal"`
	Discussion string `json:"discussion"`
}

type ArgPpslRslt struct {
	Doc     string `json:"doc"`
	Discuss string `json:"discuss"`
	Arbiter string `json:"arbiter"`
}

type ArgChkRevPpslRslt struct {
	Address string `json:"address"`
	Version uint   `json:"version"`
}

type ArgGetBasicInfo struct {
	Name     string `json:"name"`
	Proposal string `json:"proposal"`
	Arbiter  string `json:"arbiter"`
	Version  string `json:"version"`
}

type ArgGetParticipants struct {
	Persons []string `json:"persons"`
}

type ArgGetNofTokenArg struct {
	Address string `json:"address"`
	Person  string `json:"person"`
}

func NewAssembly() *Assembly {
	return &Assembly{}
}

func (self *Assembly) Create(name string, rslt *string) error {
	tx, err := alg.NewAssembly(name)
	if err != nil {
		fmt.Printf("%s:\n", err.Error())
		return err
	}
	*rslt = tx
	return nil
}

func (self *Assembly) getName(address string, rslt *string) error {
	name, err := alg.Assembly_GetName(address)
	if err != nil {
		fmt.Printf("%s:\n", err.Error())
		return err
	}
	*rslt = name
	return nil
}

func (self *Assembly) CheckMine(tx string, rslt *string) error {
	//adrs, err := alg.Assembly_CheckMine(tx)
	adrs, err := alg.Assembly_ChkCreated(tx)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return err
	}
	*rslt = adrs
	return nil
}

func (self *Assembly) GetProposal(adrs string, rslt *ArgPpslRslt) error {
	doc, discuss, arbiter, er := alg.Assembly_GetProposal(adrs)
	if er != nil {
		fmt.Printf("%s:\n", er.Error())
		return er
	}
	*rslt = ArgPpslRslt{Doc: doc, Discuss: discuss, Arbiter: arbiter}
	return nil
}

func (self *Assembly) RevisionProposal(prm ArgRevPpslParam, rslt *string) error {
	tx, err := alg.Assembly_RevisionProposal(prm.Address, prm.Proposal, prm.Discussion)
	if err != nil {
		fmt.Print("%s:\n", err.Error())
		return err
	}
	*rslt = tx
	return nil
}

func (self *Assembly) CheckRevisionProposal(tx string, rslt *ArgChkRevPpslRslt) error {
	adrs, ver, err := alg.Assembly_CheckRevision(tx)
	if err != nil {
		fmt.Print("%s:\n", err.Error())
		return err
	}

	*rslt = ArgChkRevPpslRslt{Address: adrs, Version: ver}
	return nil
}

func (self *Assembly) GetBasicInfo(address string, rslt *ArgGetBasicInfo) error {
	info, err := alg.Assembly_GetBasicInfo(address)
	if err != nil {
		fmt.Print("%s:\n", err.Error())
		return err
	}
	*rslt = ArgGetBasicInfo{Name: info[0], Proposal: info[1], Arbiter: info[2], Version: info[3]}
	return nil
}

func (self *Assembly) GetNofToken(args ArgGetNofTokenArg, rslt *int) error {
	var err error
	*rslt, err = alg.Assembly_GetNofToken(args.Address, args.Person)
	if err != nil {
		fmt.Print("%s:\n", err.Error())
		return err
	}
	return nil
}

func (self *Assembly) GetParticipants(address string, rslt *ArgGetParticipants) error {
	ret, err := alg.Assembly_GetParticipants(address)
	if err != nil {
		fmt.Print("%s:\n", err.Error())
		return err
	}
	*rslt = ArgGetParticipants{Persons: ret}
	return nil
}
