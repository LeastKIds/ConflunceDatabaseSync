package entity

// @db
// @name {프로필}
type Profile struct {
	ID     int    `json:"id" db:"id" reference:"아이디3"`
	UserID int    `json:"user_id" db:"user_id" reference:"유저아이디"`
	Phone  string `json:"phone" db:"phone" reference:"전화번호"`
	Email  string `json:"email" db:"email" reference:"이메일1"`
	Age    int    `json:"age" db:"age" reference:"나이"`
}
