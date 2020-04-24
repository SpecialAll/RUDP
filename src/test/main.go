package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

func main()  {
	packageNoBuf := strings.Split("1,2,3" , ",")
	fmt.Println(packageNoBuf)

	for index , v := range packageNoBuf {
		fmt.Println("index:c",index, "value: ", v)
	}

	i, err := strconv.Atoi("0")
	fmt.Println(i,err)

	fmt.Println("《你》字节数：", len([]byte("你")))
	fmt.Println("《你好》字节数：", len([]byte("你好")))
	fmt.Println("《a》字节数：", len([]byte("a")))
	fmt.Println("《abc》字节数：", len([]byte("abc")))
	fmt.Println("《^》字节数：", len([]byte("^")))

	fmt.Println(len(string(0)))


	sendData := make(map[int]string)
	receData := make(map[int]string)
	sendData[0] = "sds"
	sendData[1] = "jiio"
	jsonData, _ := json.Marshal(sendData)
	fmt.Println(string(jsonData))
	json.Unmarshal(jsonData, &receData)
	fmt.Println("json to map ", receData)
	fmt.Print(receData[1])



}
