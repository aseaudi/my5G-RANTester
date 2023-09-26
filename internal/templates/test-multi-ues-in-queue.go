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
	rand.Seed(time.Now().UnixNano())
	wg := sync.WaitGroup{}
	//ch := make(chan string)
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
		//time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
		go gnb.InitGnb2(cfg, int(j), &wg)

		wg.Add(1)

		
		time.Sleep(time.Duration(1) * time.Second)
		msin :=  cfg.Ue.Msin
		randNumUes := rand.Intn(numUes) + 1
		log.Info("[TESTER] TESTING Random Number of UEs = ", randNumUes)
		for i := 1; i <= randNumUes; i++ {

			imsi := imsiGenerator(i, msin)
			log.Info("[TESTER] TESTING REGISTRATION USING IMSI ", imsi, " UE")
			cfg.Ue.Msin = imsi
			go ue.RegistrationUe2(cfg, uint8(i), j, &wg)
			wg.Add(1)
			log.Info("[TESTER] Wait for next UE ", regPeriod, " seconds")

			time.Sleep(time.Duration(regPeriod) * time.Second)
		}
		imsi := imsiGenerator(numUes + 1, cfg.Ue.Msin)
		cfg.Ue.Msin = imsi

		//ch<-"ues done"

	}
	wg.Wait()

}

func imsiGenerator(i int, msin string) string {

	//msin_int, err := strconv.Atoi(msin)	
	// if err != nil {
	// 	log.Fatal("Error in get configuration")
	// }
	msin_int := rand.Intn(3000)
	base := msin_int + (i -1)

	imsi := fmt.Sprintf("%010d", base)
	return imsi
}