"use strict";

let API = "";
let root = "";

const showCreateDialog = () => {
    document.querySelector(".dialog-create").showModal();
    clearImportTimeout();
}
const closeCreateDialog = () => document.querySelector(".dialog-create").close();

const showUpdateDialog = () => {
    document.querySelector(".dialog-update").showModal();
    const noteTextArea = document.getElementById("update-note");
    noteTextArea.style.height = "auto";
    noteTextArea.style.height = noteTextArea.scrollHeight + "px";
}
const closeUpdateDialog = () => document.querySelector(".dialog-update").close();

const getOldData = async (ele) => {
    const entryID = ele.parentNode.parentNode.id;
    const fetchURL = API + "row/" + entryID;
    const res = await fetch(fetchURL);
    const oldData = await res.json();
    if (typeof Storage !== "undefined") localStorage.setItem("oldData", JSON.stringify(oldData));
    setOldData();
}

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
        oldData.archive ? document.querySelector("#update-archive-div").setAttribute("hidden", "") : document.querySelector("#update-archive-div").removeAttribute("hidden");
    }
    showUpdateDialog();
}

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
    })

    if (res.ok) {
        updateButton.disabled = false;
        document.querySelector("#update-archive-warn").setAttribute("hidden", "");
        document.querySelector("#update-checkmark").removeAttribute("hidden");
        setTimeout(() => document.querySelector("#update-checkmark").setAttribute("hidden", ""), 1000);
    }
}

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
    })

    if (res.ok) {
        clearInputs();
        addButton.disabled = false;
        document.querySelector("#create-archive-warn").setAttribute("hidden", "");
        document.querySelector("#create-checkmark").removeAttribute("hidden");
        setTimeout(() => document.querySelector("#create-checkmark").setAttribute("hidden", ""), 1000);
    }
}

const deleteEntry = async (element) => {
    const bookmarkID = element.parentNode.parentNode.id;
    let fetchURL = API + "delete/" + bookmarkID;
    const res = await fetch(fetchURL);
    if (res.ok) {
        location.reload();
    }
}

let changeToImportTimeout = "";
const showImportPageTimeout = 5 * 1000; // 5 seconds
const changeToImport = () => changeToImportTimeout = setTimeout(() => showImportPage(), showImportPageTimeout);
const showImportPage = () => document.location.href = root + "/import/";
const clearImportTimeout = () => clearTimeout(changeToImportTimeout);

const archiveCheckbox = () => {
    const createArchive = document.getElementById("create-archive");
    const createArchiveLabel = document.querySelector(".create-archive-label");
    const updateArchive = document.getElementById("update-archive");
    const updateArchiveLabel = document.querySelector(".update-archive-label");
    createArchive.checked ? createArchiveLabel.innerText = "Yes" : createArchiveLabel.innerText = "No";
    updateArchive.checked ? updateArchiveLabel.innerText = "Yes" : updateArchiveLabel.innerText = "No";
    return;
}

const clearInputs = () => {
    const input = document.querySelectorAll(".uac-input");
    if (input.length === 0) return;
    for (let i = 0; i < input.length; ++i) input[i].value = "";
}

const updatePagination = async () => {
    const fetchUrl = API + "pages/";
    const res = await fetch(fetchUrl);
    if (!res.ok) {
        console.error(res.status);
        return;
    }
    const totalPages = Number(await res.text());
    const pagesToShow = 5;
    let currentPage = 0;
    try { currentPage = Number(window.location.href.split("?")[1].split("=")[1]); } catch (err) { currentPage = 0; }
    if (currentPage <= 0) document.getElementById("prev-page").style.pointerEvents = "none";
    if (currentPage === totalPages) document.getElementById("next-page").style.pointerEvents = "none";
    document.getElementById("last-page").href = `/?page=${totalPages}`;

    const paginationNumbers = document.getElementById("pagination-numbers");
    paginationNumbers.innerHTML = "";
    if (totalPages <= pagesToShow) {
         for (let index = 0; index <= totalPages; ++index) createPaginationLink(index, currentPage);
    } else {
        if (currentPage >= totalPages - 3) {
            for (let index = totalPages - pagesToShow; index <= totalPages; ++index) createPaginationLink(index, currentPage);
        } else {
            for (let index = currentPage; index <= currentPage + pagesToShow; ++index) createPaginationLink(index, currentPage);
        }
    }

    function createPaginationLink(pageNumber, current) {
        const pageNumberLink = document.createElement("a");
        pageNumberLink.innerHTML = pageNumber;
        pageNumberLink.href = `/?page=${pageNumber}`;
        pageNumberLink.title = `Page ${pageNumber}`;
        if (pageNumber === current) pageNumberLink.classList.add("active");
        paginationNumbers.appendChild(pageNumberLink);
    }
}

window.onload = () => {
    const host = new URL(window.location.href).host;
    root = `http://${host}`;
    API = `${root}/api/`;
    updatePagination();
}
