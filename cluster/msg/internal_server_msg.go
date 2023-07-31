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

	"github.com/cr-mao/lorig/conf"
	"github.com/cr-mao/lorig/errors"
)

var (
	ErrMessageEmpty = errors.New("message empty")
)

var maxMessageSize = 5004

func init() {
	maxMessageSize = conf.GetInt("packet.bufferBytes", 5000) + 4
}

type InternalServerMsg struct {
	GateId  int32  // 在业务服 GateId_ConnId 做唯一id  用
	ConnId  int64  // 网关的连接id
	UserId  int64  // 网关,业务服务器 用户id，网关可以是0， 业务服返回的肯定是知道哪个用户id的。
	MsgData []byte // 原始数据
}

func (msg *InternalServerMsg) Pack() ([]byte, error) {
	// todo from config
	if len(msg.MsgData) > maxMessageSize {
		return nil, nil
	}
	buff := bytes.NewBuffer([]byte{})
	err := binary.Write(buff, binary.BigEndian, msg.GateId)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buff, binary.BigEndian, msg.ConnId)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buff, binary.BigEndian, msg.UserId)
	if err != nil {
		return nil, err
	}
	err = binary.Write(buff, binary.BigEndian, msg.MsgData)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func (msg *InternalServerMsg) UnPack(byteArray []byte) error {
	if nil == byteArray || len(byteArray) <= 0 {
		return ErrMessageEmpty
	}
	buff := bytes.NewBuffer(byteArray)
	err := binary.Read(buff, binary.BigEndian, &msg.GateId)
	if err != nil {
		return err
	}
	err = binary.Read(buff, binary.BigEndian, &msg.ConnId)
	if err != nil {
		return err
	}
	err = binary.Read(buff, binary.BigEndian, &msg.UserId)
	if err != nil {
		return err
	}
	msg.MsgData = buff.Bytes()
	return nil
}
