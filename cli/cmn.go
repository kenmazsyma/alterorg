// cli/cmn.go

package cli

import (
	"errors"
	"strings"
)

type ErrCode string
type Status int

func splitArgs(txt string) []string {
	prm := strings.Split(txt, " ")
	for i := range prm {
		start := 0
		end := len(prm[i])
		if prm[i][0] == '"' {
			start = 1
			end--
		}
		prm[i] = prm[i][start:end]
	}
	return prm
}

func makeError(msg ErrCode) error {
	return errors.New(string(msg))
}

func Equals(err error, tp ErrCode) bool {
	return err.Error() == string(tp)
}

const (
	LBL_ETH  string = "eth"
	LBL_IPFS string = "ipfs"
)
