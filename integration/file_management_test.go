//go:build integration

package integration_test

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikolalohinski/free-go/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("file management scenarios", func() {
	var ctx context.Context
	BeforeEach(func() {
		ctx = context.Background()

		freeboxClient = freeboxClient.WithAppID(appID).WithPrivateToken(token)

		permissions := Must(freeboxClient.Login(ctx))
		if !permissions.Settings {
			panic(fmt.Sprintf("the token for the '%s' app does not appear to have the permissions to modify freebox settings", appID))
		}
	})

	Context("full lifecycle of a file download", func() {
		It("should not return an error nor unexpected responses", func() {
			// download file
			filename := fmt.Sprintf("free-go.integration.tests.%s", uuid.New().String())
			identifier, err := freeboxClient.AddDownloadTask(ctx, types.DownloadRequest{
				DownloadURLs: []string{
					"https://raw.githubusercontent.com/NikolaLohinski/free-go/main/free-go.svg",
				},
				Filename:          filename,
				DownloadDirectory: "Freebox/Téléchargements",
			})
			Expect(err).To(BeNil())

			// Watch task
			task := new(types.DownloadTask)
			Eventually(func() types.DownloadTask {
				*task, err = freeboxClient.GetDownloadTask(ctx, identifier)
				Expect(err).To(BeNil())

				return *task
			}, "5s").Should(MatchFields(IgnoreExtras, Fields{
				"Status": BeEquivalentTo(types.DownloadTaskStatusDone),
			}))

			// read file info
			info, err := freeboxClient.GetFileInfo(ctx, string(task.DownloadDirectory)+"/"+filename)
			Expect(err).To(BeNil())
			Expect(info).To(MatchFields(IgnoreExtras, Fields{
				"Type": Equal(types.FileTypeFile),
			}))

			// delete
			removeTask, err := freeboxClient.RemoveFiles(ctx, []string{string(info.Path)})
			Expect(err).To(BeNil())
			Expect(removeTask).To(MatchFields(IgnoreExtras, Fields{
				"Type": Equal(types.FileTaskTypeRemove),
			}))
		})
	})
})
