package hash_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Hash", Label("Infrastructure"), func() {
	Context("HashPassword", func() {
		It("should return hash password", func() {
			password := "password"
			hashedPassword, err := hashInstance.HashPassword(ctx, password)

			Expect(err).ToNot(HaveOccurred())
			Expect(hashedPassword).ToNot(BeNil())
		})
	})
	Context("ComparePassword", func() {
		It("should return nil when password is valid", func() {
			password := "password"
			hashedPassword, _ := hashInstance.HashPassword(ctx, password)

			err := hashInstance.ComparePassword(ctx, password, *hashedPassword)

			Expect(err).ToNot(HaveOccurred())
		})
		It("should return error when password is invalid", func() {
			password := "password"
			hashedPassword, _ := hashInstance.HashPassword(ctx, password)

			err := hashInstance.ComparePassword(ctx, "invalid_password", *hashedPassword)

			Expect(err).To(HaveOccurred())
		})
	})
})
