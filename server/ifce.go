package server

import(
	"time"
	"io"
	"context"
	tzNet "github.com/pigfall/tzzGoUtil/net"
)


// storage all conns
type ConnsStorage interface{
	FindConnByClientIpPort(clientIpPort tzNet.IpPortFormat)(Conn)
	FindConnByTunnelIp(tunnIpnet tzNet.IpNetFormat)Conn
	ForEachConn(do func (conn Conn))
	AllocateConn(clientIpPort tzNet.IpPort,connWriter io.WriteCloser)(Conn,error)
	ReleaseConn(clientIpPort tzNet.IpPortFormat)
}

type Conn interface{
	ClientIpPort()(ClientIpPort tzNet.IpPortFormat)
	ClientTunnelIpNet() tzNet.IpNetFormat
	WriteIpPacket(ipPacket []byte)error
	UpdateHearbeat()
	GetHeartBeatTime()time.Time
	Close()
}

type TransportServerBuilder interface{
	BuildTransportServer() TransportServer
}

type TransportServer interface{
	Prepare(serverAddr *tzNet.IpPort) (cancelFuncs func(),err error)
	Serve(ctx context.Context,tunIfce tzNet.TunIfce,connsStorage  ConnsStorage)(error)
}
