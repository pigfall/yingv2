package server

import(
	"reflect"
tzNet "github.com/pigfall/tzzGoUtil/net"
	"context"
	"fmt"
"github.com/pigfall/tzzGoUtil/funcs"

"github.com/pigfall/tzzGoUtil/async"
	log "github.com/pigfall/tzzGoUtil/log/golog"
)

var (
		JOB_TUNIFCE_DATA_TO_CONNS = reflect.TypeOf( tunIfceDataToConns).String()
		JOB_TRANSPORT_SERVER_SERVE = reflect.TypeOf(TransportServer.Serve).String()
)


/*
jobs:
	tunIfceDataToConns
	transportServer.Serve
*/
func Serve(ctx context.Context,serveIpPort *tzNet.IpPort,transportServerBuilder TransportServerBuilder )(error){
	ctx,cancelCtxFunc := context.WithCancel(ctx)
	defer cancelCtxFunc()
	cleanFuncs := funcs.NewFuncs()
	defer cleanFuncs.Call()
	cancelFuncs := funcs.NewFuncs()
	defer cancelFuncs.Call()

	tunIfce,tunIpNet,err := readyTunIfce(nil)
	if err != nil {
		err = fmt.Errorf("Prepare tun interface failed %w",err)
		log.Error(err.Error())
		return err
	}
	cleanFuncs.AddFunc(func(){tunIfce.Close()})
	cancelFuncs.AddFunc(func(){tunIfce.Close()})
	cancelFuncs.AddFunc(cancelCtxFunc)
	tunIfceName,err := tunIfce.Name()
	if err != nil{
		err  = fmt.Errorf("get tunIfceName failed %v",err)
		log.Error(err.Error())
		return err
	}
	log.Info(fmt.Sprintf("Created tun interface %s which ipNet is %s",tunIfceName,tunIpNet.ToString()))

	ipPool,err := tzNet.NewIpPool(tunIpNet,[]*tzNet.IpWithMask{tunIpNet})
	if err != nil{
		log.Error(err.Error())
		return err
	}
	connsStorage := newConnStorage(ipPool)

	asyncCtrl := async.NewCtrl()
	asyncCtrl.OnRoutineQuit(func(jobName string){
		cancelFuncs.Call()
	})

	asyncCtrl.AppendCancelFuncs(
		func(){
			cancelFuncs.Call()
		},
	)
	defer func(){
		asyncCtrl.Cancel()
		asyncCtrl.Wait()
	}()

	// < readFrom tunIfce then send to all connections
	log.Info(fmt.Sprintf("Start %s",JOB_TUNIFCE_DATA_TO_CONNS))
	asyncCtrl.AsyncDo(
		ctx,
		JOB_TUNIFCE_DATA_TO_CONNS,
		func(ctx context.Context){
			err := tunIfceDataToConns(ctx,tunIfce,connsStorage)
			log.Info(fmt.Sprintf("%s func over %v",reflect.TypeOf( tunIfceDataToConns).String(),err))
		},
	)
	// >

	// <
	transportServer := transportServerBuilder.BuildTransportServer()
	transportServerCancelFunc,err := transportServer.Prepare(serveIpPort)
	if err != nil{
		err := fmt.Errorf("Prepare trasponrtServer failed :%w",err)
		log.Error(err.Error())
		return err
	}
	cancelFuncs.AddFunc(transportServerCancelFunc)
	log.Info("Start transerportServer ")
	asyncCtrl.AsyncDo(
			ctx,
			JOB_TRANSPORT_SERVER_SERVE,
			func (ctx context.Context){
				err = transportServer.Serve(ctx,tunIfce,connsStorage)
				log.Info(fmt.Sprintf("%s over %v",JOB_TRANSPORT_SERVER_SERVE,err))
			},
	)
	// >


	// if ctx done, we cancel all jobs and wait jobs quit
	asyncCtrl.WaitCtx(ctx,func(){cancelFuncs.Call()})
	return nil
}


// readFrom tunIfce then send to all connections
func tunIfceDataToConns(ctx context.Context,tunIfce tzNet.TunIfce,connsStorage ConnsStorage)error{
	var buf = make([]byte,70*1024)
	for {
		readNum,err := tunIfce.Read(buf)
		bytesReadFromTunIfce := buf[:readNum]
		if err != nil{
			log.Error(fmt.Errorf("Read from tun interface failed %w",err).Error())
			// < maybe closed? TODO
			// return if tunIfce close
			// now return ignore the err type
			log.Info("Read from tun ifce failed, the loop tunIfceDataToConns returnd")
			return err
			// >
		}
		log.Debug(fmt.Sprintf("read bytes from tun interface %v",bytesReadFromTunIfce))
		// TODO
	}
}
