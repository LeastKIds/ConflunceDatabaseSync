package main

import (
	"fmt"
	"os"

	goconfluence "github.com/virtomize/confluence-go-api"
)

func UpdateConflunce(html, parentPageID, pageID, spaceKey string) error {
	BaseURL :=
		os.Getenv("CONFLUENCE_BASE_URL")
	Username :=
		os.Getenv("CONFLUENCE_USERNAME")
	APIToken :=
		os.Getenv("CONFLUENCE_API_TOKEN")

	api, err := goconfluence.NewAPI(BaseURL, Username, APIToken)
	if err != nil {
		return err
	}

	content, err := api.GetContentByID(pageID, goconfluence.ContentQuery{
		SpaceKey: spaceKey,
		Expand:   []string{"version"},
	})
	if err != nil {
		fmt.Printf("GetContentByID error, spaceKey: %s, pageID: %s ", spaceKey, pageID)
		return err
	}
	newVersion := content.Version.Number + 1

	updateContent := &goconfluence.Content{
		ID:    pageID,
		Type:  "page",
		Title: "DB 스키마",
		Body: goconfluence.Body{
			Storage: goconfluence.Storage{
				Value:          html,
				Representation: "storage",
			},
		},
		Ancestors: []goconfluence.Ancestor{
			{
				ID: parentPageID,
			},
		},
		Version: &goconfluence.Version{
			Number: newVersion,
		},
		Space: &goconfluence.Space{
			Key: spaceKey,
		},
	}

	if _, err := api.UpdateContent(updateContent); err != nil {
		fmt.Println("UpdateContent error")
		return err
	}

	return nil
}
