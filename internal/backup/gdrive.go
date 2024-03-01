package backup

import (
	"context"
	"dalennod/internal/logger"
	"dalennod/internal/setup"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func GDrive() {
	setupFile()
}

func setupFile() {
	upF, err := os.Open(setup.GetOS() + "default_user.db")
	if err != nil {
		logger.Error.Fatalln(err)
	}
	defer upF.Close()

	creds, err := os.ReadFile("gdrivecreds/creds.json")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	config, err := google.ConfigFromJSON(creds, drive.DriveScope)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	service, err := getService(config)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	folderID, err := os.ReadFile("gdrivecreds/folderid.txt")
	if err != nil {
		logger.Error.Fatalln(err)
	}

	uploadFile(filepath.Base(upF.Name()), upF, folderID, service)
}

func uploadFile(upFName string, upF io.Reader, folderID []byte, service *drive.Service) {
	var driveF *drive.File = &drive.File{
		MimeType: "application/octet-stream",
		Name:     upFName,
		Parents:  []string{string(folderID)},
	}

	serviceF, err := service.Files.Create(driveF).Media(upF).Do()
	if err != nil {
		logger.Error.Fatalln(err)
	}

	fmt.Printf("File '%s' uploaded. File id: %s\n", driveF.Name, serviceF.Id)
}

func getService(cfg *oauth2.Config) (*drive.Service, error) {
	var tokenF string = "gdrivecreds/token.json"
	openTokenF, err := getTokenFile(tokenF)
	if err != nil {
		logger.Error.Println(openTokenF, err)
		openTokenF = getToken(cfg, tokenF)
	}

	var cfgClient *http.Client = cfg.Client(context.Background(), openTokenF)
	dService, err := drive.NewService(context.Background(), option.WithHTTPClient(cfgClient))
	if err != nil {
		return nil, err
	}
	return dService, nil
}

func getTokenFile(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var token *oauth2.Token = &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token, err
}

func getToken(cfg *oauth2.Config, path string) *oauth2.Token {
	var authURL string = cfg.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the URL below, allow access then paste the authorization code:\n%s\n", authURL)
	fmt.Println("Authorization code: ")
	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		logger.Error.Fatalln("error reading auth code", err)
	}

	token, err := cfg.Exchange(context.TODO(), authCode)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	tokenF, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		logger.Error.Fatalln(err)
	}
	defer tokenF.Close()

	err = json.NewEncoder(tokenF).Encode(token)
	if err != nil {
		logger.Error.Fatalln(err)
	}

	return token
}
