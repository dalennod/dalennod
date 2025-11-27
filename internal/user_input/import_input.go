package user_input

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"dalennod/internal/backup"
	"dalennod/internal/bookmark_import"
	"dalennod/internal/db"
	"dalennod/internal/setup"
)

func importDalennodInput(file string) {
	var fContent []byte
	var err error
	if setup.FlagVals.Crypt {
		fContent, err = os.ReadFile(file)
		if err != nil {
			log.Fatalln("ERROR: reading import file:", err)
		}
		key := backup.GetKey()
		fContent = decryptAES(fContent, key)
	} else {
		fContent, err = os.ReadFile(file)
		if err != nil {
			log.Fatalln("ERROR: reading import file:", err)
		}
	}

	var importedBookmarks []setup.Bookmark
	err = json.Unmarshal(fContent, &importedBookmarks)
	if err != nil {
		log.Fatalln("ERROR: failed to unmarshal:", err)
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
		log.Fatalln("ERROR: could not get cipher block:", err)
	}

	if len(data) < aes.BlockSize {
		log.Fatalln("ERROR: cipher text too short")
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
		log.Printf("couldn't open file. ERROR: %v\n", err)
	}

	firefoxBookmarks := &bookmark_import.Item{}
	if err := json.NewDecoder(readFile).Decode(firefoxBookmarks); err != nil {
		log.Fatalln("ERROR: failed to decode input file:", err)
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
		log.Println("ERROR: could not open file:", err)
	}

	var chromiumBookmarks bookmark_import.ChromiumItem
	if err := json.NewDecoder(readFile).Decode(&chromiumBookmarks); err != nil {
		log.Fatalln("ERROR: failed to decode input file:", err)
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
