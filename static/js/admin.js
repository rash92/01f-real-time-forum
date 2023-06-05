const renderAdminPage = (startingPage = null) => {
  fetch("/admin")
    .then(function (response) {
      if (response.ok) {
        return response.json();
      } else {
        throw new Error("Error: " + response.status);
      }
    })
    .then(function (jdata) {
      // Process the JSON data received from Go
      addAdminPageInnerHTML(jdata, startingPage)
      renderNavbar()
    })
    .catch(function (error) {
      console.log(error);
    });
}

function addAdminPageInnerHTML(data, startingPage) {
  let html = `
  <div class="container-filter">
    <div class="filter">
        <button id="mod" onclick="hideShowAdmin(this)" value="mod"><i class="fa-solid fa-hourglass-start"></i>
            Pending
            Mod Requests</button>
    </div>
    <div class="filter">
        <button id="reports" onclick="hideShowAdmin(this)" value="reports"><i class="fa-solid fa-flag"></i>
            Pending
            Reported Posts</button>
    </div>
    <div class="filter">
        <button id="users" onclick="hideShowAdmin(this)" value="users"><i class="fa-solid fa-users"></i>
            Users</button>
    </div>
    <div class="filter">
        <button id="tags" onclick="hideShowAdmin(this)" value="tags"><i class="fa-solid fa-tags"></i>
            Tags</button>
    </div>
  </div>
`;

  html += `
  <div id="mod-requests" style="display: block;">
    `;

  if (data.AdminRequests !== undefined && data.AdminRequests !== null) {
    data.AdminRequests.forEach(item => {
      if (item.Description === "this user is asking to become a moderator") {
        html += `
          <div class="container-post">
            <div class="post">
              ${item.RequestFromId} ${item.RequestFromName}
              <p>${item.Description}</p>
                <button class="delete-btn" id="delete-request" value="${item.UUID}">
                  delete request
                </button>
                <button class="delete-btn" id="set-to-moderator" value="${item.RequestFromId}">
                  Set to moderator
                </button>
            </div>
          </div>
        `;
      }
    });
  }

  html += `
  </div>
`;

  html += `
  <div id="reported-posts" style="display: none">
    `;

  if (data.AdminRequests !== undefined && data.AdminRequests !== null) {
    data.AdminRequests.forEach(item => {
      if (item.ReportedPostId !== "") {
        html += `
          <div class="container-post">
            <div class="post">
              the post ${item.ReportedPostId} has been reported by:
              ${item.RequestFromId} ${item.RequestFromName}
              <p>${item.Description}</p>
              <a href="/posts/${item.ReportedPostId}">Go to post</a>
                <input type="text" id="response-message" placeholder="response message">
                <button type="submit" id="acknowledge-report" value="${item.UUID}">
                  acknowledge report
                </button>
                <button class="delete-btn" id="delete-request" value="${item.UUID}">
                  delete request
                </button>
            </div>
          </div>
        `;
      }
    });
  }

  html += `
  </div>
`;

  html += `
  <div id="user-admin" style="display: none;">
    `;

  if (data.AllUsers !== undefined && data.AllUsers !== null) {
    data.AllUsers.forEach(item => {
      if (item.Name !== 'admin') {
        html += `
        <div class="container-post">
          <div class="post">
            <p>User name: ${item.Name}</p>
            <p>Permission: ${item.Permission}</p>
              <button class="delete-btn" id="set-to-user" value="${item.UUID}">set to user</button>
              <button class="delete-btn" id="set-to-moderator" value="${item.UUID}">
                Set to moderator
              </button>
              <button class="delete-btn" id="set-to-admin" value="${item.UUID}">set to admin</button>
              <button class="delete-btn" id="delete-user" value="${item.UUID}">delete user</button>
          </div>
        </div>
      `;
      }
    });
  }

  html += `
  </div>
`;

  html += `
  <div id="tags-admin" style="display: none;">
    <div class="container-post">
        <input type="text" id="tags-text" placeholder="Tags (#General etc.)">
        <button id="tag-submit">Submit</button>
    </div>
    `;

  if (data.TagsList !== null && data.TagsList !== undefined) {
    console.log(data.TagsList);
    data.TagsList.forEach(item => {
      html += `
        <div class="container-post">
          <div class="post">
            <p>Tag Name: ${item.TagName}</p>
              <button class="delete-btn" id="delete-tag" value="${item.UUID}">delete tag</button>
              <button class="delete-btn admin-all-tag-delete" id="delete-all-tag" value="${item.TagName}">delete all posts with this tag</button>
          </div>
        </div>
      `;
    });
  }

  html += `
  </div>
`;

  document.getElementById('container').innerHTML = html;
  if (startingPage) {
    hideShowAdmin(document.getElementById(startingPage))
  }

  const admin = async (formData, startingPage) => {
    try {
      const response = await fetch(`/admin`, {
        method: 'POST',
        body: JSON.stringify(formData),
        headers: {
          'Content-Type': 'application/json'
        }
      });

      if (response.ok) {
        renderAdminPage(startingPage)
      } else {
        // Handle error response
        console.error('Admin Action failed.');
      }
    } catch (error) {
      // Handle network or other errors
      console.error('Error occurred:', error);
    }
  }

  const setUserButton = document.querySelectorAll('#set-to-user')
  setUserButton.forEach((button) => {
    button.addEventListener('click', (event) => {
      event.preventDefault();

      const formData = {
        set_to_user: button.value,
        set_to_mod: "",
        set_to_admin: "",
        delete_user: "",
        tag_create: "",
        delete_tag: "",
        delete_all_tag: "",
        delete_request: "",
        acknowledge_report: "",
        response_message: ""
      };
      

      admin(formData, "users")
    })
  })

  const setModButton = document.querySelectorAll('#set-to-moderator')
  setModButton.forEach((button) => {
    button.addEventListener('click', (event) => {
      event.preventDefault();

      const formData = {
        set_to_user: "",
        set_to_mod: button.value,
        set_to_admin: "",
        delete_user: "",
        tag_create: "",
        delete_tag: "",
        delete_all_tag: "",
        delete_request: "",
        acknowledge_report: "",
        response_message: ""
      };
      

      admin(formData, "users")
    })
  })

  const setAdminButton = document.querySelectorAll('#set-to-admin')
  setAdminButton.forEach((button) => {
    button.addEventListener('click', (event) => {
      event.preventDefault();

      const formData = {
        set_to_user: "",
        set_to_mod: "",
        set_to_admin: button.value,
        delete_user: "",
        tag_create: "",
        delete_tag: "",
        delete_all_tag: "",
        delete_request: "",
        acknowledge_report: "",
        response_message: ""
      };
      

      admin(formData, "users")
    })
  })

  const deleteUser = document.querySelectorAll('#delete-user')
  deleteUser.forEach((button) => {
    button.addEventListener('click', (event) => {
      event.preventDefault();

      const formData = {
        set_to_user: "",
        set_to_mod: "",
        set_to_admin: "",
        delete_user: button.value,
        tag_create: "",
        delete_tag: "",
        delete_all_tag: "",
        delete_request: "",
        acknowledge_report: "",
        response_message: ""
      };
      

      admin(formData, "users")
    })
  })

  const tagCreate = document.querySelectorAll('#tag-submit')
  tagCreate.forEach((button) => {
    button.addEventListener('click', (event) => {
      event.preventDefault();

      const formData = {
        set_to_user: "",
        set_to_mod: "",
        set_to_admin: "",
        delete_user: "",
        tag_create: document.getElementById("tags-text").value,
        delete_tag: "",
        delete_all_tag: "",
        delete_request: "",
        acknowledge_report: "",
        response_message: ""
      };
      
      admin(formData, "tags")
    })
  })

  const deleteTag = document.querySelectorAll('#delete-tag')
  deleteTag.forEach((button) => {
    button.addEventListener('click', (event) => {
      event.preventDefault();

      const formData = {
        set_to_user: "",
        set_to_mod: "",
        set_to_admin: "",
        delete_user: "",
        tag_create: "",
        delete_tag: button.value,
        delete_all_tag: "",
        delete_request: "",
        acknowledge_report: "",
        response_message: ""
      };
      

      admin(formData, "tags")
    })
  })

  const deleteAllTag = document.querySelectorAll('#delete-all-tag')
  deleteAllTag.forEach((button) => {
    button.addEventListener('click', (event) => {
      event.preventDefault();

      const formData = {
        set_to_user: "",
        set_to_mod: "",
        set_to_admin: "",
        delete_user: "",
        tag_create: "",
        delete_tag: "",
        delete_all_tag: button.value,
        delete_request: "",
        acknowledge_report: "",
        response_message: ""
      };
      

      admin(formData, "tags")
    })
  })

  const deleteRequest = document.querySelectorAll('#delete-request')
  deleteRequest.forEach((button) => {
    button.addEventListener('click', (event) => {
      event.preventDefault();

      const formData = {
        set_to_user: "",
        set_to_mod: "",
        set_to_admin: "",
        delete_user: "",
        tag_create: "",
        delete_tag: "",
        delete_all_tag: "",
        delete_request: button.value,
        acknowledge_report: "",
        response_message: ""
      };
      

      admin(formData, "mod")
    })
  })

  const acknowledgeReport = document.querySelectorAll('#acknowledge-report')
  acknowledgeReport.forEach((button) => {
    button.addEventListener('click', (event) => {
      event.preventDefault();

      const formData = {
        set_to_user: "",
        set_to_mod: "",
        set_to_admin: "",
        delete_user: "",
        tag_create: "",
        delete_tag: "",
        delete_all_tag: "",
        delete_request: "",
        acknowledge_report: button.value,
        response_message: ""
      };
      

      admin(formData, "reports")
    })
  })

  const responseMessage = document.querySelectorAll('#response-message')
  responseMessage.forEach((button) => {
    button.addEventListener('click', (event) => {
      event.preventDefault();

      const formData = {
        set_to_user: "",
        set_to_mod: "",
        set_to_admin: "",
        delete_user: "",
        tag_create: "",
        delete_tag: "",
        delete_all_tag: "",
        delete_request: "",
        acknowledge_report: "",
        response_message: button.value
      };
      

      admin(formData, "reports")
    })
  })

}

