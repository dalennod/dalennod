"use strict";

let API = "";
let root = "";

let changeToImportTimeout = "";
const showImportPageTimeout = 5 * 1000; // 5 seconds
const changeToImport = () => changeToImportTimeout = setTimeout(() => showImportPage(), showImportPageTimeout);
const showImportPage = () => window.location.href = root + "/import/";
const clearImportTimeout = () => clearTimeout(changeToImportTimeout);

const groupsDatalistOptions = async () => {
    const fetchURL = API + "groups/";
    const res = await fetch(fetchURL);
    return await res.text();
};

const showCreateDialog = async () => {
    document.querySelector(".dialog-create").showModal();
    clearImportTimeout();
    const bookmarkCreateDatalist = document.getElementById("bmGroups-list-create");
    bookmarkCreateDatalist.innerHTML = await groupsDatalistOptions();
};
const closeCreateDialog = () => document.querySelector(".dialog-create").close();
document.querySelector(".dialog-create").addEventListener("close", () => {
    location.href = root;
});

const showUpdateDialog = async () => {
    document.querySelector(".dialog-update").showModal();
    const noteTextArea = document.getElementById("update-note");
    noteTextArea.style.height = "auto";
    noteTextArea.style.height = noteTextArea.scrollHeight + "px";

    const bookmarkUpdateDatalist = document.getElementById("bmGroups-list-update");
    bookmarkUpdateDatalist.innerHTML = await groupsDatalistOptions();
};
const closeUpdateDialog = () => document.querySelector(".dialog-update").close();
document.querySelector(".dialog-update").addEventListener("close", () => {
    location.reload();
});

const getOldData = async (ele) => {
    const entryID = ele.parentNode.parentNode.parentNode.id;
    const fetchURL = API + "row/" + entryID;
    const res = await fetch(fetchURL);
    const oldData = await res.json();
    if (typeof Storage !== "undefined") localStorage.setItem("oldData", JSON.stringify(oldData));
    setOldData();
};

let oldData = "";
const setOldData = () => {
    if (typeof Storage !== "undefined") {
        oldData = JSON.parse(localStorage.getItem("oldData"));
        document.querySelector("#bm-id").innerText = oldData.id;
        document.querySelector("#update-url").value = oldData.url;
        document.querySelector("#update-title").value = oldData.title;
        document.querySelector("#update-note").value = oldData.note;
        document.querySelector("#update-keywords").value = oldData.keywords;
        document.querySelector("#update-bmGroup").value = oldData.bmGroup;
        oldData.archive
            ? document.querySelector("#update-archive-div").setAttribute("hidden", "")
            : document.querySelector("#update-archive-div").removeAttribute("hidden");
    }
    showUpdateDialog();
};

const updateEntry = async () => {
    const newDataJSON = {
        url: document.querySelector("#update-url").value,
        title: document.querySelector("#update-title").value,
        note: document.querySelector("#update-note").value,
        keywords: document.querySelector("#update-keywords").value,
        bmGroup: document.querySelector("#update-bmGroup").value,
        archive: document.getElementById("update-archive").checked ? true : false
    }

    if (newDataJSON.archive) {
        document.querySelector("#update-archive-warn").removeAttribute("hidden");
    } else if (!newDataJSON.archive && oldData.archive) {
        newDataJSON.archive = true;
        newDataJSON.snapshotURL = oldData.snapshotURL;
    }

    const updateButton = document.getElementById("button-update-req");
    updateButton.disabled = true;

    const dataID = document.querySelector("#bm-id").innerText;
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
        bmGroup: document.querySelector("#create-bmGroup").value,
        archive: document.getElementById("create-archive").checked ? true : false,
    }

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
    const params = new URLSearchParams(hrefLocation.search); // .search = +rw
    params.set("page", pageNumber);

    const pageNumberLink = document.createElement("a");
    pageNumberLink.textContent = pageNumber;
    pageNumberLink.title = `Page ${pageNumber}`;
    pageNumberLink.href = "?"+params.toString();

    if (pageNumber === current) pageNumberLink.classList.add("active");
    paginationNumbers.appendChild(pageNumberLink);
};

const updateNavATags = (current, totalPages, hrefParams) => {
    const nextPage = current + 1;
    const prevPage = current - 1;
    if (hrefParams.get("search-type")) {
        document.getElementById("prev-page").href = `/?search-type=${hrefParams.get("search-type")}&search-term=${hrefParams.get("search-term")}&page=${prevPage}`;
        document.getElementById("next-page").href = `/?search-type=${hrefParams.get("search-type")}&search-term=${hrefParams.get("search-term")}&page=${nextPage}`;
        document.getElementById("last-page").href = `/?search-type=${hrefParams.get("search-type")}&search-term=${hrefParams.get("search-term")}&page=${totalPages}`;
    } else {
        document.getElementById("prev-page").href = `/?page=${prevPage}`;
        document.getElementById("next-page").href = `/?page=${nextPage}`;
        document.getElementById("last-page").href = `/?page=${totalPages}`;
    }
};

const updatePagination = async () => {
    const hrefLocation = new URL(location.href);
    const hrefParams = hrefLocation.searchParams; // .searchParams = +r
    const fetchURL = API + "pages/";
    let res;
    if (hrefParams.get("search-type")) {
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

window.onload = () => {
    const host = new URL(location.href).host;
    root = `http://${host}`;
    API = `${root}/api/`;
    updatePagination();
};
