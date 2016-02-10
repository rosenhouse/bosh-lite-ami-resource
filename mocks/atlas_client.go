package mocks

type AtlasClient struct {
	GetAMIsCall struct {
		Receives struct {
			BoxName string
			Version string
		}
		Returns struct {
			AMIMap map[string]string
			Error  error
		}
	}

	GetLatestVersionCall struct {
		Receives struct {
			BoxName string
		}
		Returns struct {
			Version string
			Error   error
		}
	}
}

func (c *AtlasClient) GetAMIs(boxName, version string) (map[string]string, error) {
	c.GetAMIsCall.Receives.BoxName = boxName
	c.GetAMIsCall.Receives.Version = version
	return c.GetAMIsCall.Returns.AMIMap, c.GetAMIsCall.Returns.Error
}

func (c *AtlasClient) GetLatestVersion(boxName string) (string, error) {
	c.GetLatestVersionCall.Receives.BoxName = boxName
	return c.GetLatestVersionCall.Returns.Version, c.GetLatestVersionCall.Returns.Error
}
