package dbmanagement

import (
	"database/sql"
	"forum/utils"
)

// add either a like or dislike or return to neutral with reaction taking values of 1,-1 or 0
/*
 */
func AddReactionToPost(userId string, postId string, reaction int) {
	if reaction > 1 || reaction < -1 || reaction == 0 {
		utils.WriteMessageToLogFile("Incorrect reaction integer")
		return
	}

	var reactionStatement string
	var postUpdateStatement string

	switch SelectReactionFromPost(postId, userId) {
	case -2:
		reactionStatement = `
		INSERT OR IGNORE INTO ReactedPosts(userId, postId, reaction) 
		VALUES (?, ?, ?)
	`
		if reaction == 1 {
			postUpdateStatement = `
	UPDATE Posts 
	SET likes = likes + 1
	WHERE uuid = ?
`
		} else if reaction == -1 {
			postUpdateStatement = `
	UPDATE Posts 
	SET dislikes = dislikes + 1 
	WHERE uuid = ?
`
		}
	case -1:
		if reaction == 1 {
			reactionStatement = `
		UPDATE ReactedPosts 
		SET reaction = ? 
		WHERE userId = ? and postId = ?
	`
			postUpdateStatement = `
		UPDATE Posts 
		SET likes = likes + 1, dislikes = dislikes - 1 
		WHERE uuid = ?
	`
		} else if reaction == -1 {
			reactionStatement = `
		UPDATE ReactedPosts 
		SET reaction = ? + 1
		WHERE userId = ? and postId = ?
	`
			postUpdateStatement = `
		UPDATE Posts 
		SET dislikes = dislikes - 1 
		WHERE uuid = ?
	`
		}
	case 0:
		reactionStatement = `
		UPDATE ReactedPosts 
		SET reaction = ? 
		WHERE userId = ? and postId = ?
	`
		if reaction == 1 {
			postUpdateStatement = `
		UPDATE Posts 
		SET likes = likes + 1
		WHERE uuid = ?
	`
		} else if reaction == -1 {
			postUpdateStatement = `
		UPDATE Posts 
		SET dislikes = dislikes + 1 
		WHERE uuid = ?
	`
		}
	case 1:
		if reaction == -1 {
			reactionStatement = `
		UPDATE ReactedPosts 
		SET reaction = ? 
		WHERE userId = ? and postId = ?
	`
			postUpdateStatement = `
		UPDATE Posts 
		SET likes = likes - 1, dislikes = dislikes + 1 
		WHERE uuid = ?
	`
		} else if reaction == 1 {
			reactionStatement = `
		UPDATE ReactedPosts 
		SET reaction = ? - 1 
		WHERE userId = ? and postId = ?
	`
			postUpdateStatement = `
		UPDATE Posts 
		SET likes = likes - 1 
		WHERE uuid = ?
	`
		}
	default:

	}

	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	statement, err := db.Prepare(reactionStatement)
	utils.HandleError("Like Prepare failed: ", err)

	ps, err := db.Prepare(postUpdateStatement)
	utils.HandleError("Post Update Prepare failed: ", err)

	if SelectReactionFromPost(postId, userId) == -2 {
		_, err = statement.Exec(userId, postId, reaction)
		utils.HandleError("statement Exec failed: ", err)
		_, err = ps.Exec(postId)
		utils.HandleError("post update Exec failed: ", err)
	} else {
		_, err = statement.Exec(reaction, userId, postId)
		utils.HandleError("statement Exec failed: ", err)
		_, err = ps.Exec(postId)
		utils.HandleError("post update Exec failed: ", err)
	}
}

func SelectReactionFromPost(postid, userid string) int {
	//-2 because you don't want to insert more entries into reaction table
	var reaction = -2
	var user string
	var post string
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	stm, err := db.Prepare("SELECT * FROM ReactedPosts WHERE userId = ? and postid = ?")
	utils.HandleError("Statement failed: ", err)

	err = stm.QueryRow(userid, postid).Scan(&user, &post, &reaction)
	utils.HandleError("Query Row failed: ", err)

	return reaction
}

// add either a like or dislike or return to neutral with reaction taking values of 1,-1 or 0
/*
 */
func AddReactionToComment(userId string, commentid string, reaction int) {
	if reaction > 1 || reaction < -1 || reaction == 0 {
		utils.WriteMessageToLogFile("Incorrect reaction integer")
		return
	}

	var reactionStatement string
	var commentUpdateStatement string

	switch SelectReactionFromComment(userId, commentid) {
	case -2:
		reactionStatement = `
		INSERT OR IGNORE INTO ReactedComments(userId, commentId, reaction) 
		VALUES (?, ?, ?)
	`
		if reaction == 1 {
			commentUpdateStatement = `
	UPDATE Comments 
	SET likes = likes + 1
	WHERE uuid = ?
`
		} else if reaction == -1 {
			commentUpdateStatement = `
	UPDATE Comments 
	SET dislikes = dislikes + 1 
	WHERE uuid = ?
`
		}
	case -1:
		if reaction == 1 {
			reactionStatement = `
		UPDATE ReactedComments 
		SET reaction = ? 
		WHERE userId = ? and commentId = ?
	`
			commentUpdateStatement = `
		UPDATE Comments 
		SET likes = likes + 1, dislikes = dislikes - 1 
		WHERE uuid = ?
	`
		} else if reaction == -1 {
			reactionStatement = `
		UPDATE ReactedComments 
		SET reaction = ? + 1
		WHERE userId = ? and commentId = ?
	`
			commentUpdateStatement = `
		UPDATE Comments 
		SET dislikes = dislikes - 1 
		WHERE uuid = ?
	`
		}
	case 0:
		reactionStatement = `
		UPDATE ReactedComments 
		SET reaction = ? 
		WHERE userId = ? and commentId = ?
	`
		if reaction == 1 {
			commentUpdateStatement = `
		UPDATE Comments 
		SET likes = likes + 1
		WHERE uuid = ?
	`
		} else if reaction == -1 {
			commentUpdateStatement = `
		UPDATE Comments 
		SET dislikes = dislikes + 1 
		WHERE uuid = ?
	`
		}
	case 1:
		if reaction == -1 {
			reactionStatement = `
		UPDATE ReactedComments 
		SET reaction = ? 
		WHERE userId = ? and commentId = ?
	`
			commentUpdateStatement = `
		UPDATE Comments 
		SET likes = likes - 1, dislikes = dislikes + 1 
		WHERE uuid = ?
	`
		} else if reaction == 1 {
			reactionStatement = `
		UPDATE ReactedComments 
		SET reaction = ? - 1 
		WHERE userId = ? and commentId = ?
	`
			commentUpdateStatement = `
		UPDATE Comments 
		SET likes = likes - 1 
		WHERE uuid = ?
	`
		}
	default:

	}

	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	statement, err := db.Prepare(reactionStatement)
	utils.HandleError("Like Prepare failed: ", err)

	ps, err := db.Prepare(commentUpdateStatement)
	utils.HandleError("Comment Update Prepare failed: ", err)

	if SelectReactionFromComment(userId, commentid) == -2 {
		_, err = statement.Exec(userId, commentid, reaction)
		utils.HandleError("statement Exec failed: ", err)
		_, err = ps.Exec(commentid)
		utils.HandleError("comment update Exec failed: ", err)
	} else {
		_, err = statement.Exec(reaction, userId, commentid)
		utils.HandleError("statement Exec failed: ", err)
		_, err = ps.Exec(commentid)
		utils.HandleError("comment update Exec failed: ", err)
	}
}

func SelectReactionFromComment(commentid, userid string) int {
	//-2 because you don't want to insert more entries into reaction table
	var reaction = -2
	var user string
	var comment string
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	stm, err := db.Prepare("SELECT * FROM ReactedComments WHERE userId = ? and commentId = ?")
	utils.HandleError("Statement failed: ", err)

	err = stm.QueryRow(commentid, userid).Scan(&user, &comment, &reaction)
	utils.HandleError("Query Row failed: ", err)

	return reaction
}
