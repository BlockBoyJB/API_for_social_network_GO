package pgmodel

type Post struct {
	Id       int    `db:"id"`
	Username string `db:"username"`
	PostId   string `db:"post_id"`
	Title    string `db:"title"`
	Text     string `db:"text"`
}
