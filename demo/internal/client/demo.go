package client

import (
	"context"
	"fmt"

	"github.com/abiosoft/ishell"
	mm "github.com/grpc-ecosystem/go-grpc-middleware/v2/metadata"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

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

	md := metadata.Pairs("authorization", fmt.Sprintf("%s %v", "bearer", "test"))
	ctx := mm.MD(md).ToOutgoing(context.Background())
	if response, err := p.client.Hi(ctx, &pb.HiRequest{
		Message: msg,
	}); err != nil {
		cshell.Warn(c, err)
	} else {
		cshell.Infof(c, "Response: %s", response.Message)
	}
}
