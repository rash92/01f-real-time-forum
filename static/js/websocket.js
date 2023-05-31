let socket = new WebSocket("wss://localhost:8080/ws");
let gotData;
console.log("Attempting Connection...");

socket.onopen = () => {
  console.log("Successfully Connected");
  socket.send("Hi From the Client!");
};

socket.onmessage = function (event) {
  gotData = JSON.parse(event.data);
  console.log("[message] Data receieved from server:", gotData);
  onlineUserInfo(gotData)
};

socket.onclose = (event) => {
  console.log("Socket Closed Connection: ", event);
};

socket.onerror = (error) => {
  console.log("Socket Error: ", error);
};

function onlineUserInfo(data) {
  const usersContainer = document.getElementById("online-users");
  data.forEach((user) => {
    const userDiv = document.createElement("div")
    userDiv.textContent = user.Name
    usersContainer.appendChild(userDiv)
  });
}
