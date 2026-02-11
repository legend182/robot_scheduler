package dto

// PCDFileUploadTokenRequest 获取点云地图上传凭证请求
type PCDFileUploadTokenRequest struct {
	FileName string `json:"fileName" binding:"required"`   // 原始文件名
	Size     int64  `json:"size" binding:"required,min=0"` // 文件大小(字节)
}

// PCDFileUploadTokenResponse 获取点云地图上传凭证响应
type PCDFileUploadTokenResponse struct {
	UploadURL string `json:"uploadUrl"` // 预签名 PUT URL
	Bucket    string `json:"bucket"`    // MinIO Bucket 名
	ObjectKey string `json:"objectKey"` // 对象 Key（前端后续要写入 MinioPath）
	ExpireAt  int64  `json:"expireAt"`  // 过期时间戳（秒）
}

