package server

import (
	"dalennod/internal/db"
	"dalennod/internal/logger"
	"dalennod/internal/setup"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

func internalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
	logger.Warn.Printf("status 500 at '%s%s'\n", r.Host, r.URL)
}

func rowHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		matchId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || matchId < 1 {
			http.NotFound(w, r)
			return
		}

		oldData, err := db.ViewSingleRow(database, matchId, true)
		if err != nil {
			logger.Error.Println(err)
		}

		writeData, err := json.Marshal(&oldData)
		if err != nil {
			logger.Error.Println(err)
			return
		}

		w.Write([]byte(writeData))
	} else {
		internalServerErrorHandler(w, r)
	}
}

func checkUrlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == http.MethodPost {
		var getData setup.Bookmark
		var err error = json.NewDecoder(r.Body).Decode(&getData)
		if err != nil {
			logger.Error.Println(err)
		}

		var searchUrl string = getData.URL
		getData, err = db.SearchByUrl(database, searchUrl)
		if err != nil {
			logger.Error.Println(err)
		}

		if getData.ID == 0 {
			http.NotFound(w, r)
			return
		}

		writeData, err := json.Marshal(&getData)
		if err != nil {
			logger.Error.Println(err)
		}
		w.Write([]byte(writeData))
	} else {
		internalServerErrorHandler(w, r)
	}
}

func refetchThumbnailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == http.MethodGet {
		matchId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || matchId < 1 {
			http.NotFound(w, r)
			return
		}

		if err := db.RefetchThumbnail(database, matchId, nil); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Thumbnail updated"))
	} else if r.Method == http.MethodPost {
		matchId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil || matchId < 1 {
			http.NotFound(w, r)
			logger.Error.Println("incorrect bookmark id. ERROR:", err)
			return
		}

		if err := r.ParseMultipartForm(10 << 19); err != nil { // ~5,23MB limit for thumbnail file
			http.Error(w, "error parsing data", http.StatusBadRequest)
			logger.Error.Println("error parsing thumbnail data. ERROR:", err)
			return
		}

		thumbnailFile, _, err := r.FormFile("thumbnail")
		if err != nil {
			http.Error(w, "error getting thumbnail field", http.StatusBadRequest)
			logger.Error.Println("error getting thumbnail field. ERROR:", err)
			return
		}
		defer thumbnailFile.Close()

		thumbnailFileBytes, err := io.ReadAll(thumbnailFile)
		if err != nil {
			http.Error(w, "error converting thumbnail file to bytes.", http.StatusInternalServerError)
			logger.Error.Println("error converting thumbnail file to bytes. ERROR:", err)
			return
		}

		if err := db.RefetchThumbnail(database, matchId, thumbnailFileBytes); err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Thumbnail updated"))
	} else {
		internalServerErrorHandler(w, r)
	}
}

func groupsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "*")
	if r.Method == http.MethodGet {
		var currGroups map[string]interface{} = make(map[string]interface{})
		listCurrGroups, err := db.GetAllGroups(database)
		if err != nil {
			logger.Error.Println(err)
		}
		currGroups["AllGroups"] = listCurrGroups

		for _, group := range listCurrGroups {
			fmt.Fprintf(w, "<option value=\"%s\"></option>", group)
		}
	} else {
		internalServerErrorHandler(w, r)
	}
}
