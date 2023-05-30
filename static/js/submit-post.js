const renderSubmitPost = () => {
    fetch("/submitpost")
        .then(function (response) {
            if (response.ok) {
                return response.json();
            } else {
                throw new Error("Error: " + response.status);
            }
        })
        .then(function (jdata) {
            // Process the JSON data received from Go
            addSubmitPostInnerHTML(jdata)
        })
        .catch(function (error) {
            console.log(error);
        });
}

function addSubmitPostInnerHTML(data) {
    let html = `
  <div class="container-post">
    <div class="submission">
        <h2>Create a post</h2>
        `;

    if (data.TagsList === null) {
        data.TagsList = []
    }

    if (data.IsEdit) {
        html += `
        <form id="submit-post-form" method="POST" enctype="multipart/form-data">
            <input type="text" name="submission-title" value="${data.EditPost.Title}" maxlength="300" autocomplete="off"
                autofocus required>
            <textarea name="post" id="" cols="30" rows="10" required>${data.EditPost.Content}</textarea>
            <select name="tags" id="tags" multiple="multiple">
                <option value="" disabled>Select a tag (Hold Ctrl for multiple tags)</option>
                ${data.TagsList.map(tag => `<option value="${tag.TagName}">${tag.TagName}</option>`).join('')}
            </select>
            <div class="submission-bottom-row">
                <label for="file-upload" class="image-upload">
                    <i class="fa-solid fa-image"></i> Upload Image
                </label>
                <input id="file-upload" name="submission-image" type="file"
                    accept="image/jpeg, image/png, image/gif, image/svg+xml" />
                <button name="editpost" value="${data.EditPost.UUID}">Edit</button>
            </div>
        </form>
        `;
    } else {
        html += `
        <form id="submit-post-form" method="POST" enctype="multipart/form-data">
            <input type="text" name="submission-title" placeholder="Title" maxlength="300" autocomplete="off" autofocus
                required>
            <textarea name="post" id="" cols="30" rows="10" placeholder="Text" required></textarea>
            <select name="tags" id="tags" multiple="multiple">
                <option value="" disabled>Select a tag (Hold Ctrl for multiple tags)</option>
                ${data.TagsList.map(tag => `<option value="${tag.TagName}">${tag.TagName}</option>`).join('')}
            </select>
            <div class="submission-bottom-row">
                <label for="file-upload" class="image-upload">
                    <i class="fa-solid fa-image"></i> Upload Image
                </label>
                <input id="file-upload" name="submission-image" type="file"
                    accept="image/jpeg, image/png, image/gif, image/svg+xml" />
                <button type="submit">Post</button>
            </div>
        </form>
        `;
    }

    html += `
    </div>
  </div>
`;

    document.getElementById('container').innerHTML = html;

    const submitPost = async (formData) => {
        try {
            const response = await fetch('/submitpost', {
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
                console.error('Submit Post failed.');
            }
        } catch (error) {
            // Handle network or other errors
            console.error('Error occurred:', error);
        }
    };

    // Submit form event handler
    document.querySelector('#submit-post-form').addEventListener('submit', async (event) => {
        event.preventDefault(); // Prevent form submission

        var editPostButton = document.querySelector('button[name="editpost"]');
        var editPostValue = "";
        if (editPostButton) {
            editPostValue = editPostButton.value;
        }
        console.log(editPostValue);

        const formData = {
            title: document.querySelector('input[name="submission-title"]').value,
            content: document.querySelector('textarea[name="post"]').value,
            tags: [""],
            editpost: editPostValue
        };

        console.log(formData);

        submitPost(formData);
    });

}