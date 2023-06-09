const renderLoginForm = (encodedData, attempted) => {
	const decodedData = decodeURIComponent(encodedData)
	const data = JSON.parse(decodedData)
	let html = `
      <div class="container">
        <div class="input-wrapper">
          <div class="input-header">
            <div>
              <h1>Log In</h1>
            </div>
            <div>
              <h3>Sign in with</h3>
            </div>
          </div>
          <form id="login-form" method="POST">
            <div class="login-input">
              <label for="user_name"><i class="fa-solid fa-user"></i></label>
              <input type="text" placeholder="Username or Email" name="user_name" autocomplete="off" maxlength="20" autofocus />
            </div>

    `

	if (data.UserInfo.isLoggedIn) {
		html += `
        <div class="login-input-incorrect">
          <p class="login-incorrect-password">You are already logged in!</p>
        </div>
      `
	}

	html += `
            <div class="login-input">
              <label for="password"><i class="fa-solid fa-lock"></i></label>
              <input type="password" placeholder="Password" name="password" autocomplete="off" required/>
            </div>
    `

	if (!data.isCorrect && attempted) {
		html += `
        <div class="login-input-incorrect">
          <p class="login-incorrect-password">Incorrect Email, Password or Username! Try Again.</p>
        </div>
      `
	}

	html += `
            <div class="login-input">
              <button type="submit">Log In</button>
            </div>
          </form>
          <p>
            Not Registered?
            <a onclick="renderRegisterForm('${encodeURIComponent(
							JSON.stringify(data)
						)}')">Click Here</a>
          </p>
        </div>
      </div>
    `

	document.getElementById("container").innerHTML = html

	const authenticateUser = async (formData) => {
		try {
			const response = await fetch("/authenticate", {
				method: "POST",
				body: JSON.stringify(formData),
				headers: {
					"Content-Type": "application/json",
				},
			})

			if (response.ok) {
				// Authentication successful
				const data = await response.json()

				if (data.IsCorrect) {
					renderForum()
					startWebSocket()
				} else {
					// Call the rendering function with the updated variables
					renderLoginForm(encodeURIComponent(JSON.stringify(data)), true)
				}
			} else {
				// Handle error response
				console.error("Authentication failed.")
			}
		} catch (error) {
			// Handle network or other errors
			console.error("Error occurred:", error)
		}
	}

	// Submit form event handler
	document
		.querySelector("#login-form")
		.addEventListener("submit", async (event) => {
			event.preventDefault() // Prevent form submission

			const formData = {
				user_name: document.querySelector('input[name="user_name"]').value,
				password: document.querySelector('input[name="password"]').value,
			}

			authenticateUser(formData)
		})
}

fetch("/forum")
	.then(function (response) {
		if (response.ok) {
			console.log("response ok, json is: ", response.json)
			return response.json()
		} else {
			throw new Error("Error: " + response.status)
		}
	})
	.then(function (jdata) {
		console.log(jdata.UserInfo.IsLoggedIn)
		if (jdata.UserInfo.IsLoggedIn) {
			renderForum()
			startWebSocket()
		} else {
			console.log("here")
			renderLoginForm(encodeURIComponent(JSON.stringify(jdata)), false)
		}
	})
	.catch(function (error) {
		console.log(error)
	})
