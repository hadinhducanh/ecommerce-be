package dto

type UploadImageRequest struct {
	Folder string `form:"folder"` // Folder trên Cloudinary (mặc định: 'ecommerce')
}

type UploadImageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		URL          string `json:"url"`
		OriginalName string `json:"originalName"`
		Size         int64  `json:"size"`
		Mimetype     string `json:"mimetype"`
	} `json:"data"`
}

type DeleteImageRequest struct {
	URL string `json:"url" binding:"required"` // Cloudinary URL cần xóa
}

type DeleteImageResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		URL      string `json:"url"`
		PublicID string `json:"publicId"`
	} `json:"data"`
}
