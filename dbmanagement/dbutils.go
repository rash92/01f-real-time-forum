package dbmanagement

import (
	"database/sql"
	"forum/utils"
	"os"
	"time"

	"github.com/gofrs/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var createUserTableStatement = `
	CREATE TABLE Users (
		uuid TEXT NOT NULL PRIMARY KEY,		
		name TEXT UNIQUE,
		email TEXT UNIQUE,
		password TEXT,
		permission TEXT,
		IsLoggedIn INTEGER,
		limitTokens INTEGER
	);`

// ADD TITLE TO POST TABLE AND THEN FIX EVERYTHING
var createPostTableStatement = `
	CREATE TABLE Posts (
		uuid TEXT NOT NULL PRIMARY KEY,
		title TEXT,
		content TEXT,		
		ownerId TEXT,
		likes INTEGER,
		dislikes INTEGER,
		time DATETIME,
		imagename TEXT,
		FOREIGN KEY (ownerId) REFERENCES Users(uuid) ON DELETE SET NULL
	);`

var createCommentTableStatement = `
	CREATE TABLE Comments (
		uuid TEXT NOT NULL PRIMARY KEY,
		content TEXT,
		postId TEXT,
		ownerId TEXT,
		likes INTEGER,
		dislikes INTEGER,
		time DATETIME,
		FOREIGN KEY (postId) REFERENCES Posts(uuid),
		FOREIGN KEY (ownerId) REFERENCES Users(uuid)
	);`

var createTagsTableStatement = `
	CREATE TABLE Tags (
		uuid TEXT NOT NULL PRIMARY KEY,
		tagname TEXT NOT NULL UNIQUE
	);`

var createTaggedPostsStatement = `
CREATE TABLE TaggedPosts (
		tagId TEXT NOT NULL,
		postId TEXT NOT NULL,
		FOREIGN KEY (tagId) REFERENCES Tags(uuid) ON DELETE CASCADE,
		FOREIGN KEY (postId) REFERENCES Posts(uuid) ON DELETE CASCADE,
		PRIMARY KEY (tagId, postId)
	);`

// can represent like as 1, dislike as -1 and neither as 0 as a single value in the reaction field
var createReactionPostsTableStatement = `
	CREATE TABLE ReactedPosts (
		userId TEXT NOT NULL,
		postId TEXT NOT NULL,
		reaction INTEGER,
		FOREIGN KEY (userId) REFERENCES Users(uuid),
		FOREIGN KEY (postId) REFERENCES Posts(uuid),
		PRIMARY KEY (userId, postId)
	);`

var createReactedCommentsTableStatement = `
	CREATE TABLE ReactedComments  (
		userId TEXT NOT NULL,
		commentId TEXT NOT NULL,
		reaction INTEGER,
		FOREIGN KEY (userId) REFERENCES Users(uuid),
		FOREIGN KEY (commentId) REFERENCES Comments(uuid),
		PRIMARY KEY (userId, commentId)
	);`

var createSessionTableStatement = `
	CREATE TABLE Sessions (
  		uuid      TEXT NOT NULL PRIMARY KEY,
  		userId    INTEGER REFERENCES Users(uuid),
  		createdAt TIMESTAMP NOT NULL   
	);`

var createAdminRequestTableStatement = `
	CREATE TABLE AdminRequests (
		uuid TEXT NOT NULL PRIMARY KEY,
		requestfromid TEXT,
		requestfromname TEXT,
		reportedpostid TEXT,
		reportedcommentid TEXT,
		reporteduserid TEXT,
		description TEXT,
		FOREIGN KEY (reportedpostid) REFERENCES Posts(uuid),
		FOREIGN KEY (reportedcommentid) REFERENCES Comments(uuid),
		FOREIGN KEY (reporteduserid) REFERENCES Users(uuid),
		FOREIGN KEY (requestfromid) REFERENCES Users(uuid),
		FOREIGN KEY (requestfromid) REFERENCES Users(name)	
	);`

var createNotificationsTableStatement = `
	CREATE TABLE Notifications (
		uuid TEXT NOT NULL PRIMARY KEY,
		receivingUserId TEXT,
		postId TEXT,
		commentId TEXT,
		sendingUserId TEXT,
		reaction INT,
		notificationStatement TEXT,
		FOREIGN KEY (receivingUserId) REFERENCES Users(uuid),
		FOREIGN KEY (postId) REFERENCES Posts(uuid),
		FOREIGN KEY (commentId) REFERENCES Comments(uuid),
		FOREIGN KEY (sendingUserId) REFERENCES Users(uuid)
	);`

/*
Only used to create brand new databases, wiping all previous data in the process.
To be used when initially implementing database or clearing data after testing.
Also inserts a user with admin permissions by default, with both username and password being 'admin'
*/
func CreateDatabaseWithTables() {
	forumDB := CreateDatabase("forum")
	defer forumDB.Close()

	CreateTable(forumDB, createUserTableStatement)
	CreateTable(forumDB, createPostTableStatement)
	CreateTable(forumDB, createCommentTableStatement)
	CreateTable(forumDB, createTagsTableStatement)
	CreateTable(forumDB, createTaggedPostsStatement)
	CreateTable(forumDB, createReactionPostsTableStatement)
	CreateTable(forumDB, createReactedCommentsTableStatement)
	CreateTable(forumDB, createSessionTableStatement)
	CreateTable(forumDB, createAdminRequestTableStatement)
	CreateTable(forumDB, createNotificationsTableStatement)

	os.Remove("logfile.txt")
	os.Create("logfile.txt")

	// had to manually reimplement hashing as get 'import cycle error' if you import auth package
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	utils.HandleError("password hashing error for default admin on database creation", err)

	InsertUser("admin", "a@a", string(hashedPassword), "admin", 0)
	e := os.RemoveAll("./static/uploads/")
	if e != nil {
		utils.HandleError("Unable to insert user", e)
	}

	utils.WriteMessageToLogFile("forum.db created successfully!")
}

/*
Creates a new database file to store tables.  If database already exists, it is REMOVED.  Beware of losing data.
*/
func CreateDatabase(name string) *sql.DB {
	os.Remove(name + ".db")
	utils.WriteMessageToLogFile("Creating " + name + ".db...")
	file, err := os.Create(name + ".db")
	utils.HandleError("", err)

	file.Close()

	forumDB, err := sql.Open("sqlite3", "./"+name+".db?_foreign_keys=on")
	utils.HandleError("", err)

	return forumDB
}

/*
Creates a table within a specified database
*/
func CreateTable(db *sql.DB, table string) {
	statement, err := db.Prepare(table)
	utils.HandleError(table, err)
	statement.Exec()
}

/*
Generates a new UUID and returns a string of that new number.
*/
func GenerateUUIDString() string {
	UUID, err := uuid.NewV4()
	utils.HandleError("problem generating uuid", err)
	return UUID.String()
}

/*
Used to provide specific information for when a user logs in by cross referencing their userID.
Creates and returns a new session when the user successfully logs in to their account.
The sessions has its own UUID, contains the usersID (user's UUID), and the time that it was created.
*/
func (user *User) CreateSession() (session Session, err error) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	statement := `INSERT INTO Sessions (uuid, userID, createdAt) values (?, ?, ?) returning uuid, userID, createdAt`

	stmt, err := db.Prepare(statement)
	utils.HandleError("session error:", err)

	defer stmt.Close()

	UUID := GenerateUUIDString()
	timeNow := time.Now()

	err = stmt.QueryRow(UUID, user.UUID, timeNow).Scan(&session.UUID, &session.UserId, &session.CreatedAt)
	return
}

func (user *User) ReturnSession(userId string) (session Session, err error) {
	db, _ := sql.Open("sqlite3", "./forum.db")
	defer db.Close()

	stm, err := db.Prepare("SELECT * FROM Sessions WHERE userId = ?")
	utils.HandleError("Statement failed: ", err)

	err = stm.QueryRow(userId).Scan(&session.UUID, &session.UserId, &session.CreatedAt)
	utils.HandleError("Query Row failed: ", err)

	return
}
