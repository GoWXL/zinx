package znet

import (
	"bytes"
	"encoding/binary"
	"zinx/ziface"
)

//拆包封包模块 解决粘包问题
type DataPack struct {
}

//初始化NEW方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

//获取头部长度 4+4
func (dp *DataPack) GetHeadLen() uint32 {
	return 8
}

//封包
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	//创建一个二进制字节缓冲
	dataBuffer := bytes.NewBuffer([]byte{})
	//把数据长度写入缓冲区
	err := binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgLen())
	if err != nil {
		//fmt.Println("databuffer datalen write err", err)
		return nil, err
	}
	//将dataId 写入缓冲区
	err = binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgId())
	if err != nil {
		//fmt.Println("databuffer dataId write err", err)
		return nil, err
	}
	//将data写入缓冲区
	err = binary.Write(dataBuffer, binary.LittleEndian, msg.GetMsgData())
	if err != nil {
		//fmt.Println("databuffer data write err", err)
		return nil, err
	}
	return dataBuffer.Bytes(), nil
}

//拆包
func (dp *DataPack) UnPack(binarydata []byte) (ziface.IMessage, error) {
	msghead := &Message{}
	//创建读取二进制流的阅读reader
	dataBuffer := bytes.NewReader(binarydata)
	//读取datalen
	err := binary.Read(dataBuffer, binary.LittleEndian, &msghead.DataLen)
	if err != nil {
		return nil, err
	}
	//读取dataId
	err = binary.Read(dataBuffer, binary.LittleEndian, &msghead.Id)
	if err != nil {
		return nil, err
	}
	return msghead, nil
}
