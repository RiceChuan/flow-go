package crypto

// #cgo CFLAGS:
// #include "dkg_include.h"
import "C"

import (
	"fmt"
)

// Implements Feldman Verifiable Secret Sharing using
// the BLS set up on the BLS12-381 curve. A complaint mechanism
// is added to qualify/disqualify the dealer if they misbehave.

// The secret is a BLS private key generated by the dealer.
// (and hence this is a centralized generation).
// The dealer generates key shares for a BLS-based
// threshold signature scheme and distributes the shares over the (n)
// participants including itself. The participants validate their shares
// using a public verification vector shared by the dealer and are able
// to broadcast complaints against a misbehaving dealer.

// The dealer has the chance to avoid being disqualified by broadcasting
// a complaint answer. The protocol ends with all honest participants
// reaching a consensus about the dealer qualification/disqualification.

// Private keys are scalar in Fr, where r is the group order of G1/G2
// Public keys are in G2.

// feldman VSS protocol, with complaint mechanism, implements DKGState
type feldmanVSSQualState struct {
	// feldmanVSSstate state
	*feldmanVSSstate
	// complaints received against the dealer:
	// the key is the origin of the complaint
	// a complaint will be created if a complaint message or an answer was
	// broadcasted, a complaint will be checked only when both the
	// complaint message and the answer were broadcasted
	complaints map[index]*complaint
	// is the dealer disqualified
	disqualified bool
	// Timeout to receive shares and verification vector
	// - if a share is not received before this timeout a complaint will be formed
	// - if the verification is not received before this timeout,
	// dealer is disqualified
	sharesTimeout bool
	// Timeout to receive complaints
	// all complaints received after this timeout are ignored
	complaintsTimeout bool
}

// these data are required to justify a slashing
type complaint struct {
	received       bool
	answerReceived bool
	answer         scalar
}

// NewFeldmanVSSQual creates a new instance of a Feldman VSS protocol
// with a qualification mechanism.
//
// An instance is run by a single participant and is usable for only one protocol.
// In order to run the protocol again, a new instance needs to be created
//
// The function returns:
//   - (nil, InvalidInputsError) if:
//   - size if not in [DKGMinSize, DKGMaxSize]
//   - threshold is not in [MinimumThreshold, size-1]
//   - myIndex is not in [0, size-1]
//   - dealerIndex is not in [0, size-1]
//   - (dkgInstance, nil) otherwise
func NewFeldmanVSSQual(size int, threshold int, myIndex int,
	processor DKGProcessor, dealerIndex int) (DKGState, error) {

	common, err := newDKGCommon(size, threshold, myIndex, processor, dealerIndex)
	if err != nil {
		return nil, err
	}

	fvss := &feldmanVSSstate{
		dkgCommon:   common,
		dealerIndex: index(dealerIndex),
	}
	fvssq := &feldmanVSSQualState{
		feldmanVSSstate: fvss,
		disqualified:    false,
	}
	fvssq.init()
	return fvssq, nil
}

func (s *feldmanVSSQualState) init() {
	s.feldmanVSSstate.init()
	s.complaints = make(map[index]*complaint)
}

// NextTimeout sets the next protocol timeout
// This function needs to be called twice by every participant in
// the Feldman VSS Qual protocol.
// The first call is a timeout for sharing the private shares.
// The second call is a timeout for broadcasting the complaints.
//
// The returned erorr is :
//   - dkgInvalidStateTransitionError if the DKG instance was not running.
//   - dkgInvalidStateTransitionError if the DKG instance already called the 2 required timeouts.
//   - nil otherwise.
func (s *feldmanVSSQualState) NextTimeout() error {
	if !s.running {
		return dkgInvalidStateTransitionErrorf("dkg protocol %d is not running", s.myIndex)
	}
	if s.complaintsTimeout {
		return dkgInvalidStateTransitionErrorf("the next timeout should be to end DKG protocol")
	}

	// if dealer is already disqualified, there is nothing to do
	if s.disqualified {
		if !s.sharesTimeout {
			s.sharesTimeout = true
			return nil
		} else {
			s.complaintsTimeout = true
			return nil
		}
	}

	if !s.sharesTimeout {
		s.setSharesTimeout()
		return nil
	} else {
		s.setComplaintsTimeout()
		return nil
	}
}

// End ends the protocol in the current participant.
// This is also a timeout to receiving all complaint answers.
// It returns the finalized public data and participant private key share:
//  1. the group public key corresponding to the group secret key
//  2. all the public key shares corresponding to the participants private key shares.
//  3. the finalized private key which is the current participant's own private key share
//  4. Error Returns:
//     - dkgFailureError if the dealer was disqualified.
//     - dkgFailureError if the public key share or group public key is identity.
//     - dkgInvalidStateTransition if Start() was not called, or NextTimeout() was not called twice
//     - nil otherwise.
func (s *feldmanVSSQualState) End() (PrivateKey, PublicKey, []PublicKey, error) {
	if !s.running {
		return nil, nil, nil, dkgInvalidStateTransitionErrorf("dkg protocol %d is not running", s.myIndex)
	}
	if !s.sharesTimeout || !s.complaintsTimeout {
		return nil, nil, nil,
			dkgInvalidStateTransitionErrorf("%d: two timeouts should be set before ending dkg", s.myIndex)
	}
	s.running = false
	// check if a complaint has remained without an answer
	// a dealer is disqualified if a complaint was never answered
	if !s.disqualified {
		for complainer, c := range s.complaints {
			if c.received && !c.answerReceived {
				s.disqualified = true
				s.processor.Disqualify(int(s.dealerIndex),
					fmt.Sprintf("complaint from (%d) was not answered",
						complainer))
				break
			}
		}
	}

	// If the dealer is disqualified, all keys are ignored
	// otherwise, the keys are valid
	if s.disqualified {
		return nil, nil, nil, dkgFailureErrorf("dealer is disqualified")
	}

	// private key of the current participant
	x := newPrKeyBLSBLS12381(&s.x)

	// Group public key
	Y := newPubKeyBLSBLS12381(&s.vA[0])

	// The participants public keys
	y := make([]PublicKey, s.size)
	for i, p := range s.y {
		y[i] = newPubKeyBLSBLS12381(&p)
	}

	// check if current public key share or group public key is identity.
	// In that case all signatures generated by the key are invalid (as stated by the BLS IETF
	//	draft) to avoid equivocation issues.
	// TODO: update generateShares to make sure no public key share is identity AND
	// update receiveVector function to disqualify the dealer if any public key share
	// is identity, only when FeldmanVSSQ is not a building primitive of Joint-Feldman
	if (&s.x).isZero() {
		s.disqualified = true
		return nil, nil, nil, dkgFailureErrorf("private key share is identity and therefore invalid")
	}
	if Y.isIdentity {
		s.disqualified = true
		return nil, nil, nil, dkgFailureErrorf("group private key is identity and is therefore invalid")
	}
	return x, Y, y, nil
}

const (
	complaintSize       = 1
	complaintAnswerSize = 1 + PrKeyLenBLSBLS12381
)

// HandleBroadcastMsg processes a new broadcasted message received by the current participant.
// orig is the message origin index
//
// The function returns:
//   - dkgInvalidStateTransitionError if the instance is not running
//   - invalidInputsError if `orig` is not valid (in [0, size-1])
//   - nil otherwise
func (s *feldmanVSSQualState) HandleBroadcastMsg(orig int, msg []byte) error {
	if !s.running {
		return dkgInvalidStateTransitionErrorf("dkg is not running")
	}

	if orig >= s.Size() || orig < 0 {
		return invalidInputsErrorf(
			"wrong origin input, should be less than %d, got %d",
			s.Size(),
			orig)
	}

	// In case a message is received by the origin participant,
	// the message is just ignored
	if s.myIndex == index(orig) {
		return nil
	}

	// if dealer is already disqualified, ignore the message
	if s.disqualified {
		return nil
	}

	if len(msg) == 0 {
		if index(orig) == s.dealerIndex {
			s.disqualified = true
		}
		s.processor.Disqualify(orig, "received broadcast is empty")
		return nil
	}

	switch dkgMsgTag(msg[0]) {
	case feldmanVSSVerifVec:
		s.receiveVerifVector(index(orig), msg[1:])
	case feldmanVSSComplaint:
		s.receiveComplaint(index(orig), msg[1:])
	case feldmanVSSComplaintAnswer:
		s.receiveComplaintAnswer(index(orig), msg[1:])
	default:
		if index(orig) == s.dealerIndex {
			s.disqualified = true
		}
		s.processor.Disqualify(orig,
			fmt.Sprintf("invalid broadcast header, got %d",
				dkgMsgTag(msg[0])))
	}
	return nil
}

// HandlePrivateMsg processes a new private message received by the current participant.
// orig is the message origin index.
//
// The function returns:
//   - dkgInvalidStateTransitionError if the instance is not running
//   - invalidInputsError if `orig` is not valid (in [0, size-1])
//   - nil otherwise
func (s *feldmanVSSQualState) HandlePrivateMsg(orig int, msg []byte) error {
	if !s.running {
		return dkgInvalidStateTransitionErrorf("dkg is not running")
	}
	if orig >= s.Size() || orig < 0 {
		return invalidInputsErrorf(
			"invalid origin, should be positive less than %d, got %d",
			s.Size(),
			orig)
	}

	// In case a private message is received by the origin participant,
	// the message is just ignored
	if s.myIndex == index(orig) {
		return nil
	}

	// if dealer is already disqualified, ignore the message
	if s.disqualified {
		return nil
	}

	// forward all the message to receiveShare because any private message
	// has to be a private share
	s.receiveShare(index(orig), msg)

	return nil
}

// ForceDisqualify forces a participant to get disqualified
// for a reason outside of the DKG protocol
// The caller should make sure all honest participants call this function,
// otherwise, the protocol can be broken
//
// The function returns:
//   - dkgInvalidStateTransitionError if the instance is not running
//   - invalidInputsError if `orig` is not valid (in [0, size-1])
//   - nil otherwise
func (s *feldmanVSSQualState) ForceDisqualify(participant int) error {
	if !s.running {
		return dkgInvalidStateTransitionErrorf("dkg is not running")
	}
	if participant >= s.Size() || participant < 0 {
		return invalidInputsErrorf(
			"invalid origin input, should be less than %d, got %d",
			s.Size(), participant)
	}
	if index(participant) == s.dealerIndex {
		s.disqualified = true
	}
	return nil
}

// The function does not check the call respects the machine
// state transition of feldmanVSSQual. The calling function must make sure this call
// is valid.
func (s *feldmanVSSQualState) setSharesTimeout() {
	s.sharesTimeout = true
	// if verif vector is not received, disqualify the dealer
	if !s.vAReceived {
		s.disqualified = true
		s.processor.Disqualify(int(s.dealerIndex),
			"verification vector was not received")
		return
	}
	// if share is not received, make a complaint
	if !s.xReceived {
		s.buildAndBroadcastComplaint()
	}
}

// The function does not check the call respects the machine
// state transition of feldmanVSSQual. The calling function must make sure this call
// is valid.
func (s *feldmanVSSQualState) setComplaintsTimeout() {
	s.complaintsTimeout = true
	// if more than t complaints are received, the dealer is disqualified
	// regardless of the answers.
	// (at this point, all answered complaints should have been already received)
	// (i.e there is no complaint with (!c.received && c.answerReceived)
	if len(s.complaints) > s.threshold {
		s.disqualified = true
		s.processor.Disqualify(int(s.dealerIndex),
			fmt.Sprintf("there are %d complaints, they exceeded the threshold %d",
				len(s.complaints), s.threshold))
	}
}

func (s *feldmanVSSQualState) receiveShare(origin index, data []byte) {
	// only accept private shares from the dealer.
	if origin != s.dealerIndex {
		return
	}

	// check the share timeout
	if s.sharesTimeout {
		s.processor.FlagMisbehavior(int(origin),
			"private share is received after the shares timeout")
		return
	}

	if s.xReceived {
		s.processor.FlagMisbehavior(int(origin),
			"private share was already received")
		return
	}

	// at this point, tag private share is received
	s.xReceived = true

	// private message general check
	if len(data) == 0 || dkgMsgTag(data[0]) != feldmanVSSShare {
		s.buildAndBroadcastComplaint()
		s.processor.FlagMisbehavior(int(origin),
			fmt.Sprintf("private share should be non-empty and first byte should be %d, received %#x",
				feldmanVSSShare, data))
		return
	}

	// consider the remaining data from message
	data = data[1:]

	if (len(data)) != shareSize {
		s.buildAndBroadcastComplaint()
		s.processor.FlagMisbehavior(int(origin),
			fmt.Sprintf("invalid share size, expects %d, got %d",
				shareSize, len(data)))
		return
	}
	// read the participant private share
	err := readScalarFrStar(&s.x, data)
	if err != nil {
		s.buildAndBroadcastComplaint()
		s.processor.FlagMisbehavior(int(origin),
			fmt.Sprintf("invalid share value %x: %s", data, err))
		return
	}

	if s.vAReceived {
		if !s.verifyShare() {
			// build a complaint
			s.buildAndBroadcastComplaint()
		}
	}
}

func (s *feldmanVSSQualState) receiveVerifVector(origin index, data []byte) {
	// only accept the verification vector from the dealer.
	if origin != s.dealerIndex {
		return
	}

	// check the share timeout
	if s.sharesTimeout {
		s.processor.FlagMisbehavior(int(origin),
			"verification vector received after the shares timeout")
		return
	}

	if s.vAReceived {
		s.processor.FlagMisbehavior(int(origin),
			"verification received was already received")
		return
	}
	s.vAReceived = true

	if len(data) != verifVectorSize*(s.threshold+1) {
		s.disqualified = true
		s.processor.Disqualify(int(origin),
			fmt.Sprintf("invalid verification vector size, expects %d, got %d",
				verifVectorSize*(s.threshold+1), len(data)))
		return
	}
	// read the verification vector
	s.vA = make([]pointE2, s.threshold+1)
	err := readVerifVector(s.vA, data)
	if err != nil {
		s.disqualified = true
		s.processor.Disqualify(int(origin),
			fmt.Sprintf("reading the verification vector failed:%s", err))
		return
	}

	s.y = make([]pointE2, s.size)
	// compute all public keys
	s.computePublicKeys()

	// check the (already) registered complaints
	for complainer, c := range s.complaints {
		if c.received && c.answerReceived {
			if s.checkComplaint(complainer, c) {
				s.disqualified = true
				s.processor.Disqualify(int(s.dealerIndex),
					fmt.Sprintf("verification vector received: a complaint answer to (%d) is invalid, answer is %s, computed key is %s",
						complainer, &c.answer, &s.y[complainer]))
				return
			}
		}
	}
	// check the private share
	if s.xReceived {
		if !s.verifyShare() {
			s.buildAndBroadcastComplaint()
		}
	}
}

// build a complaint against the dealer, add it to the local
// complaint map and broadcast it
func (s *feldmanVSSQualState) buildAndBroadcastComplaint() {
	var logMsg string
	if s.vAReceived && s.xReceived {
		logMsg = fmt.Sprintf("building a complaint, share is %s, computed public key is %s",
			&s.x, &s.y[s.myIndex])
	} else {
		logMsg = "building a complaint"
	}
	s.processor.FlagMisbehavior(int(s.dealerIndex), logMsg)
	s.complaints[s.myIndex] = &complaint{
		received:       true,
		answerReceived: false,
	}
	data := []byte{byte(feldmanVSSComplaint), byte(s.dealerIndex)}
	s.processor.Broadcast(data)
}

// build a complaint answer, add it to the local
// complaint map and broadcast it
func (s *feldmanVSSQualState) buildAndBroadcastComplaintAnswer(complainee index) {
	data := make([]byte, complaintAnswerSize+1)
	data[0] = byte(feldmanVSSComplaintAnswer)
	data[1] = byte(complainee)
	frPolynomialImage(data[2:], s.a, complainee+1, nil)
	s.complaints[complainee].answerReceived = true
	s.processor.Broadcast(data)
}

// assuming a complaint and its answer were both received, this function returns:
// - false if the complaint answer is correct
// - true if the complaint answer is not correct
func (s *feldmanVSSQualState) checkComplaint(complainer index, c *complaint) bool {
	// check y[complainer] == share.G2
	return C.G2_check_log(
		(*C.Fr)(&c.answer),
		(*C.E2)(&s.y[complainer])) == 0
}

// data = |complainee|
func (s *feldmanVSSQualState) receiveComplaint(origin index, data []byte) {
	// check the complaint timeout
	if s.complaintsTimeout {
		s.processor.FlagMisbehavior(int(origin),
			"complaint received after the complaint timeout")
		return
	}

	if len(data) != complaintSize {
		// only the dealer of the instance gets disqualified
		if origin == s.dealerIndex {
			s.disqualified = true
			s.processor.Disqualify(int(origin),
				fmt.Sprintf("invalid complaint size, expects %d, got %d",
					complaintSize, len(data)))
		}
		return
	}

	// the byte encodes the complainee
	complainee := index(data[0])

	// validate the complainee value
	if int(complainee) >= s.size {
		// only the dealer of the instance gets disqualified
		if origin == s.dealerIndex {
			s.disqualified = true
			s.processor.Disqualify(int(origin),
				fmt.Sprintf("invalid complainee, should be less than %d, got %d",
					s.size, complainee))
		}
		return
	}

	// if the complaint is coming from the dealer, ignore it
	if origin == s.dealerIndex {
		return
	}

	// if the complainee is not the dealer, ignore the complaint
	if complainee != s.dealerIndex {
		return
	}

	c, ok := s.complaints[origin]
	// if the complaint is new, add it
	if !ok {
		s.complaints[origin] = &complaint{
			received:       true,
			answerReceived: false,
		}
		// if the complainee is the current participant, prepare an answer
		if s.myIndex == s.dealerIndex {
			s.buildAndBroadcastComplaintAnswer(origin)
		}
		return
	}
	// complaint is not new in the map
	// check if the complaint has been already received
	if c.received {
		s.processor.FlagMisbehavior(int(origin),
			"complaint was already received")
		return
	}
	c.received = true
	// answerReceived flag check is a sanity check
	if s.vAReceived && c.answerReceived && s.myIndex != s.dealerIndex {
		s.disqualified = s.checkComplaint(origin, c)
		if s.disqualified {
			s.processor.Disqualify(int(s.dealerIndex),
				fmt.Sprintf("complaint received: answer to (%d) is invalid, answer is %s, computed public key is %s",
					origin, &c.answer, &s.y[origin]))
		}
		return
	}
}

// answer = |complainer| private share |
func (s *feldmanVSSQualState) receiveComplaintAnswer(origin index, data []byte) {
	// check for invalid answers
	if origin != s.dealerIndex {
		return
	}

	// check the answer format
	if len(data) != complaintAnswerSize {
		s.disqualified = true
		s.processor.Disqualify(int(s.dealerIndex),
			fmt.Sprintf("the complaint answer has an invalid length, expects %d, got %d",
				complaintAnswerSize, len(data)))
		return
	}

	// first byte encodes the complainee
	complainer := index(data[0])
	if int(complainer) >= s.size {
		s.disqualified = true
		s.processor.Disqualify(int(origin),
			fmt.Sprintf("complainer value is invalid, should be less that %d, got %d",
				s.size, int(complainer)))
		return
	}

	c, ok := s.complaints[complainer]
	// if the complaint is new, add it
	if !ok {
		s.complaints[complainer] = &complaint{
			received:       false,
			answerReceived: true,
		}

		// read the complainer private share
		err := readScalarFrStar(&s.complaints[complainer].answer, data[1:])
		if err != nil {
			s.disqualified = true
			s.processor.Disqualify(int(s.dealerIndex),
				fmt.Sprintf("invalid complaint answer value %x: %s", data, err))
			return
		}
		return
	}
	// complaint is not new in the map
	// check if the answer has been already received
	if c.answerReceived {
		s.processor.FlagMisbehavior(int(origin),
			"complaint answer was already received")
		return
	}
	c.answerReceived = true

	// flag check is a sanity check
	if c.received {
		// read the complainer private share
		err := readScalarFrStar(&c.answer, data[1:])
		if err != nil {
			s.disqualified = true
			s.processor.Disqualify(int(s.dealerIndex),
				fmt.Sprintf("invalid complaint answer value %x: %s", data, err))
			return
		}
		if s.vAReceived {
			s.disqualified = s.checkComplaint(complainer, c)
			if s.disqualified {
				s.processor.Disqualify(int(s.dealerIndex),
					fmt.Sprintf("complaint answer received: answer to (%d) is invalid, answer is %s, computed key is %s",
						complainer, &c.answer, &s.y[complainer]))
			}
		}

		// fix the share of the current participant if the complaint is invalid
		if !s.disqualified && complainer == s.myIndex {
			s.x = c.answer
		}
	}
}
