// cli/ipfs.go

package cli

import (
	ipfs "github.com/ipfs/go-ipfs-api"
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
