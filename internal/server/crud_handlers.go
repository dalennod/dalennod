package server

import (
    "encoding/json"
    "fmt"
    "html/template"
    "net/http"
    "strconv"
    "strings"

    "dalennod/internal/archive"
    "dalennod/internal/db"
    "dalennod/internal/logger"
    "dalennod/internal/setup"
)

type GotURLParams struct {
    pageNumber string
    searchType string
    searchTerm string
}

func parseURLParams(r *http.Request) GotURLParams {
    var gotURLParams GotURLParams
    err := r.ParseForm()
    if err != nil {
        logger.Error.Println("error parsing URL params. ERROR:", err)
        return gotURLParams
    }
    gotURLParams.pageNumber = r.FormValue("page")
    gotURLParams.searchType = r.FormValue("search-type")
    gotURLParams.searchTerm = r.FormValue("search-term")
    return gotURLParams
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        pageCount := 0
        var bookmarks []setup.Bookmark
        var recentlyInteracted setup.RecentInteractions

        gotURLParams := parseURLParams(r)
        if gotURLParams.pageNumber != "" {
            pageNoValid, err := strconv.Atoi(gotURLParams.pageNumber)
            if err != nil {
                logger.Warn.Println("Invalid page number. Got:", gotURLParams.pageNumber)
            }
            pageCount = pageNoValid
        }
        if gotURLParams.searchType != "" {
            switch gotURLParams.searchType {
            case "general":
                openPrefix := "o "
                searchTermAfter, prefixFound := strings.CutPrefix(gotURLParams.searchTerm, openPrefix)
                if prefixFound && searchTermAfter != "" {
                    gotURLParams.searchTerm = searchTermAfter
                    bookmark := db.OpenSesame(database, gotURLParams.searchTerm)
                    w.Header().Set("Content-Type", "application/json")
                    json.NewEncoder(w).Encode(bookmark)
                    return
                }
                fallthrough
            case "hostname":
                fallthrough
            case "keyword":
                fallthrough
            case "category":
                bookmarks, pageCountForSearch = db.SearchFor(database, gotURLParams.searchType, gotURLParams.searchTerm, pageCount)
            default:
                w.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(w, "ERROR: Unrecognized search type")
                logger.Error.Println("ERROR: Unrecognized search type. Got:", gotURLParams.searchType)
                return
            }
        } else {
            bookmarks = db.ViewAllWebUI(database, pageCount)
        }
        if gotURLParams.pageNumber == "" && gotURLParams.searchType == "" {
            recentlyInteracted = db.RecentInteractions(database)
        }

        tmpl = template.Must(template.New("index").Funcs(tmplFuncMap).ParseFS(Web, "web/index.html"))
        bookmarksMap["AllBookmarks"] = bookmarks
        bookmarksMap["RecentBookmarks"] = recentlyInteracted.Bookmarks
        if err := tmpl.ExecuteTemplate(w, "index", bookmarksMap); err != nil {
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
