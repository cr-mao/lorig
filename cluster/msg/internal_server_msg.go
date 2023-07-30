/**
User: cr-mao
Date: 2023/7/30
Time: 18:40
Desc: 网关和业务服务器到通信结构
*/
package msg

import (
	"bytes"
	"encoding/binary"
)

type InternalServerMsg struct {
	GateId  int32  // 在业务服 GateId_ConnId 做唯一id  用
	ConnId  int64  // 网关的连接id
	UserId  int64  // 网关,业务服务器 用户id，网关可以是0， 业务服返回的肯定是知道哪个用户id的。
	MsgData []byte // 原始数据
}

func (msg *InternalServerMsg) ToByteArray() []byte {
	buff := bytes.NewBuffer([]byte{})
	_ = binary.Write(buff, binary.BigEndian, msg.GateId)
	_ = binary.Write(buff, binary.BigEndian, msg.ConnId)
	_ = binary.Write(buff, binary.BigEndian, msg.UserId)
	_ = binary.Write(buff, binary.BigEndian, msg.MsgData)
	return buff.Bytes()
}

func (msg *InternalServerMsg) FromByteArray(byteArray []byte) {
	if nil == byteArray || len(byteArray) <= 0 {
		return
	}
	buff := bytes.NewBuffer(byteArray)
	_ = binary.Read(buff, binary.BigEndian, &msg.GateId)
	_ = binary.Read(buff, binary.BigEndian, &msg.ConnId)
	_ = binary.Read(buff, binary.BigEndian, &msg.UserId)
	msg.MsgData = buff.Bytes()
}
