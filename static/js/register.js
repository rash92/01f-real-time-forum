const renderRegisterForm = (encodedData) => {
    const decodedData = decodeURIComponent(encodedData);
    const data = JSON.parse(decodedData)
    let html = `
      <div class="container">
        <div class="input-wrapper">
          <div>
            <h1>Register</h1>
          </div>
          <div>
            <h3>Sign in with</h3>
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
          <form id="register-form" method="POST">
            <div class="login-input">
              <label for="user_name"><i class="fa-solid fa-user"></i></label>
              <input
                type="text"
                placeholder="Username"
                name="user_name"
                autocomplete="off"
                autofocus
                required
              />
            </div>
            <div class="login-input">
              <label for="email"><i class="fa-solid fa-envelope"></i></label>
              <input
                type="email"
                placeholder="Email"
                name="email"
                autocomplete="off"
                required
              />
            </div>
            <div class="login-input">
              <label for="password"><i class="fa-solid fa-lock"></i></label>
              <input
                type="password"
                placeholder="Password"
                name="password"
                minlength="6"
                maxlength="20"
                autocomplete="off"
                required
              />
            </div>
      `;

    // Check if there's a register error
    if (data.RegisterError !== "" && data.RegisterError !== undefined) {
        html += `
          <div class="login-input-incorrect">
            <p class="login-incorrect-password">
              Sorry, that ${data.RegisterError} is already taken.
            </p>
          </div>
        `;
    }

    html += `
            <div class="login-input">
              <button type="submit">Sign Up</button>
            </div>
          </form>
          <p>
            Registered?
            <a accesskey="a" href="/login">Sign in</a>
          </p>
        </div>
      </div>
    `;

    document.getElementById('container').innerHTML = html;

    // Add submit form event listener
    document.querySelector('#register-form').addEventListener('submit', async (event) => {
        event.preventDefault(); // Prevent form submission

        const formData = {
            user_name: document.querySelector('input[name="user_name"]').value,
            email: document.querySelector('input[name="email"]').value,
            password: document.querySelector('input[name="password"]').value
        };

        registerUser(formData);
    });
};

const registerUser = async (formData) => {
    try {
        const response = await fetch('/register_account', {
            method: 'POST',
            body: JSON.stringify(formData),
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (response.ok) {
            // Registration successful
            const data = await response.json();
            // Do something with the response data if needed
            if (data.RegisterError !== "") {
                console.log(data);
                renderRegisterForm(encodeURIComponent(JSON.stringify(data)))
            } else {
                renderLoginForm(encodeURIComponent(JSON.stringify(data)), false)
            }
        } else {
            // Handle error response
            console.error('Registration failed.');
        }
    } catch (error) {
        // Handle network or other errors
        console.error('Error occurred:', error);
    }
};
