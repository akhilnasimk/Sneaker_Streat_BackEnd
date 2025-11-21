package cloudinary

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/akhilnasimk/SS_backend/internal/config"
	"github.com/akhilnasimk/SS_backend/internal/helpers"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// UploadResult represents the result of a single upload
type UploadResult struct {
	URL     string
	AltText string
	Error   error
}

// UploadOptions defines upload configuration
type UploadOptions struct {
	Folder  string
	Timeout time.Duration
}

// DefaultUploadOptions returns default upload settings
func DefaultUploadOptions(folder string) UploadOptions {
	return UploadOptions{
		Folder:  folder,
		Timeout: 30 * time.Second,
	}
}

// UploadMultiple uploads multiple files in parallel to Cloudinary
func UploadMultiple(ctx context.Context, files []*multipart.FileHeader, opts UploadOptions) ([]UploadResult, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided")
	}

	resultsChan := make(chan UploadResult, len(files))

	// Upload files in parallel
	for _, file := range files {
		go func(f *multipart.FileHeader) {
			fileReader, err := f.Open()
			if err != nil {
				resultsChan <- UploadResult{Error: err}
				return
			}
			defer fileReader.Close()

			resp, err := config.CLD.Upload.Upload(ctx, fileReader, uploader.UploadParams{
				Folder: opts.Folder,
			})
			if err != nil {
				resultsChan <- UploadResult{Error: err}
				return
			}

			resultsChan <- UploadResult{
				URL:     resp.SecureURL,
				AltText: f.Filename,
				Error:   nil,
			}
		}(file)
	}

	// Collect results
	var results []UploadResult
	var uploadedURLs []string
	var uploadErrors []error

	for i := 0; i < len(files); i++ {
		select {
		case result := <-resultsChan:
			results = append(results, result)
			if result.Error != nil {
				uploadErrors = append(uploadErrors, result.Error)
			} else {
				uploadedURLs = append(uploadedURLs, result.URL)
			}
		case <-ctx.Done():
			// Timeout occurred - cleanup uploaded files
			DeleteMultipleAsync(uploadedURLs)
			return nil, fmt.Errorf("upload timeout exceeded")
		}
	}

	// If any upload failed, cleanup successful uploads
	if len(uploadErrors) > 0 {
		DeleteMultipleAsync(uploadedURLs)
		return nil, fmt.Errorf("failed to upload %d images: %v", len(uploadErrors), uploadErrors[0])
	}

	return results, nil
}

// UploadSingle uploads a single file to Cloudinary
func UploadSingle(ctx context.Context, file *multipart.FileHeader, opts UploadOptions) (string, error) {
	fileReader, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer fileReader.Close()

	resp, err := config.CLD.Upload.Upload(ctx, fileReader, uploader.UploadParams{
		Folder: opts.Folder,
	})
	if err != nil {
		return "", fmt.Errorf("cloudinary upload failed: %w", err)
	}

	return resp.SecureURL, nil
}

// Delete removes a single file from Cloudinary
func Delete(url string) error {
	if url == "" {
		return nil
	}

	publicID := helpers.ExtractPublicID(url)
	if publicID == "" {
		return fmt.Errorf("invalid cloudinary URL")
	}

	_, err := config.CLD.Upload.Destroy(context.Background(), uploader.DestroyParams{
		PublicID: publicID,
	})
	return err
}

// DeleteAsync removes a file asynchronously (fire and forget)
func DeleteAsync(url string) {
	go func() {
		if err := Delete(url); err != nil {
			// Log error but don't block
			fmt.Printf("Failed to delete cloudinary asset %s: %v\n", url, err)
		}
	}()
}

// DeleteMultiple removes multiple files from Cloudinary synchronously
func DeleteMultiple(urls []string) error {
	if len(urls) == 0 {
		return nil
	}

	errorsChan := make(chan error, len(urls))

	for _, url := range urls {
		go func(u string) {
			if err := Delete(u); err != nil {
				errorsChan <- err
			} else {
				errorsChan <- nil
			}
		}(url)
	}

	var errors []error
	for i := 0; i < len(urls); i++ {
		if err := <-errorsChan; err != nil {
			errors = append(errors, err)
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to delete %d files", len(errors))
	}

	return nil
}

// DeleteMultipleAsync removes multiple files asynchronously (fire and forget)
func DeleteMultipleAsync(urls []string) {
	if len(urls) == 0 {
		return
	}

	go func() {
		for _, url := range urls {
			Delete(url)
		}
	}()
}
