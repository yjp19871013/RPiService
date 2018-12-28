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

type FileInfo struct {
	ID           uint    `json:"id" binding:"required"`
	FileName     string  `json:"fileName" binding:"required"`
	CompleteDate string  `json:"completeDate" binding:"required"`
	SizeKb       float64 `json:"sizeKb" binding:"required"`
}

type GetAllFileInfosResponse struct {
	FileInfos []FileInfo `json:"fileInfos" binding:"required"`
}
