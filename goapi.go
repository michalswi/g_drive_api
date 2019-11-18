package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

const (
	driveFolder = "audycje"
	sourceDir   = "/tmp/workspace"
)

var dirID string
var dirName string

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}
	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func createDir(service *drive.Service, name string, parentID string) (*drive.File, error) {
	d := &drive.File{
		Name:     name,
		MimeType: "application/vnd.google-apps.folder",
		Parents:  []string{parentID},
	}
	file, err := service.Files.Create(d).Do()
	if err != nil {
		log.Fatalf("Could not create dir: " + err.Error())
		return nil, err
	}
	return file, nil
}

func createFile(service *drive.Service, name string, mimeType string, content io.Reader, parentID string) (*drive.File, error) {
	f := &drive.File{
		MimeType: mimeType,
		Name:     name,
		Parents:  []string{parentID},
	}
	file, err := service.Files.Create(f).Media(content).Do()
	if err != nil {
		log.Fatalf("Could not create file: " + err.Error())
		return nil, err
	}
	return file, nil
}

// https://developers.google.com/drive/api/v3/quickstart/go
func getService() (*drive.Service, error) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file. Err: %v\n", err)
		return nil, err
	}
	// If modifying these scopes, delete your previously saved token.json.
	// config, err := google.ConfigFromJSON(b, drive.DriveMetadataReadonlyScope)
	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v\n", err)
		return nil, err
	}
	client := getClient(config)
	// Retrieve Drive client
	service, err := drive.New(client)
	if err != nil {
		log.Fatalf("Cannot create the Google Drive service: %v\n", err)
		return nil, err
	}
	return service, err
}

// Upload single file to Drive
func uploadFile(service *drive.Service) {
	// Open the file
	f, err := os.Open("image.png")
	if err != nil {
		log.Fatalf("Cannot open file: %v", err)
	}
	defer f.Close()
	// Create the directory
	dir, err := createDir(service, driveFolder, "root")
	if err != nil {
		log.Fatalf("Could not create dir: %v\n", err)
	}
	// Create the file and upload its content
	file, err := createFile(service, "uploaded-image.png", "image/png", f, dir.Id)
	if err != nil {
		log.Fatalf("Could not create file: %v\n", err)
	}
	fmt.Printf("File '%s' successfully uploaded in '%s' directory\n", file.Name, dir.Name)
}

// Create folder/directory in Drive for multiple files
func multiFilesDirectory(service *drive.Service) {
	dir, err := createDir(service, driveFolder, "root")
	if err != nil {
		log.Fatalf("Could not create dir: %v\n", err)
	}
	dirID = dir.Id
	dirName = dir.Name
	fmt.Printf("ID: %s  Name: %s\n", dirID, dirName)
}

// Upload multiple files to Drive
func uploadMultiFiles(service *drive.Service, fileName string) {
	fmt.Printf("Uploading file: %s\n", fileName)
	// Open the file
	// input: '/tmp/workspace/file.mp3'
	f, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Cannot open file: %v", err)
	}

	defer f.Close()

	// Create the file and upload its content
	// input: 'file.mp3'
	file, err := createFile(service, strings.Split(fileName, "/")[3], "audio/mpeg", f, dirID)
	if err != nil {
		log.Fatalf("Could not create file: %v\n", err)
	}

	fmt.Printf("File '%s' successfully uploaded in '%s' directory\n", file.Name, dirName)
}

// Get all files from Drive (get folderID to remove)
// https://developers.google.com/drive/api/v2/reference/files/list
func getFiles(service *drive.Service) {
	r, err := service.Files.List().PageSize(10).Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v\n", err)
	}
	if len(r.Files) == 0 {
		fmt.Println("No files in Drive found.")
	} else {
		for _, value := range r.Files {
			fmt.Printf("%s (%s)\n", value.Name, value.Id)
		}
	}
}

// Get files from local directory
func getLocalFiles() []string {
	var files []string
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		log.Fatalf("Unable to get files: %v\n", err)
	}
	return files
}

// Get existing folder/directory ID from Drive
// https://developers.google.com/drive/api/v2/reference/files/list
func getFolderID(service *drive.Service) string {
	var folderID string
	r, err := service.Files.List().PageSize(10).Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v\n", err)
	}
	if len(r.Files) == 0 {
		log.Fatalf("No files in Drive found.")
	} else {
		for _, value := range r.Files {
			if value.Name == driveFolder {
				// fmt.Printf("%s (%s)\n", value.Name, value.Id)
				folderID = value.Id
				break
			}
		}
	}
	return folderID
}

// Remove Drive folder/directory
// https://developers.google.com/drive/api/v2/reference/files/delete
// func removeDir(service *drive.Service, folderID string) (*drive.File, error) {
func removeDir(service *drive.Service, folderID string) {
	err := service.Files.Delete(folderID).Do()
	if err != nil {
		log.Fatalf("An error occurred: %v\n", err)
	}
	fmt.Printf("Folder with ID: %s deleted forever.\n", folderID)
}

func doList() {
	fmt.Printf("Printing all files.\n")
	service, err := getService()
	if err != nil {
		log.Fatalf("Could not list files: %v\n", err)
	}
	getFiles(service)
}

func doRemoval() {
	fmt.Printf("Removing folder: '%s'.\n", driveFolder)
	service, err := getService()
	if err != nil {
		log.Fatalf("Could not list files: %v\n", err)
	}
	removeDir(service, getFolderID(service))
}

func doJob() {
	service, err := getService()
	if err != nil {
		log.Fatalf("Could not list files: %v\n", err)
	}

	// SINGLE
	// Upload single file (image.png) to Drive
	// uploadFile(service)
	// Get files from Drive
	// getFiles(service)

	// MULTI
	files := getLocalFiles()
	if len(files) > 1 {
		multiFilesDirectory(service)
		for _, file := range files[1:] {
			uploadMultiFiles(service, file)
		}
	} else {
		fmt.Printf("No files to be uploaded.\n")
	}
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Missing args. Available upload/delete/list")
	} else {
		switch os.Args[1] {
		case "upload":
			doJob()
		case "delete":
			doRemoval()
		case "list":
			doList()
		}
	}
}
