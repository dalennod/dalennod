const ENDPOINT = "http://localhost:41415/";

const deleteEntry = async (ele) => {
    const dataID = ele.parentNode.parentNode.id;

    const fetchURL = ENDPOINT + "delete/" + dataID;
    const res = await fetch(fetchURL, {
        method: "GET",
    });

    console.log(res.status);
}

const addEntry = async () => {
    if (document.querySelector("#input-url").value === "") {
        alert("URL is required.");
        return
    };

    const dataJSON = {
        url: document.querySelector("#input-url").value,
        title: document.querySelector("#input-title").value,
        reason: document.querySelector("#input-reason").value,
        keywords: document.querySelector("#input-keywords").value,
        group: document.querySelector("#input-group").value,
        archive: (document.querySelector("#radio-no").checked) ? 0 : 1,
    };

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
}

const resizeInput = () => {
    var input = document.querySelectorAll("input");
    for (i = 0; i < input.length; i++) {
        (input[i].type === "text") ? input[i].setAttribute("size", input[i].getAttribute("placeholder").length) : {};
        input[i].value = "";
    };
}

resizeInput();