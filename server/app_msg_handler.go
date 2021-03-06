package server

import(
	"io"
	"encoding/json"
	"fmt"

	"net"

tzNet "github.com/pigfall/tzzGoUtil/net"
"github.com/pigfall/yingv2/proto"
log "github.com/pigfall/tzzGoUtil/log/golog"
)


func handleUDPAppMsg(req *proto.ReqMsg,connsStorage ConnsStorage,clientAddr *net.UDPAddr,udpSock *tzNet.UDPSock)(res *proto.ResMsg){
	return handleAppMsg(req,connsStorage,tzNet.IpPort{IP:clientAddr.IP,Port:clientAddr.Port},newConnUDPWriter(udpSock,clientAddr))
}


func handleAppMsg(req *proto.ReqMsg,connStorage ConnsStorage,clientIpPort tzNet.IpPort,connWriter io.WriteCloser)(res *proto.ResMsg){
	res = &proto.ResMsg{}
	var protoBodyMsgIfce interface{}
	switch req.Id{
	case proto.ID_C2S_QUERY_IP_NET:
		res.Id = proto.ID_S2C_QUERY_IP_NET
		conn,err := connStorage.AllocateConn(clientIpPort,connWriter)
		if err != nil{
			log.Error(err.Error())
			res.ErrReason = err.Error()
			return res
		}
		log.Info(fmt.Sprintf("assign clientTunelIp  %s to client %s",conn.ClientTunnelIpNet(),clientIpPort.ToString()))
		protoBodyMsgIfce = &proto.MsgS2CQueryIp{
			IpNet:conn.ClientTunnelIpNet(),
		}
	case proto.ID_C2S_HEARTBEAT:
		res.Id = proto.ID_S2C_HEARTBEAT
		conn := connStorage.FindConnByClientIpPort(clientIpPort.ToIpPortFormat())
		if conn!=nil{
			conn.UpdateHearbeat()
		}else{
			// not found connection ,no reponse
			return nil
		}
	default:
		err := fmt.Errorf("Undefined appMsgReqId %v",req.Id)
		log.Error(err.Error())
		panic(err)
	}

	var bodyBytes []byte
	if protoBodyMsgIfce != nil{
		var err error
		bodyBytes,err = json.Marshal(protoBodyMsgIfce)
		if err != nil{
			log.Error(err.Error())
			return nil
		}
	}


	res.Body =string(bodyBytes)
	return res
}
