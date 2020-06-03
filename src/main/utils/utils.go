package utils

import (
	"encoding/binary"
	log "github.com/cihub/seelog"
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

func Xor(a byte, b byte) byte {
	return a ^ b
}

//自定义tcp协议内容解析
func Analysis(message []byte) {
	var msgStruct Message
	msgStruct.MsgHeadFlag = message[0]
	msgStruct.MsgFootFlag = message[len(message)-1]
	log.Infof(" 消息头部标识位 %d", message[0])
	log.Infof(" 消息尾部标识位 %d", message[len(message)-1])
	//消息类型
	msgType := binary.LittleEndian.Uint16(message[1:3])
	msgStruct.MsgType = msgType
	log.Infof(" 消息类型:0x%x", uint16(msgType))
	//消息长度
	msgLength := binary.LittleEndian.Uint16(message[3:5])
	msgStruct.MsgLength = msgLength
	log.Infof(" 消息剩余长度: %d", uint16(msgLength))
	//检验位
	msgCheck := message[len(message)-2]
	msgStruct.MsgCheck = msgCheck
	log.Infof(" 检验位: %d", msgCheck)

	//消息体内容
	msgBody := message[5 : len(message)-2]
	//发起方 1 Client发起，2 Agent发起
	sendPerson := msgBody[0]
	msgStruct.SendPerson = sendPerson
	log.Infof(" 发起方: %d", sendPerson)
	//注册用户的USERID
	userID := string(msgBody[1:33])
	msgStruct.UserID = userID
	log.Infof(" USERID: %s", userID)
	//内部分配的AGENTID
	agentID := string(msgBody[33:65])
	msgStruct.AgentID = agentID
	log.Infof(" AGENTID: %s", agentID)
	//内部分配的PLCID
	plcID := string(msgBody[65:97])
	msgStruct.PlcID = plcID
	log.Infof(" PLCID: %s", plcID)
	//内部分配的producer_id值
	producerID := msgBody[97:98][0]
	msgStruct.ProducerID = producerID
	log.Infof(" 内部分配的producer_id值: %d", producerID)
	//协议类型 1TCP，2UDP
	msgBodytype := msgBody[98:99][0]
	msgStruct.MsgBodytype = msgBodytype
	log.Infof(" 协议类型: %d", msgBodytype)
	//PLC的局域网IP地址，如：131.25.180.112
	msgBodyIp := string(msgBody[99:114])
	msgStruct.MsgBodyIp = msgBodyIp
	log.Infof(" PLC的局域网IP地址: %s", msgBodyIp)
	//PLC的监听端口号，如：80012
	msgBodyPort := binary.LittleEndian.Uint16(msgBody[94:116])
	msgStruct.MsgBodyPort = msgBodyPort
	log.Infof(" PLC的监听端口号: %d", msgBodyPort)
	//从软件发出的原始包体
	i := len(msgBody) - 1
	msgBodyInitPackage := msgBody[116:i]
	msgStruct.MsgBodyInitPackage = msgBodyInitPackage
	log.Infof(" 从软件发出的原始包体: %d", msgBodyInitPackage)
	log.Infof("%+v\n", msgStruct)
	//log.Info(msgStruct)
}
