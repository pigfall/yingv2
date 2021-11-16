package server

import(
	"net"
	"fmt"
	"context"
	tzNet "github.com/pigfall/tzzGoUtil/net"
log "github.com/pigfall/tzzGoUtil/log/golog"
)

type TransportUDPServerBuilder struct{

}

func NewTransportUDPServerBuilder()TransportServerBuilder{
	return &TransportUDPServerBuilder{}
}

func (this *TransportUDPServerBuilder) BuildTransportServer() TransportServer{
	return &transportServerUDP{}

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
	var udpSock = this.udpSock
	var buf = make([]byte,1024*70)
	for {
		readNum,clientAddr,err := udpSock.ReadFromUDP(buf)
		if err != nil{
			err = fmt.Errorf("Read from udp socket failed %w",err)
			log.Error(err.Error())
			return err
		}
		readData := buf[:readNum]
		msgType := readData[0]
		msgData:=readData[1:]
		if msgType == MSG_TYPE_IP_PACKET{
			this.handleIpPakcet(msgData,clientAddr,connsStorage)
		}else if msgType== MSG_TYPE_APP_MSG{
			this.handleAppMsg(ctx,msgData)
		}else{
			panic(fmt.Sprintf("BUG undefined msgType %v",msgType))
		}
	}
}

func (this *transportServerUDP) handleIpPakcet(ipPacket []byte,clientAddr *net.UDPAddr,connStorage ConnsStorage){
	// TODO route by clientTunnelIp
	allConns := connStorage.AllConns()
	for _,conn  := range allConns{
		err := conn.WriteIpPacket(ipPacket)
		if err != nil{
			log.Error(err.Error())
		}
	}
}

func (this *transportServerUDP) handleAppMsg(ctx context.Context,appMsg []byte){
	panic("TODO")
}


const(
		MSG_TYPE_IP_PACKET=0
 MSG_TYPE_APP_MSG= 1
)
