package types_test

import (
	"encoding/json"
	"time"

	"github.com/nikolalohinski/free-go/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("common", func() {
	returnedErr := new(error)
	Context("json marshal/unmarshal of integer timestamps", func() {
		Context("when marshaling", func() {
			var bytes []byte
			JustBeforeEach(func() {
				bytes, *returnedErr = json.Marshal(types.Timestamp{
					Time: Must(time.Parse(time.RFC3339, "2022-09-18T07:25:40Z")).(time.Time),
				})
			})
			It("should return the correct bytes", func() {
				Expect(*returnedErr).To(BeNil())
				Expect(string(bytes)).To(Equal("1663485940"))
			})
		})
		Context("when unmarshaling", func() {
			var (
				payload []byte

				timestamp *types.Timestamp
			)

			BeforeEach(func() {
				payload = make([]byte, 0)

				timestamp = new(types.Timestamp)
			})
			JustBeforeEach(func() {
				*returnedErr = json.Unmarshal(payload, timestamp)
			})

			Context("when the timestamp is an epoch integer", func() {
				BeforeEach(func() {
					payload = []byte(`1663485940`)
				})
				It("should return the correct payload", func() {
					Expect(*returnedErr).To(BeNil())
					Expect(*timestamp).To(Equal(types.Timestamp{
						Time: Must(time.Parse(time.RFC3339, "2022-09-18T07:25:40Z")).(time.Time).UTC(),
					}))
				})
			})
			Context("when the timestamp is not an epoch integer", func() {
				BeforeEach(func() {
					payload = []byte(`"foobar"`)
				})
				It("should return an error", func() {
					Expect(*returnedErr).ToNot(BeNil())
				})
			})
		})
	})
})
