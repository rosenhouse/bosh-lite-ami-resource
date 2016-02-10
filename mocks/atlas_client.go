package mocks

type AtlasClient struct {
	GetLatestAMIsCall struct {
		Receives struct {
			BoxName string
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

func (c *AtlasClient) GetLatestAMIs(boxName string) (map[string]string, error) {
	c.GetLatestAMIsCall.Receives.BoxName = boxName
	return c.GetLatestAMIsCall.Returns.AMIMap, c.GetLatestAMIsCall.Returns.Error
}

func (c *AtlasClient) GetLatestVersion(boxName string) (string, error) {
	c.GetLatestVersionCall.Receives.BoxName = boxName
	return c.GetLatestVersionCall.Returns.Version, c.GetLatestVersionCall.Returns.Error
}
