package lib

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

type jsonClient interface {
	Get(route string, outData interface{}) error
}

type AtlasClient struct {
	JSONClient jsonClient
}

func (c *AtlasClient) GetLatestVersion(boxName string) (string, error) {
	var metadata struct {
		CurrentVersion struct {
			Version string `json:"version"`
		} `json:"current_version"`
	}

	err := c.JSONClient.Get("/api/v1/box/"+boxName, &metadata)
	if err != nil {
		return "", err
	}

	ver := metadata.CurrentVersion.Version

	if ver == "" {
		return "", errors.New("missing version in JSON returned from Atlas API")
	}

	return ver, nil
}

func (c *AtlasClient) GetAMIs(boxName, version string) (map[string]string, error) {
	var metadata struct {
		Versions []struct {
			Version   string
			Providers []struct {
				Name        string
				DownloadURL string `json:"download_url"`
			}
		}
	}

	err := c.JSONClient.Get("/api/v1/box/"+boxName, &metadata)
	if err != nil {
		return nil, err
	}

	var downloadURL = ""
	for _, v := range metadata.Versions {
		if v.Version != version {
			continue
		}
		for _, provider := range v.Providers {
			if provider.Name == "aws" {
				downloadURL = provider.DownloadURL
			}
		}
	}

	if downloadURL == "" {
		return nil, fmt.Errorf("unable to find box %q version %q", boxName, version)
	}

	gzippedBoxResp, err := http.Get(downloadURL)
	if err != nil {
		return nil, fmt.Errorf("error downloading box: %s", err)
	}
	defer gzippedBoxResp.Body.Close()

	tarReader, err := gzip.NewReader(gzippedBoxResp.Body)
	if err != nil {
		return nil, fmt.Errorf("error gunzipping box: %s", err)
	}

	tarBytes, err := ioutil.ReadAll(tarReader)
	if err != nil {
		return nil, fmt.Errorf("error reading box bytes: %s", err) // not tested
	}

	// aws.region_config "eu-west-1", ami: "ami-4d8eac3a"
	amiLineParts := regexp.MustCompile(`\"([a-z,0-9,\-]*)\", ami: \"(ami-[a-z,0-9]*)\"`).FindAllSubmatch(tarBytes, -1)
	if amiLineParts == nil {
		return nil, fmt.Errorf("no AMIs found within box")
	}

	amiMap := map[string]string{}

	for _, lineParts := range amiLineParts {
		amiMap[string(lineParts[1])] = string(lineParts[2])
	}

	return amiMap, nil
}
