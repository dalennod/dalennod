"use strict";

const ENDPOINT = "http://localhost:41415/";

const existsIconPaths = {
    48: "icons/exists/dalennod-exists-48.png",
    96: "icons/exists/dalennod-exists-96.png",
};
const defaultIconPaths = {
    48: "icons/dalennod-48.png",
    96: "icons/dalennod-96.png",
};

const checkUrl = async (currTabUrl) => {
    const fetchUrl = ENDPOINT + "checkUrl/";
    const res = await fetch(fetchUrl, {
        method: "POST",
        body: JSON.stringify({ url: currTabUrl }),
    });
    if (res.status === 404) {
        await chrome.action.setIcon({ path: defaultIconPaths });
        return;
    }
    await chrome.action.setIcon({ path: existsIconPaths });
};

let conn = false;
const checkConnection = async () => {
    if (!conn) {
        try {
            const fetchURL = ENDPOINT + "add/";
            const res = await fetch(fetchURL);
            console.log(res.status, await res.text());
        } catch (err) {
            conn = false;
            console.error("server not running. ERROR:", err);
            return;
        }
        conn = true;
    }
    await checkUrl(currTab);
};

let currTab = "";
chrome.tabs.onActivated.addListener((activeInfo) => {
    chrome.tabs.get(activeInfo.tabId, async (tab) => {
        currTab = await tab.url;
        await checkConnection();
    });
});
