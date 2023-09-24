package sm_5gs

import (
	"bytes"
	"fmt"
	"my5G-RANTester/lib/nas"
	"my5G-RANTester/lib/nas/nasConvert"
	"my5G-RANTester/lib/nas/nasMessage"
	"my5G-RANTester/lib/nas/nasType"
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
	// pduSessionReleaseRequest.IntegrityProtectionMaximumDataRate.SetMaximumDataRatePerUEForUserPlaneIntegrityProtectionForDownLink(0xff)
	// pduSessionReleaseRequest.IntegrityProtectionMaximumDataRate.SetMaximumDataRatePerUEForUserPlaneIntegrityProtectionForUpLink(0xff)

	// pduSessionReleaseRequest.PDUSessionType = nasType.NewPDUSessionType(nasMessage.PDUSessionEstablishmentRequestPDUSessionTypeType)
	// pduSessionReleaseRequest.PDUSessionType.SetPDUSessionTypeValue(uint8(0x01)) //IPv4 type

	// pduSessionReleaseRequest.ExtendedProtocolConfigurationOptions = nasType.NewExtendedProtocolConfigurationOptions(nasMessage.PDUSessionEstablishmentRequestExtendedProtocolConfigurationOptionsType)
	// protocolConfigurationOptions := nasConvert.NewProtocolConfigurationOptions()
	// protocolConfigurationOptions.AddIPAddressAllocationViaNASSignallingUL()
	// protocolConfigurationOptions.AddDNSServerIPv4AddressRequest()
	// protocolConfigurationOptions.AddDNSServerIPv6AddressRequest()
	// pcoContents := protocolConfigurationOptions.Marshal()
	// pcoContentsLength := len(pcoContents)
	// pduSessionReleaseRequest.ExtendedProtocolConfigurationOptions.SetLen(uint16(pcoContentsLength))
	// pduSessionReleaseRequest.ExtendedProtocolConfigurationOptions.SetExtendedProtocolConfigurationOptionsContents(pcoContents)

	m.GsmMessage.PDUSessionEstablishmentRequest = pduSessionReleaseRequest

	data := new(bytes.Buffer)
	err := m.GsmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	nasPdu = data.Bytes()
	return
}
