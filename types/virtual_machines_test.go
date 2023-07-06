package types_test

import (
	"encoding/json"

	"github.com/nikolalohinski/free-go/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
)

var _ = Describe("virtual machines", func() {
	returnedErr := new(error)
	Context("json unmarshal of bind_usb_ports in the VirtualMachine object", func() {
		var (
			payload []byte

			virtualMachine *types.VirtualMachine
		)
		BeforeEach(func() {
			payload = make([]byte, 0)

			virtualMachine = new(types.VirtualMachine)
		})
		JustBeforeEach(func() {
			*returnedErr = json.Unmarshal(payload, virtualMachine)
		})
		Context("when bind_usb_ports is an empty string", func() {
			BeforeEach(func() {
				payload = []byte(`{
					"bind_usb_ports": ""
				}`)
			})
			It("should return the correct usb port binds", func() {
				Expect(*returnedErr).To(BeNil())
				Expect((*virtualMachine)).To(MatchFields(IgnoreExtras, Fields{
					"VirtualMachinePayload": MatchFields(IgnoreExtras, Fields{
						"BindUSBPorts": BeEmpty(),
					}),
				}))
			})
		})
		Context("when bind_usb_ports is a string with content", func() {
			BeforeEach(func() {
				payload = []byte(`{
					"bind_usb_ports": "error"
				}`)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
		Context("when bind_usb_ports is a a list of strings", func() {
			BeforeEach(func() {
				payload = []byte(`{
					"bind_usb_ports": ["foo", "bar"]
				}`)
			})
			It("should return the correct usb port binds", func() {
				Expect(*returnedErr).To(BeNil())
				Expect((*virtualMachine)).To(MatchFields(IgnoreExtras, Fields{
					"VirtualMachinePayload": MatchFields(IgnoreExtras, Fields{
						"BindUSBPorts": Equal(types.BindUSBPorts{"foo", "bar"}),
					}),
				}))
			})
		})
		Context("when bind_usb_ports is a list of integers", func() {
			BeforeEach(func() {
				payload = []byte(`{
					"bind_usb_ports": [1, 2]
				}`)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
		Context("when bind_usb_ports is neither a list of strings nor a string", func() {
			BeforeEach(func() {
				payload = []byte(`{
					"bind_usb_ports": 1
				}`)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
		Context("when bind_usb_ports is an invalid json", func() {
			JustBeforeEach(func() {
				*returnedErr = (&types.BindUSBPorts{}).UnmarshalJSON([]byte(`{`))
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
	Context("json marshal of a Base64Path", func() {
		var (
			result []byte

			path = types.Base64Path("/tmp/example/")
		)
		JustBeforeEach(func() {
			result, *returnedErr = json.Marshal(path)
		})

		It("should return the correct usb port binds", func() {
			Expect(*returnedErr).To(BeNil())
			Expect(string(result)).To(Equal("\"L3RtcC9leGFtcGxlLw==\""))
		})
	})
	Context("json unmarshal of a Base64Path", func() {
		var (
			payload []byte

			encodedPath *types.Base64Path
		)
		BeforeEach(func() {
			payload = make([]byte, 0)

			encodedPath = new(types.Base64Path)
		})
		JustBeforeEach(func() {
			*returnedErr = json.Unmarshal(payload, encodedPath)
		})
		Context("when cd_path is a base64 encoded string", func() {
			BeforeEach(func() {
				payload = []byte(`"L0ZyZWVib3gvcGF0aC90by9pbWFnZQ=="`)
			})
			It("should return the correct object", func() {
				Expect(*returnedErr).To(BeNil())
				Expect((*encodedPath)).To(Equal(types.Base64Path("/Freebox/path/to/image")))
			})
		})
		Context("when cd_path is a string but is not base64 encoded", func() {
			BeforeEach(func() {
				payload = []byte(`"\nà@"`)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
		Context("when cd_path is not a string", func() {
			BeforeEach(func() {
				payload = []byte(`123`)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
})
