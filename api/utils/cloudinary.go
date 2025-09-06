package utils

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	client *cloudinary.Cloudinary
}

func NewCloudinaryService() (*CloudinaryService, error) {
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return nil, err
	}

	return &CloudinaryService{client: cld}, nil
}

func (cs *CloudinaryService) UploadImage(file multipart.File, fileName string) (string, error) {
	ctx := context.Background()

	uploadResult, err := cs.client.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:     fileName,
		ResourceType: "image",
		Folder:       "yt-backend/avatars",
	})
	if err != nil {
		return "", err
	}
	fmt.Println("file uploaded successfully..", uploadResult.SecureURL)
	return uploadResult.SecureURL, nil
}

func (cs *CloudinaryService) UploadVideo(file multipart.File, fileName string) (string, error) {
	ctx := context.Background()

	uploadResult, err := cs.client.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:     fileName,
		ResourceType: "video",
		Folder:       "yt-backend/videos",
	})
	if err != nil {
		return "", err
	}
	fmt.Println("video uploaded successfully..", uploadResult.SecureURL)
	return uploadResult.SecureURL, nil
}

func (cs *CloudinaryService) DeleteFile(publicID string, resourceType string) error {
	ctx := context.Background()

	_, err := cs.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: resourceType,
	})

	return err
}
