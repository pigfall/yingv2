package server

import(
	"io"
	"sync"
tzNet "github.com/pigfall/tzzGoUtil/net"
log "github.com/pigfall/tzzGoUtil/log/golog"
)


type connsStorage struct{
	writeLock sync.Mutex
	clientIpPort_Conn map[tzNet.IpPortFormat]Conn
	clientTunnelIpNet_Conn map[tzNet.IpNetFormat]Conn
	ipPool *tzNet.IpPool
}

func newConnStorage(ipPool *tzNet.IpPool)ConnsStorage{
	return &connsStorage{
		ipPool :ipPool,
		clientIpPort_Conn:make(map[tzNet.IpPortFormat]Conn),
		clientTunnelIpNet_Conn:make(map[tzNet.IpNetFormat]Conn) ,
	}
}



func(this *connsStorage)	FindConnByClientIpPort(clientIpPort tzNet.IpPortFormat)(Conn){
	return this.clientIpPort_Conn[clientIpPort]
}

func(this *connsStorage)	FindConnByTunnelIp(tunnIpNet tzNet.IpNetFormat)Conn{
	return this.clientTunnelIpNet_Conn[tunnIpNet]
}

func(this *connsStorage)	ForEachConn(do func(conn Conn)){
	for _,conn  := range this.clientIpPort_Conn{
		do(conn)
	}
}

func (this *connsStorage) AllocateConn(clientIpPort tzNet.IpPort,connWriter io.WriteCloser)(Conn,error){
	this.writeLock.Lock()
	defer this.writeLock.Unlock()
	var conn Conn
	conn = this.clientIpPort_Conn[clientIpPort.ToIpPortFormat()]
	if conn == nil {
		ipNet,err := this.ipPool.Take()
		if err != nil{
			log.Error(err.Error())
			return nil,err
		}
		conn = NewConn(clientIpPort,ipNet,connWriter)
		conn.UpdateHearbeat()
		this.clientIpPort_Conn[clientIpPort.ToIpPortFormat()] = conn
		this.clientTunnelIpNet_Conn[ipNet.ToIpNetFormat()]= conn
	}
	return conn,nil
}

func (this *connsStorage) ReleaseConn(clientIpPort tzNet.IpPortFormat){
	this.writeLock.Lock()
	defer this.writeLock.Unlock()
	conn := this.clientIpPort_Conn[clientIpPort]
	if conn == nil {
		return
	}
	conn.Close()
	delete( this.clientIpPort_Conn,clientIpPort)
}
