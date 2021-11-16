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
	ipPool *tzNet.IpPool
}

func newConnStorage(ipPool *tzNet.IpPool)ConnsStorage{
	return &connsStorage{
		ipPool :ipPool,
		clientIpPort_Conn:make(map[tzNet.IpPortFormat]Conn),
	}
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

func (this *connsStorage) AllocateConn(clientIpPort tzNet.IpPort,connWriter io.Writer)(Conn,error){
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
		this.clientIpPort_Conn[clientIpPort.ToIpPortFormat()] = conn
	}
	return conn,nil
}
