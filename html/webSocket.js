/* global submitModalNewCommand */
/* global escapeHTMLcharacters */
/* global messageKindCommand */
/* global messageOperationWrite */

let webSocket;

$(document).ready(function () {
    const webSocketStatusElement = document.getElementById("webSocketStatus");
    webSocket = new WebSocket("ws://" + window.location.host + "/ws");

    // Web socket open
    webSocket.onopen = function () {
        webSocketStatusElement.innerHTML = "Web socket status: <b>Connected</b>";
    };

    // Web socket receive message
    webSocket.onmessage = function (e) {
        let data = JSON.parse(e.data);

        // Receive message of "Command" kind
        if (data.kind === messageKindCommand) {
            if (data.operation === messageOperationWrite) {
                data.data.forEach((d) => {
                    submitModalNewCommand(escapeHTMLcharacters(d.title), escapeHTMLcharacters(d.description), escapeHTMLcharacters(d.command), escapeHTMLcharacters(d.output));
                });
            }
        }
    };

    // Web socket close
    webSocket.onclose = function () {
        webSocketStatusElement.innerHTML = "Web socket status: <b>Disconnected</b>";
    };

    // Web socket error
    webSocket.onerror = function () {
        webSocketStatusElement.innerHTML = "Web socket status: <b>Error</b>";
    };

    // On exit
    window.onbeforeunload = function () {
        webSocket.close();
    };
});

// Send data to web socket
const send = function (dataToSend) {
    webSocket.send(JSON.stringify(dataToSend));
};