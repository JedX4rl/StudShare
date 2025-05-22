package domain

type FileRequest struct {
	FileName    string `json:"file_name" validate:"required,min=1"`
	ContentType string `json:"content_type" validate:"required"`
}

type FileUploadResponse struct {
	FileName     string `json:"file_name"`
	PresignedURL string `json:"presigned_url"`
	PublicURL    string `json:"public_url"`
}
