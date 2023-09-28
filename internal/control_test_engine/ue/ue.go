package ue

import (
	log "github.com/sirupsen/logrus"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/ue/context"
	"my5G-RANTester/internal/control_test_engine/ue/nas/service"
	"my5G-RANTester/internal/control_test_engine/ue/nas/trigger"
	"my5G-RANTester/internal/monitoring"
	"my5G-RANTester/lib/nas/security"
	"os"
	"os/signal"
	"sync"
	"time"
)

func RegistrationUe(conf config.Config, id uint8, wg *sync.WaitGroup) {

	// wg := sync.WaitGroup{}

	// new UE instance.
	ue := &context.UEContext{}

	// new UE context
	ue.NewRanUeContext(
		conf.Ue.Msin,
		security.AlgCiphering128NEA0,
		security.AlgIntegrity128NIA2,
		conf.Ue.Key,
		conf.Ue.Opc,
		"c9e8763286b5b9ffbdf56e1297d0887b",
		conf.Ue.Amf,
		conf.Ue.Sqn,
		conf.Ue.Hplmn.Mcc,
		conf.Ue.Hplmn.Mnc,
		conf.Ue.Dnn,
		int32(conf.Ue.Snssai.Sst),
		conf.Ue.Snssai.Sd,
		id)

	// starting communication with GNB and listen.
	err := service.InitConn(ue)
	if err != nil {
		log.Fatal("Error in", err)
	} else {
		log.Info("[UE] UNIX/NAS service is running")
		// wg.Add(1)
	}

	// registration procedure started.
	trigger.InitRegistration(ue)

	// wg.Wait()

	// control the signals
	sigUe := make(chan os.Signal, 1)
	signal.Notify(sigUe, os.Interrupt)

	// Block until a signal is received.
	<-sigUe
	ue.Terminate()
	wg.Done()
	// os.Exit(0)

}

func RegistrationUe2(conf config.Config, id uint8, id2 int, wg *sync.WaitGroup) {

	// wg := sync.WaitGroup{}


	// new channel
	ch := make(chan string)

	// new UE instance.
	ue := &context.UEContext{}

	// new UE context
	ue.NewRanUeContext(
		conf.Ue.Msin,
		security.AlgCiphering128NEA0,
		security.AlgIntegrity128NIA2,
		conf.Ue.Key,
		conf.Ue.Opc,
		"c9e8763286b5b9ffbdf56e1297d0887b",
		conf.Ue.Amf,
		conf.Ue.Sqn,
		conf.Ue.Hplmn.Mcc,
		conf.Ue.Hplmn.Mnc,
		conf.Ue.Dnn,
		int32(conf.Ue.Snssai.Sst),
		conf.Ue.Snssai.Sd,
		id,
		ch)

	// starting communication with GNB and listen.
	ue_connect := 0
	for ue_connect == 0 {
	err := service.InitConn2(ue, id2)
	if err != nil {
		log.Warn("[UE] UE ", id, " INIT CONNECT Error : ", err)
		time.Sleep(1 * time.Second)
	} else {
		log.Info("[UE] UE ", id, " UNIX/NAS service is running")
		ue_connect = 1
		// wg.Add(1)
	}
	}
	// registration procedure started.
	trigger.InitRegistration(ue)

	// wg.Wait()

	// wait for a random time, then deregister ue
	//log.Info("[UE] UE ", id, " registered, wait 60 seconds then deregister")
	//time.Sleep(60 * time.Second)
	//trigger.DeRegister(ue)
	//log.Info("[UE] UE ", id, " deregistered")

	// control the signals
	sigUe := make(chan os.Signal, 1)
	signal.Notify(sigUe, os.Interrupt)

	// Block until a signal is received.
	// <-sigUe

	select {
	case <-sigUe:
		log.Info("[UE] OS Signal Interrupt, UE will now terminate")
	case <-ch:
		log.Info("[UE] UE Signal received, UE will now terminate")
	}
	ue.Terminate()
	wg.Done()
	// os.Exit(0)

}

func RegistrationUeMonitor(conf config.Config,
	id uint8, monitor *monitoring.Monitor, wg *sync.WaitGroup, start time.Time) {

	// new UE instance.
	ue := &context.UEContext{}

	// new UE context
	ue.NewRanUeContext(
		conf.Ue.Msin,
		security.AlgCiphering128NEA0,
		security.AlgIntegrity128NIA2,
		conf.Ue.Key,
		conf.Ue.Opc,
		"c9e8763286b5b9ffbdf56e1297d0887b",
		conf.Ue.Amf,
		conf.Ue.Sqn,
		conf.Ue.Hplmn.Mcc,
		conf.Ue.Hplmn.Mnc,
		conf.Ue.Dnn,
		int32(conf.Ue.Snssai.Sst),
		conf.Ue.Snssai.Sd,
		id)

	// starting communication with GNB and listen.
	err := service.InitConn(ue)
	if err != nil {
		log.Fatal("Error in", err)
	} else {
		log.Info("[UE] UNIX/NAS service is running")
		// wg.Add(1)
	}

	// registration procedure started.
	trigger.InitRegistration(ue)

	for {

		// UE is register in network
		if ue.GetStateMM() == 0x03 {
			elapsed := time.Since(start)
			monitor.LtRegisterLocal = elapsed.Milliseconds()
			log.Warn("[TESTER][UE] UE LATENCY IN REGISTRATION: ", monitor.LtRegisterLocal, " ms")
			break
		}

		// timeout is 10 000 ms
		if time.Since(start).Milliseconds() >= 10000 {
			log.Warn("[TESTER][UE] TIME EXPIRED IN UE REGISTRATION 10 000 ms")
			break
		}
	}

	wg.Done()
	// ue.Terminate()
	// os.Exit(0)
}
