// cli/ipfs.go

package cli

import (
	//"encoding/json"
	ipfs "github.com/ipfs/go-ipfs-api"
	//ipfs "../../go-ipfs-api"
	"errors"
	"fmt"
	"io"
	"os"
)

var ipfsurl string
var shell *ipfs.Shell

func InitIpfs(url string) error {
	ipfsurl = url
	shell = ipfs.NewShell(ipfsurl)
	return nil
}

func IpfsAddFile(path string) (string, error) {
	file, er := os.Open(path)
	if er != nil {
		return "", er
	}
	return IpfsAdd(file)
}

func IpfsAdd(reader io.Reader) (string, error) {
	hash, er := shell.Add(reader)
	if er != nil {
		return "", er
	}
	return hash, nil
}

func IpfsGet(hash, path string) error {
	return shell.Get(hash, path)
}

func IpfsBlockGet(hash string) ([]byte, error) {
	return shell.BlockGet(hash)
}

var DIR_IPFS_BOARD = "boardforalterorg"
var dIR_IPFS = "/ipfs/"
var hASH_EMPTY_FILE = "QmbFMke1KXqnYyBBWxB74N4c5SBnJMVAiMNRcGu6x1AwQH"
var hASH_EMPTY_DIR = "QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn"
var ERR_IpfsCreateIpfsDir_01 = "01" // ipns link is created as a file
var ERR_IpfsCreateIpfsDir_02 = "02" // dir which is used for boards of alterorg is created as a file
func IpfsCreateBoardDir() error {
	out, err := shell.ID()
	if err != nil {
		return err
	}
	adrs, err := shell.Resolve(out.ID)
	if err != nil {
		return err
	}
	obj, err := shell.ObjectGet(adrs)
	if err != nil {
		return err
	}
	if len(obj.Links) == 0 && adrs != dIR_IPFS+hASH_EMPTY_DIR { // link to a file
		return errors.New(ERR_IpfsCreateIpfsDir_01)
	} else { // link to a directory
		lst, err := shell.List(adrs)
		if err != nil {
			return err
		}
		nodir := true
		for _, v := range lst {
			if v.Name == DIR_IPFS_BOARD {
				if v.Type == 2 { // 1:dir 2:file
					return errors.New(ERR_IpfsCreateIpfsDir_02)
				}
				nodir = false
				break
			}
		}
		// if the directory is not created yet
		if nodir {
			obj.Links = append(obj.Links, ipfs.ObjectLink{Name: DIR_IPFS_BOARD, Hash: hASH_EMPTY_DIR, Size: 3})
			newhash, err := shell.ObjectPut(obj)
			if err != nil {
				return err
			}
			err = shell.Publish("", newhash)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// dir:1, file:2
func IpfsList(path string) ([]*ipfs.LsLink, error) {
	return shell.List(path)
}

func IpfsObjGet(path string) (*ipfs.IpfsObject, error) {
	return shell.ObjectGet(path)
}
