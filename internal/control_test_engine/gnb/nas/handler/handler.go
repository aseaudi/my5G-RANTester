package handler

import (
	log "github.com/sirupsen/logrus"
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	"my5G-RANTester/internal/control_test_engine/gnb/ngap/message/ngap_control/nas_transport"
	"my5G-RANTester/internal/control_test_engine/gnb/ngap/message/sender"
)

func HandlerUeInitialized(ue *context.GNBUe, message []byte, gnb *context.GNBContext) {
	log.Info("[GNB][NAS] Handler for UE Initialized")
	// encode NAS message in NGAP.
	ngap, err := nas_transport.SendInitialUeMessage(message, ue, gnb)
	if err != nil {
		log.Info("[GNB][NAS] Error making initial UE message: ", err)
	}

	// change state of UE.
	ue.SetStateOngoing()

	// Send Initial UE Message
	//conn := ue.GetSCTP()
	amf_id := ue.GetAmfId()
	amf, _ := gnb.GetGnbAmf(amf_id)
	log.Info("[GNB][AMF] Sending Initial UE Message to AMF ID ", amf.GetAmfId())
	conn := amf.GetSCTPConn()
	log.Info("[GNB][AMF] STPConn = ", conn)
	err = sender.SendToAmF(ngap, conn)
	if err != nil {
		log.Info("[GNB][AMF] Error sending initial UE message: ", err)
	}
}

func HandlerUeOngoing(ue *context.GNBUe, message []byte, gnb *context.GNBContext) {
	log.Info("[GNB][NAS] Handler for UE Ongoing")

	ngap, err := nas_transport.SendUplinkNasTransport(message, ue, gnb)
	if err != nil {
		log.Info("[GNB][NGAP] Error making Uplink Nas Transport: ", err)
	}

	// Send Uplink Nas Transport
	conn := ue.GetSCTP()
	err = sender.SendToAmF(ngap, conn)
	if err != nil {
		log.Info("[GNB][AMF] Error sending Uplink Nas Transport: ", err)
	}
}

func HandlerUeReady(ue *context.GNBUe, message []byte, gnb *context.GNBContext) {
	log.Info("[GNB][NAS] Handler for UE Ready")

	// receive UE ip or other messages.

	ngap, err := nas_transport.SendUplinkNasTransport(message, ue, gnb)
	

	if err != nil {
		log.Info("[GNB][NGAP] Error making Uplink Nas Transport: ", err)
	}

	// Send Uplink Nas Transport
	conn := ue.GetSCTP()
	err = sender.SendToAmF(ngap, conn)
	if err != nil {
		log.Info("[GNB][AMF] Error sending Uplink Nas Transport: ", err)
	} else {
	log.Info("[GNB][NAS] Sent UL NAS Transport to AMF")
	}
}
