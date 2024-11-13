package client_test

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/nikolalohinski/free-go/client"
)

var _ = Describe("APIError", func() {
	var err *client.APIError

	BeforeEach(func() {
		err = &client.APIError{}
	})

	Context("when Code is empty", func() {
		Context("when Message is not empty", func() {
			BeforeEach(func() {
				err.Message = "message"
			})

			It("should return a string without code", func() {
				Expect(err.Error()).To(Equal(`failed with error code "": message`))
			})
		})

		Context("when Message is empty", func() {
			It("should return a string without code and message", func() {
				Expect(err.Error()).To(Equal(`failed with error code ""`))
			})
		})
	})

	Context("when Code is not empty", func() {
		BeforeEach(func() {
			err.Code = "code"
		})

		Context("when Message is empty", func() {
			It("should return a string without message", func() {
				Expect(err.Error()).To(Equal(`failed with error code "code"`))
			})
		})

		Context("when Message is not empty", func() {
			BeforeEach(func() {
				err.Message = "message"
			})

			It("should return a string without message", func() {
				Expect(err.Error()).To(Equal(`failed with error code "code": message`))
			})

			Describe("errors.IS", func() {
				var target *client.APIError

				BeforeEach(func() {
					target = &client.APIError{}
				})

				Context("when target has the same code", func() {
					BeforeEach(func() {
						target.Code = err.Code
					})

					It("should return true for APIError", func() {
						Expect(errors.Is(err, target)).To(BeTrue())
					})
				})

				Context("when target has a different code", func() {
					BeforeEach(func() {
						target.Code = "not_" + err.Code
					})

					It("should return false for APIError", func() {
						Expect(errors.Is(err, target)).To(BeFalse())
					})
				})
			})

			Describe("APIError.As", func() {
				var target *client.APIError

				BeforeEach(func() {
					target = &client.APIError{}
				})

				It("should return true for APIError", func() {
					Expect(errors.As(err, &target)).To(BeTrue())
					Expect(target.Code).To(Equal(err.Code))
				})
			})
		})
	})
})
