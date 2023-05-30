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
        console.log(jdata);
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
                            ${item.Tags.map((tag) => `<a class="tag-link" href="/categories/${tag.TagName}">${tag.TagName}</a>`).join('')}
                        </div>
                        ${(data.UserInfo.Name === item.OwnerId || data.UserInfo.Permission === 'admin' || data.UserInfo.Permission === 'moderator')
                    ? `
                            <div class="post-interactions">
                                <form method="POST">
                                    <div class="tooltip">
                                        <button class="delete-btn" name="deletepost" value="${item.UUID}">
                                            <i class="fa-solid fa-trash"></i><span class="tooltiptext">Delete</span>
                                        </button>
                                    </div>
                                </form>
                                ${data.UserInfo.Name === item.OwnerId
                        ? `
                                  <form action="/submitpost" method="post">
                                      <div class="tooltip">
                                          <button type="submit" class="delete-btn tooltip" id="editbutton" name="editpost" value="${item.UUID}">
                                              <i class="fa-solid fa-pen-to-square"></i><span class="tooltiptext">Edit</span>
                                          </button>
                                      </div>
                                  </form>
                                  `
                        : ''
                    }
                            </div>
                            `
                    : ''
                }
                    </div>
                    <div class="post-info">
                        <a href="/posts/${item.UUID}">
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
                    <form action="/forum" method="post" target="_self">
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
                    </form>
                </div>
            </div>
          `;
        });
    }

    document.getElementById('container').innerHTML = html;
}

renderForum()
