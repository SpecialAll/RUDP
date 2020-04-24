package RUDP

import (
	"../common"
	"encoding/json"
	"net"
	"strconv"
)

// 读取消息
func HandleConnection(udpConn *net.UDPConn) {

	// 读取数据
	buf := make([]byte,1024)
	len, udpAddr, err := udpConn.ReadFromUDP(buf)
	if err != nil{
		return
	}
	receData := make(map[int]string)

	//防止该错误：invalid character '\x00' after top-level value
	err = json.Unmarshal(buf[:len], &receData)
	if err != nil {
		common.Log.Error("数据解析错误，<", err, ">")
	}
	common.Log.Info("服务端已接受客户端数据长度[", len, "]: " , receData)

	//反馈数据，以字符串形式反馈
	feedbackStr := feedBackData(receData)

	// 发送数据
	len, err = udpConn.WriteToUDP([]byte(feedbackStr), udpAddr)
	if err != nil{
		common.Log.Error(err)
		return
	}

	common.Log.Info("服务端已回传数据！")
}

func feedBackData(data map[int]string) string{
	feedbackStr := ""
	start := 1
	for index := range data {
		if start == 1 {
			feedbackStr = strconv.Itoa(index)
		}else {
			feedbackStr = "," + strconv.Itoa(index)
		}
		start++
	}
	return feedbackStr
}
