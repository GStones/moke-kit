package client

import (
	"context"

	"github.com/abiosoft/ishell"
	"google.golang.org/grpc"

	pb "moke-kit/demo/api/gen/demo/api"
	"moke-kit/utility/cshell"
)

type Hello struct {
	client pb.HelloClient
	cmd    *ishell.Cmd
}

func NewHello(conn *grpc.ClientConn) *Hello {
	cmd := &ishell.Cmd{
		Name:    "demo",
		Help:    "demo interactive",
		Aliases: []string{"D"},
	}
	p := &Hello{
		client: pb.NewHelloClient(conn),
		cmd:    cmd,
	}
	p.initSubShells()
	return p
}

func (p *Hello) GetCmd() *ishell.Cmd {
	return p.cmd
}

func (p *Hello) initSubShells() {
	p.cmd.AddCmd(&ishell.Cmd{
		Name:    "hi",
		Help:    "say hi",
		Aliases: []string{"hi"},
		Func:    p.sayHi,
	})

}

func (p *Hello) sayHi(c *ishell.Context) {
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	cshell.Info(c, "Enter say hi message...")
	msg := cshell.ReadLine(c, "message: ")

	if response, err := p.client.Hi(context.Background(), &pb.HiRequest{
		Message: msg,
	}); err != nil {
		cshell.Warn(c, err)
	} else {
		cshell.Infof(c, "Response: %s", response.Message)
	}
}
