const ENDPOINT = "http://localhost:41415/";

const deleteEntry = async (ele) => {
    const dataID = ele.parentNode.id;
    await fetch(ENDPOINT+"delete/"+dataID)
        .then(res => {
            if (!res.ok) {
                console.log("Failed.", res);
            }
            else {
                console.log("Deleted entry #"+dataID);
            }
        });
}