package handler

import (
	"embeck/model"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UploadTeamLogo uploads team logo image
// @Summary Upload team logo
// @Description Upload team logo image (PNG, JPG, JPEG)
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Team logo image file"
// @Success 200 {object} model.UploadResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/upload/team-logo [post]
func UploadTeamLogo(c *fiber.Ctx) error {
	// Parse multipart form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "bad_request",
			Message: "No file uploaded",
		})
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	if !allowedExts[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_file_type",
			Message: "Only JPG, JPEG, and PNG files are allowed",
		})
	}

	// Validate file size (max 5MB)
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if file.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "file_too_large",
			Message: "File size must be less than 5MB",
		})
	}

	// Generate unique filename
	uniqueID := uuid.New().String()
	timestamp := time.Now().Format("20060102_150405")
	newFileName := fmt.Sprintf("team_logo_%s_%s%s", timestamp, uniqueID[:8], ext)

	// Create upload directory if it doesn't exist
	uploadDir := "./uploads/team_logos"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "server_error",
			Message: "Failed to create upload directory",
		})
	}

	// Save file
	filePath := filepath.Join(uploadDir, newFileName)
	if err := saveUploadedFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "save_failed",
			Message: "Failed to save uploaded file",
		})
	}

	// Generate file URL (relative path for serving)
	fileURL := fmt.Sprintf("/uploads/team_logos/%s", newFileName)

	return c.Status(fiber.StatusOK).JSON(model.UploadResponse{
		Message:  "File uploaded successfully",
		FileURL:  fileURL,
		FileName: newFileName,
	})
}

// UploadPlayerAvatar uploads player avatar image
// @Summary Upload player avatar
// @Description Upload player avatar image (PNG, JPG, JPEG)
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Player avatar image file"
// @Success 200 {object} model.UploadResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 500 {object} model.ErrorResponse
// @Router /api/upload/player-avatar [post]
func UploadPlayerAvatar(c *fiber.Ctx) error {
	// Parse multipart form
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "bad_request",
			Message: "No file uploaded",
		})
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}

	if !allowedExts[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "invalid_file_type",
			Message: "Only JPG, JPEG, and PNG files are allowed",
		})
	}

	// Validate file size (max 2MB for avatars)
	maxSize := int64(2 * 1024 * 1024) // 2MB
	if file.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(model.ErrorResponse{
			Error:   "file_too_large",
			Message: "File size must be less than 2MB",
		})
	}

	// Generate unique filename
	uniqueID := uuid.New().String()
	timestamp := time.Now().Format("20060102_150405")
	newFileName := fmt.Sprintf("player_avatar_%s_%s%s", timestamp, uniqueID[:8], ext)

	// Create upload directory if it doesn't exist
	uploadDir := "./uploads/player_avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "server_error",
			Message: "Failed to create upload directory",
		})
	}

	// Save file
	filePath := filepath.Join(uploadDir, newFileName)
	if err := saveUploadedFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.ErrorResponse{
			Error:   "save_failed",
			Message: "Failed to save uploaded file",
		})
	}

	// Generate file URL (relative path for serving)
	fileURL := fmt.Sprintf("/uploads/player_avatars/%s", newFileName)

	return c.Status(fiber.StatusOK).JSON(model.UploadResponse{
		Message:  "File uploaded successfully",
		FileURL:  fileURL,
		FileName: newFileName,
	})
}

// saveUploadedFile saves the uploaded file to disk
func saveUploadedFile(fileHeader *multipart.FileHeader, destPath string) error {
	// Open uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Create destination file
	dst, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy file content
	_, err = io.Copy(dst, src)
	return err
}
