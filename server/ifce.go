package server

tzNet "github.com/pigfall/tzzGoUtil/net"

// storage all conns
type ConnsStoreage interface{
	PutConn(Conn)
	FindConnByClientIpPort(clientIpPort )(Conn)
	FindConnByTunnelIp(tunnIp)Conn
	AllConns()[]Conn
}

type Conn interface{
	ClientIpPort() ClientIpPort
	ClientTunnelIp() IpFormat
}

type TransportServerBuilder interface{
	BuildTransportServer() TranportServer
}
