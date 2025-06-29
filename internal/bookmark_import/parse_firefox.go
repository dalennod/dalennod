// Thank you, Christopher:
// https://github.com/chrsm/

package bookmark_import

import "dalennod/internal/setup"

// Firefox's bookmark constants from:
// https://github.com/mozilla/gecko-dev/blob/599a15d3547d862764048ff62b74252dd41a56d3/toolkit/components/places/Bookmarks.jsm#L95
const (
    TypeBookmark  = 1
    TypeFolder    = 2
    TypeSeparator = 3

    DefaultIndex = -1

    MaxTagLength = 100

    GUIDRoot    = "root________"
    GUIDMenu    = "menu________"
    GUIDToolbar = "toolbar_____"
    GUIDUnfiled = "unfiled_____"
    GUIDMobile  = "mobile______"
    GUIDTag     = "tags________"

    GUIDVirtMenu    = "menu_______v"
    GUIDVirtToolbar = "toolbar____v"
    GUIDVirtUnfiled = "unfiled___v"
    GUIDVirtMobile  = "mobile____v"
)

// Firefox's bookmark properties
// https://github.com/mozilla/gecko-dev/blob/599a15d3547d862764048ff62b74252dd41a56d3/toolkit/components/places/PlacesBackups.jsm#L398-L421
// https://github.com/mozilla/gecko-dev/blob/599a15d3547d862764048ff62b74252dd41a56d3/toolkit/components/places/Bookmarks.jsm#L7
// type Item struct {
//      // The globally unique identifier of the item.
//      GUID string `json:"guid"`

//      // The globally unique identifier of the folder containing the item.
//      // This will be an empty string for the Places root folder.
//      ParentGUID string `json:"parentGuid"`

//      Title string `json:"title"`

//      // The 0-based position of the item in the parent folder.
//      Index int `json:"index"`

//      // The time at which the item was added.
//      DateAdded time.Time `json:"dateAdded"`
//      // The time at which the item was last modified.
//      LastModified time.Time `json:"lastModified"`

//      ID int `json:"id"`

//      // TypeCode designates the type of item i.e. bookmark, folder or separator
//      TypeCode int `json:"typeCode"`

//      Type string `json:"type"`
//      Root string `json:"root"`

//      // Children are the items within a TypeFolder.
//      Children []*Item `json:"children"`

//      // The following fields only apply to a subset of items.
//      Annos   []Anno `json:"annos"`
//      URI     string `json:"uri"`
//      IconURI string `json:"iconuri"`
//      Keyword string `json:"keyword"`
//      Charset string `json:"charset"`
//      Tags    string `json:"tags"`
// }

type Item struct {
    // The globally unique identifier of the item.
    GUID string `json:"guid"`

    // The globally unique identifier of the folder containing the item.
    // This will be an empty string for the Places root folder.
    ParentGUID string `json:"parentGuid"`

    Title string `json:"title"`

    // The 0-based position of the item in the parent folder.
    Index int `json:"index"`

    ID int `json:"id"`

    // TypeCode designates the type of item i.e. bookmark, folder or separator
    TypeCode int `json:"typeCode"`

    Type string `json:"type"`
    Root string `json:"root"`

    // Children are the items within a TypeFolder.
    Children []*Item `json:"children"`

    // The following fields only apply to a subset of items.
    Annos   []Anno `json:"annos"`
    URI     string `json:"uri"`
    IconURI string `json:"iconuri"`
    Keyword string `json:"keyword"`
    Charset string `json:"charset"`
    Tags    string `json:"tags"`
}

type Anno struct {
    Name    string `json:"name"`
    Value   string `json:"value"`
    Expires int    `json:"expires"`
    Flags   int    `json:"flags"`
}

var (
    currentCategory string
    parsedBookmarks []setup.Bookmark
)

func ParseFirefox(bookmarks *Item, prefix string) []setup.Bookmark {
    if bookmarks.TypeCode == TypeFolder && len(bookmarks.Children) > 0 {
        currentCategory = bookmarks.Title
        for i := range bookmarks.Children {
            ParseFirefox(bookmarks.Children[i], prefix+"\t")
        }
    } else if bookmarks.TypeCode == TypeBookmark {
        parsedBookmarks = append(parsedBookmarks, setup.Bookmark{
            URL:      bookmarks.URI,
            Title:    bookmarks.Title,
            Keywords: bookmarks.Keyword,
            Category:  currentCategory,
        })

    }
    return parsedBookmarks
}
