const renderAdminPage = () => {
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
            addAdminPageInnerHTML(jdata)
            renderNavbar()
        })
        .catch(function (error) {
            console.log(error);
        });
}

function addAdminPageInnerHTML(data) {
    console.log("admin", data);
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
              <form action="/admin" method="POST">
                <button class="delete-btn" name="delete request" value="${item.UUID}">
                  delete request
                </button>
                <button class="delete-btn" name="set to moderator" value="${item.RequestFromId}">
                  Set to moderator
                </button>
              </form>
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
              <form action="/admin" method="POST">
                <input type="text" name="response message" placeholder="response message">
                <button type="submit" name="acknowledge report" value="${item.UUID}">
                  acknowledge report
                </button>
                <button class="delete-btn" name="delete request" value="${item.UUID}">
                  delete request
                </button>
              </form>
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
            html += `
        <div class="container-post">
          <div class="post">
            <p>User name: ${item.Name}</p>
            <p>Permission: ${item.Permission}</p>
            <form action="/admin" method="POST">
              <button class="delete-btn" name="set to user" value="${item.UUID}">set to user</button>
              <button class="delete-btn" name="set to moderator" value="${item.UUID}">
                Set to moderator
              </button>
              <button class="delete-btn" name="set to admin" value="${item.UUID}">set to admin</button>
              <button class="delete-btn" name="delete user" value="${item.UUID}">delete user</button>
            </form>
          </div>
        </div>
      `;
        });
    }

    html += `
  </div>
`;

    html += `
  <div id="tags-admin" style="display: none;">
    <div class="container-post">
      <form action="/admin" method="post">
        <input type="text" name="tags" placeholder="Tags (#General etc.)">
        <input type="submit">
      </form>
    </div>
    `;

    if (data.TagsList !== null && data.TagsList !== undefined) {
        data.TagsList.forEach(item => {
            html += `
        <div class="container-post">
          <div class="post">
            <p>Tag Name: ${item.TagName}</p>
            <form action="/admin" method="POST">
              <button class="delete-btn" name="delete tag" value="${item.UUID}">delete tag</button>
              <button class="delete-btn admin-all-tag-delete" name="delete all posts with tag" value="${item.TagName}">delete all posts with this tag</button>
            </form>
          </div>
        </div>
      `;
        });
    }

    html += `
  </div>
`;

    document.getElementById('container').innerHTML = html;

}