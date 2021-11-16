package server

import(
	"io"
	"fmt"
log "github.com/pigfall/tzzGoUtil/log/golog"
	ws "github.com/gorilla/websocket"
	"context"
	"net"
"github.com/pigfall/yingv2/proto"
	"net/http"
	tzNet "github.com/pigfall/tzzGoUtil/net"
)


type TransportServerWebsocketBuilder struct{}


func NewTransportServerWebsocketBuilder() TransportServerBuilder{
	return &TransportServerWebsocketBuilder{}
}

func (this *TransportServerWebsocketBuilder) BuildTransportServer()  TransportServer{
	return &transportServerWebsocket{}
}

type transportServerWebsocket struct{
	l net.Listener
}

func (this *transportServerWebsocket) Prepare(serverAddr *tzNet.IpPort)(cancelFunc func(),err error){
	log.Info(fmt.Sprintf("Prepare to listent websocketServer at %s",serverAddr.ToString()))
	l,err := net.Listen("tcp",string(serverAddr.ToIpPortFormat()))
	if err != nil{
		err = fmt.Errorf("Failed to listent websocket server at %s %w",serverAddr.ToString(),err)
		log.Error(err.Error())
		return nil, err
	}
	this.l = l
	return func(){
		l.Close()
	},nil
}

func (this *transportServerWebsocket) Serve(ctx context.Context,tunIfce tzNet.TunIfce,connsStorage ConnsStorage)error{
	httpSvr := http.NewServeMux()
	httpSvr.HandleFunc(
		"/",
		func (w http.ResponseWriter,req *http.Request){
			up:=ws.Upgrader{}
			wsConn,err := up.Upgrade(w,req,nil)
			if err != nil{
				log.Error(err.Error())
				return
			}
			defer wsConn.Close()
			this.serveConn(ctx,wsConn,tunIfce,connsStorage)
		},
	)
	err := http.Serve(this.l,httpSvr)
	if err != nil{
		err = fmt.Errorf("HttpServer over %w",err)
		log.Error(err.Error())
	}

	return nil
}

func (this *transportServerWebsocket) serveConn(ctx context.Context,wsConn *ws.Conn,tunIfce tzNet.TunIfce,connsStorage ConnsStorage)error{
	ctx,cancel := context.WithCancel(ctx)
	defer cancel()
	go func(){
		<-ctx.Done()
		wsConn.Close()
	}()
	for {
		msgType,bytes,err := wsConn.ReadMessage()
		if err != nil{
			log.Error(err.Error())
			return err
		}
		if msgType == ws.BinaryMessage{
			_,err = tunIfce.Write(bytes)
			if err != nil{
				log.Error(err.Error())
			}
			continue
		}

		if msgType == ws.TextMessage{
			var reqMsg = proto.ReqMsg{}
			err = proto.Decode(bytes,&reqMsg)
			if err != nil{
				log.Error(err.Error())
				continue
			}
			handleWebsocketAppMsgType(&reqMsg,wsConn,&connWebsocketWriter{conn:wsConn},connsStorage,wsConn.RemoteAddr())
			continue
		}
		log.Debug(fmt.Sprintf("get websocket msgType %v",msgType))
	}
}

func handleWebsocketAppMsgType(reqMsg *proto.ReqMsg,wsConn *ws.Conn,connWriter io.Writer,connsStorage ConnsStorage,clientAddr net.Addr){
	clientIpPort,err := tzNet.IpPortFromString(clientAddr.String())
	if err != nil{
		err = fmt.Errorf("parse clientIpPort failed %w",err)
		log.Error(err.Error())
		return 
	}
	
	res := handleAppMsg(reqMsg,connsStorage,*clientIpPort,connWriter)
	resBytes,err := proto.Encode(res)
	if err != nil{
		log.Error(err.Error())
		return 
	}
	err = writeToWebsocketClient(wsConn,ws.TextMessage,resBytes)
	if err != nil{
		log.Error(err.Error())
	}
}

func writeToWebsocketClient(wsConn *ws.Conn,msgType int,bytes []byte)error{
	return wsConn.WriteMessage(msgType,bytes)
}


type connWebsocketWriter struct{
	conn *ws.Conn
}

func (this *connWebsocketWriter) Write(b []byte)(int,error){
	err := this.conn.WriteMessage(ws.BinaryMessage,b)
	return len(b),err
}


