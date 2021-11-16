package proto


type MsgBase struct{
	Id string `json:"id"`
}

type ReqMsg struct{
	MsgBase
}

type ResMsg struct{
	MsgBase
}

func Encode(msg interface{})([]byte,error){
	panic("TODO")
}

func Decode(bytes []byte,msg interface{})error{
	panic("TODO")
}


