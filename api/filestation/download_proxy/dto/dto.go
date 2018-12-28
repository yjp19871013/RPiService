package dto

type AddDownloadTaskRequest struct {
	Url          string `json:"url" binding:"required,url_validator"`
	SaveFilename string `json:"saveFilename"`
}

type DownloadTask struct {
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

type GetDownloadProgressResponse struct {
	Progresses []DownloadProgress `json:"progresses" binding:"required"`
}

type GetAllTaskResponse struct {
	Tasks []DownloadTask `json:"tasks" binding:"required"`
}
