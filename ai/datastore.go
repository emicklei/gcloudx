package ai

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/discoveryengine/v1"
	"google.golang.org/api/option"
)

type DataStoreArguments struct {
	ProjectID     string
	URLPattern    string
	DatastoreName string
	Location      string
}

// assume for now: global, advanced public site, type chat, generic
// https://cloud.google.com/generative-ai-app-builder/docs/create-data-store-es#website
func CreateDatastore(args DataStoreArguments) error {
	opts := option.WithEndpoint("global-discoveryengine.googleapis.com")
	srv, err := discoveryengine.NewService(context.Background(), opts)
	if err != nil {
		log.Println("new service", err)
		return err
	}
	parent := fmt.Sprintf("projects/%s/locations/%s/collections/default_collection", args.ProjectID, args.Location)
	ds := &discoveryengine.GoogleCloudDiscoveryengineV1DataStore{
		Name:             fmt.Sprintf("projects/%s/locations/%s/collections/default_collection/dataStores/%s", args.ProjectID, args.Location, args.DatastoreName),
		DisplayName:      args.DatastoreName,
		IndustryVertical: "GENERIC",
		SolutionTypes:    []string{"SOLUTION_TYPE_CHAT"},
		ContentConfig:    "PUBLIC_WEBSITE",
	}
	call := srv.Projects.Locations.DataStores.Create(parent, ds)
	call.CreateAdvancedSiteSearch(true)
	call.DataStoreId(args.DatastoreName)
	_, err = call.Do()
	if err != nil {
		log.Println("datastore create", err, args)
		return err
	}
	return nil
}

// assume for now:  include only, exactmatch false, siteSearchEngine.
func CreateTargetSite(args DataStoreArguments) error {
	opts := option.WithEndpoint("global-discoveryengine.googleapis.com")
	srv, err := discoveryengine.NewService(context.Background(), opts)
	if err != nil {
		log.Println("new service", err)
		return err
	}
	sess := discoveryengine.NewProjectsLocationsDataStoresSiteSearchEngineService(srv)
	site := &discoveryengine.GoogleCloudDiscoveryengineV1TargetSite{
		ExactMatch:         false,
		ProvidedUriPattern: args.URLPattern,
		Type:               "INCLUDE",
	}
	dsParent := fmt.Sprintf("projects/%s/locations/%s/dataStores/%s/siteSearchEngine", args.ProjectID, args.Location, args.DatastoreName)
	op, err := sess.TargetSites.Create(dsParent, site).Do()
	if err != nil {
		log.Println("targetsite create", err)
		return err
	}
	for {
		if op.Done {
			break
		}
		if op.Error != nil {
			err = errors.New(op.Error.Message)
			break
		}
		log.Println("wait for complete...")
		time.Sleep(1 * time.Second)
	}
	return err
}
