package dto

type FileInfo struct {
	ID           uint    `json:"id" binding:"required"`
	FileName     string  `json:"fileName" binding:"required"`
	CompleteDate string  `json:"completeDate" binding:"required"`
	SizeKb       float64 `json:"sizeKb" binding:"required"`
}

type GetAllFileInfosResponse struct {
	FileInfos []FileInfo `json:"fileInfos" binding:"required"`
}

type DownloadFileResponse struct {
	StaticUrl string `json:"staticUrl" binding:"required"`
}

type DeleteFileResponse struct {
	ID uint `json:"id" binding:"required"`
}
