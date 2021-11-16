package server

import(
tzNet "github.com/pigfall/tzzGoUtil/net"
)


type connsStorage struct{}

func newConnStorage()ConnsStorage{
	return &connsStorage{}
}



func(this *connsStorage)	PutConn(Conn){
	panic("TODO")
}
func(this *connsStorage)	FindConnByClientIpPort(clientIpPort tzNet.IpPortFormat)(Conn){
	panic("TODO")
}
func(this *connsStorage)	FindConnByTunnelIp(tunnIp tzNet.IpFormat)Conn{
	panic("TODO")
}
func(this *connsStorage)	AllConns()[]Conn{
	panic("TODO")
}

func (this *connsStorage) AllocateConn(clientIpPort tzNet.IpPort)(Conn,error){
	panic("TODO")
}
