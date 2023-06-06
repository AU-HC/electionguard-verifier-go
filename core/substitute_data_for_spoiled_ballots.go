package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateSubstituteDataForSpoiledBallots(er *deserialize.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 13, "Correctness of substitute data for spoiled ballots")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	// Split the slice of ballots into multiple slices
	ballots := er.SpoiledBallots
	chunkSize := 1
	if len(ballots) > v.verifierStrategy.getBallotSplitSize() {
		chunkSize = len(ballots) / v.verifierStrategy.getBallotSplitSize()
	}

	for i := 0; i < len(ballots); i += chunkSize {
		end := i + chunkSize

		if end > len(ballots) {
			end = len(ballots)
		}

		helper.wg.Add(1)
		go v.validateSubstituteDataForSpoiledBallotsForSlice(helper, ballots[i:end], er)
	}

	helper.wg.Wait()
	v.helpers[helper.VerificationStep] = helper
}

func (v *Verifier) validateSubstituteDataForSpoiledBallotsForSlice(helper *ValidationHelper, spoiledBallots []schema.SpoiledBallot, er *deserialize.ElectionRecord) {
	defer helper.wg.Done()

	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash

	for _, ballot := range spoiledBallots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.Selections {
				alpha := selection.Message.Pad
				beta := selection.Message.Data

				for _, share := range selection.Shares {
					for _, part := range share.RecoveredParts {
						mil := part.Share
						ail := part.Proof.Pad
						bil := part.Proof.Data
						cil := part.Proof.Challenge
						vil := part.Proof.Response

						missingGuardian := findGuardian(er, part.MissingGuardianIdentifier)
						guardian := findGuardian(er, part.GuardianIdentifier)

						sum := schema.IntToBigInt(1)
						for j, commitment := range missingGuardian.ElectionCommitments {
							temp := v.powP(&commitment, v.powP(schema.IntToBigInt(guardian.SequenceOrder), schema.IntToBigInt(j)))
							sum = v.mulP(sum, temp)
						}

						helper.addCheck(step13A, v.isInRange(vil))
						helper.addCheck(step13B1, v.isValidResidue(ail))
						helper.addCheck(step13B2, v.isValidResidue(bil))
						helper.addCheck(step13C, cil.Compare(crypto.HashElems(extendedBaseHash, alpha, beta, ail, bil, mil)))
						helper.addCheck(step13D, v.powP(v.constants.G, &vil).Compare(v.mulP(&ail, v.powP(sum, &cil))))
						helper.addCheck(step13E, v.powP(&alpha, &vil).Compare(v.mulP(&bil, v.powP(&mil, &cil))))
					}
				}
			}
		}
	}
}
