package main

import (
	"os"

	goconfluence "github.com/virtomize/confluence-go-api"
)

func updateConflunce(html, parentPageID, pageID, spaceKey string) error {
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
		return err
	}

	return nil
}
