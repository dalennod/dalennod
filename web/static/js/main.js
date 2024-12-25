"use strict";

const API_ENDPOINT = "http://localhost:41415/api/";
const ROOT_ENDPOINT = "http://localhost:41415/";

const createDialog = document.querySelector(".dialog-create");
const showCreateDialog = () => {
    createDialog.showModal();
    clearImportTimeout();
};
const closeCreateDialog = () => createDialog.close();

const updateDialog = document.querySelector(".dialog-update");
const showUpdateDialog = () => updateDialog.showModal();
const closeUpdateDialog = () => updateDialog.close();

const getOldData = async (ele) => {
    const entryID = ele.parentNode.parentNode.id;
    const fetchURL = API_ENDPOINT + "row/" + entryID;

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
        document.querySelector("#bm-id").innerText = oldData.id;
        document.querySelector("#update-url").value = oldData.url;
        document.querySelector("#update-title").value = oldData.title;
        document.querySelector("#update-note").value = oldData.note;
        document.querySelector("#update-keywords").value = oldData.keywords;
        document.querySelector("#update-bmGroup").value = oldData.bmGroup;
        oldData.archive ? document.querySelector("#update-archive-div").setAttribute("hidden", "") : document.querySelector("#update-archive-div").removeAttribute("hidden");
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
    };

    if (newDataJSON.archive) {
        document.querySelector("#update-archive-warn").removeAttribute("hidden");
    } else if (!newDataJSON.archive && oldData.archive) {
        newDataJSON.archive = true;
        newDataJSON.snapshotURL = oldData.snapshotURL;
    }
    
    const updateButton = document.getElementById("button-update-req");
    updateButton.disabled = true;

    const dataID = document.querySelector("#bm-id").innerText;
    const fetchURL = API_ENDPOINT + "update/" + dataID;
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
    };

    if (dataJSON.archive) {
        document.getElementById("create-archive-warn").removeAttribute("hidden");
    }
    
    const addButton = document.getElementById("button-add-req");
    addButton.disabled = true;

    const fetchURL = API_ENDPOINT + "add/";
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

let changeToImportTimeout = "";
const changeToImport = () => {
    changeToImportTimeout = setTimeout(() => showImportPage(), 5000);
    return;
};

const showImportPage = () => {
    document.location.href = ROOT_ENDPOINT + "import/";
    return;
};

const clearImportTimeout = () => {
    clearTimeout(changeToImportTimeout);
    return;
};

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
    for (let i = 0; i < input.length; i++) input[i].value = "";
};
