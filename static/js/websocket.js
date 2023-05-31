

let socket = new WebSocket("wss://localhost:8080/ws");
console.log("Attempting Connection...");

socket.onopen = () => {
    console.log("Successfully Connected");
    socket.send("Hi From the Client!");
};

socket.onmessage = function (event) {
    let gotData = JSON.parse(event.data)
    console.log("[message] Data receieved from server:", gotData)
}

socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
};

socket.onerror = error => {
    console.log("Socket Error: ", error);
};