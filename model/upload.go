package model

// UploadResponse represents response for file upload operations
type UploadResponse struct {
	Message  string `json:"message"`
	FileURL  string `json:"file_url"`
	FileName string `json:"file_name"`
}

// ErrorResponse represents error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
