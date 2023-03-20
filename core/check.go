package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
	"fmt"
	"time"
)

func (v *Verifier) CheckSpeed(er deserialize.ElectionRecord) {
	start := time.Now()
	for _, ballot := range er.SubmittedBallots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.BallotSelections {
				v.mulP(&selection.Ciphertext.Data, &selection.Ciphertext.Pad)
			}
		}
	}
	fmt.Println("Current approach took: " + time.Since(start).String())

	start = time.Now()
	for _, ballot := range er.SubmittedBallots {
		for _, contest := range ballot.Contests {
			for _, selection := range contest.BallotSelections {
				var xd schema.BigInt
				xd.Mod(&mul(&selection.Ciphertext.Data, &selection.Ciphertext.Pad).Int, &v.constants.P.Int)
			}
		}
	}
	fmt.Println("Another approach took: " + time.Since(start).String())
}
