package server

import (
    "dalennod/internal/archive"
    "dalennod/internal/db"
    "dalennod/internal/logger"
    "dalennod/internal/setup"
    "encoding/json"
    "html/template"
    "net/http"
    "strconv"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        pageCount = 0
        var bookmarks []setup.Bookmark
        r.ParseForm()

        var pageNo string = r.FormValue("page")
        if pageNo == "" {
            bookmarks = db.ViewAllWebUI(database, 0)
        } else {
            pageNoInt, err := strconv.Atoi(pageNo)
            if err != nil {
                logger.Error.Printf("error: invalid page no. %v", err)
            }
            pageCount = pageNoInt
            bookmarks = db.ViewAllWebUI(database, pageNoInt)
        }

        tmpl = template.Must(template.New("index").Funcs(tmplFuncMap).ParseFS(Web, "web/index.html"))
        allBookmarks["Bookmarks"] = bookmarks
        if err := tmpl.ExecuteTemplate(w, "index", allBookmarks); err != nil {
            logger.Warn.Println("error executing template for root index. ERROR:", err)
        }
    } else {
        internalServerErrorHandler(w, r)
    }
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "*")
    if r.Method == http.MethodGet {
        matchId, err := strconv.Atoi(r.PathValue("id"))
        if err != nil || matchId < 1 {
            http.NotFound(w, r)
            return
        }

        db.Remove(database, matchId)
    } else {
        internalServerErrorHandler(w, r)
    }
}

func addHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "*")
    if r.Method == http.MethodPost {
        var insData setup.Bookmark
        if err := json.NewDecoder(r.Body).Decode(&insData); err != nil {
            logger.Error.Println(err)
        }

        if !insData.Archived {
            db.Add(database, insData)
        } else {
            insData.Archived, insData.SnapshotURL = archive.SendSnapshot(insData.URL)
            if insData.Archived {
                db.Add(database, insData)
            } else {
                logger.Warn.Println("Snapshot failed")
                db.Add(database, insData)
            }
        }

        w.WriteHeader(http.StatusCreated)
    } else if r.Method == http.MethodGet {
        w.Write([]byte("Alive"))
    } else {
        internalServerErrorHandler(w, r)
    }
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "*")
    if r.Method == http.MethodPost {
        matchId, err := strconv.Atoi(r.PathValue("id"))
        if err != nil || matchId < 1 {
            http.NotFound(w, r)
            return
        }

        var newData setup.Bookmark
        err = json.NewDecoder(r.Body).Decode(&newData)
        if err != nil {
            logger.Error.Println(err)
        }
        newData.ID = matchId

        if !newData.Archived {
            db.Update(database, newData, true)
        } else if newData.Archived && newData.SnapshotURL != "" {
            db.Update(database, newData, true)
        } else {
            newData.Archived, newData.SnapshotURL = archive.SendSnapshot(newData.URL)
            if newData.Archived {
                db.Update(database, newData, true)
            } else {
                logger.Warn.Println("Snapshot failed")
                db.Update(database, newData, true)
            }
        }

        w.WriteHeader(http.StatusCreated)
    } else {
        internalServerErrorHandler(w, r)
    }
}
