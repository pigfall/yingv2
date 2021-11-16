package server

import(
	"io"
	"net"
	"fmt"
	"context"
	tzNet "github.com/pigfall/tzzGoUtil/net"
log "github.com/pigfall/tzzGoUtil/log/golog"

"github.com/pigfall/yingv2/proto"
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
			this.handleIpPakcet(msgData,tunIfce)
		}else if msgType== MSG_TYPE_APP_MSG{
			this.handleAppMsg(clientAddr,msgData,connsStorage)
		}else{
			panic(fmt.Sprintf("BUG undefined msgType %v",msgType))
		}
	}
}

func (this *transportServerUDP) handleIpPakcet(ipPacket []byte,tunIfce tzNet.TunIfce){
	_,err  := tunIfce.Write(ipPacket)
	if err != nil{
		log.Error(err.Error())
	}
}

func (this *transportServerUDP) handleAppMsg(clientAddr *net.UDPAddr,appMsgBytes []byte,connStorage ConnsStorage){
	var udpSock = this.udpSock
	var reqMsg = proto.ReqMsg{}
	err :=  proto.Decode(appMsgBytes,&reqMsg)
	if err != nil{
		log.Error(err.Error())
		return 
	}
	res :=handleUDPAppMsg(&reqMsg,connStorage,clientAddr,udpSock)
	resMsgBytes,err  := proto.Encode(res)
	if err != nil{
		log.Error(err.Error())
		return
	}
	_,err = writeToUDPClient(udpSock,resMsgBytes,clientAddr)
	if err != nil{
		log.Error(err.Error())
	}
}

func writeToUDPClient(udpSock *tzNet.UDPSock,bytes []byte,clientAddr *net.UDPAddr)(int,error){
	return udpSock.WriteToUDP(bytes,clientAddr)
}


const(
		MSG_TYPE_IP_PACKET=0
	 MSG_TYPE_APP_MSG= 1
)

type connUDPWriter struct{
	udpSock *tzNet.UDPSock
	remoteAddr *net.UDPAddr
}

func newConnUDPWriter(udpSock *tzNet.UDPSock,remoteAddr *net.UDPAddr)io.Writer{
	return &connUDPWriter{
		udpSock:udpSock,
		remoteAddr:remoteAddr,
	}
}

func (this *connUDPWriter) Write(b []byte)(int,error){
	return  writeToUDPClient(this.udpSock,b,this.remoteAddr)
}
