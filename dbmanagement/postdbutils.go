package dbmanagement

import (
	"database/sql"
	"fmt"
	"forum/utils"
	"strings"
	"time"
)

/*
Inserts post into database with the relevant data, likes and dislikes should be set to 0 for most cases.  Each post has it's own UUID.
*/
func InsertPost(title string, content string, ownerId string, likes int, dislikes int, inputtime time.Time, imageName string) (Post, error) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	utils.WriteMessageToLogFile("Inserting post record...")

	UUID := GenerateUUIDString()
	insertPostData := "INSERT INTO Posts(UUID, title, content, ownerId, likes, dislikes, time, imagename) VALUES (?, ?, ?, ?, ?, ?, ?, ?)"
	statement, err := db.Prepare(insertPostData)
	if err != nil {
		utils.HandleError("User Prepare failed: ", err)
		return Post{}, err
	}

	_, err = statement.Exec(UUID, title, content, ownerId, likes, dislikes, inputtime, imageName)
	if err != nil {
		utils.HandleError("Statement Exec failed: ", err)
		return Post{}, err
	}

	tags := SelectAllTagsFromPost(UUID)

	return Post{UUID, title, content, ownerId, likes, dislikes, tags, inputtime, strings.TrimSuffix(inputtime.Format(time.RFC822), "UTC"), 0, imageName}, err
}

func UpdatePost(postuuid string, title string, content string, ownerId string, likes int, dislikes int, edittime time.Time, imageName string) (Post, error) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()
	utils.WriteMessageToLogFile("Updating post record...")

	updatePostData := `
	UPDATE Posts
	SET title = ?, content = ?, time = ?
	WHERE uuid = ?
	`
	statement, err := db.Prepare(updatePostData)
	if err != nil {
		utils.HandleError("User Prepare failed: ", err)
		return Post{}, err
	}

	_, err = statement.Exec(title, content, edittime, postuuid)
	if err != nil {
		utils.HandleError("Statement Exec failed: ", err)
		return Post{}, err
	}

	tags := SelectAllTagsFromPost(postuuid)

	return Post{postuuid, title, content, ownerId, likes, dislikes, tags, edittime, strings.TrimSuffix(edittime.Format(time.RFC822), "UTC"), 0, imageName}, err
}

/*
Displays all posts from the database in the console.  Only for internal use.
*/
func DisplayAllPosts() error {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Posts ORDER BY time")
	utils.HandleError("Display posts query failed: ", err)

	defer row.Close()

	for row.Next() {
		var UUID string
		var title string
		var content string
		var ownerId string
		var likes int
		var dislikes int
		var time time.Time
		var imageName string
		row.Scan(&UUID, &title, &content, &ownerId, &likes, &dislikes, &time, &imageName)
		owner, err := SelectUserFromUUID(ownerId)
		utils.HandleError("unable to get user to display post", err)
		message := fmt.Sprint("Post: ", UUID, " content: ", content, " owner: ", owner.Name, " likes ", likes, " dislikes ", dislikes, " time ", time, "tags ", SelectAllTagsFromPost(UUID))
		utils.WriteMessageToLogFile(message)
	}
	return err
}

/*
Finds a specific post based on the UUID (of the post).  Could be used for when bringing up a particular post onto one page.
*/
func SelectPostFromUUID(UUID string) (Post, error) {
	var post Post
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	stm, err := db.Prepare("SELECT * FROM Posts WHERE uuid = ?")
	if err != nil {
		utils.HandleError("Statement failed: ", err)
		return Post{}, err
	}

	err = stm.QueryRow(UUID).Scan(&post.UUID, &post.Title, &post.Content, &post.OwnerId, &post.Likes, &post.Dislikes, &post.Time, &post.ImageName)
	post.FormattedTime = strings.TrimSuffix(post.Time.Format(time.RFC822), "UTC")
	post.NumOfComments = len(SelectAllCommentsFromPost(post.UUID))
	post.Tags = SelectAllTagsFromPost(post.UUID)

	if err != nil {
		utils.HandleError("Query Row failed: ", err)
		return Post{}, err
	}

	return post, err
}

/*
Gathers all the posts from the database and returns them as an array of Post struct.  This function is used when displaying all the posts on the forum website.
*/
func SelectAllPosts() ([]Post, error) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Posts")
	utils.HandleError("All posts query failed: ", err)
	defer row.Close()

	var allPosts []Post
	for row.Next() {
		var currentPost Post
		row.Scan(&currentPost.UUID, &currentPost.Title, &currentPost.Content, &currentPost.OwnerId, &currentPost.Likes, &currentPost.Dislikes, &currentPost.Time, &currentPost.ImageName)
		currentPost.FormattedTime = strings.TrimSuffix(currentPost.Time.Format(time.RFC822), "UTC")
		currentPost.NumOfComments = len(SelectAllCommentsFromPost(currentPost.UUID))
		currentPost.Tags = SelectAllTagsFromPost(currentPost.UUID)
		allPosts = append(allPosts, currentPost)
	}
	return allPosts, err
}

/*
Similar to SelectAllPosts() but for a specific user.  Uses the ownerID (users UUID) to specify which user and returns all the posts created by that user.
*/
func SelectAllPostsFromUser(ownerId string) ([]Post, error) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	row, err := db.Query("SELECT * FROM Posts WHERE ownerId = ?", ownerId)
	utils.HandleError("Post from User query failed: ", err)
	defer row.Close()

	var allPosts []Post

	for row.Next() {
		var currentPost Post
		row.Scan(&currentPost.UUID, &currentPost.Title, &currentPost.Content, &currentPost.OwnerId, &currentPost.Likes, &currentPost.Dislikes, &currentPost.Time, &currentPost.ImageName)
		currentPost.FormattedTime = strings.TrimSuffix(currentPost.Time.Format(time.RFC822), "UTC")
		currentPost.NumOfComments = len(SelectAllCommentsFromPost(currentPost.UUID))
		currentPost.Tags = SelectAllTagsFromPost(currentPost.UUID)
		allPosts = append(allPosts, currentPost)
	}
	return allPosts, err
}

func SelectAllLikedPostsFromUser(user User) ([]Post, error) {
	allPosts, err := SelectAllPosts()
	utils.HandleError("Unable to Select all liked posts", err)
	likedPosts := []Post{}
	for _, v := range allPosts {
		if SelectReactionFromPost(v.UUID, user.UUID) == 1 {
			likedPosts = append(likedPosts, v)
		}
	}
	return likedPosts, err
}

func SelectAllDislikedPostsFromUser(user User) ([]Post, error) {
	allPosts, err := SelectAllPosts()
	utils.HandleError("Unable to Select all liked posts", err)
	DislikedPosts := []Post{}
	for _, v := range allPosts {
		if SelectReactionFromPost(v.UUID, user.UUID) == -1 {
			DislikedPosts = append(DislikedPosts, v)
		}
	}
	return DislikedPosts, err
}

/*
Similar to SelectAllPosts() but for a specific user.  Uses the ownerID (users UUID) to specify which user and returns all the posts created by that user.
*/
func SelectAllPostsFromTag(tagName string) ([]Post, error) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	tag, err := SelectTagFromName(tagName)
	if err != nil {
		utils.HandleError("couldn't find tag", err)
		return []Post{}, err
	}

	row, err := db.Query("SELECT postId FROM TaggedPosts WHERE tagId = ?", tag.UUID)
	if err != nil {
		utils.HandleError("Tag query failed: ", err)
		return []Post{}, err
	}
	defer row.Close()

	var allPosts []Post

	for row.Next() {
		var currentPostId string
		var currentPost Post
		row.Scan(&currentPostId)
		currentPost, err = SelectPostFromUUID(currentPostId)
		if err != nil {
			utils.HandleError("Unable to select post from uuid from select all posts ", err)
			return []Post{}, err
		}
		currentPost.FormattedTime = strings.TrimSuffix(currentPost.Time.Format(time.RFC822), "UTC")
		currentPost.Tags = SelectAllTagsFromPost(currentPost.UUID)
		currentPost.NumOfComments = len(SelectAllCommentsFromPost(currentPost.UUID))
		message := fmt.Sprint("found post from tag: ", tagName, "the post: ", currentPost)
		utils.WriteMessageToLogFile(message)
		allPosts = append(allPosts, currentPost)
	}

	return allPosts, err
}
