package server

import (
    "database/sql"
    "embed"
    "fmt"
    "html/template"
    "io/fs"
    "net/http"

    "dalennod/internal/logger"
    "dalennod/internal/constants"
    "dalennod/internal/setup"
)

var (
    pageCountForSearch int                   = 0
    tmplFuncMap  template.FuncMap            = make(template.FuncMap)
    bookmarksMap map[string][]setup.Bookmark = make(map[string][]setup.Bookmark)
    database     *sql.DB
    tmpl         *template.Template
    Web          embed.FS
)

func Start(data *sql.DB) {
    database = data

    var mux *http.ServeMux = http.NewServeMux()

    var fsopen fs.FS = fs.FS(Web)
    webStatic, err := fs.Sub(fsopen, "web/static")
    if err != nil {
        fmt.Println("error when opening embedded 'web' directory. ERROR:", err)
        logger.Error.Fatalln("error when opening embedded 'web' directory. ERROR:", err)
    }
    mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(webStatic))))

    tmplFuncMap["getHostname"]    = getHostname
    tmplFuncMap["keywordSplit"]   = keywordSplit
    tmplFuncMap["byteConversion"] = byteConversion
    tmplFuncMap["webUIAddress"]   = webUIAddress

    mux.HandleFunc("/{$}", rootHandler)
    mux.HandleFunc("/import/", importHandler)
    mux.HandleFunc("/api/import-bookmark/", importBookmarkHandler)
    mux.HandleFunc("/api/add/", addHandler)
    mux.HandleFunc("/api/row/{id}", rowHandler)
    mux.HandleFunc("/api/update/{id}", updateHandler)
    mux.HandleFunc("/api/delete/{id}", deleteHandler)
    mux.HandleFunc("/api/categories/", categoriesHandler)
    mux.HandleFunc("/api/check-url/", checkUrlHandler)
    mux.HandleFunc("/api/refetch-thumbnail/{id}", refetchThumbnailHandler)
    mux.HandleFunc("/api/pages/", pagesHandler)

    if constants.WEBUI_ADDR[0] == 58 { // ':'
        logger.Info.Printf("Web-server starting at: http://localhost%s\n", constants.WEBUI_ADDR)
        fmt.Printf("Web-server starting at: http://localhost%s\n", constants.WEBUI_ADDR)
    } else {
        logger.Info.Printf("Web-server starting at: http://%s\n", constants.WEBUI_ADDR)
        fmt.Printf("Web-server starting at: http://%s\n", constants.WEBUI_ADDR)
    }

    if err := http.ListenAndServe(constants.WEBUI_ADDR, mux); err != nil {
        fmt.Println("Stopping. ERROR:", err)
        logger.Error.Fatalln("Stopping. ERROR:", err)
    }
}
