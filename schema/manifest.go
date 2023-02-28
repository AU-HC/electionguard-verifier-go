package schema

type Manifest struct {
	ElectionScopeID    string             `json:"election_scope_id"`
	SpecVersion        string             `json:"spec_version"`
	Type               string             `json:"type"`
	StartDate          string             `json:"start_date"`
	EndDate            string             `json:"end_date"`
	GeopoliticalUnits  []GeopoliticalUnit `json:"geopolitical_units"`
	Parties            []Party            `json:"parties"`
	Candidates         []Candidate        `json:"candidates"`
	Contests           []Contest          `json:"contests"`
	BallotStyles       []BallotStyle      `json:"ballot_styles"`
	Name               Text               `json:"name"`
	ContactInformation Contact            `json:"contact_information"`
}

type Contest struct {
	ObjectID            string                    `json:"object_id"`
	SequenceOrder       int                       `json:"sequence_order"`
	ElectoralDistrictID string                    `json:"electoral_district_id"`
	VoteVariation       string                    `json:"vote_variation"`
	NumberElected       int                       `json:"number_elected"`
	VotesAllowed        int                       `json:"votes_allowed"`
	Name                string                    `json:"name"`
	BallotSelections    []ManifestBallotSelection `json:"ballot_selections"`
	BallotTitle         Text                      `json:"ballot_title"`
	BallotSubtitle      Text                      `json:"ballot_subtitle"`
}

type ManifestBallotSelection struct {
	ObjectID      string `json:"object_id"`
	SequenceOrder int    `json:"sequence_order"`
	CandidateID   string `json:"candidate_id"`
}

type GeopoliticalUnit struct {
	ObjectID           string  `json:"object_id"`
	Name               string  `json:"name"`
	Type               string  `json:"type"`
	ContactInformation Contact `json:"contact_information"`
}

type Party struct {
	ObjectID     string `json:"object_id"`
	Name         Text   `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Color        string `json:"color"`
	LogoURI      string `json:"logo_uri"`
}

type Candidate struct {
	ObjectID  string `json:"object_id"`
	Name      Text   `json:"name"`
	PartyID   string `json:"party_id"`
	ImageURI  string `json:"image_uri"`
	IsWriteIn bool   `json:"is_write_in"`
}

type BallotStyle struct {
	ObjectID            string   `json:"object_id"`
	GeopoliticalUnitIds []string `json:"geopolitical_unit_ids"`
	PartyIds            []string `json:"party_ids"`
	ImageURI            string   `json:"image_uri"`
}

type Text struct {
	Text                string `json:"value"`
	LanguageDesignation string `json:"language"`
}

type Information struct {
	Annotation string `json:"annotation"`
	Value      string `json:"value"`
}

type Contact struct {
	Address []string      `json:"address_line"`
	Email   []Information `json:"email"`
	Phone   []Information `json:"phone"`
	Name    string        `json:"name"`
}
