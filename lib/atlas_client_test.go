package lib_test

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/rosenhouse/bosh-lite-ami-resource/lib"
	"github.com/rosenhouse/bosh-lite-ami-resource/mocks"
)

var _ = Describe("AtlasClient", func() {

	var fakeDownloadServer *ghttp.Server
	var jsonClient *mocks.JSONClient
	var atlasClient lib.AtlasClient

	Describe("#GetAMIs", func() {
		BeforeEach(func() {
			gzippedBoxData, err := ioutil.ReadFile("fixtures/test-box.gz")
			Expect(err).NotTo(HaveOccurred())
			fakeDownloadServer = ghttp.NewServer()
			fakeDownloadRoute := "/some/download/url"
			fakeDownloadServer.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", fakeDownloadRoute),
					ghttp.RespondWith(http.StatusOK, gzippedBoxData),
				),
			)
			fakeDownloadURL := fakeDownloadServer.URL() + fakeDownloadRoute

			jsonClient = &mocks.JSONClient{}
			jsonClient.GetCall.ResponseJSON = fmt.Sprintf(`{
				"versions" : [
					{
						"version": "8000.92.0",
						"providers": [
							{
								"name": "some-provider",
								"download_url": "some-other-download-url"
							},
							{
								"name": "aws",
								"download_url": "wrong-url"
							}
						]
					},
					{
						"version": "some-special-version",
						"providers": [
							{
								"name": "some-provider",
								"download_url": "some-other-download-url"
							},
							{
								"name": "aws",
								"download_url": "%s"
							}
						]
					},
					{
						"version": "9000.91.0",
						"providers": [
							{
								"name": "some-provider",
								"download_url": "some-other-download-url"
							},
							{
								"name": "aws",
								"download_url": "another-wrong-url"
							}
						]
					}
				]
			}`, fakeDownloadURL)
			atlasClient = lib.AtlasClient{jsonClient}
		})
		AfterEach(func() {
			fakeDownloadServer.Close()
		})

		It("should return the AMI used by the box in the specified region", func() {
			amiMap, err := atlasClient.GetAMIs("someuser/somebox", "some-special-version")
			Expect(err).NotTo(HaveOccurred())
			Expect(amiMap).To(Equal(map[string]string{
				"ap-northeast-1": "ami-58d24558",
				"ap-southeast-1": "ami-4a2e3b18",
				"ap-southeast-2": "ami-0dd89737",
				"eu-west-1":      "ami-4d8eac3a",
				"sa-east-1":      "ami-3370e52e",
				"us-east-1":      "ami-4f1e6a2a",
				"us-west-1":      "ami-5df23719",
				"us-west-2":      "ami-8b4956bb",
			}))

			Expect(jsonClient.GetCall.Args.Route).To(Equal("/api/v1/box/someuser/somebox"))
		})

		Context("when the desired version isn't available", func() {
			It("should return a useful error message", func() {
				_, err := atlasClient.GetAMIs("someuser/somebox", "non-existent-version")
				Expect(err).To(MatchError(`unable to find box "someuser/somebox" version "non-existent-version"`))
			})
		})

		Context("when downloading the gzipped box fails", func() {
			It("should return a useful error", func() {
				fakeDownloadServer.HTTPTestServer.Close()
				_, err := atlasClient.GetAMIs("someuser/somebox", "some-special-version")
				Expect(err).To(MatchError(HavePrefix("error downloading box: Get")))
			})
		})

		Context("when unzipping the box fails", func() {
			It("should return a useful error", func() {
				fakeDownloadServer.SetHandler(0,
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, []byte("not gzipped")),
					),
				)
				_, err := atlasClient.GetAMIs("someuser/somebox", "some-special-version")
				Expect(err).To(MatchError("error gunzipping box: gzip: invalid header"))
			})
		})

		Context("when the box doesn't include AMI definitions", func() {
			It("should return a useful error", func() {
				gzippedMsg, _ := base64.StdEncoding.DecodeString("H4sIADxwvVYAA8vI5AIAenpv7QMAAAA=")
				fakeDownloadServer.SetHandler(0,
					ghttp.CombineHandlers(
						ghttp.RespondWith(http.StatusOK, gzippedMsg),
					),
				)
				_, err := atlasClient.GetAMIs("someuser/somebox", "some-special-version")
				Expect(err).To(MatchError("no AMIs found within box"))
			})
		})
	})

	Describe("#GetLatestVersion", func() {
		It("should return the most recent box version", func() {
			jsonClient.GetCall.ResponseJSON = `
{
  "current_version" : {
    "version": "9000.92.0",
    "status": "active",
    "created_at": "2016-02-04T05:37:57.213Z",
    "updated_at": "2016-02-10T00:00:13.580Z",
    "number": "9000.92.0",
    "release_url": "https://atlas.hashicorp.com/api/v1/box/cloudfoundry/bosh-lite/version/9000.92.0/release",
    "revoke_url": "https://atlas.hashicorp.com/api/v1/box/cloudfoundry/bosh-lite/version/9000.92.0/revoke"
  }
}`
			version, err := atlasClient.GetLatestVersion("someuser/somebox")
			Expect(err).NotTo(HaveOccurred())
			Expect(version).To(Equal("9000.92.0"))

			Expect(jsonClient.GetCall.Args.Route).To(Equal("/api/v1/box/someuser/somebox"))
		})
	})
})
