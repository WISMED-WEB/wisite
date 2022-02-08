export function fetch_get_json(url) {
    const cData = fetch(url)
        .then((resp) => resp.json())
        .then((data) => {
            // console.log("fetch_get_json:", data);
            return data;
        }).catch((error) => { console.log(error) });
    return cData;
}

export function fetch_get(url) {
    const cData = fetch(url)
        .then((resp) => {
            // console.log("fetch_get:", resp);
            return resp.text()
        }).catch((error) => { console.log(error) });
    return cData;
}