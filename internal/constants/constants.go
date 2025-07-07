package constants

const (
    NAME                string = "dalennod"
    VERSION             string = "0.6.0"

    LOGS_DIRNAME        string = "logs"
    LOGS_FILENAME       string = NAME + ".log"
    DB_DIRNAME          string = "db"
    DB_FILENAME         string = NAME + "." + DB_DIRNAME
    CONFIG_FILENAME     string = "config.json"

    WEBUI_PORT          string = ":41415"
    TIME_FORMAT         string = "2006-01-02 15:04:05"
    PAGE_UPDATE_LIMIT   int    = 60
    RECENT_ENGAGE_LIMIT int    = PAGE_UPDATE_LIMIT>>2

    THUMBNAIL_FILE_SIZE int64  = 10<<19 // ~5.24MB limit on thumbnail file size
    LOG_FILE_SIZE       int64  = 10<<20 // ~10.48MB limit on log file size
    IMPORT_FILE_SIZE    int64  = 10<<21 // 10<<21 = 10*(2^21) = 20,971,520 = ~20.9MB limit on import file size

    COMMON_USERAGENT    string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36"
)

var (
    LOGS_PATH           string
    DB_PATH             string
    CONFIG_PATH         string

    WEBUI_ADDR          string
)
