package dto

type DownloadFileRequest struct {
	Url          string `json:"url" binding:"required,url_validator"`
	SaveFilename string `json:"saveFilename"`
}

type DownloadFileResponse struct {
	ID           uint   `json:"id" binding:"required"`
	Url          string `json:"url" binding:"required,url_validator"`
	SaveFilename string `json:"saveFilename" binding:"required"`
}

type DeleteDownloadFileRequest struct {
	ID uint `json:"id" binding:"required"`
}

type DeleteDownloadFileResponse struct {
	ID uint `json:"id" binding:"required"`
}

type DownloadProgressRequest struct {
	ID uint `json:"id" binding:"required"`
}

type DownloadProgressResponse struct {
	Url      string `json:"url" binding:"required,url_validator"`
	Progress string `json:"progress" binding:"required"`
}
