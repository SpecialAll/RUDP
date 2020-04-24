package main

import (
	"../RUDP"
	"../common"
	"net"
)

// udp 服务端
func main() {
	/*
	   network: "udp"、"udp4"或"udp6"
	   addr: "host:port"或"[ipv6-host%zone]:port"
	*/
	udpAddr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:9999")

	//监听端口
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		common.Log.Error(err)
	}
	defer udpConn.Close()

	common.Log.Info("udp listening to ------->  ", udpAddr.Port)

	//udp不需要Accept
	for {
		RUDP.HandleConnection(udpConn)
	}
}


