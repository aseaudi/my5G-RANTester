package mm_5gs

import (
	"bytes"
	"fmt"
	"my5G-RANTester/internal/control_test_engine/ue/context"
	"my5G-RANTester/lib/nas"
	"my5G-RANTester/lib/nas/nasMessage"
	"my5G-RANTester/lib/nas/nasType"
)

func GetDeRegistrationRequest(registrationType uint8, requestedNSSAI *nasType.RequestedNSSAI, uplinkDataStatus *nasType.UplinkDataStatus, capability bool, ue *context.UEContext) (nasPdu []byte) {

	//ueSecurityCapability := context.SetUESecurityCapability(ue)

	m := nas.NewMessage()
	m.GmmMessage = nas.NewGmmMessage()
	m.GmmHeader.SetMessageType(nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration)

	deregistrationRequest := nasMessage.NewDeregistrationRequestUEOriginatingDeregistration(0)
	deregistrationRequest.SetExtendedProtocolDiscriminator(nasMessage.Epd5GSMobilityManagementMessage)
	deregistrationRequest.SpareHalfOctetAndSecurityHeaderType.SetSecurityHeaderType(nas.SecurityHeaderTypePlainNas)
	deregistrationRequest.SpareHalfOctetAndSecurityHeaderType.SetSpareHalfOctet(0x00)
	deregistrationRequest.DeregistrationRequestMessageIdentity.SetMessageType(nas.MsgTypeDeregistrationRequestUEOriginatingDeregistration)
	// deregistrationRequest.NgksiAndRegistrationType5GS.SetTSC(nasMessage.TypeOfSecurityContextFlagNative)
	// deregistrationRequest.NgksiAndRegistrationType5GS.SetNasKeySetIdentifiler(ue.GetUeId())
	// deregistrationRequest.NgksiAndRegistrationType5GS.SetRegistrationType5GS(registrationType)
	deregistrationRequest.MobileIdentity5GS = ue.GetSuci()
	// if capability {
	// 	deregistrationRequest.Capability5GMM = &nasType.Capability5GMM{
	// 		Iei:   nasMessage.RegistrationRequestCapability5GMMType,
	// 		Len:   1,
	// 		Octet: [13]uint8{0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
	// 	}
	// } else {
	// 	deregistrationRequest.Capability5GMM = nil
	// }
	// deregistrationRequest.UESecurityCapability = ueSecurityCapability
	// deregistrationRequest.RequestedNSSAI = requestedNSSAI
	// deregistrationRequest.UplinkDataStatus = uplinkDataStatus

	//deregistrationRequest.SetFOR(1)

	m.GmmMessage.DeregistrationRequestUEOriginatingDeregistration = deregistrationRequest

	data := new(bytes.Buffer)
	err := m.GmmMessageEncode(data)
	if err != nil {
		fmt.Println(err.Error())
	}

	nasPdu = data.Bytes()
	return
}
