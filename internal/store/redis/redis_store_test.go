package redis_store_test

import (
	"context"
	"time"

	redisStore "github.com/lorenzophys/secure_share/internal/store/redis"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var _ = Describe("RedisStore", func() {
	var rs *redisStore.RedisStore

	BeforeEach(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		req := testcontainers.ContainerRequest{
			Image:        "redis:latest",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor:   wait.ForLog("Ready to accept connections"),
		}
		redisC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		Expect(err).NotTo(HaveOccurred())

		endpoint, err := redisC.Endpoint(ctx, "")
		Expect(err).NotTo(HaveOccurred())

		rs = redisStore.NewRedisStore(endpoint, "", 0)
	})

	Describe("Set and Get functions", func() {
		It("should store and retrieve data correctly", func() {
			value := "secret data"
			urlKey := rs.Set(value, 0)

			Expect(urlKey).NotTo(BeEmpty())

			retrievedValue, ok := rs.Get(urlKey)

			Expect(ok).To(BeTrue())
			Expect(retrievedValue).NotTo(BeNil())
			Expect(retrievedValue).To(Equal(value))
		})
	})

	Describe("Get with invalid key", func() {
		It("should return false for invalid keys", func() {
			secret, ok := rs.Get("invalid-key")
			Expect(ok).To(BeFalse())
			Expect(secret).To(BeEmpty())
		})
	})

	Describe("Get a secret twice", func() {
		It("should delete the secret once has been retrieved", func() {
			value := "secret data"
			urlKey := rs.Set(value, 0)

			retrievedValue, _ := rs.Get(urlKey)
			Expect(retrievedValue).NotTo(BeNil())

			retrievedValueAgain, ok := rs.Get(urlKey)
			Expect(ok).To(BeFalse())
			Expect(retrievedValueAgain).To(BeZero())
		})
	})

	Describe("Get a secret after expiration", func() {
		It("should not be possible", func() {
			// TODO: find a better way to test this
			value := "secret data"
			urlKey := rs.Set(value, 100*time.Millisecond)
			time.Sleep(200 * time.Millisecond)

			retrievedValue, _ := rs.Get(urlKey)
			Expect(retrievedValue).To(BeEmpty())
		})
	})
})
