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

// 1) should be loaded upon clicking on a friend from the "active" list
// 2) should be updated every time a text is sent/received in sql or relative to WS
// 3) should create a []ChatBox (or a map) which represents every active chat
type ChatBox struct {
	UUID string `json:"uuid"`

	Content []ChatText //sort it chronologically
}

// implement a date system
type ChatText struct {
	Content    string `json:"content"` //encrypt content for security purposes
	SenderId   string `json:"sender"`
	ReceiverId string `json:"receiver"`
	//time is relative to the server's time zone
	//the text's time will be loaded relative
	//to the user's time zone in JS
	Time time.Time `json:"time"`
}
