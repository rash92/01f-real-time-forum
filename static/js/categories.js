const renderCategory = (categoryName) => {
  fetch(`/categories/${categoryName}`)
    .then(function (response) {
      if (response.ok) {
        return response.json();
      } else {
        throw new Error("Error: " + response.status);
      }
    })
    .then(function (jdata) {
      // Process the JSON data received from Go
      addCategoryInnerHTML(jdata)
      renderNavbar()
    })
    .catch(function (error) {
      console.log(error);
    });
}

function addCategoryInnerHTML(data) {

  document.getElementById("online-users").style.visibility = "visible"

  let html = `
  <form action="/categories/${data.SubName}" method="post">
    <div class="container-filter">
      <div class="filter">
        <p>
          <i class="fa-solid fa-angles-right"></i>
          ${data.SubName}
        </p>
      </div>
      <div class="filter">
        <button name="filter" value="newest">
          <i class="fa-brands fa-hotjar"></i>
          Hot
        </button>
      </div>
      <div class="filter">
        <button name="filter" value="oldest">
          <i class="fa-solid fa-clock"></i>
          Oldest
        </button>
      </div>
    </div>
  </form>
`;

  data.ListOfData.forEach(item => {
    html += `
    <div class="container-post">
      <div class="post">
        <div class="tags-interactions">
          <div class="post-tags">
            `;

    item.Tags.forEach(tag => {
      html += `
      <a class="tag-link" href="/categories/${tag.TagName}">${tag.TagName}</a>
      `;
    });

    html += `
          </div>
          `;

    if (
      data.UserInfo.Name === item.OwnerId ||
      data.UserInfo.Permission === "admin" ||
      data.UserInfo.Permission === "moderator"
    ) {
      html += `
      <div class="post-interactions">
          <div class="tooltip">
            <button class="delete-btn" name="deletepost" value="${item.UUID}">
              <i class="fa-solid fa-trash"></i><span class="tooltiptext">Delete</span>
            </button>
          </div>
        `;

      if (data.UserInfo.Name === item.OwnerId) {
        html += `
        <form action="/submitpost" method="post">
          <div class="tooltip">
            <button type="submit" class="delete-btn tooltip" id="editbutton" name="editpost" value="${item.UUID}">
              <i class="fa-solid fa-pen-to-square"></i><span class="tooltiptext">Edit</span>
            </button>
          </div>
        </form>
        `;
      }

      html += `
      </div>
      `;
    }

    html += `
        </div>
        <div class="post-info">
          <a href="/posts/${item.UUID}">
            <div class="post-author-time">
              <p>Posted by ${item.OwnerId} on ${item.FormattedTime}</p>
            </div>
            <div class="post-title">
              <h2>${item.Title}</h2>
            </div>
            `;

    if (item.ImageName !== "") {
      html += `
      <div class="post-image">
        <img src="${item.ImageName}" alt="user submitted image" />
      </div>
      `;
    }

    html += `
            <div class="post-content">
              <div class="post-content-text">${item.Content}</div>
            </div>
            <div class="post-comments-likes-dislikes">
              <div>
                <p>
                  <i class="fa-regular fa-comment"></i>
                  ${item.NumOfComments}
                </p>
              </div>
            </div>
          </a>
        </div>
        <form action="/categories/${data.SubName}" method="post" target="_self">
          <div class="post-likes-dislikes">
            `;

    if (data.UserInfo.Name !== "") {
      html += `
      <button class="like" name="like" value="${item.UUID}">
        <i class="fa-solid fa-thumbs-up"></i>
      </button>
      `;
    } else {
      html += `
      <p><i class="fa-solid fa-thumbs-up"></i></p>
      `;
    }

    html += `
            <p>${item.Likes}</p>
            `;

    if (data.UserInfo.Name !== "") {
      html += `
      <button class="dislike" name="dislike" value="${item.UUID}">
        <i class="fa-solid fa-thumbs-down"></i>
      </button>
      `;
    } else {
      html += `
      <p><i class="fa-solid fa-thumbs-down"></i></p>
      `;
    }

    html += `
            <p>${item.Dislikes}</p>
          </div>
        </form>
      </div>
    </div>
  `;
  });

  document.getElementById('container').innerHTML = html;

  // Attach event listeners to each delete button
  const deleteButtons = document.querySelectorAll('#delete-post-button');
  deleteButtons.forEach((button) => {
    button.addEventListener('click', (event) => {
      event.preventDefault();

      const id = button.value;
      const formData = { editPost: id };

      deletePost(formData);
    });
  });


  const deletePost = async (formData) => {
    try {
      const response = await fetch(`/deletepost`, {
        method: 'POST',
        body: JSON.stringify(formData),
        headers: {
          'Content-Type': 'application/json'
        }
      });

      if (response.ok) {
        // Authentication successful
        renderForum()
      } else {
        // Handle error response
        console.error('Delete Post failed.');
      }
    } catch (error) {
      // Handle network or other errors
      console.error('Error occurred:', error);
    }
  };

}