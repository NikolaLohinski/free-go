package types_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nikolalohinski/free-go/types"
)

var _ = Describe("Events", func() {
	Describe("EventDescription", func() {
		Describe("String", func() {
			It("should return the event description", func() {
				Expect((&types.EventDescription{
					Source: types.EventSourceVM,
					Name:   types.EventStateChanged,
				}).String()).To(Equal("vm_state_changed"))
			})
		})
	})

	Describe("WebSocketNotification", func() {
		Describe("VmStateChange", func() {
			Context("when the action is not a VmStateChange", func() {
				It("should return an error", func() {
					_, err := (&types.WebSocketNotification{
						Action:  "unexpected_action",
						Success: true,
						Source:  types.EventSourceVM,
						Event:   types.EventStateChanged,
						Result:  json.RawMessage(`{}`),
					}).VmStateChange()

					Expect(err).To(MatchError("unexpected action: unexpected_action"))
				})
			})

			Context("when the result cannot be unmarshalled", func() {
				It("should return an error", func() {
					_, err := (&types.WebSocketNotification{
						Action:  "vm_state_changed",
						Success: true,
						Source:  types.EventSourceVM,
						Event:   types.EventStateChanged,
						Result:  json.RawMessage(`{`),
					}).VmStateChange()

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when the action is a VmStateChange", func() {
				It("should return the VmStateChange", func() {
					eventResult, err := json.Marshal(&types.VmStateChange{
						ID: 123,
					})
					Expect(err).ToNot(HaveOccurred())

					notification := &types.WebSocketNotification{
						Action:  "vm_state_changed",
						Success: true,
						Source:  types.EventSourceVM,
						Event:   types.EventStateChanged,
						Result:  json.RawMessage(eventResult),
					}
					Expect(notification.VmStateChange()).To(Equal(&types.VmStateChange{
						ID: 123,
					}))
					_, err = notification.VmDiskTask()
					Expect(err).To(MatchError("unexpected action: vm_state_changed"))
					_, err = notification.LanHost()
					Expect(err).To(MatchError("unexpected action: vm_state_changed"))
				})
			})
		})

		Describe("VmDiskTask", func() {
			Context("when the action is not a VmDiskTask", func() {
				It("should return an error", func() {
					_, err := (&types.WebSocketNotification{
						Action:  "unexpected_action",
						Success: true,
						Source:  types.EventSourceVMDisk,
						Event:   types.EventDiskTaskDone,
						Result:  json.RawMessage(`{}`),
					}).VmDiskTask()

					Expect(err).To(MatchError("unexpected action: unexpected_action"))
				})
			})

			Context("when the result cannot be unmarshalled", func() {
				It("should return an error", func() {
					_, err := (&types.WebSocketNotification{
						Action:  "vm_disk_task_done",
						Success: true,
						Source:  types.EventSourceVMDisk,
						Event:   types.EventDiskTaskDone,
						Result:  json.RawMessage(`{`),
					}).VmDiskTask()

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when the action is a VmDiskTask", func() {
				It("should return the VmDiskTask", func() {
					eventResult, err := json.Marshal(&types.VmDiskTask{
						ID: 123,
					})
					Expect(err).ToNot(HaveOccurred())

					notification := &types.WebSocketNotification{
						Action:  "vm_disk_task_done",
						Success: true,
						Source:  types.EventSourceVMDisk,
						Event:   types.EventDiskTaskDone,
						Result:  json.RawMessage(eventResult),
					}
					Expect(notification.VmDiskTask()).To(Equal(&types.VmDiskTask{
						ID: 123,
					}))
					_, err = notification.VmStateChange()
					Expect(err).To(MatchError("unexpected action: vm_disk_task_done"))
					_, err = notification.LanHost()
					Expect(err).To(MatchError("unexpected action: vm_disk_task_done"))
				})
			})
		})

		Describe("LanHost", func() {
			Context("when the action is not a LanHost", func() {
				It("should return an error", func() {
					_, err := (&types.WebSocketNotification{
						Action:  "unexpected_action",
						Success: true,
						Source:  types.EventSourceLANHost,
						Event:   types.EventHostL3AddrReachable,
						Result:  json.RawMessage(`{}`),
					}).LanHost()

					Expect(err).To(MatchError("unexpected action: unexpected_action"))
				})
			})

			Context("when the result cannot be unmarshalled", func() {
				It("should return an error", func() {
					_, err := (&types.WebSocketNotification{
						Action:  "lan_host_l3addr_reachable",
						Success: true,
						Source:  types.EventSourceLANHost,
						Event:   types.EventHostL3AddrReachable,
						Result:  json.RawMessage(`{`),
					}).LanHost()

					Expect(err).To(HaveOccurred())
				})
			})

			Context("when the action is a LanHost", func() {
				It("should return the LanHost", func() {
					eventResult, err := json.Marshal(&types.LanHost{
						ID: 123,
					})
					Expect(err).ToNot(HaveOccurred())

					notification := &types.WebSocketNotification{
						Action:  "lan_host_l3addr_reachable",
						Success: true,
						Source:  types.EventSourceLANHost,
						Event:   types.EventHostL3AddrReachable,
						Result:  json.RawMessage(eventResult),
					}
					Expect(notification.LanHost()).To(Equal(&types.LanHost{
						ID: 123,
					}))
					_, err = notification.VmDiskTask()
					Expect(err).To(MatchError("unexpected action: lan_host_l3addr_reachable"))
					_, err = notification.VmStateChange()
					Expect(err).To(MatchError("unexpected action: lan_host_l3addr_reachable"))
				})
			})
		})
	})
})
