package backup

import (
	"dalennod/internal/db"
	"dalennod/internal/logger"
	"dalennod/internal/setup"
	"database/sql"
	"encoding/json"
	"fmt"
)

func JSONOut(database *sql.DB) {
	var data []setup.Bookmark = db.ViewAll(database, true)
	jsonIndent, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		logger.Error.Println(err)
	}
	fmt.Println(string(jsonIndent))
}
