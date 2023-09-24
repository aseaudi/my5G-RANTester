package sm_5gs

import (
	"bytes"
	"fmt"
	"my5G-RANTester/lib/nas"

	//"my5G-RANTester/lib/nas/nasConvert"
	nasType "command-line-arguments/Users/aseaudi/src/github/my5G-RANTester/lib/nas/nasType/NAS_Cause5GSM.go"
	"my5G-RANTester/lib/nas/nasMessage"
	//"my5G-RANTester/lib/nas/nasType"
)

func GetPduSessionReleaseRequest(pduSessionId uint8) (nasPdu []byte) {

	m := nas.NewMessage()
	m.GsmMessage = nas.NewGsmMessage()
	m.GsmHeader.SetMessageType(nas.MsgTypePDUSessionReleaseRequest)

	pduSessionReleaseRequest := nasMessage.NewPDUSessionReleaseRequest(0)
	pduSessionReleaseRequest.ExtendedProtocolDiscriminator.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSSessionManagementMessage)
	pduSessionReleaseRequest.SetMessageType(nas.MsgTypePDUSessionReleaseRequest)
	pduSessionReleaseRequest.PDUSessionID.SetPDUSessionID(pduSessionId)
	pduSessionReleaseRequest.PTI.SetPTI(0x01)
	pduSessionReleaseRequest.Cause5GSM = nasMessage.NewCause5GSM(nasMessage.PDUSessionReleaseRequestCause5GSMType)
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

	nasPdu = data.Bytes()
	return
}
