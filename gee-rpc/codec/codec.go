package codec

import (
	"io"
)

// err = client.Call("Arith.Multiply", args, &reply)
// 典型的rpc调用如上所示， 请求参数为：服务名，方法名， 参数args。 响应参数为reply， err
// 将请求和响应中的参数和返回值抽象为 body， 剩余的信息放在 header 中

type Header struct {
	ServiceMethod string //服务名和方法名 format "Service.Method"
	Seq uint64 // 请求的序号
	Error string //
}

// 对消息体的编解码结构

type Codec interface {
	io.Closer
	ReadHeader(*Header) error
	ReadBody(interface{}) error
	Write(*Header, interface{}) error
}

