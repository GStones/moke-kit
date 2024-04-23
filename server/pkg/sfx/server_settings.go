package sfx

import (
	"go.uber.org/fx"

	"github.com/gstones/moke-kit/utility"
)

//Version          string // The version of the Zinx framework.(当前Zinx版本号)
//MaxPacketSize    uint32 // The maximum size of the packets that can be sent or received.(读写数据包的最大值)
//MaxConn          int    // The maximum number of connections that the server can handle.(当前服务器主机允许的最大链接个数)
//WorkerPoolSize   uint32 // The number of worker pools in the business logic.(业务工作Worker池的数量)
//MaxWorkerTaskLen uint32 // The maximum number of tasks that a worker pool can handle.(业务工作Worker对应负责的任务队列最大任务存储数量)
//WorkerMode       string // The way to assign workers to connections.(为链接分配worker的方式)
//MaxMsgChanLen    uint32 // The maximum length of the send buffer message queue.(SendBuffMsg发送消息的缓冲最大长度)
//IOReadBuffSize   uint32 // The maximum size of the read buffer for each IO operation.(每次IO最大的读取长度)

// SettingsParams All server settings module
type SettingsParams struct {
	fx.In

	Port       int32 `name:"Port"`       // grpc/http port
	Timeout    int32 `name:"Timeout"`    // tcp service heartbeat timeout
	RateLimit  int32 `name:"RateLimit"`  // all server type rate limit per second
	OtelEnable bool  `name:"OtelEnable"` // open telemetry enable

	//--------------------- zinx settings ---------------------
	// pure tcp port
	ZinxTcpPort int32 `name:"ZinxTcpPort"` // tcp port
	// websocket port
	ZinxWSPort int32 `name:"ZinxWSPort"` // websocket port
	// The maximum size of the packets that can be sent or received
	MaxPacketSize uint32 `name:"MaxPacketSize"`
	// The number of worker pools in the business logic
	WorkerPoolSize uint32 `name:"WorkerPoolSize"`
	// The maximum number of tasks that a worker pool can handle
	MaxWorkerTaskLen uint32 `name:"MaxWorkerTaskLen"`
	// The maximum length of the send buffer message queue
	MaxMsgChanLen uint32 `name:"MaxMsgChanLen"`
}

// SettingsResult loads from the environment and its members are injected into the tfx dependency graph.
type SettingsResult struct {
	fx.Out

	Port       int32 `name:"Port"  envconfig:"PORT" default:"8081"`
	Timeout    int32 `name:"Timeout" envconfig:"TIMEOUT" default:"10"`
	RateLimit  int32 `name:"RateLimit" envconfig:"RATE_LIMIT" default:"1000"`
	OtelEnable bool  `name:"OtelEnable" envconfig:"OTEL_ENABLE" default:"false"`

	// --------------------- zinx settings ---------------------
	ZinxTcpPort int32 `name:"ZinxTcpPort" envconfig:"ZINX_TCP_PORT" default:"8888"`
	// websocket port
	ZinxWSPort int32 `name:"ZinxWSPort" envconfig:"ZINX_WS_PORT" default:""`
	// The maximum size of the packets that can be sent or received
	MaxPacketSize uint32 `name:"MaxPacketSize" envconfig:"MAX_PACKET_SIZE" default:"4096"`
	// The number of worker pools in the business logic
	WorkerPoolSize uint32 `name:"WorkerPoolSize" envconfig:"WORKER_POOL_SIZE" default:"64"`
	// The maximum number of tasks that a worker pool can handle
	MaxWorkerTaskLen uint32 `name:"MaxWorkerTaskLen" envconfig:"MAX_WORKER_TASK_LEN" default:"1024"`
	// The maximum length of the send buffer message queue
	MaxMsgChanLen uint32 `name:"MaxMsgChanLen" envconfig:"MAX_MSG_CHAN_LEN" default:"1024"`
}

func (g *SettingsResult) LoadFromEnv() (err error) {
	err = utility.Load(g)
	return
}

var SettingsModule = fx.Provide(
	func() (out SettingsResult, err error) {
		err = out.LoadFromEnv()
		return
	},
)
