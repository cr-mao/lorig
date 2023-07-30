/**
User: cr-mao
Date: 2023/7/30
Time: 18:40
Desc: internalServerMsg.go
*/
package msg

import (
	"bytes"
	"encoding/binary"
)

type InternalServerMsg struct {
	GateId  int32  // 暂时没用
	ConnId  int32  // 网关的连接id
	UserId  int64  // 网关过来的用户id
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
