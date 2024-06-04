const ENDPOINT = "http://localhost:41415/";

const overlayDiv = document.querySelector("#done-overlay-div");
const overlayText = document.querySelector("#done-overlay-text");
const checkmark = document.querySelector("#checkmark");
const archiveWarn = document.querySelector("#archive-warn");
const btnCreate = document.querySelector("#button-add-req");
const btnUpdate = document.querySelector("#button-update-req");
const btnRemove = document.querySelector("#button-remove-req");
const btnArchive = document.querySelector("#radio-btn-archive");

const bmId = document.querySelector("#bm-id");
const inputUrl = document.querySelector("#input-url");
const inputTitle = document.querySelector("#input-title");
const inputNote = document.querySelector("#input-note");
const inputKeywords = document.querySelector("#input-keywords");
const inputBGroup = document.querySelector("#input-bGroup");
const archiveRadioNo = document.querySelector("#radio-no");

let conn = false;
let currTab = "";

const addEntry = async () => {
    const dataJSON = {
        url: inputUrl.value,
        title: inputTitle.value,
        note: inputNote.value,
        keywords: inputKeywords.value,
        bGroup: inputBGroup.value,
        archive: archiveRadioNo.checked ? false : true,
    };

    if (dataJSON.archive) {
        archiveWarn.removeAttribute("hidden");
        btnCreate.disabled = true;
    };

    const fetchURL = ENDPOINT + "add/";
    const res = await fetch(fetchURL, {
        method: "POST",
        body: JSON.stringify(dataJSON),
    });

    if (res.ok) {
        resizeInput();
    };

    archiveWarn.setAttribute("hidden", "");
    btnCreate.disabled = false;
    checkmark.removeAttribute("hidden");
    overlayText.innerHTML = "Created&nbsp;&check;";
    overlayDiv.style.display = "block";
    setTimeout(() => {
        checkmark.setAttribute("hidden", "");
        overlayDiv.style.display = "none";
    }, 2000);
};

const getCurrTab = () => {
    browser.tabs.query({ currentWindow: true, active: true }).then((tabs) => {
        currTab = tabs[0];
        inputUrl.value = currTab.url;
        inputTitle.value = currTab.title;
    });
};

const checkConnection = async () => {
    if (!conn) {
        try {
            const fetchURL = ENDPOINT + "add/";
            const res = await fetch(fetchURL);
            console.log(res.status, await res.text());
        } catch (e) {
            conn = false;
            document.querySelector(".centered").innerHTML = `
                <a href="https://github.com/dalennod/dalennod" target="_blank"> <span style="text-decoration: underline;">Dalennod</span> </a>
                (web-server) must be running.`;
            return;
        }
        conn = true;
    };
    checkUrl(currTab.url);
};

const checkUrl = async (currTabUrl) => {
    const fetchUrl = ENDPOINT + "checkUrl/";
    const res = await fetch(fetchUrl, {
        method: "POST",
        body: JSON.stringify({ url: currTabUrl }),
    })
    if (res.status == 404) {
        return;
    };
    const receviedData = await res.json();

    btnUpdate.removeAttribute("hidden");
    btnRemove.removeAttribute("hidden");
    btnCreate.setAttribute("hidden", "");
    fillData(JSON.parse(JSON.stringify(receviedData)));
};

const fillData = (dataFromDb) => {
    bmId.innerHTML = dataFromDb.id;
    inputUrl.value = dataFromDb.url;
    inputTitle.value = dataFromDb.title;
    inputNote.value = dataFromDb.note;
    inputKeywords.value = dataFromDb.keywords;
    inputBGroup.value = dataFromDb.bGroup;
    dataFromDb.archive ? btnArchive.setAttribute("hidden", "") : btnArchive.removeAttribute("hidden");
};

window.addEventListener("load", () => {
    getCurrTab();
    checkConnection();
});

btnCreate.addEventListener("click", () => {
    addEntry();
});

btnUpdate.addEventListener("click", () => {
    updateEntry(bmId.innerHTML);
});

btnRemove.addEventListener("click", () => {
    removeEntry(bmId.innerHTML);
});

const updateEntry = async (idInDb) => {
    const dataJSON = {
        url: inputUrl.value,
        title: inputTitle.value,
        note: inputNote.value,
        keywords: inputKeywords.value,
        bGroup: inputBGroup.value,
        archive: archiveRadioNo.checked ? false : true,
    };

    const fetchURL = ENDPOINT + "update/" + idInDb;
    await fetch(fetchURL, {
        method: "POST",
        body: JSON.stringify(dataJSON),
    });

    checkmark.removeAttribute("hidden");
    overlayText.innerHTML = "Updated&nbsp;&check;";
    overlayDiv.style.display = "block";
    setTimeout(() => {
        checkmark.setAttribute("hidden", "");
        overlayDiv.style.display = "none";
    }, 2000);
};

const removeEntry = async (idInDb) => {
    const fetchURL = ENDPOINT + "delete/" + idInDb;
    await fetch(fetchURL);
    
    checkmark.removeAttribute("hidden");
    overlayText.innerHTML = "Removed&nbsp;&check;";
    overlayDiv.style.display = "block";
    setTimeout(() => {
        checkmark.setAttribute("hidden", "");
        overlayDiv.style.display = "none";
    }, 2000);
};

const resizeInput = () => {
    const input = document.querySelectorAll("input");
    for (i = 0; i < input.length; i++) {
        (input[i].type === "text") ? input[i].setAttribute("size", input[i].getAttribute("placeholder").length) : {};
        input[i].value = "";
    };
};
resizeInput();