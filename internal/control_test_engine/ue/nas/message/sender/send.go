package sender

import (
	"fmt"
	"my5G-RANTester/internal/control_test_engine/ue/context"
	"encoding/hex"
)

func SendToGnb(ue *context.UEContext, message []byte) {
	log.Info("[UE] [NAS] UL NAS MSG = ", hex.EncodeToString(message))
	conn := ue.GetUnixConn()
	_, err := conn.Write(message)
	if err != nil {
		fmt.Println("[UE] [NAS] Error Writing to GNB UNIX Socket")
	}
}
