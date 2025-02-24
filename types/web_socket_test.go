package types_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nikolalohinski/free-go/types"
)

var _ = Describe("WebSocket", func() {
	Describe("WebSocketAction", func() {
		Describe("Is", func() {
			Context("when the action is the same", func() {
				It("should return true", func() {
					Expect(types.WebSocketAction("lan_host_l3addr_reachable").Is(&types.EventDescription{
						Source: types.EventSourceLANHost,
						Name:   types.EventHostL3AddrReachable,
					})).To(BeTrue())
				})
			})

			Context("when the action is different", func() {
				It("should return false", func() {
					Expect(types.WebSocketAction("lan_host_l3addr_reachable").Is(&types.EventDescription{
						Source: types.EventSourceLANHost,
						Name:   types.EventHostL3AddrUnreachable,
					})).To(BeFalse())
				})
			})
		})
	})

	Describe("WebSocketResponse", func() {
		Describe("GetError", func() {
			Context("when the response is a success", func() {
				It("should return nil", func() {
					Expect((&types.WebSocketResponse[interface{}]{
						Success: true,
					}).GetError()).To(BeNil())
				})
			})

			Context("when the response is a failure", func() {
				It("should return the error", func() {
					Expect((&types.WebSocketResponse[interface{}]{
						Success: false,
						ErrorCode: `error_code`,
						Message: `un message d'erreur`,
					}).GetError()).To(MatchError(And(ContainSubstring(`error_code`), ContainSubstring(`un message d'erreur`))))
				})
			})
		})
	})
})
