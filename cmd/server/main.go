package main

import (
	"fmt"
	"net"
	"flag"
	"encoding/json"
tzNet "github.com/pigfall/tzzGoUtil/net"
	"os"
 log "github.com/pigfall/tzzGoUtil/log/golog"
 "github.com/pigfall/tzzGoUtil/encoding"
 "context"
	"github.com/pigfall/yingv2/server"
)


func main() {
	var confPath string
	var outputDemoConfig bool
	flag.StringVar(&confPath,"confPath","","config path")
	flag.BoolVar(&outputDemoConfig,"demoConfig",false,"output demo config")
	flag.Parse()
	if outputDemoConfig{
		rawCfg := buildDemoConfig()
		bytes,err := json.Marshal(rawCfg)
		if err != nil{
			panic(err)
		}
		fmt.Println(string(bytes))
		os.Exit(0)
	}
	if len(confPath) == 0{
		log.Error("config file path is nil ")
		flag.Usage()
		os.Exit(1)
	}
	rawSvrCfg:=&RawServerCfg{}
	err := encoding.UnMarshalByFile(
		confPath,
		rawSvrCfg,
		json.Unmarshal,
	)
	if err != nil{
		err = fmt.Errorf("parse config file %s error %w",confPath,err)
		log.Error(err.Error())
		os.Exit(1)
	}
	ctx := context.Background()

	svrCfg,err := fromRawServerCfg(rawSvrCfg)
	if err != nil{
		log.Error(err.Error())
		os.Exit(1)
	}

	var tp server.TransportServerBuilder
	switch svrCfg.TransportType{
	case TPUDP:
		tp = server.NewTransportUDPServerBuilder() 
	case TPWebsocket:
		panic("TODO")
	default:
		panic("BUG unreachable")
	}

	server.Serve(
		ctx,
		&svrCfg.IpPort,
		tp,
	)
}


const (
		TPUDP TransportType ="udp"
		TPWebsocket TransportType="ws"
)

type TransportType string

type ServerCfg struct{
	IpPort tzNet.IpPort
	TransportType TransportType
}

func fromRawServerCfg(rawCfg *RawServerCfg)(*ServerCfg,error){
	cfgErr := make([]error,0)
	if len(rawCfg.AddrIp) == 0{
		cfgErr = append(cfgErr,fmt.Errorf("addrIp cannot be  nil"))
	}
	if rawCfg.AddrPort == 0 {
		cfgErr = append(cfgErr,fmt.Errorf("addrPort cannot be  nil"))
	}
	var tpMode TransportType
	if len(rawCfg.Mode) == 0{
		cfgErr = append(cfgErr,fmt.Errorf("transport Mode cannot be nil"))
	}else{
		switch rawCfg.Mode{
		case "udp":
			tpMode = TPUDP
		case "ws":
			tpMode = TPWebsocket
		default:
		cfgErr = append(cfgErr,fmt.Errorf("undefined transport mode %s",rawCfg.Mode))
		}
	}
	if len(cfgErr ) >  0{
		return nil,fmt.Errorf("Cfg error %v",cfgErr)
	}

	ipToListen := net.ParseIP(rawCfg.AddrIp)
	if ipToListen == nil{
		cfgErr = append(cfgErr,fmt.Errorf("Invalid addrIp %s",rawCfg.AddrIp)) 
	}

	if len(cfgErr ) >  0{
		return nil,fmt.Errorf("Cfg error %v",cfgErr)
	}
	svrCfg := &ServerCfg{
		IpPort:tzNet.IpPort{
			IP:ipToListen,
			Port:rawCfg.AddrPort,
		},
		TransportType:tpMode,
	}

	return svrCfg,nil
}
