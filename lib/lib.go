package lib

import "fmt"

type Source struct {
	Region       string `json:"region"`
	BoxName      string `json:"box_name"`
	AtlasBaseURL string `json:"atlas_base_url"`
}

type Version struct {
	BoxVersion string `json:"box_version"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type atlasClient interface {
	GetLatestVersion(boxName string) (string, error)
	GetAMIs(boxName, version string) (map[string]string, error)
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

func (r *Resource) In(ver Version) (string, error) {
	allRegions, err := r.AtlasClient.GetAMIs(r.SourceConfig.BoxName, ver.BoxVersion)
	if err != nil {
		return "", fmt.Errorf("atlas client: %s", err)
	}
	ami, ok := allRegions[r.SourceConfig.Region]
	if !ok {
		return "", fmt.Errorf("no ami found for region %q", r.SourceConfig.Region)
	}
	return ami, nil
}
