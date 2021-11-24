package main


func buildDemoConfig()*RawServerCfg{
	return &RawServerCfg{
		AddrIp:"0.0.0.0",
		AddrPort:10101,
		Mode:"udp",
	}
}
