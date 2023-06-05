
var currentUrl = window.location.href;
console.log(currentUrl.split("8080")[1]);



const renderNavbar = () => {

  fetch("/forum")
    .then(function (response) {
      if (response.ok) {
        return response.json();
      } else {
        throw new Error("Error: " + response.status);
      }
    })
    .then(function (data) {
      // Process the JSON data received from Go
      addNavBarHTML(data)
    })
    .catch(function (error) {
      console.log(error);
    });
};

function addNavBarHTML(data) {
  let html = `
      <nav class="navbar">
        <div class="logo"><a accesskey="h" href="/">F<i class="fa-regular fa-comment"></i>rum</a></div>
        <div class="searchbar">
          <input type="text" placeholder="Search tags" name="search" onkeyup="searchTags()" id="search" autocomplete="off">
          <div id="list">
    `;

  if (data.TagsList !== null && data.TagsList !== undefined) {
    data.TagsList.forEach((tag) => {
      html += `
            <ul>
              <li class="listItem"><a href="/categories/${tag.TagName}">${tag.TagName}</a></li>
            </ul>
          `;
    });
  }

  html += `
          </div>
          <button type="submit"></button>
        </div>
    `;

  if (!data || data.UserInfo.Name === '') {
    html += `
        <div class="menu">
        <li><a accesskey="a" onclick="renderLoginForm('${encodeURIComponent(JSON.stringify(data))}', false)">Sign in</a></li>
        <p>|</p>
        <li><a accesskey="a" onclick="renderLoginForm('${encodeURIComponent(JSON.stringify(data))}', false)">Sign up</a></li>
        </div>
      `;
  } else {
    html += `
        <div class="menu">
          <div class="create">
            <a onclick="renderSubmitPost()"><i class="fa-solid fa-plus"></i> Create</a>
          </div>
          <li><a onclick="renderUserPage()"><i class="fa-solid fa-user"></i> ${data.UserInfo.Name}</a></li>
          <div class="dropdown">
            <button class="dropbtn">
              <i class="fa-solid fa-bell"></i>
              ${data.UserInfo.Notifications ? '<span class="badge"><i class="fa-solid fa-circle"></i></span>' : ''}
            </button>
            <div class="dropdown-content">
              <div class="notifications-interactions">
      `;

    if (data.UserInfo.Notifications !== null) {
      data.UserInfo.Notifications.forEach((notification) => {
        html += `
           <div class="tags-interactions">
             <a class="dropdown-content-links" onclick="renderPostPage('${notification.PostId}')">
               <div>${notification.Sender} ${notification.Statement}</div>
             </a>
             <div class="tooltip">
                 <button id="delete-notification" class="delete-btn" name="delete notification" value="${notification.UUID}">
                   <i class="fa-solid fa-trash"></i><span class="tooltiptext">Delete</span>
                 </button>
             </div>
           </div>
         `;
      });
    }

    html += `
              </div>
            </div>
          </div>
          <div class="dropdown">
            <button class="dropbtn">
              <i class="fa-solid fa-chevron-down"></i>
            </button>
            <div class="dropdown-content">
              <a class="dropdown-content-links" onclick="renderUserPage()">
                <div><i class="fa-solid fa-user"></i> Profile</div>
              </a>
      `;

    if (data.UserInfo.Permission === 'user') {
      html += `
            <button class="dropdown-content-links" name="request to become moderator" value=${data.UserInfo.UUID} onclick="renderUserPage()">
              <div><i class="fa-solid fa-gavel"></i> Become a moderator?</div>
            </button>
        `;
    }

    if (data.UserInfo.Permission === 'admin') {
      html += `
          <a class="dropdown-content-links" onclick="renderAdminPage()">
            <div><i class="fa-solid fa-gear"></i> Admin tools</div>
          </a>
        `;
    }

    html += `
              <a class="dropdown-content-links" href="/logout">
                <div><i class="fa-solid fa-right-from-bracket"></i> Sign out</div>
              </a>
            </div>
          </div>
        </div>
      `;
  }

  html += `</nav>`;

  document.getElementById('navbar').innerHTML = html;

  const deleteNotificationButton = document.querySelectorAll("#delete-notification")
  deleteNotificationButton.forEach(notification => {
    notification.addEventListener('click', (event) => {
      event.preventDefault();

      const notificationid = notification.value
      const formData = {
        deleteNotification: notificationid
      }
      deleteNotification(formData)
    })
  })

  const deleteNotification = async (formData) => {
    try {
      const response = await fetch(`/notification`, {
          method: 'POST',
          body: JSON.stringify(formData),
          headers: {
              'Content-Type': 'application/json'
          }
      });

      if (response.ok) {
          // Authentication successful
          renderNavbar()
      } else {
          // Handle error response
          console.error('Delete Post failed.');
      }
  } catch (error) {
      // Handle network or other errors
      console.error('Error occurred:', error);
  }
  }
}