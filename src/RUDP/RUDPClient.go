package RUDP

import (
	"../common"
	"../utils"
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	ch = make(chan int)
	now0 int64
	now1 int64
	time1 *time.Timer
)
func RUDP(c *net.UDPConn) {
	var rtt map[int]int64
	rtt = make(map[int]int64)
	var rto int64
	//数据传输发送次数，不包括重传
	sendDataCount := 1
	//数据传输发送次数，包括重传次数
	allSendDataCount := 1
	//包位置
	index := 0

	//网络传输阈值
	ssthresh := common.SSTHRESH
	swnd := common.SWND
	//从终端获取数据（键盘输入）
	//input := bufio.NewReader(os.Stdin)
	for {
		//dataStr, _ := input.ReadString('\n')
		var dataStr string
		var allData = make(map[int]string)
		//当此次输入的所有数据都传输完成才进行下一轮
		if len(allData) == 0 {
			index = 0
			fmt.Scanln(&dataStr)
			allData = utils.SplitByPackSize(utils.InitData(dataStr))
			common.Log.Info("传输数据已分包:", allData)
		}

		//本次所要发送的数据
		sendData := make(map[int]string)
		for ; allData != nil && len(allData) != 0; {
			//如果sendData不为空说明上一次发送给服务端的数据没有得到反馈，需要重传
			if len(sendData) == 0 {
				//依据窗口大小获取本次发送的数据
				sendData = utils.GetSlice(allData, index, swnd)
				common.Log.Info("本次发送的数据包：", sendData)
				index += swnd
				if swnd > ssthresh {
					swnd += 1
				} else {
					swnd *= 2
				}
			} else {
				//发生拥塞发送窗口大小从1开始
				swnd = 1
				ssthresh /= 2
			}
			allSendDataCount += 1
			jsonData, _ := json.Marshal(sendData)
			common.Log.Info("Marshal map to json ", string(jsonData))
			var err error
			//向服务器发送数据
			_, err = c.Write(jsonData)
			if err != nil {
				common.Log.Error("向服务端发送数据出错: ", err)
				return
			} else {
				common.Log.Info("已向服务端发送数据！")
				//记录开始发送数据的时间
				now0 = time.Now().UnixNano()
				common.Log.Info("发送数据包的开始时间：", now0)

				//启动定时器
				go func() {
					common.Log.Info("定时器启动......")
					if sendDataCount == 1 {
						time1 = time.NewTimer(10000000000 * time.Nanosecond)
					} else {
						var rtoo time.Duration
						rtoo = time.Duration(rto)
						time1 = time.NewTimer((rtoo / 1e3) * time.Microsecond)
					}
					//监控通道，若通道值为1则停止定时器
					WaitChannel(ch, time1)
				}()
			}

			//接收服务器的反馈信息
			buf := make([]byte, swnd)
			l, _, err := c.ReadFromUDP(buf) //从服务器中取出反馈信息
			if err != nil {
				common.Log.Error("服务器反馈信息接收出错：%v\n", err)
				return
			} else {
				common.Log.Info("已收到服务器反馈信息！")
				//管道记录标志已收到服务器反馈消息
				ch <- 1
				go func() {
					//记录发送结束的时间
					now1 = time.Now().UnixNano()
					common.Log.Info("发送数据包的结束时间：", now1)
				}()
			}

			packageNoStr := string(buf[:l])
			packageNoBuf := strings.Split(packageNoStr, ",")
			//对反馈接收到的包从sendData清理，如果超时，发送需要重新传送的剩余数据
			for _, v := range packageNoBuf {
				common.Log.Info("发送成功数据包序列号：", v)
				int, err := strconv.Atoi(v)
				if err != nil {
					common.Log.Error("数据转换错误，<", err, ">")
					return
				}
				//查看map是否包含该元素
				_, exit := sendData[int]
				if exit {
					delete(sendData, int)
					delete(allData, int)
				}
			}

			t := now1 - now0
			common.Log.Info("本次数据传输的时间: ", t)
			rtt[allSendDataCount%5] = t

			//计算定时器的重传时间 rto
			if len(rtt) < 2 {
				rto = 10000000000
				common.Log.Info("定时器的时间: ", rto)
			} else if len(rtt) > 5 {
				rt := (rtt[0] + rtt[1] + rtt[2] + rtt[3] + rtt[4]) / int64(len(rtt))
				rto = int64(1.2 * float64(rt))
				common.Log.Info("定时器的时间: ", rto)
			} else {
				rt := (rtt[0] + rtt[1] + rtt[2] + rtt[3] + rtt[4]) / int64(len(rtt))
				rto = int64(1.4 * float64(rt))
				common.Log.Info("定时器的时间: ", rto)
			}
		}
	}
}

//定时器
func WaitChannel(conn <-chan int, timer *time.Timer)  {
	select {
	case <- conn:
		timer.Stop()
		common.Log.Info("数据传输成功，定时器停止！ch：",  <- conn)
	case <- timer.C: // 超时
		common.Log.Info("数据超时，等待重传!")
	}
}


