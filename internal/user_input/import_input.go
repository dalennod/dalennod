package user_input

import (
	"crypto/aes"
	"crypto/cipher"
	"dalennod/internal/backup"
	"dalennod/internal/bookmark_import"
	"dalennod/internal/db"
	"dalennod/internal/logger"
	"dalennod/internal/setup"
	"encoding/json"
	"fmt"
	"os"
)

func importDalennodInput(file string) {
	var fContent []byte
	var err error
	if setup.FlagVals.Crypt {
		fContent, err = os.ReadFile(file)
		if err != nil {
			fmt.Println("error reading file. ERROR:", err)
			logger.Error.Fatalln("error reading file. ERROR:", err)
		}
		key := backup.GetKey()
		fContent = decryptAES(fContent, key)
	} else {
		fContent, err = os.ReadFile(file)
		if err != nil {
			fmt.Println("error reading file. ERROR:", err)
			logger.Error.Fatalln("error reading file. ERROR:", err)
		}
	}

	var importedBookmarks []setup.Bookmark
	err = json.Unmarshal(fContent, &importedBookmarks)
	if err != nil {
		fmt.Println("error with json unmarshal. ERROR:", err)
		logger.Error.Fatalln("error with json unmarshal. ERROR:", err)
	}
	importedBookmarksCount := len(importedBookmarks)
	for i, importedData := range importedBookmarks {
		db.Add(database, importedData)
		fmt.Printf("\rAdded %d / %d", i+1, importedBookmarksCount)
	}
	fmt.Println()
}

func decryptAES(data, key []byte) []byte {
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("error getting cipher block. ERROR:", err)
		logger.Error.Fatalln("error getting cipher block. ERROR:", err)
	}

	if len(data) < aes.BlockSize {
		fmt.Println("ERROR: cipher text too short")
		logger.Error.Fatalln("cipher text too short to import")
	}

	iv := data[:aes.BlockSize]
	stream := cipher.NewCFBDecrypter(cipherBlock, iv)
	plaintext := make([]byte, len(data)-aes.BlockSize)
	stream.XORKeyStream(plaintext, data[aes.BlockSize:])

	return plaintext
}

func importFirefoxInput(file string) {
	readFile, err := os.Open(file)
	if err != nil {
		logger.Error.Printf("couldn't open file. ERROR: %v", err)
	}

	firefoxBookmarks := &bookmark_import.Item{}
	if err := json.NewDecoder(readFile).Decode(firefoxBookmarks); err != nil {
		logger.Error.Println(err)
		fmt.Println(err)
		return
	}

	parsedBookmarks := bookmark_import.ParseFirefox(firefoxBookmarks, "")
	parsedBookmarksLength := len(parsedBookmarks)

	for i, parsedBookmark := range parsedBookmarks {
		db.Add(database, parsedBookmark)
		fmt.Printf("\rAdded %d / %d", i+1, parsedBookmarksLength)
	}
	fmt.Println()
}

func importChromiumInput(file string) {
	readFile, err := os.Open(file)
	if err != nil {
		logger.Error.Printf("couldn't open file. ERROR: %v", err)
	}

	var chromiumBookmarks bookmark_import.ChromiumItem
	if err := json.NewDecoder(readFile).Decode(&chromiumBookmarks); err != nil {
		logger.Error.Println(err)
		fmt.Println(err)
		return
	}

	parsedBookmarks := bookmark_import.ParseChromium(chromiumBookmarks)
	parsedBookmarksLength := len(parsedBookmarks)

	for i, parsedBookmark := range parsedBookmarks {
		db.Add(database, parsedBookmark)
		fmt.Printf("\rAdded %d / %d", i+1, parsedBookmarksLength)
	}
	fmt.Println()
}
