const webSocketStatusElement = document.getElementById("webSocketStatus");
const webSocket = new WebSocket("ws://" + window.location.host + "/ws");

// Web socket open
webSocket.onopen = function () {
    webSocketStatusElement.innerHTML = "Web socket status: <b>Connected</b>";
};

// Web socket receive message
webSocket.onmessage = function (e) {
    let data = JSON.parse(e.data);
    console.log("Message received:", data);
    data.forEach(d => {
       submitModalNewCommand(d.Title, d.Description, d.Command, d.Output)
    });
};

// Web socket close
webSocket.onclose = function () {
    webSocketStatusElement.innerHTML = "Web socket status: <b>Disconnected</b>";
};

// Web socket error
webSocket.onerror = function () {
    webSocketStatusElement.innerHTML = "Web socket status: <b>Error</b>";
};

// Send data to web socket
function send(dataToSend) {
    webSocket.send(JSON.stringify(dataToSend));
}

// On exit
window.onbeforeunload = function () {
    webSocket.close();
};