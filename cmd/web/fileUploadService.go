package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

func createDriveService(ctx context.Context, credentialsFile string) (*drive.Service, error) {
	// Read service account credentials
	credBytes, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file: %w", err)
	}

	// Create credentials from JSON
	creds, err := google.CredentialsFromJSON(ctx, credBytes, drive.DriveScope)
	if err != nil {
		return nil, fmt.Errorf("unable to create credentials from JSON: %w", err)
	}

	// Create Drive service client
	srv, err := drive.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		return nil, fmt.Errorf("unable to create Drive service: %w", err)
	}

	fmt.Println("successfully created Drive service")
	return srv, nil
}

// uploadImageToDrive uploads a specified file to Google Drive service
func uploadImageToDrive(srv *drive.Service, file multipart.File, fileName string) (string, error) {
	// Metadata for the uploaded file
	fileMetadata := &drive.File{
		Name:     fileName,
		MimeType: "image/*", // or image/png, image/jpeg, etc.
	}

	// Create a reader for the file content
	fileReader := file

	// Upload file to Google Drive
	uploadedFile, err := srv.Files.Create(fileMetadata).Media(fileReader).Do()
	if err != nil {
		return "", fmt.Errorf("failed to upload file to Google Drive: %w", err)
	}

	// Make the file public (so the link can be used in <img>)
	_, err = srv.Permissions.Create(uploadedFile.Id, &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}).Do()
	if err != nil {
		return "", fmt.Errorf("failed to make file public: %w", err)
	}
	return uploadedFile.Id, nil
}

// downloadFile downloads a specific uploaded to the drive by fileID
func downloadFile(srv *drive.Service, fileID, outputFilename string) error {
	// Get the file from Drive
	resp, err := srv.Files.Get(fileID).Download()
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	// Create the output file
	outFile, err := os.Create(outputFilename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Copy the content
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file content: %w", err)
	}

	return nil
}

// deleteFileFromDrive deletes a file from Google Drive given its file ID
func deleteFileFromDrive(srv *drive.Service, fileID string) error {
	err := srv.Files.Delete(fileID).Do()
	if err != nil {
		return fmt.Errorf("failed to delete file with ID %s: %w", fileID, err)
	}
	fmt.Printf("Successfully deleted file with ID: %s\n", fileID)
	return nil
}

func getCredentialsPath() (string, error) {
	if encoded := os.Getenv("GOOGLE_CREDENTIALS_BASE64"); encoded != "" {
		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return "", err
		}

		tmpFile := filepath.Join(os.TempDir(), "credentials.json")
		err = os.WriteFile(tmpFile, decoded, 0644)
		if err != nil {
			return "", err
		}
		return tmpFile, nil
	}

	return "credentials.json", nil
}
