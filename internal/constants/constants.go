package constants

const (
	NAME    string = "dalennod"
	VERSION string = "0.17.3"

	DB_DIRNAME         string = "db"
	DB_FILENAME        string = NAME + "." + DB_DIRNAME
	THUMBNAILS_DIRNAME string = "thumbnails"

	WEBUI_PORT          string = ":41415"
	SECONDARY_PORT      string = ":41417"
	TIME_FORMAT         string = "2006-01-02 15:04"
	PAGE_UPDATE_LIMIT   int    = 60
	RECENT_ENGAGE_LIMIT int    = PAGE_UPDATE_LIMIT >> 2

	THUMBNAIL_FILE_SIZE int64 = 10 << 19 // ~5.24MB limit on thumbnail file size
	IMPORT_FILE_SIZE    int64 = 10 << 21 // 10<<21 = 10*(2^21) = 20,971,520 = ~20.9MB limit on import file size

	COMMON_USERAGENT string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/143.0.0.0 Safari/537.36"
)

var (
	DATA_PATH       string
	DB_PATH         string
	THUMBNAILS_PATH string
	WEBUI_ADDR      string
)
