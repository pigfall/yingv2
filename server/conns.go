package server

import(
	"time"
		"io"
tzNet "github.com/pigfall/tzzGoUtil/net"
)

type conn struct{
	writer io.WriteCloser
	clientTunnelIpNet *tzNet.IpWithMask
	clientIpPort tzNet.IpPort
	hearbeat time.Time
}


func NewConn(clientIpPort tzNet.IpPort,clientTunnelIpNet *tzNet.IpWithMask,writer io.WriteCloser)Conn{
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

func (this *conn) UpdateHearbeat(){
	this.hearbeat = time.Now()
}

func (this *conn)GetHeartBeatTime()time.Time{
	return this.hearbeat
}

func (this *conn) Close(){
	this.writer.Close()
}
