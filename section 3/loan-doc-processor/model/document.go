package model

type DocumentJob struct {
	ApplicationID string
	DocumentType  string
	FilePath      string
	Priority      int
	Callback      func(result ProcessingResult)
}

type ProcessingResult struct {
	ApplicationID string
	DocumentType  string
	Status        string
	ExtractedData map[string]string
	Error         error
}
