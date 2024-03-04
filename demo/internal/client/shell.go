package client

import (
	"net"

	"github.com/abiosoft/ishell"

	"github.com/gstones/moke-kit/logging/slogger"
	"github.com/gstones/moke-kit/server/tools"
)

func RunGrpc(url string) {
	sh := ishell.New()
	slogger.Info(sh, "interactive demo connect to "+url)

	if conn, err := tools.DialInsecure(url); err != nil {
		slogger.Die(sh, err)
	} else {
		demoGrpc := NewDemoGrpc(conn)
		sh.AddCmd(demoGrpc.GetCmd())

		sh.Interrupt(func(c *ishell.Context, count int, input string) {
			if count >= 2 {
				c.Stop()
			}
			if count == 1 {
				conn.Close()
				slogger.Done(c, "interrupted, press again to exit")
			}
		})
	}
	sh.Run()
}

func RunZinx(url string) {
	sh := ishell.New()
	slogger.Info(sh, "interactive demo zinx connect to "+url)
	if conn, err := net.Dial("tcp", url); err != nil {
		slogger.Die(sh, err)
	} else {
		demoZinx := NewZinxDemo(conn)
		sh.AddCmd(demoZinx.GetCmd())

		sh.Interrupt(func(c *ishell.Context, count int, input string) {
			if count >= 2 {
				c.Stop()
			}
			if count == 1 {
				conn.Close()
				slogger.Done(c, "interrupted, press again to exit")
			}
		})
	}
	sh.Run()
}
