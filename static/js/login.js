const renderLoginForm = (encodedData, attempted) => {
    const decodedData = decodeURIComponent(encodedData);
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
          <div class="alternative-sign-up">
            <div class="social">
              <form action="/google/login" method="POST">
                <button class="social-button" type="submit">
                  <i class="fa-brands fa-google"></i>
                </button>
              </form>
            </div>
            <div class="social">
              <form action="/github/login" method="POST">
                <button class="social-button" type="submit">
                  <i class="fa-brands fa-github"></i>
                </button>
              </form>
            </div>
            <!-- <div class="social">
              <form action="/facebook/login" method="POST">
                <button class="social-button" type="submit">
                  <i class="fa-brands fa-facebook"></i>
                </button>
              </form>
            </div> -->
          </div>
          <form id="login-form" method="POST">
            <div class="login-input">
              <label for="user_name"><i class="fa-solid fa-user"></i></label>
              <input type="text" placeholder="Username" name="user_name" autocomplete="off" maxlength="20" autofocus required/>
            </div>
    `;

    if (data.UserInfo.isLoggedIn) {
        html += `
        <div class="login-input-incorrect">
          <p class="login-incorrect-password">You are already logged in!</p>
        </div>
      `;
    }

    html += `
            <div class="login-input">
              <label for="password"><i class="fa-solid fa-lock"></i></label>
              <input type="password" placeholder="Password" name="password" autocomplete="off" required/>
            </div>
    `;

    if (!data.isCorrect && attempted) {
        html += `
        <div class="login-input-incorrect">
          <p class="login-incorrect-password">Incorrect Password or Username! Try Again.</p>
        </div>
      `;
    }

    html += `
            <div class="login-input">
              <button type="submit">Log In</button>
            </div>
          </form>
          <p>
            Not Registered?
            <a onclick="renderRegisterForm('${encodeURIComponent(JSON.stringify(data))}')">Click Here</a>
          </p>
        </div>
      </div>
    `;

    document.getElementById('container').innerHTML = html;

    const authenticateUser = async (formData) => {
        try {
            const response = await fetch('/authenticate', {
                method: 'POST',
                body: JSON.stringify(formData),
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (response.ok) {
                // Authentication successful
                const data = await response.json();

                if (data.IsCorrect) {
                    renderForum()
                } else {
                    // Call the rendering function with the updated variables
                    renderLoginForm(encodeURIComponent(JSON.stringify(data)), true);
                }
            } else {
                // Handle error response
                console.error('Authentication failed.');
            }
        } catch (error) {
            // Handle network or other errors
            console.error('Error occurred:', error);
        }
    };

    // Submit form event handler
    document.querySelector('#login-form').addEventListener('submit', async (event) => {
        event.preventDefault(); // Prevent form submission

        const formData = {
            user_name: document.querySelector('input[name="user_name"]').value,
            password: document.querySelector('input[name="password"]').value
        };

        authenticateUser(formData);
    });

};
