package client

import (
	"github.com/abiosoft/ishell"

	"moke-kit/utility/cshell"
	"moke-kit/utility/ugrpc"
)

func Shell(url string) {
	sh := ishell.New()
	cshell.Info(sh, "interactive demo connect to "+url)

	if conn, err := ugrpc.DialWithOptions(url, false); err != nil {
		cshell.Die(sh, err)
	} else {
		hello := NewHello(conn)
		sh.AddCmd(hello.GetCmd())

		sh.Interrupt(func(c *ishell.Context, count int, input string) {
			if count >= 2 {
				c.Stop()
			}
			if count == 1 {
				conn.Close()
				cshell.Done(c, "interrupted, press again to exit")
			}
		})
	}
	sh.Run()
}
