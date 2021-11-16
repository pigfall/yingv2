package server

import(
	"io"
	"context"
	tzNet "github.com/pigfall/tzzGoUtil/net"
)


// storage all conns
type ConnsStorage interface{
	FindConnByClientIpPort(clientIpPort tzNet.IpPortFormat)(Conn)
	FindConnByTunnelIp(tunnIpnet tzNet.IpNetFormat)Conn
	ForEachConn(do func (conn Conn))
	AllocateConn(clientIpPort tzNet.IpPort,connWriter io.Writer)(Conn,error)
}

type Conn interface{
	ClientIpPort()(ClientIpPort tzNet.IpPortFormat)
	ClientTunnelIpNet() tzNet.IpNetFormat
	WriteIpPacket(ipPacket []byte)error
}

type TransportServerBuilder interface{
	BuildTransportServer() TransportServer
}

type TransportServer interface{
	Prepare(serverAddr *tzNet.IpPort) (cancelFuncs func(),err error)
	Serve(ctx context.Context,tunIfce tzNet.TunIfce,connsStorage  ConnsStorage)(error)
}
