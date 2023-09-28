package client

import (
	"io"
	"net"

	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zpack"
	"google.golang.org/protobuf/proto"

	"github.com/abiosoft/ishell"

	pb "github.com/gstones/moke-kit/demo/api/gen/demo/api"
	"github.com/gstones/moke-kit/logging/slogger"
)

type DemoZinx struct {
	conn net.Conn
	cmd  *ishell.Cmd
}

func NewZinxDemo(conn net.Conn) *DemoZinx {
	cmd := &ishell.Cmd{
		Name:    "demo",
		Help:    "demo interactive",
		Aliases: []string{"D"},
	}
	p := &DemoZinx{
		conn: conn,
		cmd:  cmd,
	}
	p.initSubShells()
	return p
}

func (dz *DemoZinx) GetCmd() *ishell.Cmd {
	return dz.cmd
}

func (dz *DemoZinx) initSubShells() {
	dz.cmd.AddCmd(&ishell.Cmd{
		Name:    "hi",
		Help:    "say hi",
		Aliases: []string{"hi"},
		Func:    dz.sayHi,
	})
	dz.cmd.AddCmd(&ishell.Cmd{
		Name:    "watch",
		Help:    "watch topic",
		Aliases: []string{"w"},
		Func:    dz.watch,
	})

}

func (dz *DemoZinx) sayHi(c *ishell.Context) {
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	slogger.Info(c, "Enter say hi message...")
	msg := slogger.ReadLine(c, "message: ")

	req := &pb.HiRequest{
		Message: msg,
		Uid:     "10000",
	}
	data, err := proto.Marshal(req)
	if err != nil {
		slogger.Warn(c, err)
		return
	}
	dp := zpack.NewDataPack()
	sendData, _ := dp.Pack(zpack.NewMsgPackage(1, data))
	_, err = dz.conn.Write(sendData)
	if err != nil {
		return
	}
}

func (dz *DemoZinx) watch(c *ishell.Context) {
	c.ShowPrompt(false)
	defer c.ShowPrompt(true)

	slogger.Info(c, "Enter watch topic...")
	topic := slogger.ReadLine(c, "topic: ")

	req := &pb.WatchRequest{
		Topic: topic,
		Uid:   "10000",
	}
	data, err := proto.Marshal(req)
	if err != nil {
		slogger.Warn(c, err)
		return
	}
	dp := zpack.NewDataPack()
	sendData, _ := dp.Pack(zpack.NewMsgPackage(2, data))
	_, err = dz.conn.Write(sendData)
	if err != nil {
		return
	}
	watchResponse(c, dz.conn)
}

func watchResponse(c *ishell.Context, conn net.Conn) {
	dp := zpack.NewDataPack()
	go func() {
		for {
			id, data, err := unPackResponse(dp, conn)
			if err != nil {
				slogger.Warn(c, err)
				return
			}
			if id == 1 {
				resp := &pb.HiResponse{}
				err := proto.Unmarshal(data, resp)
				if err != nil {
					slogger.Warn(c, err)
					return
				}
				slogger.Infof(c, "response: %s \r\n", resp.String())
				continue
			} else if id == 2 {
				resp := &pb.WatchResponse{}
				err := proto.Unmarshal(data, resp)
				if err != nil {
					slogger.Warn(c, err)
					return
				}
				slogger.Infof(c, "watching:%d %s \r\n", id, resp.String())
				continue
			}
		}
	}()

}

func unPackResponse(dp ziface.IDataPack, conn io.Reader) (uint32, []byte, error) {
	zconf.GlobalObject.MaxPacketSize = 0
	//先读出流中的head部分
	headData := make([]byte, dp.GetHeadLen())
	_, err := conn.Read(headData)
	if err != nil {
		return 0, nil, err
	}
	//将headData字节流 拆包到msg中
	msgHead, err := dp.Unpack(headData)
	if err != nil {
		return 0, nil, err
	}
	if msgHead.GetDataLen() > 0 {
		//msg是有data数据的，需要再次读取data数据
		data := make([]byte, msgHead.GetDataLen())
		//根据dataLen从io中读取字节流
		_, err := io.ReadFull(conn, data)
		if err != nil {
			return 0, nil, err
		}
		return msgHead.GetMsgID(), data, nil
	}

	return msgHead.GetMsgID(), nil, nil
}
