package entity

// @db
// @name {유저}
type User struct {
	ID    int    `json:"id" db:"id" reference:"아이디3"`
	Name  string `json:"name" db:"name" reference:"이름d"`
	Email string `json:"email" db:"email" reference:"이메일1"`
	Age   int    `json:"age" db:"age" reference:"나이"`
}

// 일반 struct
type Temp struct {
	Something string
}
