package trigger

import (
	"my5G-RANTester/internal/control_test_engine/ue/context"
	"my5G-RANTester/internal/control_test_engine/ue/nas/message/nas_control/mm_5gs"
	"my5G-RANTester/internal/control_test_engine/ue/nas/message/sender"
	"my5G-RANTester/lib/nas/nasMessage"
)

func InitRegistration(ue *context.UEContext) {

	// registration procedure started.
	registrationRequest := mm_5gs.GetRegistrationRequest(
		nasMessage.RegistrationType5GSInitialRegistration,
		nil,
		nil,
		false,
		ue)

	// send to GNB.
	sender.SendToGnb(ue, registrationRequest)

	// change the state of ue for deregistered
	ue.SetStateMM_DEREGISTERED()
}

func DeRegister(ue *context.UEContext) {

	// deregistration procedure started.
	deRegistrationRequest := mm_5gs.GetDeRegistrationRequest(
		0x09,
		nil,
		nil,
		false,
		ue)

	// send to GNB.
	sender.SendToGnb(ue, deRegistrationRequest)

	// change the state of ue for deregistered
	ue.SetStateMM_DEREGISTERED()
}