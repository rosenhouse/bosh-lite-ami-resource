package lib_test

import (
	"errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"github.com/rosenhouse/bosh-lite-ami-resource/lib"
	"github.com/rosenhouse/bosh-lite-ami-resource/mocks"
)

var _ = Describe("Resource", func() {
	var resource *lib.Resource
	var atlasClient *mocks.AtlasClient

	BeforeEach(func() {
		atlasClient = &mocks.AtlasClient{}
		resource = &lib.Resource{
			AtlasClient: atlasClient,
			SourceConfig: lib.Source{
				BoxName: "some-box-name",
				Region:  "some-region",
			},
		}
		atlasClient.GetLatestVersionCall.Returns.Version = "1.2.3"
	})

	Describe("In", func() {
		XIt("should return the AMI for the given version and region", func() {
			// ami, err := resource.In(lib.Version{"3.4.5"})
			// Expect(err).NotTo(HaveOccurred())

			// Expect(atlasClient.GetAMICall.Receives.BoxName).To(Equal("some-box-name"))
			// Expect(atlasClient.GetAMICall.Receives.Version).To(Equal("3.4.5"))
			// Expect(atlasClient.GetAMICall.Receives.Region).To(Equal("some-region"))
		})
	})

	Describe("Check", func() {
		DescribeTable("version comparison",
			func(oldVersion, currentVersion lib.Version, expectedResult []lib.Version) {
				atlasClient.GetLatestVersionCall.Returns.Version = currentVersion.BoxVersion
				versionList, err := resource.Check(oldVersion)
				Expect(err).NotTo(HaveOccurred())

				Expect(versionList).To(Equal(expectedResult))
			},
			Entry("old is empty", lib.Version{}, lib.Version{"1.2.3"}, []lib.Version{{"1.2.3"}}),
			Entry("old is older (major)", lib.Version{"0.3.5"}, lib.Version{"1.2.3"}, []lib.Version{{"1.2.3"}}),
			Entry("old is older (minor)", lib.Version{"1.1.6"}, lib.Version{"1.2.3"}, []lib.Version{{"1.2.3"}}),
			Entry("old is older (patch)", lib.Version{"1.2.2"}, lib.Version{"1.2.3"}, []lib.Version{{"1.2.3"}}),
			Entry("old is equal", lib.Version{"1.2.3"}, lib.Version{"1.2.3"}, []lib.Version{}),
			Entry("old is newer (major)", lib.Version{"2.1.0"}, lib.Version{"1.2.3"}, []lib.Version{{"1.2.3"}}),
			Entry("old is newer (minor)", lib.Version{"1.3.0"}, lib.Version{"1.2.3"}, []lib.Version{{"1.2.3"}}),
			Entry("old is newer (patch)", lib.Version{"1.2.4"}, lib.Version{"1.2.3"}, []lib.Version{{"1.2.3"}}),
		)

		It("should check for the latest version", func() {
			_, err := resource.Check(lib.Version{})
			Expect(err).NotTo(HaveOccurred())

			Expect(atlasClient.GetLatestVersionCall.Receives.BoxName).To(Equal("some-box-name"))
		})

		Context("when checking for the latest version errors", func() {
			It("should wrap and return the error", func() {
				atlasClient.GetLatestVersionCall.Returns.Error = errors.New("some error")

				_, err := resource.Check(lib.Version{})
				Expect(err).To(MatchError("atlas client: some error"))
			})
		})
	})
})
