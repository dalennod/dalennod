"use strict";

const ENDPOINT = "http://localhost:41415/";

const [overlayDiv, overlayText, checkmark, archiveWarn] = [document.querySelector("#done-overlay-div"), document.querySelector("#done-overlay-text"), document.querySelector("#checkmark"), document.querySelector("#archive-warn")];
const [btnCreate, btnUpdate, btnRemove, btnArchive] = [document.querySelector("#button-add-req"), document.querySelector("#button-update-req"), document.querySelector("#button-remove-req"), document.querySelector("#radio-btn-archive")];
const [bmId, inputUrl, inputTitle, inputNote, inputKeywords, inputBmGroup, bmGroupsList, archiveRadioNo] = [document.querySelector("#bm-id"), document.querySelector("#input-url"), document.querySelector("#input-title"), document.querySelector("#input-note"), document.querySelector("#input-keywords"), document.querySelector("#input-bmGroup"), document.querySelector("#bmGroups-list"), document.querySelector("#radio-no")];

window.addEventListener("load", () => { getCurrTab(); checkConnection(); });

let currTab = "";
const getCurrTab = () => {
    chrome.tabs.query({ currentWindow: true, active: true }).then((tabs) => {
        currTab = tabs[0];
        inputUrl.value = currTab.url;
        inputTitle.value = currTab.title;
    });
}

let conn = false;
const checkConnection = async () => {
    if (!conn) {
        try {
            const fetchURL = ENDPOINT + "add/";
            const res = await fetch(fetchURL);
            console.log(res.status, await res.text());
        } catch (e) {
            conn = false;
            document.querySelector(".centered").innerHTML = `<a href="https://github.com/dalennod/dalennod" target="_blank"> <span style="text-decoration: underline;">Dalennod</span> </a>(web-server) must be running.`;
            return;
        }
        conn = true;
    }
    checkUrl(currTab.url);
    await fillAllGroups(getGroups());
}

const checkUrl = async (currTabUrl) => {
    const fetchUrl = ENDPOINT + "checkUrl/";
    const res = await fetch(fetchUrl, {
        method: "POST",
        body: JSON.stringify({ url: currTabUrl }),
    })
    if (res.status == 404) {
        return;
    }
    const receviedData = await res.json();

    btnUpdate.removeAttribute("hidden");
    btnRemove.removeAttribute("hidden");
    btnCreate.setAttribute("hidden", "");
    fillData(JSON.parse(JSON.stringify(receviedData)));
}

const fillData = (dataFromDb) => {
    bmId.innerHTML = dataFromDb.id;
    inputUrl.value = dataFromDb.url;
    inputTitle.value = dataFromDb.title;
    inputNote.value = dataFromDb.note;
    inputKeywords.value = dataFromDb.keywords;
    inputBmGroup.value = dataFromDb.bmGroup;
    dataFromDb.archive ? btnArchive.setAttribute("hidden", "") : btnArchive.removeAttribute("hidden");
}

const getGroups = async () => {
    const fetchUrl = ENDPOINT + "groups/";
    const res = await fetch(fetchUrl);
    return await res.text();
}
const fillAllGroups = async (allGroups) => bmGroupsList.innerHTML = await allGroups;

btnCreate.addEventListener("click", () => addEntry());
const addEntry = async () => {
    const dataJSON = {
        url: inputUrl.value,
        title: inputTitle.value,
        note: inputNote.value,
        keywords: inputKeywords.value,
        bmGroup: inputBmGroup.value,
        archive: archiveRadioNo.checked ? false : true,
    }

    if (dataJSON.archive) {
        archiveWarn.removeAttribute("hidden");
        btnCreate.disabled = true;
    }

    const fetchURL = ENDPOINT + "add/";
    const res = await fetch(fetchURL, {
        method: "POST",
        body: JSON.stringify(dataJSON),
    });

    if (res.ok) {
        resizeInput();
    }

    archiveWarn.setAttribute("hidden", "");
    btnCreate.disabled = false;
    checkmark.removeAttribute("hidden");
    overlayText.innerHTML = "Created&nbsp;&check;";
    overlayDiv.style.display = "block";
    setTimeout(() => {
        checkmark.setAttribute("hidden", "");
        overlayDiv.style.display = "none";
        window.close();
    }, 2000);
}

btnUpdate.addEventListener("click", () => updateEntry(bmId.innerHTML));
const updateEntry = async (idInDb) => {
    const dataJSON = {
        url: inputUrl.value,
        title: inputTitle.value,
        note: inputNote.value,
        keywords: inputKeywords.value,
        bmGroup: inputBmGroup.value,
        archive: archiveRadioNo.checked ? false : true,
    }

    const fetchURL = ENDPOINT + "update/" + idInDb;
    const res = await fetch(fetchURL, {
        method: "POST",
        body: JSON.stringify(dataJSON),
    });

    if (res.ok) {
        checkmark.removeAttribute("hidden");
        overlayText.innerHTML = "Updated&nbsp;&check;";
        overlayDiv.style.display = "block";
        setTimeout(() => {
            checkmark.setAttribute("hidden", "");
            overlayDiv.style.display = "none";
            window.close();
        }, 2000);
    }
}

btnRemove.addEventListener("click", () => removeEntry(bmId.innerHTML));
const removeEntry = async (idInDb) => {
    const fetchURL = ENDPOINT + "delete/" + idInDb;
    const res = await fetch(fetchURL);

    if (res.ok) resizeInput();

    checkmark.removeAttribute("hidden");
    overlayText.innerHTML = "Removed&nbsp;&check;";
    overlayDiv.style.display = "block";
    setTimeout(() => {
        checkmark.setAttribute("hidden", "");
        overlayDiv.style.display = "none";
        window.close();
    }, 2000);
}

const resizeInput = () => {
    const input = document.querySelectorAll("input");
    for (let i = 0; i < input.length; i++) {
        (input[i].type === "text") ? input[i].setAttribute("size", input[i].getAttribute("placeholder").length) : {};
        input[i].value = "";
    }
}
resizeInput();
