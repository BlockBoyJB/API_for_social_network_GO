package pgmodel

type Comment struct {
	Id        int    `db:"id"`
	Username  string `db:"username"`
	PostId    string `db:"post_id"`
	CommentId string `db:"comment_id"`
	Comment   string `db:"comment"`
}
