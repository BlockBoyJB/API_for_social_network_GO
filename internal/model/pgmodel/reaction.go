package pgmodel

type Reaction struct {
	Id         int    `db:"id"`
	PostId     string `db:"post_id"`
	ReactionId string `db:"reaction_id"`
	Reaction   string `db:"reaction"`
}
