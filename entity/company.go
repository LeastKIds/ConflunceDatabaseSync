package entity

// @db
// @name {회사}
type Company struct {
	ID   int    `no:"1" column:"id" name:"아이디" type:"BIGINT" pk:"O" fk:"" null:"X" index:"" reference:"자동생성"`
	Name string `no:"2" column:"name" name:"회사명" type:"VARCHAR(255)" pk:"" fk:"" null:"X" index:"" reference:"회사명"`
}
