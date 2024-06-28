"use strict";

const ENDPOINT = "http://localhost:41415/";

const createDialog = document.querySelector(".dialog-create");
const showCreateDialog = () => createDialog.showModal();
const closeCreateDialog = () => createDialog.close();

const updateDialog = document.querySelector(".dialog-update");
const showUpdateDialog = () => updateDialog.showModal();
const closeUpdateDialog = () => updateDialog.close();

const getOldData = async (ele) => {
    const entryID = ele.parentNode.parentNode.id;

    const fetchURL = ENDPOINT + "getRow/" + entryID;
    const res = await fetch(fetchURL);

    const oldData = await res.json();
    if (typeof (Storage) !== "undefined") {
        localStorage.setItem("oldData", JSON.stringify(oldData));
    }
}

const updateEntry = async () => {
    const newDataJSON = {
        url: document.querySelector("#input-url").value,
        title: document.querySelector("#input-title").value,
        note: document.querySelector("#input-note").value,
        keywords: document.querySelector("#input-keywords").value,
        bmGroup: document.querySelector("#input-bmGroup").value,
        archive: document.querySelector("#radio-no").checked ? false : true,
    }

    const dataID = document.querySelector("#bm-id").innerText;
    const fetchURL = ENDPOINT + "update/" + dataID;
    await fetch(fetchURL, {
        method: "POST",
        headers: {
            "Content-Type": "application/json; charset=utf-8"
        },
        body: JSON.stringify(newDataJSON),
    });

    document.querySelector("#checkmark").removeAttribute("hidden");
    setTimeout(() => {
        document.querySelector("#checkmark").setAttribute("hidden", "");
    }, 2000);
}

const addEntry = async () => {
    if (document.querySelector("#input-url").value === "") {
        alert("URL is required");
        return;
    }

    const dataJSON = {
        url: document.querySelector("#input-url").value,
        title: document.querySelector("#input-title").value,
        note: document.querySelector("#input-note").value,
        keywords: document.querySelector("#input-keywords").value,
        bmGroup: document.querySelector("#input-bmGroup").value,
        archive: (document.querySelector("#radio-no").checked) ? false : true,
    }

    if (dataJSON.archive) {
        document.querySelector("#archive-warn").removeAttribute("hidden");
    }

    const fetchURL = ENDPOINT + "add/";
    const res = await fetch(fetchURL, {
        method: "POST",
        headers: {
            "Content-Type": "application/json; charset=utf-8"
        },
        body: JSON.stringify(dataJSON),
    });

    if (res.ok) {
        clearInputs();
    }

    document.querySelector("#archive-warn").setAttribute("hidden", "");
    document.querySelector("#checkmark").removeAttribute("hidden");
    setTimeout(() => {
        document.querySelector("#checkmark").setAttribute("hidden", "");
    }, 2000);
}

const inputEventKey = () => {
    const inputSearch = document.querySelector("#input-search");
    if (inputSearch === null) {
        return;
    }
    inputSearch.addEventListener("keypress", function (event) {
        if (event.key === "Enter") {
            event.preventDefault();
            document.querySelector("#button-search-html").click();
        }
    });
}
inputEventKey();

const clearInputs = () => { 
    const input = document.querySelectorAll(".uac-input");
    if (input.length === 0) return;
    for (let i = 0; i < input.length; i++) input[i].value = ""; 
}