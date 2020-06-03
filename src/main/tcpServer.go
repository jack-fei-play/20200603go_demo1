package main

import (
	"encoding/binary"
	"fmt"
	log "github.com/cihub/seelog"
	"net"
)

type Message struct {
	MsgHeadFlag byte   //消息头部标识位
	MsgType     uint16 //消息类型（消息ID）
	MsgLength   uint16 //剩余长度
	MsgCheck    byte   //校验位
	MsgFootFlag byte   //发起方
	//MsgBody 消息体中内容
	SendPerson         byte   //发起方
	UserID             string //USERID
	AgentID            string //AGENTID
	PlcID              string //PLCID
	ProducerID         byte   //内部分配的producer_id值
	MsgBodytype        byte   //协议类型
	MsgBodyIp          string //PLC的局域网IP地址
	MsgBodyPort        uint16 //PLC的监听端口号
	MsgBodyInitPackage []byte //从软件发出的原始包体

}

func parallelVisit(conn net.Conn) {
	defer conn.Close()
	if conn == nil {
		fmt.Println("并发访问传递连接失败：nil")
	}
	//tcp长连接循环保持读取tcp客户端消息
	for {
		//切片缓冲
		buf := make([]byte, 1024)
		//每次读取读取的字节长度个数
		n, err := conn.Read(buf)
		if err != nil || n == 0 {
			fmt.Println("\r\ntcp客户端发送消息读取失败，错误：", err)
			return
		}
		message := buf[:n]
		fmt.Printf("\r\n %v", message)
		if message == nil {
			fmt.Println("message为nil")
			return
		}
		//校验位
		flag := checkMsg(message[1 : len(message)-2])
		if flag == true {
			//自定义解析tcp
			analysis(message)
		}

	}

}

//异或校验
func checkMsg(message []byte) bool {
	return true
}

//自定义tcp协议内容解析
func analysis(message []byte) {
	var msgStruct Message
	msgStruct.MsgHeadFlag = message[0]
	msgStruct.MsgFootFlag = message[len(message)-1]
	fmt.Println("\r\n---------------------------------------------")
	fmt.Printf("\r\n 消息头部标识位 %d", message[0])
	fmt.Printf("\r\n 消息尾部标识位 %d", message[len(message)-1])
	//消息类型
	msgType := binary.LittleEndian.Uint16(message[1:3])
	msgStruct.MsgType = msgType
	fmt.Printf("\r\n 消息类型:0x%x", uint16(msgType))
	//消息长度
	msgLength := binary.LittleEndian.Uint16(message[3:5])
	msgStruct.MsgLength = msgLength
	fmt.Printf("\r\n 消息剩余长度: %d", uint16(msgLength))
	//检验位
	msgCheck := message[len(message)-2]
	msgStruct.MsgCheck = msgCheck
	fmt.Printf("\r\n 检验位: %d", msgCheck)

	//消息体内容
	msgBody := message[5 : len(message)-2]
	//发起方 1 Client发起，2 Agent发起
	sendPerson := msgBody[0]
	msgStruct.SendPerson = sendPerson
	fmt.Printf("\r\n 发起方: %d", sendPerson)
	//注册用户的USERID
	userID := string(msgBody[1:33])
	msgStruct.UserID = userID
	fmt.Printf("\r\n USERID: %s", userID)
	//内部分配的AGENTID
	agentID := string(msgBody[33:65])
	msgStruct.AgentID = agentID
	fmt.Printf("\r\n AGENTID: %s", agentID)
	//内部分配的PLCID
	plcID := string(msgBody[65:97])
	msgStruct.PlcID = plcID
	fmt.Printf("\r\n PLCID: %s", plcID)
	//内部分配的producer_id值
	producerID := msgBody[97:98][0]
	msgStruct.ProducerID = producerID
	fmt.Printf("\r\n 内部分配的producer_id值: %d", producerID)
	//协议类型 1TCP，2UDP
	msgBodytype := msgBody[98:99][0]
	msgStruct.MsgBodytype = msgBodytype
	fmt.Printf("\r\n 协议类型: %d", msgBodytype)
	//PLC的局域网IP地址，如：131.25.180.112
	msgBodyIp := string(msgBody[99:114])
	msgStruct.MsgBodyIp = msgBodyIp
	fmt.Printf("\r\n PLC的局域网IP地址: %s", msgBodyIp)
	//PLC的监听端口号，如：80012
	msgBodyPort := binary.LittleEndian.Uint16(msgBody[94:116])
	msgStruct.MsgBodyPort = msgBodyPort
	fmt.Printf("\r\n PLC的监听端口号: %d", msgBodyPort)
	//从软件发出的原始包体
	i := len(msgBody) - 1
	msgBodyInitPackage := msgBody[116:i]
	msgStruct.MsgBodyInitPackage = msgBodyInitPackage
	fmt.Printf("\r\n 从软件发出的原始包体: %d", msgBodyInitPackage)
	fmt.Println("\r\n---------------------------------------------")
	fmt.Printf("\r\n%+v\n", msgStruct)
}

func main() {
	//初始化日志
	defer log.Flush()
	logger, err := log.LoggerFromConfigAsFile("seelog.xml")
	if err != nil {
		log.Error(err)
		return
	}
	log.ReplaceLogger(logger)
	//启动tcp服务监听本地8088端口
	listen, err := net.Listen("tcp", "127.0.0.1:8088")
	if err != nil {
		fmt.Println("启动监听8088端口失败，错误:", err)
	}
	log.Info("socket服务8088启动成功!")
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
