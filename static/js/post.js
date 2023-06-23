const renderPostPage = (postid) => {
    fetch(`/posts/${postid}`)
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
            addPostInnerHTML(data)
            renderNavbar()
        })
        .catch(function (error) {
            console.log(error);
        });
}

function addPostInnerHTML(data) {

    document.getElementById("online-users").style.visibility = "visible"

    const html = `
        <div class="container-post">
            <div class="single-post">
                <div class="tags-interactions">
                    <div class="post-tags">
                        ${data.Post.Tags !== null ? data.Post.Tags.map(tag => `
                        <a class="tag-link" href="/categories/${tag.TagName}">${tag.TagName}</a>
                        `).join('') : ''}
                    </div>
                    ${((data.UserInfo.Name === data.Post.OwnerId) || (data.UserInfo.Permission === "admin") || (data.UserInfo.Permission === "moderator")) ? `
                    <div class="post-interactions">
                            <div class="tooltip">
                                <button id="delete-post-button" class="delete-btn" name="deletepost" value="${data.Post.UUID}">
                                    <i class="fa-solid fa-trash"></i><span class="tooltiptext">Delete</span>
                                </button>
                            </div>
                        ${(data.UserInfo.Name === data.Post.OwnerId) ? `
                            <div class="tooltip">
                                <button type="submit" class="delete-btn tooltip" id="editbutton" name="editpost" value="${data.Post.UUID}" onclick="renderEditPost('${data.Post.UUID}')">
                                    <i class="fa-solid fa-pen-to-square"></i><span class="tooltiptext">Edit</span>
                                </button>
                            </div>
                        ` : ''}
                    </div>
                    ` : ''}

                    ${((data.UserInfo.Permission === "moderator") || (data.UserInfo.Permission === "admin")) ? `
                        <div class="tooltip">
                            <button id="report-post-button" class="report-btn tooltip" name="reportpost" value="${data.Post.UUID}">
                                <i class="fa-solid fa-flag"></i><span class="tooltiptext">Report</span>
                            </button>
                        </div>
                    ` : ''}
                </div>
                <div class="single-post-title">
                    <h2>${data.Post.Title}</h2>
                </div>
                <div class="single-post-author">
                    <p>Posted by ${data.Post.OwnerId} at ${data.Post.FormattedTime}</p>
                </div>
                ${(data.Post.ImageName !== "") ? `
                <div class="post-image">
                    <img src="${data.Post.ImageName}" alt="user submitted image">
                </div>
                ` : ''}
                <div class="single-post-content">
                    <p>${data.Post.Content}</p>
                </div>
                <div class="post-comments-likes-dislikes">
                    <div>
                        <p><i class="fa-regular fa-comment"></i> ${data.NumOfComments}</p>
                    </div>
                        <div class="post-likes-dislikes">
                            ${(data.UserInfo.Name !== "") ? `
                            <button id="like-button" class="like" name="like" value="${data.Post.UUID}"><i class="fa-solid fa-thumbs-up"></i></button>
                            ` : `
                            <p><i class="fa-solid fa-thumbs-up"></i></p>
                            `}
                            <p>${data.Post.Likes}</p>
                            ${(data.UserInfo.Name !== "") ? `
                            <button id="dislike-button" class="dislike" name="dislike" value="${data.Post.UUID}"><i class="fa-solid fa-thumbs-down"></i></button>
                            ` : `
                            <p><i class="fa-solid fa-thumbs-down"></i></p>
                            `}
                            <p>${data.Post.Dislikes}</p>
                        </div>
                </div>
                ${(data.UserInfo.Name !== "") ? `
                    <textarea name="comment" id="comment-text" cols="30" rows="5" placeholder="Comment" autofocus required></textarea>
                    <button id="comment-btn">Comment</button>
                ` : `
                <p>Click <a onclick="renderLoginForm('${encodeURIComponent(JSON.stringify(data))}', false)">here</a> to register to make a comment</p>
                `}
            </div>
            ${data.Comments ? data.Comments.map(comment => `
            <div class="single-post">
                <div class="commentor-name-time">
                    <p>${comment.OwnerName} - ${comment.FormattedTime}</p>
                </div>
                <div class="comment">
                    <p>${comment.Content}</p>
                </div>

                    <div class="post-likes-dislikes">
                        ${(data.UserInfo.Name !== "") ? `
                        <button id="comment-like-button" class="like" name="commentlike" value="${comment.UUID}"><i class="fa-solid fa-thumbs-up"></i></button>
                        ` : `
                        <p><i class="fa-solid fa-thumbs-up"></i></p>
                        `}
                        <p>${comment.Likes}</p>
                        ${(data.UserInfo.Name !== "") ? `
                        <button id="comment-dislike-button" class="dislike" name="commentdislike" value="${comment.UUID}"><i class="fa-solid fa-thumbs-down"></i></button>
                        ` : `
                        <p><i class="fa-solid fa-thumbs-down"></i></p>
                        `}
                        <p>${comment.Dislikes}</p>
                    </div>
                ${((data.UserInfo.Name === comment.OwnerName) || (data.UserInfo.Permission === "admin") || (data.UserInfo.Permission === "moderator")) ? `
                <div class="post-interactions">
                        <div class="tooltip">
                            <button id="comment-delete-button" class="delete-btn" name="deletecomment" value="${comment.UUID}">
                                <i class="fa-solid fa-trash"></i><span class="tooltiptext">Delete</span>
                            </button>
                        </div>
                        ${(data.UserInfo.Name === comment.OwnerName) ? `
                        <div class="tooltip">
                            <button type="button" class="delete-btn tooltip" id="editbutton" name="showedit" value="showComment" onclick="showCommentEdit(this)" style="display: block;">
                                <i class="fa-solid fa-pen-to-square"></i><span class="tooltiptext">Edit</span>
                            </button>
                        </div>
                        ` : ''}
                </div>
                ` : ''}
                ${(data.UserInfo.Name === comment.OwnerName) ? `
                    <textarea name="editcomment" id="commentEditor" cols="30" rows="5" placeholder="Comment" style="display: none;">${comment.Content}</textarea>
                    </textarea>
                    <button class="tooltip" type="submit" name="commentuuid" id="finalEditButton" style="display: none;" value="${comment.UUID}"><i class="fa-solid fa-pen-to-square"></i></button><span class="tooltiptext" id="spanToolTip"> </span>
                ` : ''}
            </div>
            `).join('') : ''}
        </div>
    `;

    document.getElementById("container").innerHTML = html;

    const comment = async (formData) => {
        try {
            const response = await fetch(`/comments/${data.Post.UUID}`, {
                method: 'POST',
                body: JSON.stringify(formData),
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            if (response.ok) {
                // Authentication successful
                renderPostPage(data.Post.UUID)
            } else {
                // Handle error response
                console.error('Comment Post failed.');
            }
        } catch (error) {
            // Handle network or other errors
            console.error('Error occurred:', error);
        }
    }

    const commentButton = document.getElementById("comment-btn");
    commentButton.addEventListener('click', (event) => {
        event.preventDefault();

        const commentText = document.getElementById("comment-text").value
        const formData = {
            commentContent: commentText,
            deleteComment: "",
            editComment: "",
            commentUUID: ""
        }

        comment(formData)
    })

    const deleteCommentButton = document.getElementById("comment-delete-button");
    deleteCommentButton.addEventListener('click', (event) => {
        event.preventDefault();

        const commentid = document.getElementById("comment-delete-button").value
        const formData = {
            commentContent: "",
            deleteComment: commentid,
            editComment: "",
            commentUUID: ""
        }

        comment(formData)
    })

    const editCommentButton = document.getElementById("finalEditButton");
    editCommentButton.addEventListener('click', (event) => {
        event.preventDefault();

        const commentText = document.getElementById("commentEditor").value
        const formData = {
            commentContent: "",
            deleteComment: "",
            editComment: commentText,
            commentUUID: document.getElementById("finalEditButton").value
        }

        comment(formData)
    })

    //add button functionality below
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
                renderPostPage(data.Post.UUID)
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
    const likeButtons = document.querySelectorAll('#like-button');
    likeButtons.forEach((button) => {
        button.addEventListener('click', (event) => {
            event.preventDefault();

            const id = button.value;
            const formData = { isComment: false, id, like: true };

            reactPost(formData);
        });
    });

    // Attach event listeners to each dislike button
    const dislikeButtons = document.querySelectorAll('#dislike-button');
    dislikeButtons.forEach((button) => {
        button.addEventListener('click', (event) => {
            event.preventDefault();

            const id = button.value;
            const formData = { isComment: false, id, dislike: true };

            reactPost(formData);
        });
    });

    // Attach event listeners to each comment like button
    const commentLikeButtons = document.querySelectorAll('#comment-like-button');
    commentLikeButtons.forEach((button) => {
        button.addEventListener('click', (event) => {
            event.preventDefault();

            const id = button.value;
            const formData = { isComment: true, id, like: true };

            reactPost(formData);
        });
    });

    // Attach event listeners to each comment dislike button
    const commentDislikeButtons = document.querySelectorAll('#comment-dislike-button');
    commentDislikeButtons.forEach((button) => {
        button.addEventListener('click', (event) => {
            event.preventDefault();

            const id = button.value;
            const formData = { isComment: true, id, dislike: true };

            reactPost(formData);
        });
    });

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
                renderPostPage(data.Post.UUID)
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