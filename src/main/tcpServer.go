package main

import (
	"fmt"
	"net"
)

func parallelVisit(conn net.Conn) {
	if conn == nil {
		fmt.Println("并发访问传递连接失败：nil")
	}
	//tcp长连接循环保持读取tcp客户端消息
	for {
		//切片缓冲
		buf := make([]byte, 4096)
		//每次读取读取的字节长度个数
		n, err := conn.Read(buf)
		if err != nil || n == 0 {
			fmt.Println("tcp客户端发送消息读取失败，错误：", err)
			break
		}
		fmt.Println(buf[:n])
	}

}

func main() {
	//启动tcp服务监听本地8088端口
	listen, err := net.Listen("tcp", "127.0.0.1:8088")
	if err != nil {
		fmt.Println("启动监听8088端口失败，错误:", err)
	}
	for {
		//获取空闲连接
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("获取空闲连接失败，错误:", err)
		}
		//开启多协程
		go parallelVisit(conn)
	}

}
