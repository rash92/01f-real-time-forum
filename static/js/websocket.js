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
				if (isTyping) {
					typingProgressDiv.innerText = sender + " is typing..."
				} else {
					typingProgressDiv.innerText = ""
				}
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
									// recipient: user.Name,
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
						socket.send(JSON.stringify(messageToSend))
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

					const userName = {
						type: "recipientSelect",
						info: { name: user.Name },
					}
					// console.log(userName)
					socket.send(JSON.stringify(userName))
				})
			}

			// Update textContent based on user's online status
			if (user.LoggedInStatus === 0) {
				userDiv.textContent = `🔴 ${
					user.Name.charAt(0).toUpperCase() + user.Name.slice(1)
				}`
			} else {
				userDiv.textContent = `🟢 ${
					user.Name.charAt(0).toUpperCase() + user.Name.slice(1)
				}`
			}
		})
	}
}
