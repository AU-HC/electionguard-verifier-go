package core

import (
	"electionguard-verifier-go/crypto"
	"electionguard-verifier-go/schema"
	"time"
)

func (v *Verifier) validateParameters(er *schema.ElectionRecord) {
	helper := MakeValidationHelper(v.logger, 1, "Election parameters are correct")
	defer v.wg.Done()
	defer helper.measureTimeToValidateStep(time.Now())

	p := er.ElectionConstants.LargePrime
	q := er.ElectionConstants.SmallPrime
	r := er.ElectionConstants.Cofactor
	g := er.ElectionConstants.Generator

	// Verifying specification version and election parameters
	helper.addCheck("(1.A) Specification is not the same as verifier specification.", er.Manifest.SpecVersion == "1.0")
	helper.addCheck("(1.B) Large prime is not correct.", v.constants.P.Compare(&p))
	helper.addCheck("(1.C) Small prime is not correct.", v.constants.Q.Compare(&q))
	helper.addCheck("(1.D) Cofactor is not correct.", v.constants.C.Compare(&r))
	helper.addCheck("(1.E) Generator is not correct.", v.constants.G.Compare(&g))

	// Verifying election base hash
	ver := schema.MakeBigIntFromString("76322E3000000000000000000000000000000000000000000000000000000000", 16) // hardcoded value (ver 2.0.0) with 27 empty bytes appended
	hashP := crypto.Hash(&g, ver, "00", &p, &q, &g)
	hashM := crypto.Hash(&g, hashP, "01", schema.MakeBigIntFromString("2388DAE2DAD9E79ED64BD2C5401C361BDA77B9019EA8C6F416568185679B4854", 16))
	hashQ := crypto.Hash(&g, hashP, "02", 5, 3, hashM)
	helper.addCheck("(1.I) Base hash has not been computed correctly.", hashQ.Compare(&er.CiphertextElectionRecord.CryptoBaseHash))

	v.helpers[helper.VerificationStep] = helper
}
