package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

//函数名Test开头 后面函数名自定义 形参
func TestDataPack(t *testing.T) {
	fmt.Println("tes tdatapack...")
	listener, err := net.Listen("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("net listen err", err)
		return
	}
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("conn accept err", err)
			}
			//创建读写业务
			go func(conn *net.Conn) {
				dp := NewDataPack()
				for {
					//第一次从conn读 把head读出来
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(*conn, headData)
					if err != nil {
						fmt.Println("read head err")
						break
					}
					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack err", err)
						return
					}
					//数据区有内容 需要第二次读取
					if msgHead.GetMsgLen() > 0 {
						//将message转换 从iMessage转换到message
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						_, err := io.ReadFull(*conn, msg.Data)
						if err != nil {
							fmt.Println("server unpach err", err)
							return
						}
						fmt.Println("---> Recv MsgID = ", msg.Id, " datalen = ", msg.Datalen, "data = ", string(msg.Data))
					}
				}
			}(&conn)
		}
	}()
	//模拟写一个client 进行封包拆包
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err", err)
		return
	}
	dp := NewDataPack()
	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client send data err", err)
		return
	}
	msg2 := &Message{
		Id:      2,
		DataLen: 5,
		Data:    []byte{'h', 'e', 'l', 'l','o'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client send data2 err", err)
		return
	}
	//将两个包黏在一起
	sendData1 = append(sendData1,sendData2...) //[4][1]zinx[5][2]hello
	//发送
	conn.Write(sendData1)


	//让test不结束
	select{}
}