package core

import (
	"electionguard-verifier-go/deserialize"
	"electionguard-verifier-go/schema"
)

func contains(slice []schema.Contest, s string) bool {
	for _, con := range slice {
		if s == con.ObjectID {
			return true
		}
	}
	return false
}

func getGuardianPublicKey(id string, guardians []schema.Guardian) *schema.BigInt {
	for _, guardian := range guardians {
		if guardian.GuardianId == id {
			return &guardian.ElectionPublicKey
		}
	}

	return &schema.BigInt{}
}

func doesContestExistInManifest(objectID string, contests []schema.Contest) bool {
	for _, contest := range contests {
		if objectID == contest.ObjectID {
			return true
		}
	}

	return false
}

func getSpoiledBallot(objectID string, ballots []schema.SpoiledBallot) schema.SpoiledBallot {
	for _, ballot := range ballots {
		if ballot.ObjectId == objectID {
			return ballot
		}
	}

	return schema.SpoiledBallot{}
}

func isPlaceholderSelection(ballot schema.SubmittedBallot, objectID string) bool {
	for _, contest := range ballot.Contests {
		for _, selection := range contest.BallotSelections {
			if selection.ObjectId == objectID {
				return selection.IsPlaceholderSelection
			}
		}
	}

	// TODO: Log error
	return false
}

func doesManifestSelectionExist(objectID string, selections []schema.ManifestBallotSelection) bool {
	for _, selection := range selections {
		if selection.ObjectID == objectID {
			return true
		}
	}

	return false
}

func doesContestContainSelection(slice []schema.Contest, contestID, selectionID string) bool {
	for _, contest := range slice {
		if contest.ObjectID == contestID {
			for _, selection := range contest.BallotSelections {
				if selection.ObjectID == selectionID {
					return true
				}
			}
		}
	}

	return false
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

func findGuardian(er *deserialize.ElectionRecord, guardianId string) schema.Guardian {
	for _, guardian := range er.Guardians {
		if guardian.GuardianId == guardianId {
			return guardian
		}
	}
	panic("Guardian not found!")
}
