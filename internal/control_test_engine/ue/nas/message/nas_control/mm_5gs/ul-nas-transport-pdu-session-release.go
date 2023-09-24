package mm_5gs

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"my5G-RANTester/internal/control_test_engine/ue/context"
	"my5G-RANTester/internal/control_test_engine/ue/nas/message/nas_control"
	"my5G-RANTester/internal/control_test_engine/ue/nas/message/nas_control/sm_5gs"
	"my5G-RANTester/lib/nas"
	"my5G-RANTester/lib/nas/nasMessage"
	"my5G-RANTester/lib/nas/nasType"
	"my5G-RANTester/lib/openapi/models"
	log "github.com/sirupsen/logrus"


)

func UlNasTransport2(ue *context.UEContext, requestType uint8) ([]byte, error) {
	log.Info("UL NAS Transport 2")
	pdu := getUlNasTransport_PduSessionReleaseRequest(ue.PduSession.Id, requestType, ue.PduSession.Dnn, &ue.PduSession.Snssai)
	if pdu == nil {
		return nil, fmt.Errorf("Error encoding %s IMSI UE PduSession Release Request Msg", ue.UeSecurity.Supi)
	}
	log.Info("UL NAS Transport 2 - created pdu")
	pdu, err := nas_control.EncodeNasPduWithSecurity(ue, pdu, nas.SecurityHeaderTypeIntegrityProtectedAndCiphered, true, false)
	if err != nil {
		return nil, fmt.Errorf("Error encoding %s IMSI UE PduSession Release Request Msg", ue.UeSecurity.Supi)
	}
	log.Info("UL NAS Transport 2 - encoded pdu with security")
	return pdu, nil
}

func getUlNasTransport_PduSessionReleaseRequest(pduSessionId uint8, requestType uint8, dnnString string, sNssai *models.Snssai) (nasPdu []byte) {
	log.Info("getULNASTRansport Pdu session release request")
	pduSessionReleaseRequest := sm_5gs.GetPduSessionReleaseRequest(pduSessionId)

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeULNASTransport)

	ulNasTransport := nasMessage.NewULNASTransport(0)
	ulNasTransport.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	ulNasTransport.SetMessageType(nas.MsgTypeULNASTransport)
	ulNasTransport.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	ulNasTransport.PduSessionID2Value = new(nasType.PduSessionID2Value)
	ulNasTransport.PduSessionID2Value.SetIei(nasMessage.ULNASTransportPduSessionID2ValueType)
	ulNasTransport.PduSessionID2Value.SetPduSessionID2Value(pduSessionId)
	ulNasTransport.RequestType = new(nasType.RequestType)
	ulNasTransport.RequestType.SetIei(nasMessage.ULNASTransportRequestTypeType)
	ulNasTransport.RequestType.SetRequestTypeValue(requestType)

	if dnnString != "" {
		dnn := []byte(dnnString)
		ulNasTransport.DNN = new(nasType.DNN)
		ulNasTransport.DNN.SetIei(nasMessage.ULNASTransportDNNType)
		ulNasTransport.DNN.SetLen(uint8(len(dnn)))
		ulNasTransport.DNN.SetDNN(dnn)
	}

	if sNssai != nil {
		ulNasTransport.SNSSAI = nasType.NewSNSSAI(nasMessage.ULNASTransportSNSSAIType)
		if sNssai.Sd == "" {
			ulNasTransport.SNSSAI.SetLen(1)
		} else {
			ulNasTransport.SNSSAI.SetLen(4)
			var sdTemp [3]uint8
			sd, _ := hex.DecodeString(sNssai.Sd)
			copy(sdTemp[:], sd)
			ulNasTransport.SNSSAI.SetSD(sdTemp)
		}
		ulNasTransport.SNSSAI.SetSST(uint8(sNssai.Sst))
	}

	ulNasTransport.SpareHalfOctetAndPayloadContainerType.SetPayloadContainerType(nasMessage.PayloadContainerTypeN1SMInfo)
	ulNasTransport.PayloadContainer.SetLen(uint16(len(pduSessionReleaseRequest)))
	ulNasTransport.PayloadContainer.SetPayloadContainerContents(pduSessionReleaseRequest)

	m.GmmMessage.ULNASTransport = ulNasTransport

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}
	log.Info("gmm encoded data pdu session release")

	nasPdu = data.Bytes()
	return
}
