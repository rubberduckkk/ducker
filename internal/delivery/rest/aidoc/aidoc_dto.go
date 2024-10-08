package aidoc

type AddDocumentRequest struct {
	Texts []string `json:"texts"`
}

type QueryDocumentsRequest struct {
	Content string `json:"content"`
}

type QueryDocumentsResponse struct {
	Content string `json:"content"`
}
