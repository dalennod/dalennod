package backup

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"dalennod/internal/db"
	"dalennod/internal/logger"
	"dalennod/internal/setup"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

func JSONOut(database *sql.DB) {
	bookmarkData := db.BackupViewAll(database)
	if len(bookmarkData) == 0 {
		fmt.Println("database is empty")
		logger.Warn.Println("database is empty. exiting to prevent panic")
		return
	}

	var jsonData []byte
	var err error
	if setup.FlagVals.Crypt {
		jsonData, err = json.Marshal(bookmarkData)
		if err != nil {
			logger.Error.Fatalln("error with json marshal. ERROR:", err)
		}
		key := GetKey()
		jsonData = encryptAES(jsonData, key)
	} else {
		jsonData, err = json.MarshalIndent(bookmarkData, "", "\t")
		if err != nil {
			logger.Error.Fatalln("error with json marshal & indent. ERROR:", err)
		}
	}

	jsonBackupFile, err := os.Create("dalennod-backup.json")
	if err != nil {
		logger.Error.Fatalln("error creating backup file. ERROR:", err)
	}
	defer jsonBackupFile.Close()

	if _, err = jsonBackupFile.Write(jsonData); err != nil {
		logger.Error.Fatalln("error writing to backup file. ERROR:", err)
	}

	fmt.Printf("backup \"%s\" created at current directory\n", jsonBackupFile.Name())
}

func GetKey() []byte {
	fmt.Print("Enter lock/unlock key (min 16 chars): ")
	key, err := term.ReadPassword(0)
	if err != nil {
		fmt.Println("error reading key. ERROR:", err)
		logger.Error.Fatalln("error reading key. ERROR:", err)
	}
	fmt.Println()
	return key
}

func encryptAES(plaintext, key []byte) []byte {
	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("error getting cipher block. ERROR:", err)
		logger.Error.Fatalln("error getting cipher block. ERROR:", err)
	}

	out := make([]byte, aes.BlockSize+len(plaintext))
	iv := out[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		fmt.Println("error reading cipher. ERROR:", err)
		logger.Error.Fatalln("error reading cipher. ERROR:", err)
	}
	stream := cipher.NewCFBEncrypter(cipherBlock, iv)
	stream.XORKeyStream(out[aes.BlockSize:], plaintext)

	return out
}
