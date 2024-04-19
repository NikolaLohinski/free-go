//go:build integration

package integration_test

import (
	"context"
	"fmt"
	"path"
	"strconv"

	"github.com/google/uuid"
	"github.com/nikolalohinski/free-go/client"
	"github.com/nikolalohinski/free-go/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("virtual machines", Ordered, func() {
	const (
		rootDirectoryName = "free-go"
		rootDirectory     = rootFS + "/" + rootDirectoryName
		diskImagePath     = rootDirectory + "/" + "free-go.integration.tests.qcow2"
	)
	var (
		ctx context.Context
	)
	BeforeAll(func() {
		ctx = context.Background()

		freeboxClient = freeboxClient.WithAppID(appID).WithPrivateToken(token)
		permissions := Must(freeboxClient.Login(ctx))
		Expect(permissions.Settings).To(BeTrue(), fmt.Sprintf("the token for the '%s' app does not appear to have the permissions to modify freebox settings", appID))

		_, err := freeboxClient.CreateDirectory(ctx, rootFS, rootDirectoryName)
		Expect(err).To(Or(BeNil(), Equal(client.ErrDestinationConflict)))

		// Check that the image exists and if not pre-download it
		_, err = freeboxClient.GetFileInfo(ctx, diskImagePath)
		if err == nil {
			return
		}
		if err != nil && err == client.ErrPathNotFound {
			taskID, err := freeboxClient.AddDownloadTask(ctx, types.DownloadRequest{
				DownloadURLs: []string{
					"https://cloud.debian.org/images/cloud/bullseye/daily/latest/debian-11-generic-arm64-daily.qcow2",
				},
				Hash:              "https://cloud.debian.org/images/cloud/bullseye/daily/latest/SHA512SUMS",
				DownloadDirectory: path.Dir(diskImagePath),
				Filename:          path.Base(diskImagePath),
			})
			Expect(err).To(BeNil())
			Eventually(func() types.DownloadTask {
				downloadTask, err := freeboxClient.GetDownloadTask(ctx, taskID)
				Expect(err).To(BeNil())

				return downloadTask
			}, "5m").Should(MatchFields(IgnoreExtras, Fields{
				"Status": BeEquivalentTo(types.DownloadTaskStatusDone),
			}))
		} else {
			Expect(err).To(BeNil())
		}
	})
	Context("full lifecycle for a virtual machine", func() {
		It("should not return an error nor unexpected responses", func() {
			// create the virtual machine
			virtualMachine, err := freeboxClient.CreateVirtualMachine(ctx, types.VirtualMachinePayload{
				Name:     fmt.Sprintf("free-go.integration.tests.%s", uuid.New().String())[:30],
				DiskPath: types.Base64Path(diskImagePath),
				DiskType: types.QCow2Disk,
				Memory:   512,
				OS:       types.DebianOS,
				VCPUs:    1,
			})
			Expect(err).To(BeNil())

			// register for state change
			eventStream, err := freeboxClient.ListenEvents(ctx, []types.EventDescription{
				{
					Source: "vm",
					Name:   "state_changed",
				},
			})
			Expect(err).To(BeNil())

			// watch events
			isRunning := false
			go func() {
				defer GinkgoRecover()
				event := <-eventStream
				Expect(event.Error).To(BeNil())
				Expect(event.Notification.Source).To(Equal("vm"))
				Expect(event.Notification.Event).To(Equal("state_changed"))
				Expect(event.Notification.Result).To(MatchJSON(`{
					"id": ` + strconv.Itoa(int(virtualMachine.ID)) + `,
					"status": "` + types.RunningStatus + `"
				}`))
				isRunning = true
			}()

			// start the virtual machine
			Expect(freeboxClient.StartVirtualMachine(ctx, virtualMachine.ID)).To(BeNil())

			// wait for the virtual machine to be running
			Eventually(func() bool { return isRunning }, "30s").Should(BeTrue())

			// watch events
			go func() {
				defer GinkgoRecover()
				event := <-eventStream
				Expect(event.Error).To(BeNil())
				Expect(event.Notification.Source).To(Equal("vm"))
				Expect(event.Notification.Event).To(Equal("state_changed"))
				Expect(event.Notification.Result).To(MatchJSON(`{
					"id": ` + strconv.Itoa(int(virtualMachine.ID)) + `,
					"status": "` + types.StoppedStatus + `"
				}`))
				isRunning = false
			}()

			// stop the virtual machine
			Expect(freeboxClient.StopVirtualMachine(ctx, virtualMachine.ID)).To(BeNil())

			// wait for the virtual machine to be stopped
			Eventually(func() bool { return isRunning }, "30s").Should(BeFalse())

			// delete the virtual machine
			Expect(freeboxClient.DeleteVirtualMachine(ctx, virtualMachine.ID)).To(BeNil())
		})
	})
})
