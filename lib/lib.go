package lib

import "fmt"

type Source struct {
	Region string `json:"region"`
}

type Version struct {
	BoxVersion string `json:"boxversion"`
}

type atlasClient interface {
	GetLatestVersion(boxName string) (string, error)
}

type Resource struct {
	AtlasClient atlasClient
	BoxName     string
}

const atlasBaseURL = "https://atlas.hashicorp.com"

func NewResource() *Resource {
	return &Resource{
		AtlasClient: &AtlasClient{
			JSONClient: &JSONClient{
				BaseURL: atlasBaseURL,
			},
		},
	}
}

func (r *Resource) Check(source Source, old Version) ([]Version, error) {
	latest, err := r.AtlasClient.GetLatestVersion(r.BoxName)
	if err != nil {
		return nil, fmt.Errorf("atlas client: %s", err)
	}

	if old.BoxVersion == latest {
		return []Version{}, nil
	}
	return []Version{Version{BoxVersion: latest}}, nil
}
