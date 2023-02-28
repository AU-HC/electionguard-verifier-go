package core

import (
	"electionguard-verifier-go/schema"
)

func getGuardianPublicKey(id string, guardians []schema.Guardian) *schema.BigInt {
	for _, guardian := range guardians {
		if guardian.GuardianId == id {
			return &guardian.ElectionPublicKey
		}
	}

	return &schema.BigInt{}
}

// TODO: Fix these to return nothing if two of the same exist, and log a message (Also move them somewhere else -> helper_methods.go)
func getContest(objectID string, contests []schema.Contest) schema.Contest {
	for _, contest := range contests {
		if objectID == contest.ObjectID {
			return contest
		}
	}
	return schema.Contest{}
}

func getBallotContest(objectID string, ballot schema.SubmittedBallot) schema.BallotContest {
	for _, contest := range ballot.Contests {
		if objectID == contest.ObjectId {
			return contest
		}
	}
	return schema.BallotContest{}
}

func getSelection(ballot schema.SubmittedBallot, contestId string, selectionId string) schema.Ciphertext {
	contest := getBallotContest(contestId, ballot)
	for _, selection := range contest.BallotSelections {
		if selection.ObjectId == selectionId {
			if !selection.IsPlaceholderSelection {
				return selection.Ciphertext
			}
		}
	}
	return makeOneCiphertext()
}

func makeOneCiphertext() schema.Ciphertext {
	var result schema.Ciphertext
	result.Pad = *schema.MakeBigIntFromString("1", 10)
	result.Data = *schema.MakeBigIntFromString("1", 10)
	return result
}
