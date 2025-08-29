package user_input

import (
    "bufio"
    "fmt"
    "os"
    "strconv"

    "dalennod/internal/archive"
    "dalennod/internal/db"
    "dalennod/internal/logger"
    "dalennod/internal/setup"
    "dalennod/internal/thumb_url"
)

func addInput(bkmStruct setup.Bookmark, callToUpdate bool) {
    var archiveUrl string
    bkmStruct, archiveUrl = getBKMInfo(bkmStruct)

    if !callToUpdate {
        addBKM(bkmStruct, archiveUrl)
    } else {
        updateBKM(bkmStruct, archiveUrl)
    }
}

func getBKMInfo(bkmStruct setup.Bookmark) (setup.Bookmark, string) {
    var scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
    var err error

    fmt.Print("URL to save: ")
    scanner.Scan()
    bkmStruct.URL = scanner.Text()

    bkmStruct.ThumbURL, err = thumb_url.GetPageThumb(bkmStruct.URL)
    if err != nil {
        bkmStruct.ThumbURL = bkmStruct.URL
    }

    fmt.Print("Title for the bookmark: ")
    scanner.Scan()
    bkmStruct.Title = scanner.Text()

    fmt.Print("Notes/log reason for bookmark: ")
    scanner.Scan()
    bkmStruct.Note = scanner.Text()

    fmt.Print("Keywords for searching later: ")
    scanner.Scan()
    bkmStruct.Keywords = scanner.Text()

    fmt.Print("Category to store the bookmark into: ")
    scanner.Scan()
    bkmStruct.Category = scanner.Text()

    fmt.Print("Archive URL? (y/N): ")
    scanner.Scan()
    var archiveUrl string = scanner.Text()

    return bkmStruct, archiveUrl
}

func updateBKM(bkmStruct setup.Bookmark, archiveUrl string) {
    switch archiveUrl {
    case "y", "Y":
        bkmStruct.Archived, bkmStruct.SnapshotURL = archive.SendSnapshot(bkmStruct.URL)
        if bkmStruct.Archived {
            db.Update(database, bkmStruct, false)
        } else {
            logger.Warn.Println("Snapshot failed")
            db.Update(database, bkmStruct, false)
        }
    case "n", "N", "":
        db.Update(database, bkmStruct, false)
    default:
        db.Update(database, bkmStruct, false)
        logger.Warn.Println("Invalid input for archive request. URL has not been archived")
    }
}

func addBKM(bkmStruct setup.Bookmark, archiveUrl string) {
    switch archiveUrl {
    case "y", "Y":
        bkmStruct.Archived, bkmStruct.SnapshotURL = archive.SendSnapshot(bkmStruct.URL)
        if bkmStruct.Archived {
            db.Add(database, bkmStruct)
        } else {
            logger.Warn.Println("Snapshot failed")
            db.Add(database, bkmStruct)
        }
    case "n", "N", "":
        db.Add(database, bkmStruct)
    default:
        db.Add(database, bkmStruct)
        logger.Warn.Println("Invalid input for archive request. URL has not been archived")
    }
}

func updateInput(updateID string) {
    var (
        confirm string
        scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
    )

    idToINT, err := strconv.ParseInt(updateID, 10, 64)
    if err != nil {
        fmt.Println("Invalid input")
        logger.Error.Println("error: invalid input on update")
        return
    }

    bkmAtID, err := db.ViewSingleRow(database, idToINT)
    if err != nil {
        fmt.Println(err)
        logger.Error.Println(err)
        return
    }
    db.PrintRow(bkmAtID)

    fmt.Print("Update this entry? (Y/n): ")
    scanner.Scan()
    confirm = scanner.Text()

    switch confirm {
    case "y", "Y", "":
        fmt.Println("Leave empty to retain old information")
        addInput(bkmAtID, true)
    case "n", "N":
        return
    default:
        logger.Info.Println("Invalid input. Received:", confirm)
        fmt.Println("Invalid input. Exiting")
        return
    }
}

func removeInput(removeID string) {
    var (
        confirm string
        scanner *bufio.Scanner = bufio.NewScanner(os.Stdin)
    )

    idToINT, err := strconv.ParseInt(removeID, 10, 64)
    if err != nil {
        logger.Error.Println("Invalid input")
        return
    }

    bkmAtID, err := db.ViewSingleRow(database, idToINT)
    if err != nil {
        fmt.Println(err)
        logger.Error.Println(err)
        return
    }
    db.PrintRow(bkmAtID)

    fmt.Print("Remove this entry? (Y/n): ")
    scanner.Scan()
    confirm = scanner.Text()

    switch confirm {
    case "y", "Y", "":
        db.Remove(database, idToINT)
    case "n", "N":
        return
    default:
        logger.Info.Println("Invalid input received:", confirm)
        fmt.Println("Invalid input. Exiting")
        return
    }
}

func viewInput(viewID string) {
    idToINT, err := strconv.ParseInt(viewID, 10, 64)
    if err != nil {
        logger.Error.Println("Invalid input")
        return
    }
    bkmAtID, err := db.ViewSingleRow(database, idToINT)
    if err != nil {
        fmt.Println(err)
        logger.Error.Println(err)
        return
    }
    db.PrintRow(bkmAtID)
}

func viewAllInput(bookmarks []setup.Bookmark) {
    if len(bookmarks) == 0 {
        fmt.Println("database is empty")
        logger.Warn.Println("database empty when trying to view all")
        return
    }

    for _, bookmark := range bookmarks {
        db.PrintRow(bookmark)
    }
}
