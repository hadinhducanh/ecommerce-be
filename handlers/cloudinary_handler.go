package handlers

import (
	"net/http"

	"ecommerce-be/dto"
	"ecommerce-be/services"

	"github.com/gin-gonic/gin"
)

type CloudinaryHandler struct {
	cloudinaryService *services.CloudinaryService
}

func NewCloudinaryHandler() (*CloudinaryHandler, error) {
	service, err := services.NewCloudinaryService()
	if err != nil {
		return nil, err
	}
	return &CloudinaryHandler{
		cloudinaryService: service,
	}, nil
}

// UploadImage xử lý request upload image
// @Summary Upload image lên Cloudinary
// @Description Upload image lên Cloudinary và trả về URL (Chỉ admin)
// @Tags cloudinary
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file"
// @Param folder formData string false "Folder trên Cloudinary (default: ecommerce)"
// @Success 200 {object} dto.UploadImageResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/cloudinary/upload-image [post]
func (h *CloudinaryHandler) UploadImage(c *gin.Context) {
	// Lấy file từ form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "không có file được upload",
		})
		return
	}

	// Lấy folder từ form (optional)
	folder := c.PostForm("folder")
	if folder == "" {
		folder = "ecommerce"
	}

	// Mở file
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "không thể mở file",
		})
		return
	}
	defer src.Close()

	// Upload image
	response, err := h.cloudinaryService.UploadImageWithResponse(
		src,
		file.Size,
		file.Filename,
		file.Header.Get("Content-Type"),
		folder,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteImage xử lý request xóa image
// @Summary Xóa image từ Cloudinary
// @Description Xóa image từ Cloudinary bằng URL (Chỉ admin)
// @Tags cloudinary
// @Accept json
// @Produce json
// @Param delete body dto.DeleteImageRequest true "Cloudinary URL"
// @Success 200 {object} dto.DeleteImageResponse
// @Failure 400 {object} map[string]string
// @Router /api/v1/cloudinary/delete-image [delete]
func (h *CloudinaryHandler) DeleteImage(c *gin.Context) {
	var req dto.DeleteImageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Dữ liệu không hợp lệ",
			"details": err.Error(),
		})
		return
	}

	response, err := h.cloudinaryService.DeleteImageWithResponse(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

