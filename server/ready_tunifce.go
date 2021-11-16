package server

import(
	"fmt"
	water_wrap "github.com/pigfall/tzzGoUtil/net/water_tun_wrap"
	log "github.com/pigfall/tzzGoUtil/log/golog"
tzNet "github.com/pigfall/tzzGoUtil/net"
)


func readyTunIfce(specTunIpNet *tzNet.IpWithMask)(tunIfce tzNet.TunIfce,tunIpNetRealUsed *tzNet.IpWithMask,err error){
	// { collect all net ip, and select the ip which not conflict with been used
	allIpV4Addrs,err :=tzNet.ListIpV4Addrs()
	if err != nil{
		log.Error(err.Error())
		return nil, nil, err
	}

	if specTunIpNet != nil{
		tunIpNetRealUsed = specTunIpNet
	}else{
		tunCidr,err := findSuitableIp(allIpV4Addrs)
		if err != nil{
			log.Error(err.Error())
			return nil, nil, err
		}
		tunIpNetRealUsed = tunCidr
	}
	log.Info(fmt.Sprintf("tun cidr %s",tunIpNetRealUsed.String()))
	// }

	// { create tun ifce 
	tunIfce,err = water_wrap.NewTun()
	if err != nil{
		log.Error(err.Error())
		return  nil, nil, err
	}
	// }

	// { set ip to tun ifce and enable it
	err = tunIfce.SetIp(tunIpNetRealUsed.String())
	if err != nil{
		err = fmt.Errorf("Set ip to tun interface failed: %w",err)
		log.Error(err.Error())
		return nil, nil, err
	}
	// }
	return tunIfce,tunIpNetRealUsed,nil
}



func findSuitableIp(allIpV4IsUsed []tzNet.IpWithMask)(*tzNet.IpWithMask,error){
	var subnet = 8
	var subnet2 = 0
	encodeIpNet := func(subNet2,subNet int)string{
		return fmt.Sprintf("10.%d.%d.1/16",subNet2,subNet)
	}
	OUT:
	for{
		subnet2++
		if subnet2 >=255{
				return nil,fmt.Errorf("Over then 255 , not found unconflict ip to tun ifce")
		}
		subnet = 1
		for{
			retIp ,err :=tzNet.FromIpSlashMask(encodeIpNet(subnet2,subnet))
			if err != nil{
				return nil,err
			}
			if tzNet.IpSubnetCoincideOrCoinCided(retIp,allIpV4IsUsed){
				subnet++
			}else{
				return retIp,nil
			}
			if subnet >= 255{
				continue OUT
			}
		}
	}



}
