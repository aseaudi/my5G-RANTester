package gnb

import (
	log "github.com/sirupsen/logrus"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	serviceNas "my5G-RANTester/internal/control_test_engine/gnb/nas/service"
	serviceNgap "my5G-RANTester/internal/control_test_engine/gnb/ngap/service"
	"my5G-RANTester/internal/control_test_engine/gnb/ngap/trigger"
	"my5G-RANTester/internal/monitoring"
	"os"
	"os/signal"
	"sync"
	"time"
)

func InitGnb(conf config.Config, wg *sync.WaitGroup) {

	// instance new gnb.
	gnb := &context.GNBContext{}

	// new gnb context.
	gnb.NewRanGnbContext(
		conf.GNodeB.PlmnList.GnbId,
		conf.GNodeB.PlmnList.Mcc,
		conf.GNodeB.PlmnList.Mnc,
		conf.GNodeB.PlmnList.Tac,
		conf.GNodeB.SliceSupportList.Sst,
		conf.GNodeB.SliceSupportList.Sd,
		conf.GNodeB.ControlIF.Ip,
		conf.GNodeB.DataIF.Ip,
		conf.GNodeB.ControlIF.Port,
		conf.GNodeB.DataIF.Port)

	// start communication with AMF (server SCTP).

	// new AMF context.
	amf := gnb.NewGnBAmf(conf.AMF.Ip, conf.AMF.Port)

	// start communication with AMF(SCTP).
	if err := serviceNgap.InitConn(amf, gnb); err != nil {
		log.Fatal("Error in", err)
	} else {
		log.Info("[GNB] SCTP/NGAP service is running")
		// wg.Add(1)
	}

	// start communication with UE (server UNIX sockets).
	if err := serviceNas.InitServer(gnb); err != nil {
		log.Fatal("Error in ", err)
	} else {
		log.Info("[GNB] UNIX/NAS service is running")
	}

	trigger.SendNgSetupRequest(gnb, amf)

	// control the signals
	sigGnb := make(chan os.Signal, 1)
	signal.Notify(sigGnb, os.Interrupt)

	// Block until a signal is received.
	<-sigGnb
	gnb.Terminate()
	wg.Done()
	// os.Exit(0)

}

func InitGnb2(conf config.Config, id int, wg *sync.WaitGroup) {

	log.Info("[GNB] Init GNB ", id)

	// instance new gnb.
	gnb := &context.GNBContext{}

	// new gnb context.
	gnb.NewRanGnbContext(
		conf.GNodeB.PlmnList.GnbId,
		conf.GNodeB.PlmnList.Mcc,
		conf.GNodeB.PlmnList.Mnc,
		conf.GNodeB.PlmnList.Tac,
		conf.GNodeB.SliceSupportList.Sst,
		conf.GNodeB.SliceSupportList.Sd,
		conf.GNodeB.ControlIF.Ip,
		conf.GNodeB.DataIF.Ip,
		conf.GNodeB.ControlIF.Port + id,
		conf.GNodeB.DataIF.Port)
	
	// start communication with AMF (server SCTP).

	// new AMF context.
	amf := gnb.NewGnBAmf(conf.AMF.Ip, conf.AMF.Port)


	// start communication with AMF(SCTP).
	sctp_ok := 0;
	for sctp_ok == 0 {
		if err := serviceNgap.InitConn(amf, gnb); err != nil {
			log.Fatal("Error in", err)
		} else {
			log.Info("[GNB] SCTP/NGAP service is running")
			sctp_ok = 1
			// wg.Add(1)
		}
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}


	// start communication with UE (server UNIX sockets).
	if err := serviceNas.InitServer2(gnb, id); err != nil {
		log.Fatal("Error in ", err)
	} else {
		log.Info("[GNB] UNIX/NAS service is running")
	}


	trigger.SendNgSetupRequest(gnb, amf)

	// control the signals
	sigGnb := make(chan os.Signal, 1)
	signal.Notify(sigGnb, os.Interrupt)

	// Block until a signal is received.
	<-sigGnb
	gnb.Terminate()
	wg.Done()
	// os.Exit(0)

}

func InitGnbForUeLatency(conf config.Config, sigGnb chan bool, synch chan bool) {

	// instance new gnb.
	gnb := &context.GNBContext{}

	// new gnb context.
	gnb.NewRanGnbContext(
		conf.GNodeB.PlmnList.GnbId,
		conf.GNodeB.PlmnList.Mcc,
		conf.GNodeB.PlmnList.Mnc,
		conf.GNodeB.PlmnList.Tac,
		conf.GNodeB.SliceSupportList.Sst,
		conf.GNodeB.SliceSupportList.Sd,
		conf.GNodeB.ControlIF.Ip,
		conf.GNodeB.DataIF.Ip,
		conf.GNodeB.ControlIF.Port,
		conf.GNodeB.DataIF.Port)

	// start communication with AMF (server SCTP).

	// new AMF context.
	amf := gnb.NewGnBAmf(conf.AMF.Ip, conf.AMF.Port)

	// start communication with AMF(SCTP).
	if err := serviceNgap.InitConn(amf, gnb); err != nil {
		log.Info("Error in", err)

		synch <- false

		return
	} else {
		log.Info("[GNB] SCTP/NGAP service is running")
		// wg.Add(1)
	}

	// start communication with UE (server UNIX sockets).
	if err := serviceNas.InitServer(gnb); err != nil {
		log.Info("Error in", err)

		synch <- false
	} else {
		log.Info("[GNB] UNIX/NAS service is running")

	}

	trigger.SendNgSetupRequest(gnb, amf)

	synch <- true

	// Block until a signal is received.
	<-sigGnb
	gnb.Terminate()
}

func InitGnbForLoadSeconds(conf config.Config, wg *sync.WaitGroup,
	monitor *monitoring.Monitor) {

	// instance new gnb.
	gnb := &context.GNBContext{}

	// new gnb context.
	gnb.NewRanGnbContext(
		conf.GNodeB.PlmnList.GnbId,
		conf.GNodeB.PlmnList.Mcc,
		conf.GNodeB.PlmnList.Mnc,
		conf.GNodeB.PlmnList.Tac,
		conf.GNodeB.SliceSupportList.Sst,
		conf.GNodeB.SliceSupportList.Sd,
		conf.GNodeB.ControlIF.Ip,
		conf.GNodeB.DataIF.Ip,
		conf.GNodeB.ControlIF.Port,
		conf.GNodeB.DataIF.Port)

	// start communication with AMF (server SCTP).

	// new AMF context.
	amf := gnb.NewGnBAmf(conf.AMF.Ip, conf.AMF.Port)

	// start communication with AMF(SCTP).
	if err := serviceNgap.InitConn(amf, gnb); err != nil {
		log.Info("Error in ", err)

		time.Sleep(1000 * time.Millisecond)

		wg.Done()

		return
	} else {
		log.Info("[GNB] SCTP/NGAP service is running")
		// wg.Add(1)
	}

	trigger.SendNgSetupRequest(gnb, amf)

	// timeout is 1 second for receive NG Setup Response
	time.Sleep(1000 * time.Millisecond)

	// AMF responds message sends by Tester
	// means AMF is available
	if amf.GetState() == 0x01 {
		monitor.IncRqs()
	}

	gnb.Terminate()
	wg.Done()
	// os.Exit(0)
}

func InitGnbForAvaibility(conf config.Config,
	monitor *monitoring.Monitor) {

	// instance new gnb.
	gnb := &context.GNBContext{}

	// new gnb context.
	gnb.NewRanGnbContext(
		conf.GNodeB.PlmnList.GnbId,
		conf.GNodeB.PlmnList.Mcc,
		conf.GNodeB.PlmnList.Mnc,
		conf.GNodeB.PlmnList.Tac,
		conf.GNodeB.SliceSupportList.Sst,
		conf.GNodeB.SliceSupportList.Sd,
		conf.GNodeB.ControlIF.Ip,
		conf.GNodeB.DataIF.Ip,
		conf.GNodeB.ControlIF.Port,
		conf.GNodeB.DataIF.Port)

	// start communication with AMF (server SCTP).

	// new AMF context.
	amf := gnb.NewGnBAmf(conf.AMF.Ip, conf.AMF.Port)

	// start communication with AMF(SCTP).
	if err := serviceNgap.InitConn(amf, gnb); err != nil {
		log.Info("Error in ", err)

		return

	} else {
		log.Info("[GNB] SCTP/NGAP service is running")

	}

	trigger.SendNgSetupRequest(gnb, amf)

	// timeout is 1 second for receive NG Setup Response
	time.Sleep(1000 * time.Millisecond)

	// AMF responds message sends by Tester
	// means AMF is available
	if amf.GetState() == 0x01 {
		monitor.IncAvaibility()

	}

	gnb.Terminate()
	// os.Exit(0)
}
