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
	return this.clientIpPort.ToIpPortFormat()
}
func (this *conn)	ClientTunnelIpNet() tzNet.IpNetFormat{
	return this.clientTunnelIpNet.ToIpNetFormat()
}
func (this *conn)	WriteIpPacket(ipPacket []byte)error{
	_,err := this.writer.Write(ipPacket)
	return err
}
