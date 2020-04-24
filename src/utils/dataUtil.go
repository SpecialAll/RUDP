package utils
/**
全局常量定义
 */
import (
	"../common"
)
/**
初始化数据
 */
func InitData(dataStr string) []byte {
	return []byte(dataStr)
}

/**
按照数据包大小548将字节数组分割存储到map映射
 */
func SplitByPackSize(dataBytes []byte) map[int]string {
	dataMap := make(map[int]string)
	//向上取整
	len := len(dataBytes)
	num := len / common.PackageSize
	if len%common.PackageSize != 0 {
		num++
	}
	if num <= 1 {
		dataMap[0] = string(dataBytes)
	} else {
		for i := 1; i < num; i++ {
			//解释一下这个问题
			dataMap[i-1] = string(dataBytes[(i-1) * common.PackageSize : i * common.PackageSize - 1 : len - 1])
		}
		dataMap[num-1] = string(dataBytes[(num-1) * common.PackageSize : (len - 1) : (len - 1)])
	}
	return dataMap
}

/**
获取map部分切片
dataMap : 所有数据map映射表
start : 从dataMap开始获取的位置
swnd ： 发送窗口大小
 */
func GetSlice(dataMap map[int]string, start int ,swnd int ) map[int]string {
	sliceMap := make(map[int]string)
	for i := 0; i < swnd; i++ {
		//若当前数据小于发送窗口大小则跳过空白内容
		if dataMap[start + i] != "" {
			sliceMap[start + i] = dataMap[start + i]
		}
	}
	return sliceMap
}


