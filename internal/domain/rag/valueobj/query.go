package valueobj

type QueryResult struct {
	Summary     string `json:"summary"`
	OriginalDoc string `json:"orig_doc"`
}
