package entity

// @db
// @name {테스트}
type Test struct {
	ID   int    `no:"1" column:"id" name:"아이디" type:"BIGINT" pk:"O" fk:"" null:"X" index:"" reference:"자동생성"`
	Name string `no:"2" column:"name" name:"이름" type:"VARCHAR(255)" pk:"" fk:"" null:"X" index:"" reference:"이름"`
}
