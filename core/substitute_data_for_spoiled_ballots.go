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

	extendedBaseHash := er.CiphertextElectionRecord.CryptoExtendedBaseHash

	for _, ballot := range er.SpoiledBallots {
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

						sum := schema.MakeBigIntFromInt(0)

						for i, commitment := range missingGuardian.ElectionCommitments {
							temp := v.powP(&commitment, v.powP(schema.MakeBigIntFromInt(guardian.SequenceOrder), schema.MakeBigIntFromInt(i)))
							sum = v.addP(sum, temp)
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

	v.helpers[helper.VerificationStep] = helper

}

func findGuardian(er *deserialize.ElectionRecord, guardianId string) schema.Guardian {
	for _, guardian := range er.Guardians {
		if guardian.GuardianId == guardianId {
			return guardian
		}
	}
	panic("Guardian not found!")
}
