package dbmanagement

import "time"

type User struct {
	UUID          string
	Name          string
	Email         string
	Password      string
	Permission    string
	IsLoggedIn    int
	Notifications []Notification
	LimitTokens   int
}

type Post struct {
	UUID          string
	Title         string
	Content       string
	OwnerId       string
	Likes         int
	Dislikes      int
	Tags          []Tag
	Time          time.Time
	FormattedTime string
	NumOfComments int
	ImageName     string
}

type Comment struct {
	UUID          string
	Content       string
	PostId        string
	OwnerId       string
	OwnerName     string
	Likes         int
	Dislikes      int
	Time          time.Time
	FormattedTime string
}

type Session struct {
	UUID      string
	UserId    string
	CreatedAt time.Time
}

type Tag struct {
	UUID    string
	TagName string
}

type AdminRequest struct {
	UUID              string
	RequestFromId     string
	RequestFromName   string
	ReportedPostId    string
	ReportedCommentId string
	ReportedUserId    string
	Description       string
}

type Notification struct {
	UUID      string
	Receiver  string
	PostId    string
	CommentId string
	Sender    string
	Reaction  int
	Statement string
}

type ChatBox struct {
	UUID     string
	Sender   User
	Receiver User
	Content  []ChatText //sort it chronologically
}

//implement a date system
type ChatText struct {
	Content string //encrypt content for security purposes

	//time is relative to the server's time zone
	//the text's time will be loaded relative
	//to the user's time zone in JS
	Time time.Time
}
