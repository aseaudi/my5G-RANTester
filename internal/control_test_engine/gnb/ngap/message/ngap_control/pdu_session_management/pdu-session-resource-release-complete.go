package pdu_session_management

import (
	"encoding/binary"
	"my5G-RANTester/internal/control_test_engine/gnb/context"
	"my5G-RANTester/lib/aper"
	"my5G-RANTester/lib/ngap"
	"my5G-RANTester/lib/ngap/ngapConvert"
	"my5G-RANTester/lib/ngap/ngapType"
	"net"
)

func PDUSessionResourceReleaseComplete(ue *context.GNBUe, ipv4 string) ([]byte, error) {

	// check hostname(Error in docker if using hostname)
	nameIp, err := net.LookupHost(ipv4)
	if err != nil {
		return nil, err
	}
	message := buildPDUSessionResourceSetupResponseForRegistrationTest(ue.GetAmfUeId(), ue.GetRanUeId(),
		nameIp[0], ue.GetPduSessionId(), ue.GetTeidDownlink(), ue.GetQosId())
	return ngap.Encoder(message)
}

func buildPDUSessionResourceSetupResponseForRegistrationTest(amfUeNgapID, ranUeNgapID int64, ipv4 string, pduId int64, teid uint32, qosId int64) (pdu ngapType.NGAPPDU) {

	pdu.Present = ngapType.NGAPPDUPresentSuccessfulOutcome
	pdu.SuccessfulOutcome = new(ngapType.SuccessfulOutcome)

	successfulOutcome := pdu.SuccessfulOutcome
	successfulOutcome.ProcedureCode.Value = ngapType.ProcedureCodePDUSessionResourceRelease
	successfulOutcome.Criticality.Value = ngapType.CriticalityPresentReject

	successfulOutcome.Value.Present = ngapType.SuccessfulOutcomePresentPDUSessionResourceReleaseResponse
	successfulOutcome.Value.ProcedureCodePDUSessionResourceRelease = new(ngapType.PDUSessionResourceRelease)

	pDUSessionResourceSetupResponse := successfulOutcome.Value.PDUSessionResourceReleaseResponse
	pDUSessionResourceSetupResponseIEs := &pDUSessionResourceSetupResponse.ProtocolIEs

	// AMF UE NGAP ID
	ie := ngapType.PDUSessionResourceReleaseResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDAMFUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceReleaseResponseIEsPresentAMFUENGAPID
	ie.Value.AMFUENGAPID = new(ngapType.AMFUENGAPID)

	aMFUENGAPID := ie.Value.AMFUENGAPID
	aMFUENGAPID.Value = amfUeNgapID

	pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

	// RAN UE NGAP ID
	ie = ngapType.PDUSessionResourceReleaseResponseIEs{}
	ie.Id.Value = ngapType.ProtocolIEIDRANUENGAPID
	ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	ie.Value.Present = ngapType.PDUSessionResourceReleaseResponseIEsPresentRANUENGAPID
	ie.Value.RANUENGAPID = new(ngapType.RANUENGAPID)

	rANUENGAPID := ie.Value.RANUENGAPID
	rANUENGAPID.Value = ranUeNgapID

	pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

	// // PDU Session Resource Setup Response List
	// ie = ngapType.PDUSessionResourceReleaseResponseIEs{}
	// ie.Id.Value = ngapType.ProtocolIEIDPDUSessionResourceReleasedListRelRes
	// ie.Criticality.Value = ngapType.CriticalityPresentIgnore
	// ie.Value.Present = ngapType.PDUSessionResourceReleaseResponseIEsPresentPDUSessionResourceReleasedListRelRes
	// ie.Value.PDUSessionResourceReleasedListRelRes = new(ngapType.PDUSessionResourceReleasedListRelRes)

	// pDUSessionResourceReleasedListRelRes := ie.Value.PDUSessionResourceReleasedListRelRes

	// // PDU Session Resource Setup Response Item in PDU Session Resource Setup Response List
	// pDUSessionResourceReleasedItemRelRes := ngapType.PDUSessionResourceReleasedItemRelRes{}

	// // PDU Session ID : This is an unique identifier generated by UE. Can’t be same as any existing PDU session.
	// pDUSessionResourceReleasedItemRelRes.PDUSessionID.Value = pduId

	// //pDUSessionResourceReleasedItemRelRes.PDUSessionResourceReleaseResponseTransfer = GetPDUSessionResourceReleaseResponseTransfer(ipv4, teid, qosId)

	// pDUSessionResourceReleasedListRelRes.List = append(pDUSessionResourceReleasedListRelRes.List, pDUSessionResourceReleasedItemRelRes)

	// pDUSessionResourceSetupResponseIEs.List = append(pDUSessionResourceSetupResponseIEs.List, ie)

	return
}

func GetPDUSessionResourceReleaseResponseTransfer(ipv4 string, teid uint32, qosId int64) []byte {
	data := buildPDUSessionResourceReleaseResponseTransfer(ipv4, teid, qosId)
	encodeData, _ := aper.MarshalWithParams(data, "valueExt")
	return encodeData
}

func buildPDUSessionResourceReleaseResponseTransfer(ipv4 string, teid uint32, qosId int64) (data ngapType.PDUSessionResourceSetupResponseTransfer) {

	// QoS Flow per TNL Information
	qosFlowPerTNLInformation := &data.QosFlowPerTNLInformation
	qosFlowPerTNLInformation.UPTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel

	// UP Transport Layer Information in QoS Flow per TNL Information
	upTransportLayerInformation := &qosFlowPerTNLInformation.UPTransportLayerInformation
	upTransportLayerInformation.Present = ngapType.UPTransportLayerInformationPresentGTPTunnel
	upTransportLayerInformation.GTPTunnel = new(ngapType.GTPTunnel)

	dowlinkTeid := make([]byte, 4)
	binary.BigEndian.PutUint32(dowlinkTeid, teid)
	upTransportLayerInformation.GTPTunnel.GTPTEID.Value = dowlinkTeid
	upTransportLayerInformation.GTPTunnel.TransportLayerAddress = ngapConvert.IPAddressToNgap(ipv4, "")

	// Associated QoS Flow List in QoS Flow per TNL Information
	associatedQosFlowList := &qosFlowPerTNLInformation.AssociatedQosFlowList

	associatedQosFlowItem := ngapType.AssociatedQosFlowItem{}
	associatedQosFlowItem.QosFlowIdentifier.Value = qosId
	associatedQosFlowList.List = append(associatedQosFlowList.List, associatedQosFlowItem)

	return
}
