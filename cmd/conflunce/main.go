package main

const (
	parentPageID = "65809"
	pageID       = "65827"
	spaceKey     = "96QjrC0InFzz"
)

func main() {
	html, err := getStructTable()
	if err != nil {
		panic(err)
	}

	if err := updateConflunce(html, parentPageID, pageID, spaceKey); err != nil {
		panic(err)
	}

}
