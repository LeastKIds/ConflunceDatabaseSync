package entity

// @db
// @name {유저}
type User struct {
	ID    int    `no:"1" column:"id" name:"아이디" type:"BIGINT" pk:"O" fk:"" null:"X" index:"" reference:"자동생성"`
	Name  string `no:"2" column:"name" name:"이름" type:"VARCHAR(255)" pk:"" fk:"" null:"X" index:"" reference:"이름"`
	Email string `no:"3" column:"email" name:"이메일" type:"VARCHAR(255)" pk:"" fk:"" null:"X" index:"" reference:"이메일"`
	Age   int    `no:"4" column:"age" name:"나이" type:"INT" pk:"" fk:"" null:"X" index:"" reference:"나이"`
}

// 일반 struct
type Temp struct {
	Something string
}
