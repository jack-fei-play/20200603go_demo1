package main

import (
	"../main/utils"
	log "github.com/cihub/seelog"
	"net"
)

//type Message struct {
//	MsgHeadFlag byte   //消息头部标识位
//	MsgType     uint16 //消息类型（消息ID）
//	MsgLength   uint16 //剩余长度
//	MsgCheck    byte   //校验位
//	MsgFootFlag byte   //发起方
//	//MsgBody 消息体中内容
//	SendPerson         byte   //发起方
//	UserID             string //USERID
//	AgentID            string //AGENTID
//	PlcID              string //PLCID
//	ProducerID         byte   //内部分配的producer_id值
//	MsgBodytype        byte   //协议类型
//	MsgBodyIp          string //PLC的局域网IP地址
//	MsgBodyPort        uint16 //PLC的监听端口号
//	MsgBodyInitPackage []byte //从软件发出的原始包体
//
//}

func parallelVisit(conn net.Conn) {
	defer conn.Close()
	if conn == nil {
		log.Error("并发访问传递连接失败：nil")
	}
	//tcp长连接循环保持读取tcp客户端消息
	for {
		//切片缓冲
		buf := make([]byte, 1024)
		//每次读取读取的字节长度个数
		n, err := conn.Read(buf)
		if err != nil || n == 0 {
			log.Error("tcp客户端发送消息读取失败，错误：", err)
			return
		}
		message := buf[:n]
		//fmt.Printf("\r\n %v", message)
		if message == nil {
			log.Error("message为nil")
			return
		}
		//校验位
		flag := checkMsg(message)
		if flag == true {
			//自定义解析tcp
			utils.Analysis(message)
		}

	}

}

//异或校验
func checkMsg(message []byte) bool {
	//for i := 1; i < len(message)-3; i++ {
	//	message[i+1] = utils.Xor(message[i],message[i+1])
	//}
	//log.Info("校验值:", message[len(message)-3])
	//if message[len(message)-3] != message[len(message)-2] {
	//	log.Error("异或校验出错")
	//	return false
	//}
	temp := message[1]
	for i := 2; i < len(message)-2; i++ {
		temp = utils.Xor(temp, message[i])
	}
	//log.Info("校验值:", temp)
	if temp != message[len(message)-2] {
		log.Error("异或校验出错")
		return false
	}
	return true
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
		log.Error("启动监听8088端口失败，错误:", err)
	}
	log.Info("socket服务8088启动成功!")
	for {
		//获取空闲连接
		conn, err := listen.Accept()
		if err != nil {
			log.Error("获取空闲连接失败，错误:", err)
		}
		//开启多协程
		go parallelVisit(conn)
	}

}
