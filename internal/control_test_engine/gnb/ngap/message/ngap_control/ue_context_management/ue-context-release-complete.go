package ue_context_management

import (
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	"my5G-RANTester/lib/ngap"
	"my5G-RANTester/lib/ngap/ngapType"
)

/*
func initialContextSetupResponse(connN2 *sctp.SCTPConn, amfUeNgapID int64, ranUeNgapID int64, supi string) error {

	sendMsg, err := InitialContextSetupResponse(amfUeNgapID, ranUeNgapID)
	if err != nil {
		return fmt.Errorf("Error getting %s ue ngap Initial Context Setup Response Msg", supi)
	}
	_, err = connN2.Write(sendMsg)
	if err != nil {
		return fmt.Errorf("Error sending %s ue Initial Context Setup Response Msg", supi)
	}

	return nil
}
*/

func UEContextReleaseComplete(ue *context.GNBUe) ([]byte, error) {
	message := BuildInitialContextSetupResponseForRegistraionTest(ue.GetAmfUeId(), ue.GetRanUeId())

	return ngap.Encoder(message)
}

func BuildInitialContextSetupResponseForRegistraionTest(amfUeNgapID, ranUeNgapID int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodeInitialContextSetup
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentUEContextReleaseComplete
	successfulOutcome.Value.UEContextReleaseComplete = new(ngapType.UEContextReleaseComplete)

	uEContextReleaseComplete := successfulOutcome.Value.UEContextReleaseComplete
	uEContextReleaseCompleteIEs := &uEContextReleaseComplete.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.UEContextReleaseCompleteIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextReleaseCompleteIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	uEContextReleaseCompleteIEs.List = append(uEContextReleaseCompleteIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.UEContextReleaseCompleteIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.UEContextReleaseCompleteIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	uEContextReleaseCompleteIEs.List = append(uEContextReleaseCompleteIEs.List, ie)

	return
}
