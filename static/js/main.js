const ENDPOINT = "http://localhost:41415/";

const deleteEntry = async (ele) => {
    const dataID = ele.parentNode.parentNode.id;

    const fetchURL = ENDPOINT + "delete/" + dataID;
    const res = await fetch(fetchURL, {
        method: "GET",
    });

    console.log(res.status);
};

const getOldData = async (ele) => {
    const entryID = ele.parentNode.parentNode.id;

    const fetchURL = ENDPOINT + "getRow/" + entryID;
    const res = await fetch(fetchURL, {
        method: "GET",
    });

    console.log(res.status);

    const oldData = await res.json();
    if (typeof (Storage) !== "undefined") {
        localStorage.setItem('old-data', JSON.stringify(oldData));
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

    // console.log(JSON.stringify(newDataJSON));

    const dataID = document.querySelector("#data-id").innerText;
    const fetchURL = ENDPOINT + "update/" + dataID;
    const res = await fetch(fetchURL, {
        method: "POST",
        headers: {
            "Content-Type": "application/json; charset=utf-8"
        },
        body: JSON.stringify(newDataJSON),
    });

    console.log(res.status);

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

    // console.log(JSON.stringify(dataJSON));

    const fetchURL = ENDPOINT + "add/";
    const res = await fetch(fetchURL, {
        method: "POST",
        headers: {
            "Content-Type": "application/json; charset=utf-8"
        },
        body: JSON.stringify(dataJSON),
    });

    console.log(res.status);
    if (res.ok) {
        resizeInput();
    };
};

const resizeInput = () => {
    const input = document.querySelectorAll("input");
    for (i = 0; i < input.length; i++) {
        (input[i].type === "text") ? input[i].setAttribute("size", input[i].getAttribute("placeholder").length) : {};
        input[i].value = "";
    };
};

resizeInput();