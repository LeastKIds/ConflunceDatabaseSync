package entity

// @db
// @name {프로필}
type Profile struct {
	ID     int    `no:"1" column:"id" name:"아이디" type:"BIGINT" pk:"O" fk:"" null:"X" index:"" reference:"자동생성"`
	UserID int    `no:"2" column:"user_id" name:"유저 아이디" type:"BIGINT" pk:"" fk:"O" null:"X" index:"" reference:"유저 아이디"`
	Phone  string `no:"3" column:"phone" name:"전화번호" type:"VARCHAR(255)" pk:"" fk:"" null:"X" index:"" reference:"전화번호"`
	Email  string `no:"4" column:"email" name:"이메일" type:"VARCHAR(255)" pk:"" fk:"" null:"X" index:"" reference:"이메일"`
	Age    int    `no:"5" column:"age" name:"나이" type:"INT" pk:"" fk:"" null:"X" index:"" reference:"나이"`
}
