const ENDPOINT = "http://localhost:41415/";

const deleteEntry = async (ele) => {
    const dataID = ele.parentNode.parentNode.id;

    const fetchURL = ENDPOINT + "delete/" + dataID;
    const res = await fetch(fetchURL, {
        method: "GET",
    });

    console.log(res.status);

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

    console.log(res.status);

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

const dmCheck = async () => {
    const fetchURL = ENDPOINT + "dm/";
    const res = await fetch(fetchURL, {
        method: "GET",
    });
    const dm = await res.json();
    if (dm) {
        dmToggle();
    } else {
        return;
    }
};
dmCheck();

const tapCheck = () => {
    const gridChild = document.querySelector(".grid-child");
    if (gridChild === null) {
        return;
    }
    let count = 0;
    gridChild.addEventListener("click", (e) => {
        if (count === 6) {
            dmToggle();
            dmStateChange();
            count = 0;
        } else {
            count++;
        }
        return;
    })
};
tapCheck();

const dmToggle = () => {
    document.body.classList.toggle("dark");
    document.querySelector("nav").classList.toggle("dark");
    document.querySelectorAll("button").forEach(b => b.classList.toggle("dark"));
    try { document.querySelectorAll("input").forEach(i => i.classList.toggle("dark")); } catch (err) { console.log(err); }
    try { document.querySelectorAll(".grid-child").forEach(gc => gc.classList.toggle("dark")); } catch (err) { console.log(err); }
    try { document.querySelectorAll(".grid-child a").forEach(gca => gca.classList.toggle("dark")); } catch (err) { console.log(err); }
    try { document.querySelector(".centered").classList.toggle("dark"); } catch (err) { console.log(err); }
    try { document.querySelectorAll(".update-main div").forEach(um => um.classList.toggle("dark")); } catch (err) { console.log(err); }
    // location.reload();
};

const dmStateChange = async () => {
    const fetchURL = ENDPOINT + "dm/";
    const res = await fetch(fetchURL, {
        method: "GET",
    });
    const dm = await res.json();
    if (!dm) {
        dmPostReq(fetchURL, true);
    } else {
        dmPostReq(fetchURL, false);
    };
};

const dmPostReq = async (fetchURL, s) => {
    await fetch(fetchURL, {
        method: "POST",
        headers: {
            "Content-Type": "application/json; charset=utf-8"
        },
        body: JSON.stringify({ darkMode: s }),
    });
}

const resizeInput = () => {
    const input = document.querySelectorAll("input");
    for (i = 0; i < input.length; i++) {
        (input[i].type === "text") ? input[i].setAttribute("size", input[i].getAttribute("placeholder").length) : {};
        input[i].value = "";
    };
};
resizeInput();