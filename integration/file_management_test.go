//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"

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

		workingDirectory = new(string)
	)
	BeforeAll(func() {
		ctx = context.Background()

		freeboxClient = freeboxClient.WithAppID(appID).WithPrivateToken(token)

		permissions := MustReturn(freeboxClient.Login(ctx))
		Expect(permissions.Settings).To(BeTrue(), fmt.Sprintf("the token for the '%s' app does not appear to have the permissions to modify freebox settings", appID))

		_, err := freeboxClient.CreateDirectory(ctx, root, "Logiciels")
		Expect(err).To(Or(BeNil(), Equal(client.ErrDestinationConflict)))

		*workingDirectory = MustReturn(freeboxClient.CreateDirectory(ctx, root+"/Logiciels", fmt.Sprintf("free-go.integration.tests.%s", uuid.New().String())))
	})
	AfterAll(func() {
		task := MustReturn(freeboxClient.RemoveFiles(ctx, []string{*workingDirectory}))

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
				DownloadDirectory: *workingDirectory,
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

	Context("uploading and extracting a ZIP file", func() {
		It("should not return an error nor unexpected responses", func() {
			_, currentFile, _, _ := runtime.Caller(0)
			path := path.Join(filepath.Dir(currentFile), "test_data", "archive.zip")
			size := int(MustReturn(os.Stat(path)).Size())

			filename := fmt.Sprintf("free-go.integration.tests.%s.zip", uuid.New().String())
			writer, _, err := freeboxClient.FileUploadStart(ctx, types.FileUploadStartActionInput{
				Filename: filename,
				Dirname:  types.Base64Path(*workingDirectory),
				Force:    types.FileUploadStartActionForceOverwrite,
				Size:     size,
			})
			Expect(err).To(BeNil())
			Expect(writer).ToNot(BeNil())

			Expect(writer.Write(MustReturn(ioutil.ReadFile(path)))).ToNot(BeZero())
			Expect(writer.Close()).To(Succeed())

			info, err := freeboxClient.GetFileInfo(ctx, *workingDirectory+"/"+filename)
			Expect(err).To(BeNil())
			Expect(info).To(MatchFields(IgnoreExtras, Fields{
				"Type": Equal(types.FileTypeFile),
			}))

			extractTask, err := freeboxClient.ExtractFile(ctx, types.ExtractFilePayload{
				Src:           types.Base64Path(*workingDirectory + "/" + filename),
				Dst:           types.Base64Path(*workingDirectory),
				Password:      "free-go",
				DeleteArchive: false,
				Overwrite:     false,
			})
			Expect(err).To(BeNil())
			Eventually(func() types.FileSystemTask {
				extractTask, err = freeboxClient.GetFileSystemTask(ctx, extractTask.ID)
				Expect(err).To(BeNil())

				return extractTask
			}, "2s").Should(MatchFields(IgnoreExtras, Fields{
				"State": BeEquivalentTo(types.FileTaskStateDone),
				"Error": BeEquivalentTo(types.DiskTaskErrorNone),
			}))
		})
	})
})
