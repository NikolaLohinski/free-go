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

var _ = Describe("virtual machines", func() {
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
	Context("getting virtual machine info", func() {
		returnedInfo := new(types.VirtualMachinesInfo)
		JustBeforeEach(func() {
			*returnedInfo, *returnedErr = freeboxClient.GetVirtualMachineInfo(context.Background())
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vm/info/", version)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"usb_used": false,
								"sata_used": false,
								"sata_ports": [
									"sata-internal-p0",
									"sata-internal-p1",
									"sata-internal-p2",
									"sata-internal-p3"
								],
								"used_memory": 0,
								"usb_ports": [
									"usb-external-type-a",
									"usb-external-type-c"
								],
								"used_cpus": 0,
								"total_memory": 1024,
								"total_cpus": 2
							}
						}`),
					),
				)
			})
			It("should return the correct virtual machine info", func() {
				Expect(*returnedErr).To(BeNil())
				Expect((*returnedInfo)).To(Equal(types.VirtualMachinesInfo{
					USBUsed:  false,
					SATAUsed: false,
					SATAPorts: []string{
						"sata-internal-p0",
						"sata-internal-p1",
						"sata-internal-p2",
						"sata-internal-p3",
					},
					UsedMemory: 0,
					USBPorts: []string{
						"usb-external-type-a",
						"usb-external-type-c",
					},
					UsedCPUs:    0,
					TotalMemory: 1024,
					TotalCPUs:   2,
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
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vm/info/", version)),
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
	Context("getting virtual machine distributions", func() {
		returnedDistros := new([]types.VirtualMachineDistribution)
		JustBeforeEach(func() {
			*returnedDistros, *returnedErr = freeboxClient.GetVirtualMachineDistributions(context.Background())
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vm/distros/", version)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": [
								{
									"hash": "http://ftp.free.fr/.private/ubuntu-cloud/releases/jammy/release/SHA256SUMS",
									"os": "ubuntu",
									"url": "http://ftp.free.fr/.private/ubuntu-cloud/releases/jammy/release/ubuntu-22.04-server-cloudimg-arm64.img",
									"name": "Ubuntu 22.04 LTS (Jammy)"
								},
								{
									"hash": "http://ftp.free.fr/.private/ubuntu-cloud/releases/impish/release/SHA256SUMS",
									"os": "ubuntu",
									"url": "http://ftp.free.fr/.private/ubuntu-cloud/releases/impish/release/ubuntu-21.10-server-cloudimg-arm64.img",
									"name": "Ubuntu 21.10 (Impish)"
								}
							]
						}`),
					),
				)
			})
			It("should return the correct virtual machine info", func() {
				Expect(*returnedErr).To(BeNil())
				Expect((*returnedDistros)).To(Equal([]types.VirtualMachineDistribution{
					{
						Hash: "http://ftp.free.fr/.private/ubuntu-cloud/releases/jammy/release/SHA256SUMS",
						OS:   "ubuntu",
						URL:  "http://ftp.free.fr/.private/ubuntu-cloud/releases/jammy/release/ubuntu-22.04-server-cloudimg-arm64.img",
						Name: "Ubuntu 22.04 LTS (Jammy)",
					},
					{
						Hash: "http://ftp.free.fr/.private/ubuntu-cloud/releases/impish/release/SHA256SUMS",
						OS:   "ubuntu",
						URL:  "http://ftp.free.fr/.private/ubuntu-cloud/releases/impish/release/ubuntu-21.10-server-cloudimg-arm64.img",
						Name: "Ubuntu 21.10 (Impish)",
					},
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
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vm/distros/", version)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {}
						}`),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
	Context("listing virtual machines", func() {
		returnedMachines := new([]types.VirtualMachine)
		JustBeforeEach(func() {
			*returnedMachines, *returnedErr = freeboxClient.ListVirtualMachines(context.Background())
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vm/", version)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": [
								{
									"mac": "f6:69:9c:d9:4f:3d",
									"cloudinit_userdata": "\n#cloud-config\n\n\nsystem_info:\n  default_user:\n    name: freemind\n",
									"cd_path": "L0ZyZWVib3gvcGF0aC90by9pbWFnZQ==",
									"id": 0,
									"os": "debian",
									"enable_cloudinit": true,
									"disk_path": "L0ZyZWVib3gvZGlzay1wYXRo",
									"vcpus": 1,
									"memory": 300,
									"name": "testing",
									"cloudinit_hostname": "testing",
									"status": "stopped",
									"bind_usb_ports": "",
									"enable_screen": false,
									"disk_type": "qcow2"
								}
							]
						}`),
					),
				)
			})
			It("should return the correct virtual machine", func() {
				Expect(*returnedErr).To(BeNil())
				Expect((*returnedMachines)).To(Equal([]types.VirtualMachine{{
					ID:     0,
					Mac:    "f6:69:9c:d9:4f:3d",
					Status: types.StoppedStatus,
					VirtualMachinePayload: types.VirtualMachinePayload{
						Name:              "testing",
						DiskPath:          "/Freebox/disk-path",
						DiskType:          types.QCow2Disk,
						CDPath:            "/Freebox/path/to/image",
						Memory:            300,
						OS:                types.DebianOS,
						VCPUs:             1,
						EnableScreen:      false,
						BindUSBPorts:      []string{},
						EnableCloudInit:   true,
						CloudInitUserData: "\n#cloud-config\n\n\nsystem_info:\n  default_user:\n    name: freemind\n",
						CloudHostName:     "testing",
					},
				}}))
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
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vm/", version)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {}
						}`),
					),
				)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
	Context("creating a virtual machine", func() {
		var (
			payload         = new(types.VirtualMachinePayload)
			returnedMachine = new(types.VirtualMachine)
		)
		BeforeEach(func() {
			*payload = types.VirtualMachinePayload{
				CloudInitUserData: "\n#cloud-config\n\n\nsystem_info:\n  default_user:\n    name: freemind\n",
				CDPath:            "/Freebox/path/to/image",
				OS:                types.DebianOS,
				EnableCloudInit:   true,
				DiskPath:          "/Freebox/disk-path",
				VCPUs:             1,
				Memory:            300,
				Name:              "testing",
				CloudHostName:     "testing",
				BindUSBPorts:      []string{},
				EnableScreen:      false,
				DiskType:          types.QCow2Disk,
			}
		})
		JustBeforeEach(func() {
			*returnedMachine, *returnedErr = freeboxClient.CreateVirtualMachine(context.Background(), *payload)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/vm/", version)),
						ghttp.VerifyJSON(`{
							"cloudinit_userdata": "\n#cloud-config\n\n\nsystem_info:\n  default_user:\n    name: freemind\n",
							"cd_path": "L0ZyZWVib3gvcGF0aC90by9pbWFnZQ==",
							"os": "debian",
							"enable_cloudinit": true,
							"disk_path": "L0ZyZWVib3gvZGlzay1wYXRo",
							"vcpus": 1,
							"memory": 300,
							"name": "testing",
							"cloudinit_hostname": "testing",
							"disk_type": "qcow2"
						}`),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"mac": "f6:69:9c:d9:4f:3d",
								"cloudinit_userdata": "\n#cloud-config\n\n\nsystem_info:\n  default_user:\n    name: freemind\n",
								"cd_path": "L0ZyZWVib3gvcGF0aC90by9pbWFnZQ==",
								"id": 0,
								"os": "debian",
								"enable_cloudinit": true,
								"disk_path": "L0ZyZWVib3gvZGlzay1wYXRo",
								"vcpus": 1,
								"memory": 300,
								"name": "testing",
								"cloudinit_hostname": "testing",
								"status": "stopped",
								"bind_usb_ports": "",
								"enable_screen": false,
								"disk_type": "qcow2"
							}
						}`),
					),
				)
			})
			It("should return the correct virtual machine", func() {
				Expect(*returnedErr).To(BeNil())
				Expect((*returnedMachine)).To(Equal(types.VirtualMachine{
					ID:     0,
					Mac:    "f6:69:9c:d9:4f:3d",
					Status: types.StoppedStatus,
					VirtualMachinePayload: types.VirtualMachinePayload{
						Name:              "testing",
						DiskPath:          "/Freebox/disk-path",
						DiskType:          types.QCow2Disk,
						CDPath:            "/Freebox/path/to/image",
						Memory:            300,
						OS:                types.DebianOS,
						VCPUs:             1,
						EnableScreen:      false,
						BindUSBPorts:      []string{},
						EnableCloudInit:   true,
						CloudInitUserData: "\n#cloud-config\n\n\nsystem_info:\n  default_user:\n    name: freemind\n",
						CloudHostName:     "testing",
					},
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
		Context("when the name in the payload is too long", func() {
			BeforeEach(func() {
				payload.Name = "this is way more than 30 characters and should fail"
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrVirtualMachineNameTooLong))
			})
		})
		Context("when the server returns an unexpected payload", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/vm/", version)),
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
	Context("getting a virtual machine", func() {
		returnedMachine := new(types.VirtualMachine)
		JustBeforeEach(func() {
			*returnedMachine, *returnedErr = freeboxClient.GetVirtualMachine(context.Background(), 1234)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vm/1234", version)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"mac": "f6:69:9c:d9:4f:3d",
								"cloudinit_userdata": "\n#cloud-config\n\n\nsystem_info:\n  default_user:\n    name: freemind\n",
								"cd_path": "L0ZyZWVib3gvcGF0aC90by9pbWFnZQ==",
								"id": 1234,
								"os": "debian",
								"enable_cloudinit": true,
								"disk_path": "L0ZyZWVib3gvZGlzay1wYXRo",
								"vcpus": 1,
								"memory": 300,
								"name": "testing",
								"cloudinit_hostname": "testing",
								"status": "stopped",
								"bind_usb_ports": "",
								"enable_screen": false,
								"disk_type": "qcow2"
							}
						}`),
					),
				)
			})
			It("should return the correct virtual machine", func() {
				Expect(*returnedErr).To(BeNil())
				Expect((*returnedMachine)).To(Equal(types.VirtualMachine{
					ID:     1234,
					Mac:    "f6:69:9c:d9:4f:3d",
					Status: types.StoppedStatus,
					VirtualMachinePayload: types.VirtualMachinePayload{
						Name:              "testing",
						DiskPath:          "/Freebox/disk-path",
						DiskType:          types.QCow2Disk,
						CDPath:            "/Freebox/path/to/image",
						Memory:            300,
						OS:                types.DebianOS,
						VCPUs:             1,
						EnableScreen:      false,
						BindUSBPorts:      []string{},
						EnableCloudInit:   true,
						CloudInitUserData: "\n#cloud-config\n\n\nsystem_info:\n  default_user:\n    name: freemind\n",
						CloudHostName:     "testing",
					},
				}))
			})
		})
		Context("when the virtual machine does not exist", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vm/1234", version)),
						ghttp.RespondWith(http.StatusOK, `{
						    "msg": "Impossible de supprimer cette VM : La VM n’existe pas",
						    "success": false,
						    "error_code": "no_such_vm"
						}`),
					),
				)
			})
			It("should return the correct virtual machine", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrVirtualMachineNotFound))
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
						ghttp.VerifyRequest(http.MethodGet, fmt.Sprintf("/api/%s/vm/1234", version)),
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
	Context("updating a virtual machine", func() {
		const (
			identifier = int64(1234)
		)
		var (
			payload         = new(types.VirtualMachinePayload)
			returnedMachine = new(types.VirtualMachine)
		)
		BeforeEach(func() {
			*payload = types.VirtualMachinePayload{
				DiskPath: "/Freebox/disk-path",
			}
		})
		JustBeforeEach(func() {
			*returnedMachine, *returnedErr = freeboxClient.UpdateVirtualMachine(context.Background(), identifier, *payload)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/vm/%d", version, identifier)),
						ghttp.VerifyJSON(`{
							"disk_path": "L0ZyZWVib3gvZGlzay1wYXRo"
						}`),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true,
							"result": {
								"mac": "f6:69:9c:d9:4f:3d",
								"cloudinit_userdata": "\n#cloud-config\n\n\nsystem_info:\n  default_user:\n    name: freemind\n",
								"cd_path": "L0ZyZWVib3gvcGF0aC90by9pbWFnZQ==",
								"id": 0,
								"os": "debian",
								"enable_cloudinit": true,
								"disk_path": "L0ZyZWVib3gvZGlzay1wYXRo",
								"vcpus": 1,
								"memory": 300,
								"name": "testing",
								"cloudinit_hostname": "testing",
								"status": "stopped",
								"bind_usb_ports": "",
								"enable_screen": false,
								"disk_type": "qcow2"
							}
						}`),
					),
				)
			})
			It("should return the correct virtual machine", func() {
				Expect(*returnedErr).To(BeNil())
				Expect((*returnedMachine)).To(Equal(types.VirtualMachine{
					ID:     0,
					Mac:    "f6:69:9c:d9:4f:3d",
					Status: types.StoppedStatus,
					VirtualMachinePayload: types.VirtualMachinePayload{
						Name:              "testing",
						DiskPath:          "/Freebox/disk-path",
						DiskType:          types.QCow2Disk,
						CDPath:            "/Freebox/path/to/image",
						Memory:            300,
						OS:                types.DebianOS,
						VCPUs:             1,
						EnableScreen:      false,
						BindUSBPorts:      []string{},
						EnableCloudInit:   true,
						CloudInitUserData: "\n#cloud-config\n\n\nsystem_info:\n  default_user:\n    name: freemind\n",
						CloudHostName:     "testing",
					},
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
						ghttp.VerifyRequest(http.MethodPut, fmt.Sprintf("/api/%s/vm/%d", version, identifier)),
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
	Context("deleting a virtual machine", func() {
		JustBeforeEach(func() {
			*returnedErr = freeboxClient.DeleteVirtualMachine(context.Background(), 1234)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, fmt.Sprintf("/api/%s/vm/1234", version)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true
						}`),
					),
				)
			})
			It("should not return an error", func() {
				Expect(*returnedErr).To(BeNil())
			})
		})
		Context("when the virtual machine does not exist", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodDelete, fmt.Sprintf("/api/%s/vm/1234", version)),
						ghttp.RespondWith(http.StatusOK, `{
						    "msg": "Impossible de supprimer cette VM : La VM n’existe pas",
						    "success": false,
						    "error_code": "no_such_vm"
						}`),
					),
				)
			})
			It("should return the correct virtual machine", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrVirtualMachineNotFound))
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
	})
	Context("starting a virtual machine", func() {
		JustBeforeEach(func() {
			*returnedErr = freeboxClient.StartVirtualMachine(context.Background(), 1234)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/vm/1234/start", version)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true
						}`),
					),
				)
			})
			It("should not return an error", func() {
				Expect(*returnedErr).To(BeNil())
			})
		})
		Context("when the virtual machine does not exist", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/vm/1234/start", version)),
						ghttp.RespondWith(http.StatusOK, `{
						    "msg": "Impossible de supprimer cette VM : La VM n’existe pas",
						    "success": false,
						    "error_code": "no_such_vm"
						}`),
					),
				)
			})
			It("should return the correct virtual machine", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrVirtualMachineNotFound))
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
	})
	Context("killing a virtual machine", func() {
		JustBeforeEach(func() {
			*returnedErr = freeboxClient.KillVirtualMachine(context.Background(), 1234)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/vm/1234/stop", version)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true
						}`),
					),
				)
			})
			It("should not return an error", func() {
				Expect(*returnedErr).To(BeNil())
			})
		})
		Context("when the virtual machine does not exist", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/vm/1234/stop", version)),
						ghttp.RespondWith(http.StatusOK, `{
						    "msg": "Impossible de supprimer cette VM : La VM n’existe pas",
						    "success": false,
						    "error_code": "no_such_vm"
						}`),
					),
				)
			})
			It("should return the correct virtual machine", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrVirtualMachineNotFound))
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
	})
	Context("stopping a virtual machine", func() {
		JustBeforeEach(func() {
			*returnedErr = freeboxClient.StopVirtualMachine(context.Background(), 1234)
		})
		Context("default", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/vm/1234/powerbutton", version)),
						ghttp.RespondWith(http.StatusOK, `{
							"success": true
						}`),
					),
				)
			})
			It("should not return an error", func() {
				Expect(*returnedErr).To(BeNil())
			})
		})
		Context("when the virtual machine does not exist", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest(http.MethodPost, fmt.Sprintf("/api/%s/vm/1234/powerbutton", version)),
						ghttp.RespondWith(http.StatusOK, `{
						    "msg": "Impossible de supprimer cette VM : La VM n’existe pas",
						    "success": false,
						    "error_code": "no_such_vm"
						}`),
					),
				)
			})
			It("should return the correct virtual machine", func() {
				Expect(*returnedErr).ToNot(BeNil())
				Expect(*returnedErr).To(Equal(client.ErrVirtualMachineNotFound))
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
	})
})
