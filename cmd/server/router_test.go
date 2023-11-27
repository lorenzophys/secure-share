package main_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	secureShare "github.com/lorenzophys/secure_share/cmd/server"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) Set(secret string, ttl time.Duration) string {
	args := m.Called(secret, ttl)
	return args.String(0)
}

func (m *MockStore) Get(key string) (string, bool) {
	args := m.Called(key)
	return args.String(0), args.Bool(1)
}

func (m *MockStore) RemoveExpiredSecrets() {
	m.Called()
}

var _ = Describe("Test Handlers", func() {
	var (
		e   *echo.Echo
		req *http.Request
		rec *httptest.ResponseRecorder
		c   echo.Context
		app *secureShare.Application
	)

	BeforeEach(func() {
		e = echo.New()
		rec = httptest.NewRecorder()

		mockStore := new(MockStore)
		mockStore.On("Set", "test secret", time.Hour).Return("mockedURLKey").Once()
		mockStore.On("Get", "mockedURLKey").Return("test secret", true).Once()

		app = &secureShare.Application{
			Store: mockStore,
		}

	})

	It("should handle valid secret posting", func() {
		req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader("textareaContent=test secret&menuSelection=1h"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

		c = e.NewContext(req, rec)

		data, err := app.HandlePostSecret(c)

		Expect(err).ToNot(HaveOccurred())
		Expect(data).To(Equal("mockedURLKey"))
	})

	It("should handle invalid duration format", func() {
		req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader("textareaContent=test secret&menuSelection=invalid"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

		c = e.NewContext(req, rec)

		data, err := app.HandlePostSecret(c)

		Expect(err).To(HaveOccurred())
		Expect(data).To(BeEmpty())
	})

	It("should handle correct retrieval", func() {
		req = httptest.NewRequest(http.MethodGet, "/mockedURLKey", nil)

		c = e.NewContext(req, rec)
		c.SetParamNames("key")
		c.SetParamValues("mockedURLKey")

		password, err := app.HandleGetSecret(c)

		Expect(err).ToNot(HaveOccurred())
		Expect(password).To(Equal("test secret"))
	})
})
