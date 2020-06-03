package main

import (
	"fmt"
	"net"
)

func main() {
	//tcpClient访问8088端口
	conn, err := net.Dial("tcp", "127.0.0.1:8088")
	if err != nil {
		fmt.Println("tcpClient访问8088端口失败，错误:", err)
	}
	arr := []byte{10, 24, 65, 88, 45}
	n, err := conn.Write(arr)
	if err != nil {
		fmt.Println("tcpClient写入内容失败，错误:", err)
	}
	fmt.Println("tcpClient写入字节长度为:", n)
}
