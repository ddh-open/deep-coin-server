package hostShell

import (
	"fmt"
	"testing"
)

func TestNewContext(t *testing.T) {
	c := NewContext("45.136.184.165", 22, "root")
	err := c.InitCommonTerminal("dou.190824")
	if err != nil {
		fmt.Println(err.Error())
	}

	go func() {
		for {
			cc := <-c.Logs
			fmt.Print(cc)
		}
	}()

	c.SendCmd("top")

}
