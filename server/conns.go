package server

import(
		"io"
tzNet "github.com/pigfall/tzzGoUtil/net"
)

type conn struct{
	writer io.Writer
	clientTunnelIpNet *tzNet.IpWithMask
	clientIpPort tzNet.IpPort
}


func NewConn(clientIpPort tzNet.IpPort,clientTunnelIpNet *tzNet.IpWithMask,writer io.Writer)Conn{
	return &conn{
		writer:writer,
		clientTunnelIpNet:clientTunnelIpNet,
		clientIpPort:clientIpPort,
	}
}


func (this *conn)	ClientIpPort()(ClientIpPort tzNet.IpPortFormat){
	panic("TODO")
}
func (this *conn)	ClientTunnelIpNet() tzNet.IpNetFormat{
	panic("TODO")
}
func (this *conn)	WriteIpPacket(ipPacket []byte)error{
	panic("TODO")
}
