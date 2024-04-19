//go:build integration

package integration_test

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/nikolalohinski/free-go/client"
	"github.com/nikolalohinski/free-go/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("file management scenarios", Ordered, func() {
	var (
		ctx context.Context

		rootDirectory = new(string)
	)
	BeforeAll(func() {
		ctx = context.Background()

		freeboxClient = freeboxClient.WithAppID(appID).WithPrivateToken(token)

		permissions := Must(freeboxClient.Login(ctx))
		Expect(permissions.Settings).To(BeTrue(), fmt.Sprintf("the token for the '%s' app does not appear to have the permissions to modify freebox settings", appID))

		*rootDirectory = Must(freeboxClient.CreateDirectory(ctx, rootFS, fmt.Sprintf("free-go.integration.tests.%s", uuid.New().String())))
	})
	AfterAll(func() {
		task := Must(freeboxClient.RemoveFiles(ctx, []string{*rootDirectory}))

		Eventually(func() interface{} {
			task, err := freeboxClient.GetFileSystemTask(ctx, task.ID)
			Expect(err).To(BeNil())

			return task.State
		}, "30s").Should(Equal(types.FileTaskStateDone))
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
				DownloadDirectory: *rootDirectory,
			})
			Expect(err).To(BeNil())

			// watch download task
			task := new(types.DownloadTask)
			Eventually(func() types.DownloadTask {
				*task, err = freeboxClient.GetDownloadTask(ctx, identifier)
				Expect(err).To(BeNil())

				return *task
			}, "2s").Should(MatchFields(IgnoreExtras, Fields{
				"Status": BeEquivalentTo(types.DownloadTaskStatusDone),
			}))

			// read file info
			info, err := freeboxClient.GetFileInfo(ctx, string(task.DownloadDirectory)+"/"+filename)
			Expect(err).To(BeNil())
			Expect(info).To(MatchFields(IgnoreExtras, Fields{
				"Type": Equal(types.FileTypeFile),
			}))

			// remove file
			removeTask, err := freeboxClient.RemoveFiles(ctx, []string{string(info.Path)})
			Expect(err).To(BeNil())
			Expect(removeTask).To(MatchFields(IgnoreExtras, Fields{
				"Type": Equal(types.FileTaskTypeRemove),
			}))

			// watch remove task
			Eventually(func() types.FileSystemTask {
				removeTask, err = freeboxClient.GetFileSystemTask(ctx, removeTask.ID)
				Expect(err).To(BeNil())

				return removeTask
			}, "2s").Should(MatchFields(IgnoreExtras, Fields{
				"State": BeEquivalentTo(types.FileTaskStateDone),
			}))

			// confirm file is removed
			_, err = freeboxClient.GetFileInfo(ctx, string(task.DownloadDirectory)+"/"+filename)
			Expect(err).ToNot(BeNil())
			Expect(err).To(Equal(client.ErrPathNotFound))
		})
	})
})
