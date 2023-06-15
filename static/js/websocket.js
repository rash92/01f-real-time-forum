const usersContainer = document.getElementById("online-users")
const typingProgressDiv = document.createElement("div")

const startWebSocket = () => {
  let socket = new WebSocket("wss://localhost:8080/ws")
  let onlineUsersData
  let clientInfo
  let chatQueue = []
  console.log("Attempting Connection...")

  socket.onopen = () => {
    console.log("Successfully Connected")
    socket.send("Hi From the Client!")
  }

  socket.onmessage = function (event) {
    let message = JSON.parse(event.data)
    // console.log("TYPE: ", typeof message, "\nDATA: ", message)
    switch (message.type) {
      case "onlineUsers":
        onlineUsersData = message.data
        onlineUserInfo(onlineUsersData)
        break
      case "userInfo":
        clientInfo = message.data
        if (clientInfo.IsLoggedIn === 1) {
          usersContainer.style.display = "block"
        }
        break
      case "typing":
        clientInfo = {
          typing: false,
        }
        let sender = message.data.username
        let isTyping = message.data.isTyping
        console.log(isTyping);
        if (isTyping) {
          typingProgressDiv.innerText = sender + " is typing..."
        } else {
          typingProgressDiv.innerText = ""
        }
        break
      case "chatSelect":
      case "private":
        console.log(`THIS IS A ${message.type} MESSAGE: \n`, message.data)
        renderChat(message)
        break
    }
  }

  socket.onclose = (event) => {
    console.log("Socket Closed Connection: ", event)
  }

  socket.onerror = (error) => {
    console.log("Socket Error: ", error)
  }

  function onlineUserInfo(data) {
    const userArr = Array.from(document.getElementsByClassName("users"))
    const chatBoxesContainer = document.getElementById("chat-boxes-container")
    data.forEach((user) => {
      // Try to find existing userDiv for the user
      const existingUserDiv = userArr.find(
        (onlineUser) =>
          onlineUser.textContent.split(" ")[1] ===
          user.Name.charAt(0).toUpperCase() + user.Name.slice(1)
      )

      let userDiv
      if (existingUserDiv) {
        // If userDiv already exists, use it
        userDiv = existingUserDiv
      } else {
        // If userDiv does not exist, create it
        userDiv = document.createElement("div")
        userDiv.classList.add("users")
        usersContainer.appendChild(userDiv)
        userDiv.addEventListener("click", function () {
          const chatDiv = document.createElement("div")
          if (chatQueue.length === 1) {
            const oldestChatDiv = chatQueue.shift()
            chatBoxesContainer.removeChild(oldestChatDiv)
          }
          chatDiv.classList.add("chat-box")

          const hideChat = document.createElement("div")
          hideChat.textContent = "hide"

          const chatTitle = document.createElement("div")
          chatTitle.classList.add("chat-title")
          chatTitle.textContent =
            user.Name.charAt(0).toUpperCase() + user.Name.slice(1)

          // const typingProgressDiv = document.createElement("div")
          // typingProgressDiv.classList.add("")

          const chatContentDiv = document.createElement("div")
          chatContentDiv.classList.add("chat-content")

          const chatInputDiv = document.createElement("div")
          chatInputDiv.classList.add("chat-input-div")

          const chatInput = document.createElement("textarea")
          chatInput.rows = 1
          chatInput.classList.add("chat-input")
          chatInput.placeholder = "Type your message here..."

          chatInput.addEventListener("input", () => {
            chatInput.style.height = "auto"
            chatInput.style.height = `${chatInput.scrollHeight}px`

            if (!clientInfo.typing) {
              const typingMessage = {
                type: "typing",
                info: {
                  //recipient: user.Name,
                  isTyping: true,
                },
              }
              // console.log(typingMessage)
              socket.send(JSON.stringify(typingMessage))
            }
          })

          // Send a typing notification when the user stops typing
          chatInput.addEventListener("blur", () => {
            if (clientInfo.typing) {
              const typingMessage = {
                type: "typing",
                info: {
                  // recipient: user.Name,
                  isTyping: false,
                },
              }
              socket.send(JSON.stringify(typingMessage))
            }
          })

          chatInputDiv.append(chatInput)
          const sendButton = document.createElement("button")
          sendButton.textContent = "Send"
          sendButton.classList.add("send-button")
          sendButton.addEventListener("click", function () {
            const text = document.querySelector(".chat-input")

            let messageToSend = {
              type: "private",
              info: {
                recipient: user.Name,
                text: text.value,
              },
            }

            console.log("this is the private message: ", messageToSend)
            socket.send(JSON.stringify(messageToSend))
            text.value = ""
          })
          chatInputDiv.append(sendButton)

          chatDiv.append(chatTitle)
          chatTitle.append(hideChat)
          chatDiv.append(typingProgressDiv)
          chatDiv.append(chatContentDiv)
          chatDiv.append(chatInputDiv)
          chatBoxesContainer.appendChild(chatDiv)
          chatQueue.push(chatDiv)

          hideChat.addEventListener("click", () => {
            chatDiv.style.display = "none"
          })

          const userName = { type: "recipientSelect", info: { name: user.Name } };

          // console.log(userName)
          socket.send(JSON.stringify(userName))
        })
      }

      // Update textContent based on user's online status
      if (user.LoggedInStatus === 0) {
        userDiv.textContent = `🔴 ${user.Name.charAt(0).toUpperCase() + user.Name.slice(1)
          }`
      } else {
        userDiv.textContent = `🟢 ${user.Name.charAt(0).toUpperCase() + user.Name.slice(1)
          }`
      }
    })
  }
}


const renderChat = (obj, size = 1) => {
  const chatBox = document.getElementsByClassName("chat-content")[0];
  const recipientBox = document.getElementsByClassName("chat-title")[0];
  const recipientName = recipientBox.innerText.split('\n')[0].toLowerCase();

  //delete everything within the chatBox
  chatBox.innerHTML = ""

  let totalChatSize = (obj.data.Content).length

  //re-render the most recent  
  for (let index = totalChatSize - size * 10; index <= totalChatSize - 1; index++) {
    let value = obj.data.Content[index]

    const text = document.createElement("div");
    text.innerText = value.time + ": " + value.content;

    if (recipientName === value.receiver) {
      text.classList.add("sent")
    } else {
      text.classList.add("received")
    }
    chatBox.appendChild(text);

  }

  const sentElements = document.getElementsByClassName("sent");
  for (let i = 0; i < sentElements.length; i++) {
    const sentElement = sentElements[i];
    sentElement.style.backgroundColor = "blue";
    sentElement.style.flexDirection = "row-reverse";

  }
}


