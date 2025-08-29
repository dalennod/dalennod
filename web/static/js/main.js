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

const showCreateDialog = async () => {
    document.querySelector(".dialog-create").showModal();
    const noteTextArea = document.getElementById("create-note");
    adjustTextarea(noteTextArea);

    const bookmarkCreateDatalist = document.getElementById("categories-list-create");
    bookmarkCreateDatalist.innerHTML = await categoriesOptions();
};
const closeCreateDialog = () => document.querySelector(".dialog-create").close();
document.querySelector(".dialog-create").addEventListener("close", () => {
    location.href = root;
});

const showUpdateDialog = async () => {
    document.querySelector(".dialog-update").showModal();
    const noteTextArea = document.getElementById("update-note");
    adjustTextarea(noteTextArea);

    const bookmarkUpdateDatalist = document.getElementById("categories-list-update");
    bookmarkUpdateDatalist.innerHTML = await categoriesOptions();
};
const closeUpdateDialog = () => document.querySelector(".dialog-update").close();
document.querySelector(".dialog-update").addEventListener("close", () => {
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
        document.querySelector("#bkm-id").innerText = oldData.id;
        document.querySelector("#update-url").value = oldData.url;
        document.querySelector("#update-title").value = oldData.title;
        document.querySelector("#update-note").value = oldData.note;
        document.querySelector("#update-keywords").value = oldData.keywords;
        document.querySelector("#update-category").value = oldData.category;
        if (oldData.archive) {
            document.querySelector("#update-archive-div").setAttribute("hidden", "")
            document.getElementById("update-snapshotURL").value = oldData.snapshotURL;
        } else {
            document.querySelector("#update-archive-div").removeAttribute("hidden");
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
        url: document.querySelector("#update-url").value,
        title: document.querySelector("#update-title").value,
        note: document.querySelector("#update-note").value,
        keywords: document.querySelector("#update-keywords").value,
        category: document.querySelector("#update-category").value,
        archive: document.getElementById("update-archive").checked ? true : false,
        thumbURL: document.getElementById("update-thumbURL").value
    };

    const updateSnapshotURL = document.getElementById("update-snapshotURL");

    if (updateSnapshotURL.value != "") {
        newDataJSON.archive = true;
        newDataJSON.snapshotURL = updateSnapshotURL.value;
        document.querySelector("#update-archive-warn").removeAttribute("hidden");
    } else if (updateSnapshotURL.value === "") {
        newDataJSON.archive = false;
    } else if (!newDataJSON.archive && oldData.archive) {
        newDataJSON.archive = true;
        newDataJSON.snapshotURL = oldData.snapshotURL;
        document.querySelector("#update-archive-warn").removeAttribute("hidden");
    } else {
        document.querySelector("#update-archive-warn").removeAttribute("hidden");
    }

    const updateButton = document.getElementById("button-update-req");
    updateButton.disabled = true;

    const dataID = document.querySelector("#bkm-id").innerText;
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
        document.querySelector("#update-archive-warn").setAttribute("hidden", "");
        document.querySelector("#update-checkmark").removeAttribute("hidden");
        setTimeout(() => document.querySelector("#update-checkmark").setAttribute("hidden", ""), 1000);
    }
};

const addEntry = async () => {
    if (document.querySelector("#create-url").value === "") {
        alert("ERROR: an URL is required");
        return;
    }

    const dataJSON = {
        url: document.querySelector("#create-url").value,
        title: document.querySelector("#create-title").value,
        note: document.querySelector("#create-note").value,
        keywords: document.querySelector("#create-keywords").value,
        category: document.querySelector("#create-category").value,
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
        document.querySelector("#create-archive-warn").setAttribute("hidden", "");
        document.querySelector("#create-checkmark").removeAttribute("hidden");
        setTimeout(() => document.querySelector("#create-checkmark").setAttribute("hidden", ""), 1000);
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
    const createArchiveLabel = document.querySelector(".create-archive-label");
    const updateArchive = document.getElementById("update-archive");
    const updateArchiveLabel = document.querySelector(".update-archive-label");
    createArchive.checked ? createArchiveLabel.innerText = "Yes" : createArchiveLabel.innerText = "No";
    updateArchive.checked ? updateArchiveLabel.innerText = "Yes" : updateArchiveLabel.innerText = "No";
};

const clearInputs = () => {
    const input = document.querySelectorAll(".uac-input");
    if (input.length === 0) return;
    for (let i = 0; i < input.length; ++i) input[i].value = "";
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
        const paginationFooter = document.querySelector(".pagination");
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

const openSearchDialog = () => {
    document.querySelector(".dialog-search").showModal();

    document.querySelector(".dialog-search").addEventListener("click", () => {
        document.querySelector(".dialog-search").close();
    });

    document.getElementById("dialog-search-div").addEventListener("click", (event) => {
        event.stopPropagation();
    });

    document.getElementById("general-search-term").focus();
    document.getElementById("general-search-term").value = "";

    const openPrefix = "o ";
    document.getElementById("general-search-term").addEventListener("keydown", async (event) => {
        if (event.key === "Enter") {
            document.querySelector(".dialog-search").close();
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
    if (document.querySelector(".dialog-create").open ||
        document.querySelector(".dialog-update").open ||
        document.querySelector(".dialog-search").open ||
        event.ctrlKey || event.altKey || event.metaKey) {
        return;
    }

    if (event.Key === "/" || event.code === "Slash") {
        event.preventDefault();
        openSearchDialog();
    }
});

const adjustTextarea = (tar) => {
    tar.value.includes("\n")
        ? tar.style.height = (tar.scrollHeight + 8) + "px"
        : tar.style.height = tar.scrollHeight + "px";
    tar.addEventListener("input", () => {
        tar.style.height = "auto";
        tar.value.includes("\n")
            ? tar.style.height = (tar.scrollHeight + 8) + "px"
            : tar.style.height = tar.scrollHeight + "px";
    });
};

/*
// const observer = new IntersectionObserver((entries) => {
//     entries.forEach((entry) => {
//         if (entry.isIntersecting) {
//             entry.target.classList.add("show");
//         } else {
//             entry.target.classList.remove("show");
//         }
//     })
// }, {});
// const gridChildren = document.querySelectorAll(".grid-child");
// gridChildren.forEach(el => observer.observe(el));
*/

const toggleMoreOptions = () => {
    const updateMoreOptions = document.querySelectorAll(".update-more-options");
    updateMoreOptions.forEach((updateOption) => {
        updateOption.hidden = !updateOption.hidden;
    });
};

window.onload = () => {
    root = new URL(location.href).origin;
    API = `${root}/api/`;
    updatePagination();
};
