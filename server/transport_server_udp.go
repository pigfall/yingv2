package server

import(
	"fmt"
	"context"
	tzNet "github.com/pigfall/tzzGoUtil/net"
log "github.com/pigfall/tzzGoUtil/log/golog"
)

type TransportUDPServerBuilder struct{


}


type transportServerUDP struct{
	udpSock *tzNet.UDPSock
}

func (this *transportServerUDP) Prepare(serverAddr *tzNet.IpPort)(cancelFunc func(),err error){
	// TODO FROM HERE
	log.Info(fmt.Sprintf("Prepare to listent udpServer at %s",serverAddr.ToString()))
	udpSock,err := tzNet.UDPListen(serverAddr.IP,serverAddr.Port)
	if err != nil{
		err = fmt.Errorf("Failed to listent udp server at %s %w",serverAddr.ToString(),err)
		log.Error(err.Error())
		return nil,err
	}

	this.udpSock = udpSock
	return func(){
		udpSock.Close()
	},nil
}

/* 
jobs:
   udpReadLoop
*/
func (this *transportServerUDP) Serve(ctx context.Context,tunIfce tzNet.TunIfce,connsStorage ConnsStorage)error{
	panic("TODOj")
}

