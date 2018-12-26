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

type DeleteDownloadFileResponse struct {
	ID uint `json:"id" binding:"required"`
}

type DownloadProgress struct {
	ID       uint `json:"id" binding:"required"`
	Progress uint `json:"progress" binding:"required"`
}

type DownloadProgressResponse struct {
	Progresses []DownloadProgress
}
