const usersContainer = document.getElementById("online-users");
const nonMessagedContainer = document.getElementById("non-messaaged-users");
const typingProgressDiv = document.createElement("div");

const startWebSocket = () => {
  let socket = new WebSocket("wss://localhost:8080/ws");
  let onlineUsersData;
  let clientInfo;
  let chatQueue = [];
  // let size = 1
  console.log("Attempting Connection...");

  socket.onopen = () => {
    console.log("Successfully Connected");
    socket.send("Hi From the Client!");
  };

  socket.onmessage = function (event) {
    let message = JSON.parse(event.data);
    // console.log("TYPE: ", typeof message, "\nDATA: ", message)
    switch (message.type) {
      case "onlineUsers":
        console.log("Online users triggered in onmessage");
        onlineUserInfo(message.data.messagedUsers, usersContainer);
        onlineUserInfo(message.data.nonMessagedUsers, nonMessagedContainer);
        break;
      case "userInfo":
        clientInfo = message.data;
        if (clientInfo.IsLoggedIn === 1) {
          usersContainer.style.display = "block";
        }
        break;
      case "typing":
        clientInfo = {
          typing: false,
        };
        let sender = message.data.username;
        let isTyping = message.data.isTyping;
        console.log(isTyping);
        if (isTyping) {
          typingProgressDiv.innerText = sender + " is typing...";
        } else {
          typingProgressDiv.innerText = "";
        }
        break;
      case "chatSelect":
        renderChat(message);
        break;
      case "private":
        console.log(message);
        renderChat(message);
        messageNotification(message);
        break;
    }
  };

  socket.onclose = (event) => {
    console.log("Socket Closed Connection: ", event);
  };

  socket.onerror = (error) => {
    console.log("Socket Error: ", error);
  };

  function onlineUserInfo(data, container) {
    const userArr = [];
    container.innerHTML = "";
    const chatBoxesContainer = document.getElementById("chat-boxes-container");
    data.forEach((user) => {
      // Try to find existing userDiv for the user
      const existingUserDiv = userArr.find(
        (onlineUser) =>
          onlineUser.textContent.split(" ")[1] ===
          user.Name.charAt(0).toUpperCase() + user.Name.slice(1)
      );

      let userDiv;
      if (existingUserDiv) {
        // If userDiv already exists, use it
        userDiv = existingUserDiv;
      } else {
        // If userDiv does not exist, create it
        userDiv = document.createElement("div");
        userDiv.classList.add("users");
        container.appendChild(userDiv);
        userDiv.addEventListener("click", function () {
          const chatDiv = document.createElement("div");
          if (chatQueue.length === 1) {
            const oldestChatDiv = chatQueue.shift();
            chatBoxesContainer.removeChild(oldestChatDiv);
          }
          chatDiv.classList.add("chat-box");

          const hideChat = document.createElement("div");
          hideChat.textContent = "hide";

          const chatTitle = document.createElement("div");
          chatTitle.classList.add("chat-title");
          chatTitle.textContent =
            user.Name.charAt(0).toUpperCase() + user.Name.slice(1);

          // const typingProgressDiv = document.createElement("div")
          // typingProgressDiv.classList.add("")

          const chatContentDiv = document.createElement("div");
          chatContentDiv.classList.add("chat-content");

          const chatInputDiv = document.createElement("div");
          chatInputDiv.classList.add("chat-input-div");

          const chatInput = document.createElement("textarea");
          chatInput.rows = 1;
          chatInput.classList.add("chat-input");
          chatInput.placeholder = "Type your message here...";

          chatInput.addEventListener("input", () => {
            chatInput.style.height = "auto";
            chatInput.style.height = `${chatInput.scrollHeight}px`;

            if (!clientInfo.typing) {
              const typingMessage = {
                type: "typing",
                info: {
                  //recipient: user.Name,
                  isTyping: true,
                },
              };
              // console.log(typingMessage)
              socket.send(JSON.stringify(typingMessage));
            }
          });

          // Send a typing notification when the user stops typing
          chatInput.addEventListener("blur", () => {
            if (clientInfo.typing) {
              const typingMessage = {
                type: "typing",
                info: {
                  // recipient: user.Name,
                  isTyping: false,
                },
              };
              socket.send(JSON.stringify(typingMessage));
            }
          });

          chatInputDiv.append(chatInput);
          const sendButton = document.createElement("button");
          sendButton.textContent = "Send";
          sendButton.classList.add("send-button");
          sendButton.addEventListener("click", function () {
            const text = document.querySelector(".chat-input");

            let messageToSend = {
              type: "private",
              info: {
                recipient: user.Name,
                text: text.value,
              },
            };

            console.log("this is the private message: ", messageToSend);
            socket.send(JSON.stringify(messageToSend));
            text.value = "";
          });
          chatInputDiv.append(sendButton);

          chatDiv.append(chatTitle);
          chatTitle.append(hideChat);
          chatDiv.append(typingProgressDiv);
          chatDiv.append(chatContentDiv);
          chatDiv.append(chatInputDiv);
          chatBoxesContainer.appendChild(chatDiv);
          chatQueue.push(chatDiv);

          hideChat.addEventListener("click", () => {
            chatDiv.style.display = "none";
          });

          const userName = {
            type: "recipientSelect",
            info: { name: user.Name },
          };

          // console.log(userName)
          socket.send(JSON.stringify(userName));
        });
      }

      // Update textContent based on user's online status
      if (user.LoggedInStatus === 0) {
        userDiv.textContent = `🔴 ${
          user.Name.charAt(0).toUpperCase() + user.Name.slice(1)
        }`;
      } else {
        userDiv.textContent = `🟢 ${
          user.Name.charAt(0).toUpperCase() + user.Name.slice(1)
        }`;
      }
    });
  }
};

const renderChat = (obj, size = 1) => {
  const chatBox = document.getElementsByClassName("chat-content")[0];
  const recipientBox = document.getElementsByClassName("chat-title")[0];
  if (!recipientBox) {
    return;
  }
  const recipientName = recipientBox.innerText.split("\n")[0].toLowerCase();
  console.log(recipientName);
  //delete everything within the chatBox
  chatBox.innerHTML = "";

  let totalChatSize = obj.data.Content.length;
  let index = totalChatSize - size * 10;
  console.log("Expected number of lines: ", size * 10);
  console.log("Total number of messages: ", totalChatSize);

  //re-render the most recent
  if (index <= 0) {
    index = 0;
  }
  // obj.data.Content.reverse();
  for (index; index <= totalChatSize - 1; index++) {
    let value = obj.data.Content[index];

    const chatBubble = document.createElement("div");
    chatBubble.classList.add("chat-bubble");

    const text = document.createElement("div");
    text.classList.add("text");
    text.innerText = value.content;

    const time = document.createElement("div");
    time.classList.add("time");

    if (recipientName === value.receiver) {
      time.innerText = `${value.sender} ${value.time}`;
      chatBubble.classList.add("sent");
    } else {
      time.innerText = `${recipientName} ${value.time}`;
      chatBubble.classList.add("received");
    }
    chatBubble.appendChild(text);
    chatBubble.appendChild(time);

    chatBox.insertBefore(chatBubble, chatBox.firstChild);
  }

  let reachedMax = totalChatSize % (size * 10) != 0;
  let isThrottled = false;

  chatBox.addEventListener("scroll", function () {
    if (!isThrottled && chatBox.scrollTop === 0 && reachedMax) {
      isThrottled = true;

      setTimeout(function () {
        renderChat(obj, size + 1);
        isThrottled = false;
      }, 1000); // Adjust the throttle delay as needed (e.g., 500 milliseconds)
    }
  });
};

function messageNotification(message) {
  const currentUser = document
    .getElementById("current-user")
    .textContent.split(" ")[1];
  const usersArr = Array.from(document.getElementsByClassName("users"));
  const messageArr = message.data.Content;
  const chatBox = document.querySelector(".chat-box");
  const currentlyTalkingTo = document.querySelector(".chat-title");

  let sender;
  if (currentUser === messageArr[messageArr.length - 1].sender) {
    sender = messageArr[messageArr.length - 1].reciever;
  } else {
    sender = messageArr[messageArr.length - 1].sender;
  }
  const senderDiv = usersArr.find(
    (onlineUser) => onlineUser.textContent.split(" ")[1] === sender
  );
  if (
    !(
      currentlyTalkingTo &&
      currentlyTalkingTo.textContent.toLowerCase().split("hide")[0] ===
        sender &&
      chatBox &&
      chatBox.style.display !== "none"
    ) &&
    messageArr[messageArr.length - 1].sender !== currentUser
  ) {
    alert(sender + " has sent you a message");
  }
}
