package templates

import (
	// "math/rand"
	"fmt"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	"strconv"

	// "my5G-RANTester/internal/control_test_engine/ue"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestMultiUesMultiGNBs(numUes int, numGNBs int) {

	log.Info("Num of UEs = ", numUes)
	log.Info("Num of GNBs = ", numGNBs)

	wg := sync.WaitGroup{}

	cfg, err := config.GetConfig()
	if err != nil {
		//return nil
		log.Fatal("Error in get configuration")
	}

	gnbControlPort := cfg.GNodeB.ControlIF.Port

	for i := 0; i < numGNBs; i++ {
		gnbID, err := strconv.Atoi(string(cfg.GNodeB.PlmnList.GnbId))
		if err != nil {
			log.Error("Failed to extract gnbID")
		}
		gnbID++
		newGnbID := fmt.Sprintf("%d", gnbID)

		cfg.GNodeB.PlmnList.GnbId = newGnbID
		cfg.GNodeB.ControlIF.Port = gnbControlPort + i
		log.Info("Initializing gnb with GnbId = ", cfg.GNodeB.PlmnList.GnbId)
		log.Info("Initializing gnb with gnbControlPort = ", cfg.GNodeB.ControlIF.Port)

		go gnb.InitGnb(cfg, &wg)
		wg.Add(1)
		time.Sleep(1 * time.Second)
	}

	//time.Sleep(1 * time.Second)

	// msin := cfg.Ue.Msin
	//startTime := time.Now()
	// for i := 1; i <= numUes; i++ {

	// 	portOffset := rand.Intn(numGNBs)
	// 	cfg.GNodeB.ControlIF.Port = gnbControlPort + portOffset
	// 	cfg.GNodeB.DataIF.Port = gnbDataPort + portOffset

	// 	log.Info("Registering ue with gnbControlPort = ", cfg.GNodeB.ControlIF.Port)
	// 	log.Info("Registering ue with gnbDataPort = ", cfg.GNodeB.DataIF.Port)

	// 	imsi := imsiGenerator(i, msin)
	// 	log.Info("[TESTER] TESTING REGISTRATION USING IMSI ", imsi, " UE")
	// 	cfg.Ue.Msin = imsi
	// 	go ue.RegistrationUe(cfg, uint8(i), &wg)
	// 	wg.Add(1)

	// 	sleepTime := time.Duration(rand.Intn(100)+1) * time.Millisecond
	// 	time.Sleep(sleepTime)
	// }

	wg.Wait()
	// endTime := time.Now()
	// executionTime := endTime.Sub(startTime)
	// log.Info("Total Registeration Time =", executionTime)
}
