const renderForum = () => {

    fetch("/forum")
        .then(function (response) {
            if (response.ok) {
                return response.json();
            } else {
                throw new Error("Error: " + response.status);
            }
        })
        .then(function (jdata) {
            // Process the JSON data received from Go
            addForumInnerHTML(jdata)
            renderNavbar()
        })
        .catch(function (error) {
            console.log(error);
        });

};

function addForumInnerHTML(data) {
    let html = `
      <form action="/" method="post">
          <div class="container-filter">
              <div class="filter">
                  <button name="filter" value="newest">
                      <i class="fa-brands fa-hotjar"></i>
                      Newest
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

    if (data.ListOfData !== undefined && data.ListOfData !== null) {
        data.ListOfData.forEach((item) => {
            if (item.Tags === null) {
                item.Tags = []
            }
            html += `
            <div class="container-post">
                <div class="post">
                    <div class="tags-interactions">
                        <div class="post-tags">
                            ${item.Tags.map((tag) => `<a class="tag-link" onclick="renderCategory(${tag.TagName})">${tag.TagName}</a>`).join('')}
                        </div>
                        ${(data.UserInfo.Name === item.OwnerId || data.UserInfo.Permission === 'admin' || data.UserInfo.Permission === 'moderator')
                    ? `
                            <div class="post-interactions">
                                    <div class="tooltip">
                                        <button id="post-delete-button" class="delete-btn" name="deletepost" value="${item.UUID}">
                                            <i class="fa-solid fa-trash"></i><span class="tooltiptext">Delete</span>
                                        </button>
                                    </div>
                                ${data.UserInfo.Name === item.OwnerId
                        ? `
                                      <div class="tooltip">
                                          <button class="delete-btn tooltip" id="editbutton" name="editpost" value="${item.UUID}">
                                              <i class="fa-solid fa-pen-to-square"></i><span class="tooltiptext">Edit</span>
                                          </button>
                                      </div>
                                  `
                        : ''
                    }
                            </div>
                            `
                    : ''
                }
                    </div>
                    <div class="post-info">
                        <a onclick="renderPostPage('${item.UUID}')">
                            <div class="post-author-time">
                                <p>Posted by ${item.OwnerId} on ${item.FormattedTime}</p>
                            </div>
                            <div class="post-title">
                                <h2>${item.Title}</h2>
                            </div>
                            ${item.ImageName
                    ? `<div class="post-image"><img src="${item.ImageName}" alt="user submitted image" /></div>`
                    : ''
                }
                            <div class="post-content">
                                <div class="post-content-text">${item.Content}</div>
                            </div>
                            <div class="post-comments-likes-dislikes">
                                <div>
                                    <p><i class="fa-regular fa-comment"></i>${item.NumOfComments}</p>
                                </div>
                            </div>
                        </a>
                    </div>
                        <div class="post-likes-dislikes">
                            ${data.UserInfo.Name
                    ? `
                              <button class="like" name="like" value="${item.UUID}">
                                  <i class="fa-solid fa-thumbs-up"></i>
                              </button>
                              `
                    : `
                              <p><i class="fa-solid fa-thumbs-up"></i></p>
                              `
                }
                            <p>${item.Likes}</p>
                            ${data.UserInfo.Name
                    ? `
                              <button class="dislike" name="dislike" value="${item.UUID}">
                                  <i class="fa-solid fa-thumbs-down"></i>
                              </button>
                              `
                    : `
                              <p><i class="fa-solid fa-thumbs-down"></i></p>
                              `
                }
                            <p>${item.Dislikes}</p>
                        </div>
                </div>
            </div>
          `;
        });
    }

    document.getElementById('container').innerHTML = html;

    const reactPost = async (formData) => {
        console.log(formData);
        try {
            const response = await fetch('/react', {
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
                console.error('React Post failed.');
            }
        } catch (error) {
            // Handle network or other errors
            console.error('Error occurred:', error);
        }
    };

    // Attach event listeners to each like button
    const likeButtons = document.querySelectorAll('.like');
    likeButtons.forEach((button) => {
        button.addEventListener('click', (event) => {
            event.preventDefault();

            const id = button.value;
            const formData = { isComment: false, id, like: true };

            reactPost(formData);
        });
    });

    // Attach event listeners to each dislike button
    const dislikeButtons = document.querySelectorAll('.dislike');
    dislikeButtons.forEach((button) => {
        button.addEventListener('click', (event) => {
            event.preventDefault();

            const id = button.value;
            const formData = { isComment: false, id, dislike: true };

            reactPost(formData);
        });
    });

    // Attach event listeners to each delete button
    const deleteButtons = document.querySelectorAll('#post-delete-button');
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

    const editPost = async (formData) => {
        try {
            const response = await fetch(`/editpost`, {
                method: 'POST',
                body: JSON.stringify(formData),
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (response.ok) {
                const data = await response.json();
                renderEditPost(data.EditPost.UUID)
            } else {
                // Handle error response
                console.error('Edit Post failed.');
            }
        } catch (error) {
            // Handle network or other errors
            console.error('Error occurred:', error);
        }
    };

    const editButtons = document.querySelectorAll('#editbutton');
    // Submit form event handler
    editButtons.forEach((button) => {
        button.addEventListener('click', async (event) => {
            event.preventDefault(); // Prevent form submission

            const formData = {
                editPostID: document.querySelector('button[name="editpost"]').value
            };

            editPost(formData);
        })
    });


};