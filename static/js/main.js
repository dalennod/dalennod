const ENDPOINT = "http://localhost:41415/";

const deleteEntry = async (ele) => {
    const dataID = ele.parentNode.parentNode.id;

    const fetchURL = ENDPOINT + "delete/" + dataID;
    await fetch(fetchURL, {
        method: "GET",
    });

    document.querySelector("#delete-checkmark-" + dataID).removeAttribute("hidden")
    setTimeout(() => {
        document.querySelector("#delete-checkmark-" + dataID).setAttribute("hidden", "");
    }, 2000);
};

const getOldData = async (ele) => {
    const entryID = ele.parentNode.parentNode.id;

    const fetchURL = ENDPOINT + "getRow/" + entryID;
    const res = await fetch(fetchURL, {
        method: "GET",
    });

    const oldData = await res.json();
    if (typeof (Storage) !== "undefined") {
        localStorage.setItem("oldData", JSON.stringify(oldData));
    };
};

const updateEntry = async () => {
    const newDataJSON = {
        url: (document.querySelector("#input-url").value === "") ? document.querySelector("#old-url").value : document.querySelector("#input-url").value,
        title: (document.querySelector("#input-title").value === "") ? document.querySelector("#old-title").value : document.querySelector("#input-title").value,
        note: (document.querySelector("#input-note").value === "") ? document.querySelector("#old-note").value : document.querySelector("#input-note").value,
        keywords: (document.querySelector("#input-keywords").value === "") ? document.querySelector("#old-keywords").value : document.querySelector("#input-keywords").value,
        bGroup: (document.querySelector("#input-bGroup").value === "") ? document.querySelector("#old-bGroup").value : document.querySelector("#input-bGroup").value,
        archive: (document.querySelector("#radio-no").checked) ? false : true,
    };

    const dataID = document.querySelector("#data-id").innerText;
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
};

const addEntry = async () => {
    if (document.querySelector("#input-url").value === "") {
        alert("URL is required.");
        return
    };

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
        headers: {
            "Content-Type": "application/json; charset=utf-8"
        },
        body: JSON.stringify(dataJSON),
    });

    if (res.ok) {
        resizeInput();
    };

    document.querySelector("#archive-warn").setAttribute("hidden", "");
    document.querySelector("#checkmark").removeAttribute("hidden");
    setTimeout(() => {
        document.querySelector("#checkmark").setAttribute("hidden", "");
    }, 2000);
};

const setValue = (term) => {
    document.querySelector("#hidden-url-param").value = term;
};

const inputEventKey = () => {
    const inputSearch = document.querySelector("#input-search");
    if (inputSearch === null) {
        return;
    };
    inputSearch.addEventListener("keypress", function (event) {
        if (event.key === "Enter") {
            event.preventDefault();
            document.querySelector("#button-search-html").click();
        };
    });
};
inputEventKey();

const resizeInput = () => {
    const input = document.querySelectorAll("input");
    for (i = 0; i < input.length; i++) {
        (input[i].type === "text") ? input[i].setAttribute("size", input[i].getAttribute("placeholder").length) : {};
        input[i].value = "";
    };
};
resizeInput();