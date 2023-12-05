package schema

type Manifest struct {
	ElectionScopeID string            `json:"election_scope_id"`
	SpecVersion     string            `json:"spec_version"`
	StartDate       string            `json:"start_date"`
	EndDate         string            `json:"end_date"`
	Contests        []ManifestContest `json:"contests"`
}

type ManifestContest struct {
	ObjectID         string                    `json:"object_id"`
	VotesAllowed     int                       `json:"votes_allowed"`
	BallotSelections []ManifestBallotSelection `json:"ballot_selections"`
}

type ManifestBallotSelection struct {
	ObjectID string `json:"object_id"`
}
