package config

type LLM struct {
	OpenAIAPIKey       string   `yaml:"open_ai_api_key"`
	OpenAIModelName    string   `json:"open_ai_model_name"`
	EmbeddingModelName string   `yaml:"embedding_model_name"`
	VectorDBHost       string   `yaml:"vector_db_host"`
	Weaviate           Weaviate `yaml:"weaviate"`
}

type Weaviate struct {
	Scheme string `yaml:"scheme"`
	Host   string `yaml:"host"`
	Index  string `yaml:"index"`
}
