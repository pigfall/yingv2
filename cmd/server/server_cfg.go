package main

type RawServerCfg struct{
	AddrIp string `json:"addr_ip"`
	AddrPort int `json:"addr_port"`
	Mode string `json:"mode"`
}
