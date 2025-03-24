package main

const (
	pageID   = "65827"
	spaceKey = "96QjrC0InFzz"
)

func main() {
	html, err := getStructTable()
	if err != nil {
		panic(err)
	}

}
