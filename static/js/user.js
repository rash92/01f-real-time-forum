const renderUserPage = () => {
    fetch("/user")
        .then(function (response) {
            if (response.ok) {
                return response.json();
            } else {
                throw new Error("Error: " + response.status);
            }
        })
        .then(function (data) {
            // Process the JSON data received from Go
            console.log(data);
            addUserInnerHTML(data)
            renderNavbar()
        })
        .catch(function (error) {
            console.log(error);
        });
}

function addUserInnerHTML(data) {

    const html = `
      <div class="container-filter">
          <div class="filter">
              <button id="posts" onclick="hideShow(this)" value="posts"><i class="fa-solid fa-signs-post"></i> Your
                  posts</button>
          </div>
          <div class="filter">
              <button id="comments" onclick="hideShow(this)" value="comments"><i class="fa-solid fa-message"></i> Your
                  comments</button>
          </div>
          <div class="filter">
              <button id="likes" onclick="hideShow(this)" value="liked-posts"><i class="fa-solid fa-thumbs-up"></i> Your liked
                  posts</button>
          </div>
          <div class="filter">
              <button id="Dislikes" onclick="hideShow(this)" value="disliked-posts"><i class="fa-solid fa-thumbs-down"></i>
                  Your disliked
                  posts</button>
          </div>
      </div>
      
      <div id="user-posts" style="display: block;">
          ${data.UserPosts !== null ? data.UserPosts.map(post => `
          <div class="container-post">
              <div class="post">
                  <div class="tags-interactions">
                      <div class="post-tags">
                          ${post.Tags !== null ? post.Tags.map(tag => `
                          <a class="tag-link" href="/categories/${tag.TagName}">${tag.TagName}</a>
                          `).join('') : ''}
                      </div>
                      ${((data.UserInfo.Name === post.OwnerId) || (data.UserInfo.Permission === "admin") || (data.UserInfo.Permission === "moderator")) ? `
                      <div class="post-interactions">
                              <div class="tooltip">
                                  <button id="delete-post-button" class="delete-btn" name="deletepost" value="${post.UUID}">
                                      <i class="fa-solid fa-trash"></i><span class="tooltiptext">Delete</span>
                                  </button>
                              </div>
                          ${data.UserInfo.Name === post.OwnerId ? `
                              <div class="tooltip">
                                  <button class="delete-btn tooltip" id="editbutton" name="editpost"
                                      value="${post.UUID}">
                                      <i class="fa-solid fa-pen-to-square"></i><span class="tooltiptext">Edit</span>
                                  </button>
                              </div>
                          ` : ''}
                      </div>
                      ` : ''}
                  </div>
                  <div class="post-info">
                      <a onclick="renderPostPage('${post.UUID}')>
                          <div class="post-author-time">
                              <p>Posted by ${post.OwnerId} on ${post.FormattedTime}</p>
                          </div>
                          <div class="post-title">
                              <h2>${post.Title}</h2>
                          </div>
                          ${post.ImageName ? `
                          <div class="post-image">
                              <img src="${post.ImageName}" alt="user submitted image" />
                          </div>
                          ` : ''}
                          <div class="post-content">
                              <div class="post-content-text">${post.Content}</div>
                          </div>
                          <div class="post-like-comment-count">
                              <p>${post.Likes} likes | ${post.NumOfComments} comments</p>
                          </div>
                      </a>
                  </div>
              </div>
          </div>
          `).join('') : ''}
      </div>
      
      <div id="user-comments" style="display: none;">
          ${data.UserComments !== null ? data.UserComments.map(comment => `
          <div class="container-comment">
              <div class="comment">
                  <div class="comment-info">
                      <a onclick="renderPostPage('${comment.PostUUID}')>
                          <div class="comment-author-time">
                              <p>Posted by ${comment.Author} on ${comment.FormattedTime}</p>
                          </div>
                          <div class="comment-content">
                              <div class="comment-content-text">${comment.Content}</div>
                          </div>
                          <div class="comment-like-count">
                              <p>${comment.Likes} likes</p>
                          </div>
                      </a>
                  </div>
              </div>
          </div>
          `).join('') : ''}
      </div>
      
      <div id="liked-user-posts" style="display: none;">
          ${data.LikedUserPosts !== null ? data.LikedUserPosts.map(likedPost => `
          <div class="container-post">
              <div class="post">
                  <div class="tags-interactions">
                      <div class="post-tags">
                          ${likedPost.Tags !== null ? likedPost.Tags.map(tag => `
                          <a class="tag-link" href="/categories/${tag.TagName}">${tag.TagName}</a>
                          `).join('') : ''}
                      </div>
                      ${((data.UserInfo.Name === likedPost.OwnerId) || (data.UserInfo.Permission === "admin") || (data.UserInfo.Permission === "moderator")) ? `
                      <div class="post-interactions">
                              <div class="tooltip">
                                  <button id="delete-post-button" class="delete-btn" name="deletepost" value="${likedPost.UUID}">
                                      <i class="fa-solid fa-trash"></i><span class="tooltiptext">Delete</span>
                                  </button>
                              </div>
                          ${data.UserInfo.Name === likedPost.OwnerId ? `
                              <div class="tooltip">
                                  <button class="delete-btn tooltip" id="editbutton" name="editpost"
                                      value="${likedPost.UUID}">
                                      <i class="fa-solid fa-pen-to-square"></i><span class="tooltiptext">Edit</span>
                                  </button>
                              </div>
                          ` : ''}
                      </div>
                      ` : ''}
                  </div>
                  <div class="post-info">
                      <a href="/posts/${likedPost.UUID}">
                          <div class="post-author-time">
                              <p>Posted by ${likedPost.OwnerId} on ${likedPost.FormattedTime}</p>
                          </div>
                          <div class="post-title">
                              <h2>${likedPost.Title}</h2>
                          </div>
                          ${likedPost.ImageName ? `
                          <div class="post-image">
                              <img src="${likedPost.ImageName}" alt="user submitted image" />
                          </div>
                          ` : ''}
                          <div class="post-content">
                              <div class="post-content-text">${likedPost.Content}</div>
                          </div>
                          <div class="post-like-comment-count">
                              <p>${likedPost.Likes} likes | ${likedPost.NumOfComments} comments</p>
                          </div>
                      </a>
                  </div>
              </div>
          </div>
          `).join('') : ''}
      </div>
      
      <div id="disliked-user-posts" style="display: none;">
          ${data.DislikedUserPosts !== null ? data.DislikedUserPosts.map(dislikedPost => `
          <div class="container-post">
              <div class="post">
                  <div class="tags-interactions">
                      <div class="post-tags">
                          ${dislikedPost.Tags !== null ? dislikedPost.Tags.map(tag => `
                          <a class="tag-link" href="/categories/${tag.TagName}">${tag.TagName}</a>
                          `).join('') : ''}
                      </div>
                      ${((data.UserInfo.Name === dislikedPost.OwnerId) || (data.UserInfo.Permission === "admin") || (data.UserInfo.Permission === "moderator")) ? `
                      <div class="post-interactions">
                              <div class="tooltip">
                                  <button id="delete-post-button" class="delete-btn" name="deletepost" value="${dislikedPost.UUID}">
                                      <i class="fa-solid fa-trash"></i><span class="tooltiptext">Delete</span>
                                  </button>
                              </div>
                          ${data.UserInfo.Name === dislikedPost.OwnerId ? `
                              <div class="tooltip">
                                  <button class="delete-btn tooltip" id="editbutton" name="editpost"
                                      value="${dislikedPost.UUID}">
                                      <i class="fa-solid fa-pen-to-square"></i><span class="tooltiptext">Edit</span>
                                  </button>
                              </div>
                          ` : ''}
                      </div>
                      ` : ''}
                  </div>
                  <div class="post-info">
                      <a href="/posts/${dislikedPost.UUID}">
                          <div class="post-author-time">
                              <p>Posted by ${dislikedPost.OwnerId} on ${dislikedPost.FormattedTime}</p>
                          </div>
                          <div class="post-title">
                              <h2>${dislikedPost.Title}</h2>
                          </div>
                          ${dislikedPost.ImageName ? `
                          <div class="post-image">
                              <img src="${dislikedPost.ImageName}" alt="user submitted image" />
                          </div>
                          ` : ''}
                          <div class="post-content">
                              <div class="post-content-text">${dislikedPost.Content}</div>
                          </div>
                          <div class="post-like-comment-count">
                              <p>${dislikedPost.Likes} likes | ${dislikedPost.NumOfComments} comments</p>
                          </div>
                      </a>
                  </div>
              </div>
          </div>
          `).join('') : ''}
      </div>
      `;

    document.getElementById("container").innerHTML = html;

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
                renderUserPage()
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

}

