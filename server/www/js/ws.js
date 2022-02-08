export function local_ws(path) {
    var loc = window.location;
    var uri = 'ws:';
    if (loc.protocol === 'https:') {
        uri = 'wss:';
    }
    uri += '//' + loc.host;
    uri += loc.pathname + path; // path registered in server reg_api
    return new WebSocket(uri)
}