package proto

const(
	ID_C2S_QUERY_IP_NET=1
	ID_S2C_QUERY_IP_NET=2

)


type MsgBase struct{
	Id int `json:"id"`
}

type ReqMsg struct{
	MsgBase
}

type ResMsg struct{
	MsgBase
	ErrReason string `json:"err_reason"` 
	Body interface{} `json:"body"`
}

func Encode(msg interface{})([]byte,error){
	panic("TODO")
}

func Decode(bytes []byte,msg interface{})error{
	panic("TODO")
}


