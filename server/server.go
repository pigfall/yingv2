package server

import(
	"fmt"
"github.com/pigfall/tzzGoUtil/funcs"

"github.com/pigfall/tzzGoUtil/async"
	log "github.com/pigfall/tzzGoUtil/log/golog"
	"reflect"
)

const(
		JOB_TUNIFCE_DATA_TO_CONNS = reflect.TypeOf( tunIfceDataToConns).String()
		JOB_TRANSPORT_SERVER_SERVE = reflect.TypeOf(TransportServer.Serve).String()
)


/*
jobs:
	tunIfceDataToConns
	transportServer.Serve
*/
func Serve(ctx context.Context,transportServerBuilder)(error){
	cleanFuncs := funcs.NewFuncs()
	defer cleanFuncs.Call()
	cancelFuncs = funcs.NewFuncs()
	tunIfce,tunIpNet,err := readyTunIfce()
	if err != nil {
		log.Error(fmt.Errorf("Prepare tun interface failed"))
		return err
	}
	log.Info(fmt.Sprintf("Created tun interface %s which ipNet is %s",tunIfce.Name(),tunIpNet.ToString()))
	cleanFuncs.AddFunc(func(){tunIfce.Close()})
	cancelFuncs.AddFunc(func(){tunIfce.Close()})
	asyncCtrl := async.NewCtrl()
	// < readFrom tunIfce then send to all connections
	log.Info(fmt.Sprintf("Start %s",JOB_TUNIFCE_DATA_TO_CONNS))
	asyncCtrl.Do(
		ctx,
		JOB_TUNIFCE_DATA_TO_CONNS,
		func(ctx context.Context){
			err = tunIfceDataToConns(ctx,tunIfce,connsStorage)
			log.Info(fmt.Sprintf("%s func over %w",reflect.TypeOf( tunIfceDataToConns).String(),err))
		},
	)
	// >

	// <
	transportServer := transportServerBuilder.BuildTransportServer
	log.Info("Start transerportServer ")
	asyncCtrl.Do(
			ctx,
			JOB_TRANSPORT_SERVER_SERVE,
			func (ctx context.Context){
				err = transportServer.Serve(ctx,tunIfce,connsStorage)
				log.Info(fmt.Sprintf("%s over %w",JOB_TRANSPORT_SERVER_SERVE,err))
			}
	)
	// >


	// if ctx done, we cancel all jobs and wait jobs quit
	asyncCtrl.WaitCtx(ctx,func(cancelFuncs.Call()))
}


// readFrom tunIfce then send to all connections
func tunIfceDataToConns(ctx context.Context,tunIfce tzNet.Tunifce,connsStorage ConnsStoreage){
	for {
		bytesFromTunIfce,err := tunIfce.Read()
		if err != nil{
			log.Error(fmt.Errorf("Read from tun interface failed %w",err))
			// < maybe closed? TODO
			// return if tunIfce close
			// now return ignore the err type
			log.Info("Read from tun ifce failed, the loop tunIfceDataToConns returnd")
			return
			// >
		}
		// TODO
	}
}
