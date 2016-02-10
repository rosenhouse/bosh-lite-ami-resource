package lib

import "fmt"

type Source struct {
	Region       string `json:"region"`
	BoxName      string `json:"box_name"`
	AtlasBaseURL string `json:"atlas_base_url"`
}

type Version struct {
	BoxVersion string `json:"boxversion"`
}

type atlasClient interface {
	GetLatestVersion(boxName string) (string, error)
}

type Resource struct {
	AtlasClient  atlasClient
	SourceConfig Source
}

const default_AtlasBaseURL = "https://atlas.hashicorp.com"
const default_BoxName = "cloudfoundry/bosh-lite"
const default_Region = "us-east-1"

func applyDefaults(source Source) Source {
	if source.Region == "" {
		source.Region = default_Region
	}
	if source.BoxName == "" {
		source.BoxName = default_BoxName
	}
	if source.AtlasBaseURL == "" {
		source.AtlasBaseURL = default_AtlasBaseURL
	}
	return source
}

func NewResource(source Source) *Resource {
	source = applyDefaults(source)
	return &Resource{
		SourceConfig: source,
		AtlasClient: &AtlasClient{
			JSONClient: &JSONClient{
				BaseURL: source.AtlasBaseURL,
			},
		},
	}
}

func (r *Resource) Check(old Version) ([]Version, error) {
	latest, err := r.AtlasClient.GetLatestVersion(r.SourceConfig.BoxName)
	if err != nil {
		return nil, fmt.Errorf("atlas client: %s", err)
	}

	if old.BoxVersion == latest {
		return []Version{}, nil
	}
	return []Version{Version{BoxVersion: latest}}, nil
}
