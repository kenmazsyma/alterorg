package cmn

import (
	"fmt"
)

func Log(lbl string, txt string, args ...interface{}) {
	fmt.Printf("["+lbl+"]"+txt+"\n", args...)
}
