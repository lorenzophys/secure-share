package memory_store_test

import (
	memoryStore "github.com/lorenzophys/secure_share/internal/store/in-memory"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("MemoryStore", func() {
	var ms *memoryStore.MemoryStore

	BeforeEach(func() {
		ms = memoryStore.NewMemoryStore()
	})

	Describe("Set and Get functions", func() {
		It("should store and retrieve data correctly", func() {
			value := "secret data"
			urlKey := ms.Set(value)

			Expect(urlKey).NotTo(BeEmpty())

			retrievedValue, ok := ms.Get(urlKey)

			Expect(ok).To(BeTrue())
			Expect(retrievedValue).NotTo(BeNil())
			Expect(retrievedValue).To(Equal(value))
		})

		It("should return false for non-existent keys", func() {
			_, ok := ms.Get("non-existent-key")
			Expect(ok).To(BeFalse())
		})
	})

	Describe("Get with invalid key", func() {
		It("should return false for invalid keys", func() {
			_, ok := ms.Get("invalid-key")
			Expect(ok).To(BeFalse())
		})
	})

	Describe("Get a secret twice", func() {
		It("should delete the secret once has been retrieved", func() {
			value := "secret data"
			urlKey := ms.Set(value)

			retrievedValue, _ := ms.Get(urlKey)
			Expect(retrievedValue).NotTo(BeNil())

			retrievedValueAgain, ok := ms.Get(urlKey)
			Expect(ok).To(BeFalse())
			Expect(retrievedValueAgain).To(BeZero())
		})
	})
})
