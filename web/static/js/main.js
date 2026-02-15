"use strict";

let API = "";
let root = "";

const showImportPage = () => {
    window.location.href = root + "/import/";
};

const homeButton = () => {
    window.location.href = root;
};

const categoriesOptions = async () => {
    const fetchURL = API + "categories/";
    const res = await fetch(fetchURL);
    return await res.text();
};

const adjustTextarea = (tar) => {
    tar.style.height = (tar.scrollHeight + 4) + "px";
    tar.addEventListener("input", (v) => {
        if (tar.classList.contains("textarea-no-resize") && v.inputType === "insertLineBreak") {
            const pos = tar.selectionStart;
            const value = tar.value;
            const newValue = value.replace(/[\r\n]+/g, '');
            tar.value = newValue;
            const newPos = Math.max(0, pos - 1);
            tar.setSelectionRange(newPos, newPos);
        }
        tar.style.height = "auto";
        tar.style.height = (tar.scrollHeight + 4) + "px";
    });
};

const showCreateDialog = async () => {
    document.getElementById("dialog-create").showModal();

    const allTextarea = document.querySelectorAll(".dialog-create textarea");
    allTextarea.forEach((ta) => adjustTextarea(ta));

    const bookmarkCreateDatalist = document.getElementById("categories-list-create");
    bookmarkCreateDatalist.innerHTML = await categoriesOptions();
};
const closeCreateDialog = () => document.getElementById("dialog-create").close();
document.getElementById("dialog-create").addEventListener("close", () => {
    location.href = root;
});

const showUpdateDialog = async () => {
    document.getElementById("dialog-update").showModal();

    const allTextarea = document.querySelectorAll(".dialog-update textarea");
    allTextarea.forEach((ta) => adjustTextarea(ta));
    toggleMoreOptions();

    const bookmarkUpdateDatalist = document.getElementById("categories-list-update");
    bookmarkUpdateDatalist.innerHTML = await categoriesOptions();
};
const closeUpdateDialog = () => document.getElementById("dialog-update").close();
document.getElementById("dialog-update").addEventListener("close", () => {
    location.reload();
});

const getOldData = async (e) => {
    const entryID = e.parentNode.parentNode.parentNode.id;
    const fetchURL = API + "row/" + entryID;
    const res = await fetch(fetchURL);
    const oldData = await res.json();
    if (typeof Storage !== "undefined") {
        localStorage.setItem("oldData", JSON.stringify(oldData));
    }
    setOldData();
};

let oldData = "";
const setOldData = () => {
    if (typeof Storage !== "undefined") {
        oldData = JSON.parse(localStorage.getItem("oldData"));
        document.getElementById("bkm-id").innerText = oldData.id;
        document.getElementById("update-url").value = oldData.url;
        document.getElementById("update-title").value = oldData.title;
        document.getElementById("update-note").value = oldData.note;
        document.getElementById("update-keywords").value = oldData.keywords;
        document.getElementById("update-category").value = oldData.category;
        if (oldData.archive) {
            document.getElementById("update-archive-div").setAttribute("hidden", "")
            document.getElementById("update-snapshotURL").value = oldData.snapshotURL;
        } else {
            document.getElementById("update-archive-div").removeAttribute("hidden");
            document.getElementById("update-snapshotURL").value = "";
        }
        oldData.thumbURL
            ? document.getElementById("update-thumbURL").value = oldData.thumbURL
            : document.getElementById("update-thumbURL").value = "";
    }
    showUpdateDialog();
};

const updateEntry = async () => {
    const newDataJSON = {
        url: document.getElementById("update-url").value,
        title: document.getElementById("update-title").value,
        note: document.getElementById("update-note").value,
        keywords: document.getElementById("update-keywords").value,
        category: document.getElementById("update-category").value,
        archive: document.getElementById("update-archive").checked ? true : false,
        thumbURL: document.getElementById("update-thumbURL").value
    };

    const updateSnapshotURL = document.getElementById("update-snapshotURL");

    if (updateSnapshotURL.value != "") {
        newDataJSON.archive = true;
        newDataJSON.snapshotURL = updateSnapshotURL.value;
        document.getElementById("update-archive-warn").removeAttribute("hidden");
    } else if (updateSnapshotURL.value === "") {
        newDataJSON.archive = false;
    } else if (!newDataJSON.archive && oldData.archive) {
        newDataJSON.archive = true;
        newDataJSON.snapshotURL = oldData.snapshotURL;
        document.getElementById("update-archive-warn").removeAttribute("hidden");
    } else {
        document.getElementById("update-archive-warn").removeAttribute("hidden");
    }

    const updateButton = document.getElementById("button-update-req");
    updateButton.disabled = true;

    const dataID = document.getElementById("bkm-id").innerText;
    const fetchURL = API + "update/" + dataID;
    const res = await fetch(fetchURL, {
        method: "POST",
        headers: {
            "Content-Type": "application/json; charset=utf-8",
        },
        body: JSON.stringify(newDataJSON),
    });

    if (res.ok) {
        updateButton.disabled = false;
        document.getElementById("update-archive-warn").setAttribute("hidden", "");
        document.getElementById("update-checkmark").removeAttribute("hidden");
        setTimeout(() => document.getElementById("update-checkmark").setAttribute("hidden", ""), 1000);
    }
};

const addEntry = async () => {
    if (document.getElementById("create-url").value === "") {
        alert("ERROR: an URL is required");
        return;
    }

    const dataJSON = {
        url: document.getElementById("create-url").value,
        title: document.getElementById("create-title").value,
        note: document.getElementById("create-note").value,
        keywords: document.getElementById("create-keywords").value,
        category: document.getElementById("create-category").value,
        archive: document.getElementById("create-archive").checked ? true : false,
    };

    if (dataJSON.archive) document.getElementById("create-archive-warn").removeAttribute("hidden");

    const addButton = document.getElementById("button-add-req");
    addButton.disabled = true;

    const fetchURL = API + "add/";
    const res = await fetch(fetchURL, {
        method: "POST",
        headers: {
            "Content-Type": "application/json; charset=utf-8",
        },
        body: JSON.stringify(dataJSON),
    });

    if (res.ok) {
        clearInputs();
        addButton.disabled = false;
        document.getElementById("create-archive-warn").setAttribute("hidden", "");
        document.getElementById("create-checkmark").removeAttribute("hidden");
        setTimeout(() => document.getElementById("create-checkmark").setAttribute("hidden", ""), 1000);
    }
};

const deleteEntry = async (element) => {
    const bookmarkID = element.parentNode.parentNode.parentNode.id;
    let fetchURL = API + "delete/" + bookmarkID;
    const res = await fetch(fetchURL);
    if (!res.ok) {
        console.error(res.status);
        return;
    }
    location.reload();
};

const archiveCheckbox = () => {
    const createArchive = document.getElementById("create-archive");
    const createArchiveLabel = document.getElementById("create-archive-label");
    const updateArchive = document.getElementById("update-archive");
    const updateArchiveLabel = document.getElementById("update-archive-label");
    createArchive.checked ? createArchiveLabel.innerText = "Yes" : createArchiveLabel.innerText = "No";
    updateArchive.checked ? updateArchiveLabel.innerText = "Yes" : updateArchiveLabel.innerText = "No";
};

const clearInputs = () => {
    const input = document.querySelectorAll(".uac-input");
    if (input.length === 0) return;
    for (let i = 0; i < input.length; ++i) input[i].value = "";
};

const openSearchDialog = () => {
    document.getElementById("dialog-search").showModal();

    document.getElementById("dialog-search").addEventListener("click", () => {
        document.getElementById("dialog-search").close();
    });

    document.getElementById("dialog-search-div").addEventListener("click", (event) => {
        event.stopPropagation();
    });

    document.getElementById("general-search-term").focus();
    document.getElementById("general-search-term").value = "";

    const openPrefix = "o ";
    document.getElementById("general-search-term").addEventListener("keydown", async (event) => {
        if (event.key === "Enter") {
            document.getElementById("dialog-search").close();
            event.preventDefault();
            const searchContent = document.getElementById("general-search-term").value;
            if (searchContent.startsWith("::import")) {
                document.getElementById("search-button").disabled = true;
                showImportPage();
                return;
            } else if (searchContent.startsWith(openPrefix)) {
                event.preventDefault();
                const searchFormData = new FormData(document.getElementById("search-form"));
                const fetchURL = `${root}/?${formKeySearchType}=${searchFormData.get(formKeySearchType)}&${formKeySearchTerm}=${searchContent}`;
                const res = await fetch(fetchURL);
                if (!res.ok) {
                    console.error(res.status);
                    return;
                }
                const responseData = await res.json();
                if (responseData.id != 0) {
                    window.open(responseData.url, "_blank");
                } else {
                    console.error("WARN: Did not find anything to open with input:", searchContent);
                }
            } else {
                document.getElementById("search-form").submit();
            }
        }
    });
};

document.addEventListener("keydown", (event) => {
    if (document.getElementById("dialog-create").open ||
        document.getElementById("dialog-update").open ||
        document.getElementById("dialog-search").open ||
        event.ctrlKey || event.altKey || event.metaKey) {
        return;
    }

    if (event.Key === "/" || event.code === "Slash") {
        event.preventDefault();
        openSearchDialog();
    }
});

const toggleMoreOptions = () => {
    const updateMoreOptions = document.querySelectorAll(".update-more-options");
    updateMoreOptions.forEach((updateOption) => {
        updateOption.hidden = !updateOption.hidden;
    });
};

const recentlyInteractedHiddenKey = "recentlyInteractedHidden";

const toggleRecentlyInteracted = (e) => {
    const hiddenText = " (Hidden)";
    let gridViewList;
    try {
        gridViewList = e.parentElement.getElementsByClassName("grid-view-list")[0];
    } catch (e) {
        console.log("WARN: caught error:", e.message);
        return;
    }
    if (gridViewList.style.display == "none") {
        gridViewList.style.display = "";
        e.innerHTML = e.innerHTML.slice(0, -(hiddenText.length));
        e.outerHTML = e.outerHTML.replace("h5", "h2");
        localStorage.setItem(recentlyInteractedHiddenKey, false);
    } else {
        gridViewList.style.display = "none";
        e.innerHTML = e.innerHTML + hiddenText;
        e.outerHTML = e.outerHTML.replace("h2", "h5");
        localStorage.setItem(recentlyInteractedHiddenKey, true);
    }
};

const createPaginationLink = (pageNumber, current, paginationNumbers, hrefLocation) => {
    const params = new URLSearchParams(hrefLocation.search);
    params.set("page", pageNumber);

    const pageNumberLink = document.createElement("a");
    pageNumberLink.textContent = pageNumber;
    pageNumberLink.title = `Page ${pageNumber}`;
    pageNumberLink.href = "?"+params.toString();

    if (pageNumber === current) pageNumberLink.classList.add("active");
    paginationNumbers.appendChild(pageNumberLink);
};

const formKeySearchType = "search-type";
const formKeySearchTerm = "search-term";

const updateNavATags = (current, totalPages, hrefParams) => {
    const nextPage = current + 1;
    const prevPage = current - 1;
    if (hrefParams.get(formKeySearchType)) {
        const searchType = hrefParams.get(formKeySearchType);
        const searchTerm = hrefParams.get(formKeySearchTerm);
        document.getElementById("root-first-page").href = `/?${formKeySearchType}=${searchType}&${formKeySearchTerm}=${searchTerm}`;
        document.getElementById("prev-page").href = `/?${formKeySearchType}=${searchType}&${formKeySearchTerm}=${searchTerm}&page=${prevPage}`;
        document.getElementById("next-page").href = `/?${formKeySearchType}=${searchType}&${formKeySearchTerm}=${searchTerm}&page=${nextPage}`;
        document.getElementById("last-page").href = `/?${formKeySearchType}=${searchType}&${formKeySearchTerm}=${searchTerm}&page=${totalPages}`;

        const allBookmarksList = document.getElementById("all-bookmarks-list");
        searchType === "general"
            ? allBookmarksList.innerHTML = `<img class="svg-img" src="static/imgs/search_button.svg" alt="Search all bookmark icon" />All Bookmarks with '${searchTerm}'`
            : allBookmarksList.innerHTML = `<img class="svg-img" src="static/imgs/search_button.svg" alt="Search all bookmark icon" />Bookmarks with '${searchTerm}' in '${searchType}'`;
    } else {
        document.getElementById("prev-page").href = `/?page=${prevPage}`;
        document.getElementById("next-page").href = `/?page=${nextPage}`;
        document.getElementById("last-page").href = `/?page=${totalPages}`;
    }
};

const updatePagination = async () => {
    const hrefLocation = new URL(location.href);
    const hrefParams = hrefLocation.searchParams;
    const fetchURL = API + "pages/";
    let res;
    if (hrefParams.get(formKeySearchType)) {
        res = await fetch(fetchURL, {
            method: "POST",
        });
    } else {
        res = await fetch(fetchURL);
    }
    if (!res.ok) {
        console.error(res.status);
        return;
    }
    const totalPages = Number(await res.text());

    const pagesToShow = 5;
    const currentPage = Number(hrefParams.get("page"));
    updateNavATags(currentPage, totalPages, hrefParams);

    if (totalPages <= 0) {
        const paginationFooter = document.getElementById("pagination");
        paginationFooter.hidden = true;
        return;
    }

    if (currentPage <= 0) document.getElementById("prev-page").style.pointerEvents = "none";
    if (currentPage === totalPages) document.getElementById("next-page").style.pointerEvents = "none";

    const paginationNumbers = document.getElementById("pagination-numbers");
    paginationNumbers.innerHTML = "";

    const start = Math.max(0, currentPage - pagesToShow);
    const end = Math.min(totalPages, currentPage + pagesToShow);

    for (let i = start; i <= end; ++i) {
        createPaginationLink(i, currentPage, paginationNumbers, hrefLocation);
    }
};

const checkRecentlyInteractedVisibility = () => {
    const recentlyInteractedHiddenState = JSON.parse(localStorage.getItem(recentlyInteractedHiddenKey));
    if (recentlyInteractedHiddenState) {
        toggleRecentlyInteracted(document.getElementById("special-bookmarks"));
    }
};

window.onload = () => {
    root = new URL(location.href).origin;
    API = `${root}/api/`;
    checkRecentlyInteractedVisibility();
    updatePagination();
};
