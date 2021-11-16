package server

import(
	tzNet "github.com/pigfall/tzzGoUtil/net"
)

type TransportUDPServerBuilder struct{


}


type transportServerUDP struct{}

/* 
jobs:
   udpReadLoop
*/
func (this *transportServerUDP) Serve(ctx context.Context,serveAddr ,tunIfce,connsStorage){
	// TODO FROM HERE
	udpSock,err := tzNet.UDPListen(serverAddr.IP,serverAddr.Port)
	if err != nil
}

