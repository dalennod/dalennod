package backup

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/term"

	"dalennod/internal/db"
	"dalennod/internal/setup"
)

func JSONOut(database *sql.DB) {
	bookmarkData := db.BackupViewAll(database)
	if len(bookmarkData) == 0 {
		log.Println("WARN: database is empty")
		return
	}

	var jsonData []byte
	var err error
	if setup.FlagVals.Crypt {
		jsonData, err = json.Marshal(bookmarkData)
		if err != nil {
			log.Fatalln("ERROR: error with json marshal:", err)
		}
		key := GetKey()
		jsonData = encryptAES(jsonData, key)
	} else {
		jsonData, err = json.MarshalIndent(bookmarkData, "", "\t")
		if err != nil {
			log.Fatalln("ERROR: error with json marshal & indent:", err)
		}
	}

	jsonBackupFile, err := os.Create("dalennod-backup.json")
	if err != nil {
		log.Fatalln("ERROR: error creating backup file:", err)
	}
	defer jsonBackupFile.Close()

	if _, err = jsonBackupFile.Write(jsonData); err != nil {
		log.Fatalln("ERROR: error writing to backup file:", err)
	}

	fmt.Printf("Backup \"%s\" created at current directory\n", jsonBackupFile.Name())
}

func GetKey() []byte {
	fmt.Print("Enter lock/unlock key: ")
	key, err := term.ReadPassword(0)
	if err != nil {
		log.Fatalln("ERROR: reading key:", err)
	}
	fmt.Println()

	const (
		maxBytes    int  = 32
		zeroDecimal byte = 0x30
	)
	keyLength := len(key)
	if keyLength < maxBytes && keyLength != 0 {
		for i := keyLength; i < maxBytes; i++ {
			key = append(key, zeroDecimal)
		}
	}

	return key
}

func encryptAES(plaintext, key []byte) []byte {
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		log.Fatalln("ERROR: getting cipher block:", err)
	}

	out := make([]byte, aes.BlockSize+len(plaintext))
	iv := out[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Fatalln("ERROR: reading cipher:", err)
	}
	stream := cipher.NewCFBEncrypter(cipherBlock, iv)
	stream.XORKeyStream(out[aes.BlockSize:], plaintext)

	return out
}
