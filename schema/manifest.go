package schema

type Manifest struct {
	ElectionScopeID   string `json:"election_scope_id"`
	SpecVersion       string `json:"spec_version"`
	Type              string `json:"type"`
	StartDate         string `json:"start_date"`
	EndDate           string `json:"end_date"`
	GeopoliticalUnits []struct {
		ObjectID           string `json:"object_id"`
		Name               string `json:"name"`
		Type               string `json:"type"`
		ContactInformation struct {
			AddressLine []string `json:"address_line"`
			Email       []struct {
				Annotation string `json:"annotation"`
				Value      string `json:"value"`
			} `json:"email"`
			Phone []struct {
				Annotation string `json:"annotation"`
				Value      string `json:"value"`
			} `json:"phone"`
			Name string `json:"name"`
		} `json:"contact_information"`
	} `json:"geopolitical_units"`
	Parties []struct {
		ObjectID string `json:"object_id"`
		Name     struct {
			Text []struct {
				Value    string `json:"value"`
				Language string `json:"language"`
			} `json:"text"`
		} `json:"name"`
		Abbreviation string `json:"abbreviation"`
		Color        string `json:"color"`
		LogoURI      string `json:"logo_uri"`
	} `json:"parties"`
	Candidates []struct {
		ObjectID string `json:"object_id"`
		Name     struct {
			Text []struct {
				Value    string `json:"value"`
				Language string `json:"language"`
			} `json:"text"`
		} `json:"name"`
		PartyID   string      `json:"party_id"`
		ImageURI  interface{} `json:"image_uri"`
		IsWriteIn interface{} `json:"is_write_in"`
	} `json:"candidates"`
	Contests []struct {
		ObjectID            string `json:"object_id"`
		SequenceOrder       int    `json:"sequence_order"`
		ElectoralDistrictID string `json:"electoral_district_id"`
		VoteVariation       string `json:"vote_variation"`
		NumberElected       int    `json:"number_elected"`
		VotesAllowed        int    `json:"votes_allowed"`
		Name                string `json:"name"`
		BallotSelections    []struct {
			ObjectID      string `json:"object_id"`
			SequenceOrder int    `json:"sequence_order"`
			CandidateID   string `json:"candidate_id"`
		} `json:"ballot_selections"`
		BallotTitle struct {
			Text []struct {
				Value    string `json:"value"`
				Language string `json:"language"`
			} `json:"text"`
		} `json:"ballot_title"`
		BallotSubtitle struct {
			Text []struct {
				Value    string `json:"value"`
				Language string `json:"language"`
			} `json:"text"`
		} `json:"ballot_subtitle"`
	} `json:"contests"`
	BallotStyles []struct {
		ObjectID            string      `json:"object_id"`
		GeopoliticalUnitIds []string    `json:"geopolitical_unit_ids"`
		PartyIds            interface{} `json:"party_ids"`
		ImageURI            interface{} `json:"image_uri"`
	} `json:"ballot_styles"`
	Name struct {
		Text []struct {
			Value    string `json:"value"`
			Language string `json:"language"`
		} `json:"text"`
	} `json:"name"`
	ContactInformation struct {
		AddressLine []string `json:"address_line"`
		Email       []struct {
			Annotation string `json:"annotation"`
			Value      string `json:"value"`
		} `json:"email"`
		Phone []struct {
			Annotation string `json:"annotation"`
			Value      string `json:"value"`
		} `json:"phone"`
		Name string `json:"name"`
	} `json:"contact_information"`
}
