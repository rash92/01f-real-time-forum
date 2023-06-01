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
  onlineUserInfo(gotData);
};

socket.onclose = (event) => {
  console.log("Socket Closed Connection: ", event);
};

socket.onerror = (error) => {
  console.log("Socket Error: ", error);
};

function onlineUserInfo(data) {
  const usersContainer = document.getElementById("online-users");
  const userArr = Array.from(document.getElementsByClassName("users"));

  data.forEach((user) => {
    // Try to find existing userDiv for the user
    const existingUserDiv = userArr.find(
      (onlineUser) => onlineUser.textContent.split(" ")[1] === user.Name
    );

    let userDiv;
    if (existingUserDiv) {
      // If userDiv already exists, use it
      userDiv = existingUserDiv;
    } else {
      // If userDiv does not exist, create it
      userDiv = document.createElement("div");
      userDiv.classList.add("users");
      usersContainer.appendChild(userDiv);
      userDiv.addEventListener("click", function () {
        const userName = { type: "recipientSelect", info: {name: user.Name} };
        console.log(userName);
        socket.send(JSON.stringify(userName));
      });
    }

    // Update textContent based on user's online status
    if (user.LoggedInStatus === 0) {
      userDiv.textContent = `🔴 ${user.Name}`;
    } else {
      userDiv.textContent = `🟢 ${user.Name}`;
    }
  });
}
