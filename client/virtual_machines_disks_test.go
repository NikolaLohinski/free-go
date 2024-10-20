package client_test

import (
	"context"
	"fmt"
	"net/http"

	//
	"github.com/nikolalohinski/free-go/client"
	"github.com/nikolalohinski/free-go/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("virtual machines disks", func() {
	var (
		freeboxClient client.Client

		server   *ghttp.Server
		endpoint = new(string)

		sessionToken = new(string)

		returnedErr = new(error)
	)
	BeforeEach(func() {
		server = ghttp.NewServer()
		*endpoint = server.Addr()

		freeboxClient = Must(client.New(*endpoint, version)).
			WithAppID(appID).
			WithPrivateToken(privateToken)

		*sessionToken = setupLoginFlow(server)
	})
	AfterEach(func() {
		server.Close()
	})
	Context("creating a virtual disk", func() {
		var (
			payload         = new(types.VirtualDisksCreatePayload)
			returnedID = new(int64)
		)
		BeforeEach(func() {
			*payload = types.VirtualDisksCreatePayload{
				DiskType: types.QCow2Disk,
				DiskPath: "/Freebox/disk-path",
				Size: 	  1000000000,
			}
		})
		JustBeforeEach(func() {
			*returnedID, *returnedErr = freeboxClient.CreateVirtualDisk(context.Background(), *payload)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/vm/disk/create/", version)),
						ghttp.VerifyJSON(`{
							"disk_type": "qcow2",
							"disk_path": "L0ZyZWVib3gvZGlzay1wYXRo",
							"size": 1000000000
						}`),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"id": 1234
							}
						}`),
					),
				)
			})
			It("should create a new disk", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(*returnedID).To(BeEquivalentTo(1234))
			})
		})
		Context("when server fails to respond", func() {
			BeforeEach(func() {
				server.Close()
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
		Context("when the size is invalid", func() {
			BeforeEach(func() {
				payload.Size = -1
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrVMDiskSizeInvalid))
			})
		})
		Context("when the server returns an unexpected payload", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/vm/disk/create/", version)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": []
						}`),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
	Context("getting virtual disk info", func() {
		returnedInfo := new(types.VirtualDiskInfo)
		JustBeforeEach(func() {
			*returnedInfo, *returnedErr = freeboxClient.GetVirtualDiskInfo(context.Background(), "/Freebox/disk-path")
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/vm/disk/info/", version)),
						ghttp.VerifyJSON(`{
							"disk_path": "L0ZyZWVib3gvZGlzay1wYXRo"
						}`),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"type": "qcow2",
								"actual_size": 2000000000,
								"virtual_size": 5000000000
							}
						}`),
					),
				)
			})
			It("should return the correct virtual machine info", func() {
				Expect(*returnedErr).To(BeNil())
				Expect((*returnedInfo)).To(Equal(types.VirtualDiskInfo{
					Type: 	     types.QCow2Disk,
					ActualSize:  2000000000,
					VirtualSize: 5000000000,
				},
				))
			})
		})
		Context("when server fails to respond", func() {
			BeforeEach(func() {
				server.Close()
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
		Context("when the server returns an unexpected payload", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/vm/disk/info/", version)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": []
						}`),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
	Context("resizing a virtual disk", func() {
		var (
			payload         = new(types.VirtualDisksResizePayload)
			returnedTaskIdentifier = new(int64)
		)
		BeforeEach(func() {
			*payload = types.VirtualDisksResizePayload{
				DiskPath: "/Freebox/disk-path",
				NewSize: 2000000000,
				ShrinkAllow: false,
			}
		})
		JustBeforeEach(func() {
			*returnedTaskIdentifier, *returnedErr = freeboxClient.ResizeVirtualDisk(context.Background(), *payload)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/vm/disk/resize/", version)),
						ghttp.VerifyJSON(`{
							"disk_path": "L0ZyZWVib3gvZGlzay1wYXRo",
							"size": 2000000000,
							"shrink_allow": false
						}`),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"id": 1234
							}
						}`),
					),
				)
			})
			It("should return the correct virtual machine", func() {
				Expect(*returnedErr).To(BeNil())
				Expect((*returnedTaskIdentifier)).To(BeEquivalentTo(1234))
			})
		})
		Context("when server fails to respond", func() {
			BeforeEach(func() {
				server.Close()
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
		Context("when the server returns an unexpected payload", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/vm/disk/resize/", version)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": []
						}`),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
	Context("getting a virtual disk task", func() {
		const identifier int64 = 42
		var returnedTask = new(types.VirtualMachineDiskTask)
		JustBeforeEach(func(ctx SpecContext) {
			*returnedTask, *returnedErr = freeboxClient.GetVirtualDiskTask(ctx, identifier)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vm/disk/task/%d", version, identifier)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"id": 42,
								"type": "resize",
								"done": false
							}
						}`),
					),
				)
			})
			It("should return the correct virtual machine", func() {
				Expect(*returnedErr).To(BeNil())
				Expect((*returnedTask)).To(Equal(types.VirtualMachineDiskTask{
					ID: 42,
					Type: "resize",
					Done: false,
					Error: false,
				}))
			})
		})
		Context("when server fails to respond", func() {
			BeforeEach(func() {
				server.Close()
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
		Context("when the server returns an unexpected payload", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vm/disk/task/%d", version, identifier)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": []
						}`),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
	Context("deleting a virtual disk task", func() {
		var (
			identifier = new(int64)
		)
		BeforeEach(func() {
			*identifier = 1234
		})
		JustBeforeEach(func() {
			*returnedErr = freeboxClient.DeleteVirtualDiskTask(context.Background(), *identifier)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, fmt.Sprintf("/api/%s/vm/disk/task/%d", version, *identifier)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"id": 1234
							}
						}`),
					),
				)
			})
			It("should return the correct virtual machine", func() {
				Expect(*returnedErr).To(BeNil())
			})
		})
		Context("when server fails to respond", func() {
			BeforeEach(func() {
				server.Close()
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
		Context("when the server returns an unexpected payload", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, fmt.Sprintf("/api/%s/vm/disk/task/%d", version, *identifier)),
						ghttp.RespondWith(http.StatusNotFound, nil),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
})
