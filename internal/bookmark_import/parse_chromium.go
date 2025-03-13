package bookmark_import

import "dalennod/internal/setup"

// type BookmarkNode struct {
//      Name         string         `json:"name"`
//      URL          string         `json:"url,omitempty"`
//      DateAdded    string         `json:"date_added"`
//      DateModified string         `json:"date_modified"`
//      Children     []BookmarkNode `json:"children,omitempty"`
// }

type BookmarkNode struct {
    Name     string         `json:"name"`
    URL      string         `json:"url,omitempty"`
    Children []BookmarkNode `json:"children,omitempty"`
}

type ChromiumItem struct {
    Roots struct {
        BookmarkBar BookmarkNode `json:"bookmark_bar"`
        Other       BookmarkNode `json:"other"`
        Synced      BookmarkNode `json:"synced"`
    } `json:"roots"`
}

func ParseChromium(workingData ChromiumItem) []setup.Bookmark {
    if len(workingData.Roots.BookmarkBar.Children) > 0 {
        deeper(workingData.Roots.BookmarkBar.Children, 1, 10)
    }

    if len(workingData.Roots.Other.Children) > 0 {
        deeper(workingData.Roots.Other.Children, 1, 10)
    }

    if len(workingData.Roots.Synced.Children) > 0 {
        deeper(workingData.Roots.Synced.Children, 1, 10)
    }

    return parsedBookmarks
}

func deeper(array []BookmarkNode, depth int, maxDepth int) {
    if maxDepth > 0 && depth == maxDepth {
        return
    }

    for _, node := range array {
        if len(node.Children) > 0 {
            currentGroup = node.Name
            deeper(node.Children, depth+1, maxDepth)
        } else if node.URL != "" {
            parsedBookmarks = append(parsedBookmarks, setup.Bookmark{
                URL:     node.URL,
                Title:   node.Name,
                BmGroup: currentGroup,
            })
        }
    }
}
