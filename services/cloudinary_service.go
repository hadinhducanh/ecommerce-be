package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"ecommerce-be/dto"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

const (
	MaxFileSize   = 5 * 1024 * 1024 // 5MB
	DefaultFolder = "ecommerce"
)

type CloudinaryService struct {
	cld *cloudinary.Cloudinary
}

func NewCloudinaryService() (*CloudinaryService, error) {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		return nil, fmt.Errorf("cloudinary credentials không đầy đủ. Vui lòng kiểm tra .env file")
	}

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		return nil, fmt.Errorf("không thể khởi tạo Cloudinary: %w", err)
	}

	return &CloudinaryService{
		cld: cld,
	}, nil
}

// UploadImage upload image lên Cloudinary
func (s *CloudinaryService) UploadImage(file io.Reader, folder string, publicID string) (string, error) {
	if folder == "" {
		folder = DefaultFolder
	}

	ctx := context.Background()
	overwrite := true
	uploadParams := uploader.UploadParams{
		Folder:         folder,
		ResourceType:   "image",
		Overwrite:      &overwrite,
		Transformation: "q_auto,f_auto", // Quality auto, format auto (WebP, AVIF)
	}

	if publicID != "" {
		uploadParams.PublicID = publicID
	}

	result, err := s.cld.Upload.Upload(ctx, file, uploadParams)
	if err != nil {
		return "", fmt.Errorf("cloudinary upload failed: %w", err)
	}

	return result.SecureURL, nil
}

// DeleteImage xóa image từ Cloudinary
func (s *CloudinaryService) DeleteImage(publicID string) error {
	ctx := context.Background()
	_, err := s.cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete image from Cloudinary: %w", err)
	}
	return nil
}

// ExtractPublicId extract public ID từ Cloudinary URL
func (s *CloudinaryService) ExtractPublicId(url string) (string, error) {
	// Pattern: /v\d+/(.+)\.(jpg|jpeg|png|gif|webp|avif)
	re := regexp.MustCompile(`/v\d+/(.+)\.(jpg|jpeg|png|gif|webp|avif)`)
	matches := re.FindStringSubmatch(url)
	if len(matches) < 2 {
		return "", fmt.Errorf("URL không hợp lệ hoặc không phải Cloudinary URL")
	}
	return matches[1], nil
}

// ValidateFile validate file upload
func (s *CloudinaryService) ValidateFile(fileSize int64, mimetype string) error {
	if fileSize == 0 {
		return fmt.Errorf("không có file được upload")
	}

	if fileSize > MaxFileSize {
		return fmt.Errorf("file quá lớn. Kích thước tối đa là 5MB")
	}

	allowedTypes := []string{"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"}
	mimetypeLower := strings.ToLower(mimetype)
	for _, allowedType := range allowedTypes {
		if mimetypeLower == allowedType {
			return nil
		}
	}

	return fmt.Errorf("định dạng file không được hỗ trợ. Chỉ chấp nhận: jpg, jpeg, png, gif, webp")
}

// UploadImageWithResponse upload image với validation và response format chuẩn
func (s *CloudinaryService) UploadImageWithResponse(
	file io.Reader,
	fileSize int64,
	originalName string,
	mimetype string,
	folder string,
) (*dto.UploadImageResponse, error) {
	// Validate file
	if err := s.ValidateFile(fileSize, mimetype); err != nil {
		return nil, err
	}

	if folder == "" {
		folder = DefaultFolder
	}

	// Upload image
	url, err := s.UploadImage(file, folder, "")
	if err != nil {
		return nil, fmt.Errorf("upload thất bại: %w", err)
	}

	response := &dto.UploadImageResponse{
		Success: true,
		Message: "Upload ảnh thành công",
	}
	response.Data.URL = url
	response.Data.OriginalName = originalName
	response.Data.Size = fileSize
	response.Data.Mimetype = mimetype

	return response, nil
}

// DeleteImageWithResponse delete image với validation và response format chuẩn
func (s *CloudinaryService) DeleteImageWithResponse(url string) (*dto.DeleteImageResponse, error) {
	publicID, err := s.ExtractPublicId(url)
	if err != nil {
		return nil, fmt.Errorf("URL không hợp lệ hoặc không phải Cloudinary URL")
	}

	if err := s.DeleteImage(publicID); err != nil {
		return nil, fmt.Errorf("xóa ảnh thất bại: %w", err)
	}

	response := &dto.DeleteImageResponse{
		Success: true,
		Message: "Xóa ảnh thành công",
	}
	response.Data.URL = url
	response.Data.PublicID = publicID

	return response, nil
}
