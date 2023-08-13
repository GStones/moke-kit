package client

import (
	"github.com/abiosoft/ishell"
	"net"

	"moke-kit/utility/cshell"
	"moke-kit/utility/ugrpc"
)

func RunGrpc(url string) {
	sh := ishell.New()
	cshell.Info(sh, "interactive demo connect to "+url)

	if conn, err := ugrpc.DialWithOptions(url, false); err != nil {
		cshell.Die(sh, err)
	} else {
		demoGrpc := NewDemoGrpc(conn)
		sh.AddCmd(demoGrpc.GetCmd())

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

func RunZinx(url string) {
	sh := ishell.New()
	cshell.Info(sh, "interactive demo zinx connect to "+url)
	if conn, err := net.Dial("tcp", url); err != nil {
		cshell.Die(sh, err)
	} else {
		demoZinx := NewZinxDemo(conn)
		sh.AddCmd(demoZinx.GetCmd())

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
