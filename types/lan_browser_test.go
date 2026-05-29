package types_test

import (
	"encoding/json"

	"github.com/nikolalohinski/free-go/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LanInterfaceHost", func() {
	returnedErr := new(error)

	Context("json unmarshal of primary_name_manual", func() {
		var (
			payload []byte
			host    *types.LanInterfaceHost
		)

		BeforeEach(func() {
			host = new(types.LanInterfaceHost)
		})

		JustBeforeEach(func() {
			*returnedErr = json.Unmarshal(payload, host)
		})

		Context("when primary_name_manual is false (bool)", func() {
			BeforeEach(func() {
				payload = []byte(`{"primary_name_manual": false}`)
			})
			It("should set PrimaryNameManual to false with no error", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(host.PrimaryNameManual).To(BeFalse())
			})
		})

		Context("when primary_name_manual is true (bool)", func() {
			BeforeEach(func() {
				payload = []byte(`{"primary_name_manual": true}`)
			})
			It("should set PrimaryNameManual to true with no error", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(host.PrimaryNameManual).To(BeTrue())
			})
		})

		Context("when primary_name_manual is an empty string", func() {
			BeforeEach(func() {
				payload = []byte(`{"primary_name_manual": ""}`)
			})
			It("should set PrimaryNameManual to false with no error", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(host.PrimaryNameManual).To(BeFalse())
			})
		})

		Context("when primary_name_manual is a non-empty string", func() {
			BeforeEach(func() {
				payload = []byte(`{"primary_name_manual": "unexpected"}`)
			})
			It("should return an error", func() {
				Expect(*returnedErr).ToNot(BeNil())
			})
		})
	})
})
