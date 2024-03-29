const ENDPOINT = "http://localhost:41415/";
let conn = false;

const addEntry = async () => {
    const dataJSON = {
        url: document.querySelector("#input-url").value,
        title: document.querySelector("#input-title").value,
        note: document.querySelector("#input-note").value,
        keywords: document.querySelector("#input-keywords").value,
        bGroup: document.querySelector("#input-bGroup").value,
        archive: (document.querySelector("#radio-no").checked) ? false : true,
    };

    if (dataJSON.archive) {
        document.querySelector("#archive-warn").removeAttribute("hidden");
    };

    const fetchURL = ENDPOINT + "add/";
    const res = await fetch(fetchURL, {
        method: "POST",
        body: JSON.stringify(dataJSON),
    });

    console.log(res.status);
    if (res.ok) {
        resizeInput();
    };

    document.querySelector("#archive-warn").setAttribute("hidden", "");
    document.querySelector("#checkmark").removeAttribute("hidden");
    setTimeout(() => {
        document.querySelector("#checkmark").setAttribute("hidden", "");
    }, 2000);
};

const checkConnection = async () => {
    if (!conn) {
        try {
            const fetchURL = ENDPOINT + "add/";
            const res = await fetch(fetchURL, {
                method: "GET"
            });
            console.log(res.status, await res.text());
        } catch (e) {
            conn = false;
            document.querySelector(".centered").innerHTML = "Dalennod (web-server) must be running.";
            return;
        }
        conn = true;
    };
    browser.tabs.query({ currentWindow: true, active: true }).then((tabs) => {
        const currTab = tabs[0];
        document.querySelector("#input-url").value = currTab.url;
        document.querySelector("#input-title").value = currTab.title;
    });
};

window.addEventListener("load", (event) => {
    checkConnection();
});

document.querySelector("#button-add-req").addEventListener("click", (event) => {
    addEntry();
});

const resizeInput = () => {
    const input = document.querySelectorAll("input");
    for (i = 0; i < input.length; i++) {
        (input[i].type === "text") ? input[i].setAttribute("size", input[i].getAttribute("placeholder").length) : {};
        input[i].value = "";
    };
};
resizeInput();