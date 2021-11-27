package proto

import(
		"encoding/json"
)

const(
	ID_C2S_QUERY_IP_NET=1
	ID_S2C_QUERY_IP_NET=2
	ID_C2S_HEARTBEAT=3
	ID_S2C_HEARTBEAT=4
)


type MsgBase struct{
	Id int `json:"id"`
}

type ReqMsg struct{
	MsgBase
	Body []byte `json:"body"`
}

type ResMsg struct{
	MsgBase
	ErrReason string `json:"err_reason"` 
	Body string `json:"body"`
}

func Encode(msg interface{})([]byte,error){
	return json.Marshal(msg)
}

func Decode(bytes []byte,msg interface{})error{
	return json.Unmarshal(bytes,msg)
}


