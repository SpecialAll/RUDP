package main

import (
	"../RUDP"
	"../common"
	"net"
	"os"
)

//客户端
func main() {

	udpAddr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:9999")

	//连接udpAddr，返回 udpConn
	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		common.Log.Error(err)
		os.Exit(2)
	}
	common.Log.Info("客户端连接成功，开始发送数据...")
	// 发送数据
	RUDP.RUDP(udpConn)
	if err != nil{
		return
	}

	//读取数据
	buf := make([]byte, 1024)
	var len int
	len, _ = udpConn.Read(buf)
	common.Log.Info("客户端接受数据长度:", len)
	common.Log.Info("客户端接受数据:", string(buf))
}

