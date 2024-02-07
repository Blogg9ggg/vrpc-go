package client

// Call 代表一次 RPC.
type Call struct {
	Seq           uint64
	ServiceMethod string      // 指示服务名, 格式: "<service>.<method>"
	Args          interface{} // 输入参数
	Reply         interface{} // 服务的输出
	Error         error       // if error occurs, it will be set
	Done          chan *Call  // Strobes when call is complete.
}

func (call *Call) done() {
	call.Done <- call
}
