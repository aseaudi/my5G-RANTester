package templates

import (
	log "github.com/sirupsen/logrus"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	"my5G-RANTester/internal/control_test_engine/ue"
	"strconv"
	"sync"
	"time"
	"fmt"
)

func TestMultiUesInQueue(numUes int) {

	wg := sync.WaitGroup{}

	cfg, err := config.GetConfig()
	if err != nil {
		//return nil
		log.Fatal("Error in get configuration")
	}
    for j:= 1; j<=10; j++{
		log.Info("[TESTER] INIT GNB ", j)
		go gnb.InitGnb2(cfg, int(j), &wg)

		wg.Add(1)

		time.Sleep(1 * time.Second)
		msin :=  cfg.Ue.Msin
		for i := 1; i <= numUes; i++ {

			imsi := imsiGenerator(i, msin)
			log.Info("[TESTER] TESTING REGISTRATION USING IMSI ", imsi, " UE")
			cfg.Ue.Msin = imsi
			go ue.RegistrationUe2(cfg, uint8(i), j, &wg)
			wg.Add(1)

			time.Sleep(4 * time.Second)
		}
	}
	wg.Wait()

}

func imsiGenerator(i int, msin string) string {

	msin_int, err := strconv.Atoi(msin)
	if err != nil {
		log.Fatal("Error in get configuration")
	}
	base := msin_int + (i -1)

	imsi := fmt.Sprintf("%010d", base)
	return imsi
}