package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateCorrectnessOfDecryptionsForSpoiledBallots(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 12, "Correctness of decryptions for spoiled ballots")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	g := er.ElectionConstants.Generator
	q := &er.ElectionConstants.SmallPrime
	ctx := er.CiphertextElectionRecord
	k := ctx.ElgamalPublicKey
	extendedBaseHash := ctx.CryptoExtendedBaseHash

	for _, spoiledBallot := range er.SpoiledBallots {
		ballot, foundBallot := findBallot(spoiledBallot.Name, er)

		if !foundBallot {
			helper.addCheck("(12) Could not locate submitted ballot for spoiled ballotID: "+spoiledBallot.Name, false)
		}

		for _, contest := range spoiledBallot.Contests {
			for _, selection := range contest.Selections {
				encryptedSelection := findEncryptedSelectionForBallot(contest.ObjectId, selection.ObjectId, ballot)
				alpha := encryptedSelection.Ciphertext.Pad
				beta := encryptedSelection.Ciphertext.Data

				helper.addCheck("(12.A) The challenge is not valid.", v.isInRange(selection.Proof.Response))

				// Computing values needed for 12.C
				m := v.mulP(&beta, v.invP(&selection.Value))
				a := v.mulP(v.powP(&g, &selection.Proof.Response), v.powP(&k, &selection.Proof.Challenge))
				b := v.mulP(v.powP(&alpha, &selection.Proof.Response), v.powP(m, &selection.Proof.Challenge))
				hash := crypto.Hash(q, "30")
				hash = crypto.Hash(q, hash, extendedBaseHash)
				hash = crypto.Hash(q, hash, k)
				hash = crypto.Hash(q, hash, alpha)
				hash = crypto.Hash(q, hash, beta)
				hash = crypto.Hash(q, hash, a)
				hash = crypto.Hash(q, hash, b)
				hash = crypto.Hash(q, hash, m)

				helper.addCheck("(12.C) The challenge is not computed correctly.", selection.Proof.Challenge.Compare(hash))
			}
		}
	}

	v.helpers[helper.VerificationStep] = helper
}

func findEncryptedSelectionForBallot(contestID, selectionID string, ballot schema.SubmittedBallot) schema.BallotSelection {
	for _, contest := range ballot.Contests {
		if contestID == contest.ObjectId {
			for _, selection := range contest.BallotSelections {
				if selectionID == selection.ObjectId {
					return selection
				}
			}
		}
	}
	return schema.BallotSelection{}
}

func findBallot(name string, er *deserialize.ElectionRecord) (schema.SubmittedBallot, bool) {
	for _, ballot := range er.SubmittedBallots {
		if name == ballot.Code {
			return ballot, true
		}
	}

	return schema.SubmittedBallot{}, false
}
