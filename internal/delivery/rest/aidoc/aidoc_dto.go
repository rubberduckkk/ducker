package aidoc

type AddDocumentRequest struct {
	Texts []string `json:"texts" binding:"required"`
}

type QueryDocumentsRequest struct {
	Content string `json:"content" binding:"required"`
}
