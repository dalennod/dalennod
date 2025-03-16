package server

import (
    "dalennod/internal/logger"
    "dalennod/internal/constants"
    "database/sql"
    "embed"
    "fmt"
    "html/template"
    "io/fs"
    "net/http"
)

var (
    pageCount    int                    = 0
    tmplFuncMap  template.FuncMap       = make(template.FuncMap)
    allBookmarks map[string]interface{} = make(map[string]interface{})
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

    tmplFuncMap["getHostname"] = getHostname
    tmplFuncMap["keywordSplit"] = keywordSplit
    tmplFuncMap["byteConversion"] = byteConversion
    tmplFuncMap["pageCountUp"] = pageCountUp
    tmplFuncMap["pageCountDown"] = pageCountDown
    tmplFuncMap["pageCountNowUpdate"] = pageCountNowUpdate
    tmplFuncMap["pageCountNowDelete"] = pageCountNowDelete

    mux.HandleFunc("/{$}", rootHandler)
    mux.HandleFunc("/import/", importHandler)
    mux.HandleFunc("/api/import-bookmark/", importBookmarkHandler)
    mux.HandleFunc("/api/add/", addHandler)
    mux.HandleFunc("/api/row/{id}", rowHandler)
    mux.HandleFunc("/api/update/{id}", updateHandler)
    mux.HandleFunc("/api/delete/{id}", deleteHandler)
    mux.HandleFunc("/api/groups/", groupsHandler)
    mux.HandleFunc("/api/search/", searchHandler)
    mux.HandleFunc("/api/search-keyword/", searchKeywordHandler)
    mux.HandleFunc("/api/search-group/", searchGroupHandler)
    mux.HandleFunc("/api/search-hostname/", searchHostnameHandler)
    mux.HandleFunc("/api/check-url/", checkUrlHandler)
    mux.HandleFunc("/api/refetch-thumbnail/{id}", refetchThumbnailHandler)
    mux.HandleFunc("/api/pages/", pagesHandler)

    logger.Info.Printf("Web-server starting on http://localhost%s/\n", constants.WEBUI_PORT)
    fmt.Printf("Web-server starting on http://localhost%s/\n", constants.WEBUI_PORT)

    if err := http.ListenAndServe(constants.WEBUI_PORT, mux); err != nil {
        fmt.Println("Stopping. ERROR:", err)
        logger.Error.Fatalln("Stopping. ERROR:", err)
    }
}
