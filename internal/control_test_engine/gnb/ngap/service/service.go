package service

import (
	"fmt"
	"github.com/ishidawataru/sctp"
	log "github.com/sirupsen/logrus"
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	"my5G-RANTester/internal/control_test_engine/gnb/ngap"
	"time"
)

func InitConn(amf *context.GNBAmf, gnb *context.GNBContext) error {

	// check AMF IP and AMF port.
	remote := fmt.Sprintf("%s:%d", amf.GetAmfIp(), amf.GetAmfPort())
	local := fmt.Sprintf("%s:%d", gnb.GetGnbIp(), gnb.GetGnbPort())

	rem, err := sctp.ResolveSCTPAddr("sctp", remote)
	if err != nil {
		return err
	}
	loc, err := sctp.ResolveSCTPAddr("sctp", local)
	if err != nil {
		return err
	}

	// streams := amf.GetTNLAStreams()

	conn, err := sctp.DialSCTPExt(
		"sctp",
		loc,
		rem,
		sctp.InitMsg{NumOstreams: 2, MaxInstreams: 2})
	if err != nil {
		amf.SetSCTPConn(nil)
		return err
	}

	// set streams and other information about TNLA

	// successful established SCTP (TNLA - N2)
	amf.SetSCTPConn(conn)
	gnb.SetN2(conn)

	conn.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)

	go GnbListen(amf, gnb)

	return nil
}

func GnbListen(amf *context.GNBAmf, gnb *context.GNBContext) {

	buf := make([]byte, 65535)
	conn := amf.GetSCTPConn()

	/*
		defer func() {
			err := conn.Close()
			if err != nil {
				log.Info("[GNB][SCTP] Error in closing SCTP association for %d AMF\n", amf.GetAmfId())
			}
		}()
	*/

	for {

		n, info, err := conn.SCTPRead(buf[:])
		if err != nil {
			conn.Close()
			log.Info("[GNB][NGAP] Error reading SCTP : ", err)
			log.Info("[GNB][NGAP] Reconnecting with AMF SCTP")
			//break
			// check AMF IP and AMF port.
			remote := fmt.Sprintf("%s:%d", amf.GetAmfIp(), amf.GetAmfPort())
			local := fmt.Sprintf("%s:%d", gnb.GetGnbIp(), gnb.GetGnbPort())
			rem, err := sctp.ResolveSCTPAddr("sctp", remote)
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			loc, err := sctp.ResolveSCTPAddr("sctp", local)
			if err != nil {
				time.Sleep(1 * time.Second)
				continue
			}
			sctp_connect := 0
			for sctp_connect == 0 {
				conn_new, err := sctp.DialSCTPExt(
					"sctp",
					loc,
					rem,
					sctp.InitMsg{NumOstreams: 2, MaxInstreams: 2})
				if err != nil {
					log.Info("[GNB][NGAP] Error reconnecting with AMF SCTP: ", err)
					log.Info("[GNB][NGAP] Retry reconnect with AMF SCTP again ...")
					//amf.SetSCTPConn(nil)
					time.Sleep(1 * time.Second)
					continue
				}
				// set streams and other information about TNLA
				// successful established SCTP (TNLA - N2)
				amf.SetSCTPConn(conn_new)
				gnb.SetN2(conn_new)
				conn_new.SubscribeEvents(sctp.SCTP_EVENT_DATA_IO)
				sctp_connect = 1
				log.Info("[GNB][NGAP] SCTP Reconnected OK with AMF ", amf.GetAmfId())
			}
		if sctp_connect == 1 {
			conn = amf.GetSCTPConn()
			continue
		}
		}

		log.Info("[GNB][SCTP] Receive message in ", info.Stream, " stream")

		forwardData := make([]byte, n)
		copy(forwardData, buf[:n])

		// handling NGAP message.
		go ngap.Dispatch(amf, gnb, forwardData)

	}

}
