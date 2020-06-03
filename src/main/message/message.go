package message

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
