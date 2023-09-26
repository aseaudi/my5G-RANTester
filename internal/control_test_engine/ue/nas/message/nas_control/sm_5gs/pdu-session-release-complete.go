package sm_5gs

import (
	"bytes"
	"fmt"
	"my5G-RANTester/lib/nas"

	//"my5G-RANTester/lib/nas/nasConvert"
	"my5G-RANTester/lib/nas/nasMessage"
	"my5G-RANTester/lib/nas/nasType"
	log "github.com/sirupsen/logrus"
)

func GetPduSessionReleaseComplete(pduSessionId uint8) (nasPdu []byte) {
	log.Info("GetPduSessionReleaseComplete")
	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionReleaseComplete)

	pduSessionReleaseRequest := nasMessage.NewPDUSessionReleaseComplete(0)
	pduSessionReleaseRequest.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	pduSessionReleaseRequest.SetMessageType(nas.MsgTypePDUSessionReleaseRequest)
	pduSessionReleaseRequest.PDUSessionID.SetPDUSessionID(pduSessionId)
	pduSessionReleaseRequest.PTI.SetPTI(0x01)
	pduSessionReleaseRequest.PDUSESSIONRELEASEREQUESTMessageIdentity = *nasType.NewPDUSESSIONRELEASEREQUESTMessageIdentity(209)
	//pduSessionReleaseRequest.PDUSESSIONRELEASEREQUESTMessageIdentity.SetMessageType(209)
	pduSessionReleaseRequest.Cause5GSM = nasType.NewCause5GSM(nasMessage.PDUSessionReleaseRequestCause5GSMType)
	pduSessionReleaseRequest.Cause5GSM.SetCauseValue(36)
	// pduSessionReleaseRequest.IntegrityProtectionMaximumDataRate.SetMaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink(0xff)
	// pduSessionReleaseRequest.IntegrityProtectionMaximumDataRate.SetMaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink(0xff)

	// pduSessionReleaseRequest.PDUSessionType = nasType.NewPDUSessionType(nasMessage.PDUSessionReleaseRequestPDUSessionTypeType)
	// pduSessionReleaseRequest.PDUSessionType.SetPDUSessionTypeValue(uint8(0x01)) //IPv4 type

	// pduSessionReleaseRequest.ExtendedProtocolConfigurationOptions = nasType.NewExtendedProtocolConfigurationOptions(nasMessage.PDUSessionReleaseRequestExtendedProtocolConfigurationOptionsType)
	// protocolConfigurationOptions := nasConvert.NewProtocolConfigurationOptions()
	// protocolConfigurationOptions.AddIPAddressAllocationViaNASSignallingUL()
	// protocolConfigurationOptions.AddDNSServerIPv4AddressRequest()
	// protocolConfigurationOptions.AddDNSServerIPv6AddressRequest()
	// pcoContents := protocolConfigurationOptions.Marshal()
	// pcoContentsLength := len(pcoContents)
	// pduSessionReleaseRequest.ExtendedProtocolConfigurationOptions.SetLen(uint16(pcoContentsLength))
	// pduSessionReleaseRequest.ExtendedProtocolConfigurationOptions.SetExtendedProtocolConfigurationOptionsContents(pcoContents)

	m.GsmMessage.PDUSessionReleaseRequest = pduSessionReleaseRequest

	data := new(bytes.Buffer)
	err := m.GsmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}
	log.Info("gsm encoded pdu session release request")

	nasPdu = data.Bytes()
	return
}
