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

func NewAssembly() *Assembly {
	return &Assembly{}
}

func (self *Assembly) Create(prm ArgPpslParam, rslt *string) error {
	tx, e := alg.NewAssembly(prm.Proposal, prm.Discussion)
	if e != nil {
		fmt.Printf("%s:\n", e.Error())
		return e
	}
	*rslt = tx
	return nil
}

func (self *Assembly) CheckMine(tx string, rslt *string) error {
	adrs, e := alg.Assembly_CheckMine(tx)
	if e != nil {
		fmt.Printf("%s\n", e.Error())
		return e
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
	tx, e := alg.Assembly_RevisionProposal(prm.Address, prm.Proposal, prm.Discussion)
	if e != nil {
		fmt.Print("%s:\n", e.Error())
		return e
	}
	*rslt = tx
	return nil
}

func (self *Assembly) CheckRevisionProposal(tx string, rslt *ArgChkRevPpslRslt) error {
	adrs, ver, e := alg.Assembly_CheckRevision(tx)
	if e != nil {
		fmt.Print("%s:\n", e.Error())
		return e
	}

	*rslt = ArgChkRevPpslRslt{Address: adrs, Version: ver}
	return nil
}
