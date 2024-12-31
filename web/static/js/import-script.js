"use strict";

const API_ENDPOINT = "http://localhost:41415/api/";

const handleFileSelect = async (event) => {
    event.preventDefault();
    const file = fileInput.files[0];
    updateLog("Selected file: " + file.name);

    const url = API_ENDPOINT + "import-bookmark/";
    const formData = new FormData();
    formData.append("importFile", file);
    formData.append("selectedBrowser", selectedBrowser.value);
    try {
        const response = await fetch(url, {
            method: "POST",
            body: formData,
        });
        const result = await response.text();
        if (!response.ok) {
            updateLog(`http error. Status: ${response.status}. HTTP reply: ${result}`);
            return;
        }
        updateLog("File uploaded successfully.");
        updateLog(result);
    } catch (err) {
        updateLog("Error uploading file. ERROR:", err);
    }
};

const handleDragOver = (event) => {
    console.log("File in drop zone");
    event.preventDefault();
};

const handleDrop = (event) => {
    console.log("File dropped");
    event.preventDefault();
};

const updateLog = (message) => {
    const logTextArea = document.getElementById("logTextArea");
    logTextArea.value += message + "\n";
    logTextArea.scrollTop = logTextArea.scrollHeight;
};

let selectedBrowser = document.querySelector("input[name='browser']:checked");
document.getElementById("browserForm").addEventListener("change", () => {
    selectedBrowser = document.querySelector("input[name='browser']:checked");
    if (selectedBrowser) updateLog(`Selected Browser: ${selectedBrowser.value}`);
});

const fileInput = document.getElementById("fileInput");
fileInput.addEventListener("change", handleFileSelect, false);
