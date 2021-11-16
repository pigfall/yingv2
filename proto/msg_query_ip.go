package proto

import(
	tzNet "github.com/pigfall/tzzGoUtil/net"
)


type MsgS2CQueryIp struct{
	IpNet tzNet.IpNetFormat `json:"ip_net"`
}
