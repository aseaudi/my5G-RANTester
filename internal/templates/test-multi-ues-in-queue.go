package templates

import (
	"fmt"
	"math/rand"
	"my5G-RANTester/config"
	"my5G-RANTester/internal/control_test_engine/gnb"
	"my5G-RANTester/internal/control_test_engine/ue"
	"strconv"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

func TestMultiUesInQueue(numUes int, numGnbs int, msinOffset int, regPeriod int) {

	wg := sync.WaitGroup{}

	cfg, err := config.GetConfig()
	if err != nil {
		//return nil
		log.Fatal("Error in get configuration")
	}
	var msin_int int
	var new_msin_int int
	msin_int, err = strconv.Atoi(cfg.Ue.Msin)
	new_msin_int = msin_int + msinOffset
	cfg.Ue.Msin = strconv.Itoa(new_msin_int)

    for j:= 1; j<=numGnbs; j++{
		log.Info("[TESTER] INIT GNB ", j)
		//go gnb.InitGnb2(cfg, int(j), &wg)

		//wg.Add(1)

		gnb.InitGnb2(cfg, int(j), &wg)

		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		msin :=  cfg.Ue.Msin
		
		for i := 1; i <= numUes; i++ {

			imsi := imsiGenerator(i, msin)
			log.Info("[TESTER] TESTING REGISTRATION USING IMSI ", imsi, " UE")
			cfg.Ue.Msin = imsi
			go ue.RegistrationUe2(cfg, uint8(i), j, &wg)
			wg.Add(1)

			time.Sleep(time.Duration(regPeriod) * time.Second)
		}
		imsi := imsiGenerator(numUes + 1, cfg.Ue.Msin)
		cfg.Ue.Msin = imsi
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