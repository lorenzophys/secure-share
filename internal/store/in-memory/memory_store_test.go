package memory_store_test

import (
	"time"

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
			urlKey := ms.Set(value, 0)

			Expect(urlKey).NotTo(BeEmpty())

			retrievedValue, ok := ms.Get(urlKey)

			Expect(ok).To(BeTrue())
			Expect(retrievedValue).NotTo(BeNil())
			Expect(retrievedValue).To(Equal(value))
		})
	})

	Describe("Get with invalid key", func() {
		It("should return false for invalid keys", func() {
			secret, ok := ms.Get("invalid-key")
			Expect(ok).To(BeFalse())
			Expect(secret).To(BeEmpty())
		})
	})

	Describe("Get a secret twice", func() {
		It("should delete the secret once has been retrieved", func() {
			value := "secret data"
			urlKey := ms.Set(value, 0)

			retrievedValue, _ := ms.Get(urlKey)
			Expect(retrievedValue).NotTo(BeNil())

			retrievedValueAgain, ok := ms.Get(urlKey)
			Expect(ok).To(BeFalse())
			Expect(retrievedValueAgain).To(BeZero())
		})
	})

	Describe("Get a secret after expiration", func() {
		It("should not be possible", func() {
			// TODO: find a better way to test this
			value := "secret data"
			urlKey := ms.Set(value, 100*time.Millisecond)
			time.Sleep(200 * time.Millisecond)

			retrievedValue, _ := ms.Get(urlKey)
			Expect(retrievedValue).To(BeEmpty())
		})
	})
})
