package mocks

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"

	"github.com/authelia/authelia/v4/internal/authorization"
	"github.com/authelia/authelia/v4/internal/configuration/schema"
	"github.com/authelia/authelia/v4/internal/middlewares"
	"github.com/authelia/authelia/v4/internal/regulation"
	"github.com/authelia/authelia/v4/internal/session"
	"github.com/authelia/authelia/v4/internal/templates"
)

// MockAutheliaCtx a mock of AutheliaCtx.
type MockAutheliaCtx struct {
	// Logger hook.
	Hook *test.Hook
	Ctx  *middlewares.AutheliaCtx
	Ctrl *gomock.Controller

	// Providers.
	UserProviderMock *MockUserProvider
	StorageMock      *MockStorage
	NotifierMock     *MockNotifier
	TOTPMock         *MockTOTP

	UserSession *session.UserSession

	Clock TestingClock
}

// TestingClock implementation of clock for tests.
type TestingClock struct {
	now time.Time
}

// Now return the stored clock.
func (dc *TestingClock) Now() time.Time {
	return dc.now
}

// After return a channel receiving the time after duration has elapsed.
func (dc *TestingClock) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}

// Set set the time of the clock.
func (dc *TestingClock) Set(now time.Time) {
	dc.now = now
}

// NewMockAutheliaCtx create an instance of AutheliaCtx mock.
func NewMockAutheliaCtx(t *testing.T) *MockAutheliaCtx {
	mockAuthelia := new(MockAutheliaCtx)
	mockAuthelia.Clock = TestingClock{}

	datetime, _ := time.Parse("2006-Jan-02", "2013-Feb-03")
	mockAuthelia.Clock.Set(datetime)

	config := schema.Configuration{}
	config.Session.RememberMeDuration = schema.DefaultSessionConfiguration.RememberMeDuration
	config.Session.Name = "authelia_session"
	config.Session.Domain = "example.com"
	config.AccessControl.DefaultPolicy = "deny"
	config.AccessControl.Rules = []schema.ACLRule{{
		Domains: []string{"bypass.example.com"},
		Policy:  "bypass",
	}, {
		Domains: []string{"one-factor.example.com"},
		Policy:  "one_factor",
	}, {
		Domains: []string{"two-factor.example.com"},
		Policy:  "two_factor",
	}, {
		Domains: []string{"deny.example.com"},
		Policy:  "deny",
	}, {
		Domains:  []string{"admin.example.com"},
		Policy:   "two_factor",
		Subjects: [][]string{{"group:admin"}},
	}, {
		Domains:  []string{"grafana.example.com"},
		Policy:   "two_factor",
		Subjects: [][]string{{"group:grafana"}},
	}}

	providers := middlewares.Providers{}

	mockAuthelia.Ctrl = gomock.NewController(t)
	mockAuthelia.UserProviderMock = NewMockUserProvider(mockAuthelia.Ctrl)
	providers.UserProvider = mockAuthelia.UserProviderMock

	mockAuthelia.StorageMock = NewMockStorage(mockAuthelia.Ctrl)
	providers.StorageProvider = mockAuthelia.StorageMock

	mockAuthelia.NotifierMock = NewMockNotifier(mockAuthelia.Ctrl)
	providers.Notifier = mockAuthelia.NotifierMock

	providers.Authorizer = authorization.NewAuthorizer(
		&config)

	providers.SessionProvider = session.NewProvider(
		config.Session, nil)

	providers.Regulator = regulation.NewRegulator(config.Regulation, providers.StorageProvider, &mockAuthelia.Clock)

	mockAuthelia.TOTPMock = NewMockTOTP(mockAuthelia.Ctrl)
	providers.TOTP = mockAuthelia.TOTPMock

	var err error

	if providers.Templates, err = templates.New(templates.Config{}); err != nil {
		panic(err)
	}

	request := &fasthttp.RequestCtx{}
	// Set a cookie to identify this client throughout the test.
	// request.Request.Header.SetCookie("authelia_session", "client_cookie").

	ctx := middlewares.NewAutheliaCtx(request, config, providers)
	mockAuthelia.Ctx = ctx

	logger, hook := test.NewNullLogger()
	mockAuthelia.Hook = hook

	mockAuthelia.Ctx.Logger = logrus.NewEntry(logger)

	return mockAuthelia
}

// NewMockAutheliaCtxWithUserSession create an instance of AutheliaCtx mock with predefined user session.
func NewMockAutheliaCtxWithUserSession(t *testing.T, userSession session.UserSession) *MockAutheliaCtx {
	mock := NewMockAutheliaCtx(t)
	err := mock.Ctx.SaveSession(userSession)
	require.NoError(t, err)

	return mock
}

// Close close the mock.
func (m *MockAutheliaCtx) Close() {
	m.Hook.Reset()
	m.Ctrl.Finish()
}

// SetRequestBody set the request body from a struct with json tags.
func (m *MockAutheliaCtx) SetRequestBody(t *testing.T, body interface{}) {
	bodyBytes, err := json.Marshal(body)
	require.NoError(t, err)
	m.Ctx.Request.SetBody(bodyBytes)
}

// Assert401KO assert an error response from the service.
func (m *MockAutheliaCtx) Assert401KO(t *testing.T, message string) {
	assert.Equal(t, 401, m.Ctx.Response.StatusCode())
	assert.Equal(t, fmt.Sprintf("{\"status\":\"KO\",\"message\":\"%s\"}", message), string(m.Ctx.Response.Body()))
}

// Assert200KO assert an error response from the service.
func (m *MockAutheliaCtx) Assert200KO(t *testing.T, message string) {
	assert.Equal(t, 200, m.Ctx.Response.StatusCode())
	assert.Equal(t, fmt.Sprintf("{\"status\":\"KO\",\"message\":\"%s\"}", message), string(m.Ctx.Response.Body()))
}

// Assert200OK assert a successful response from the service.
func (m *MockAutheliaCtx) Assert200OK(t *testing.T, data interface{}) {
	assert.Equal(t, 200, m.Ctx.Response.StatusCode())

	response := middlewares.OKResponse{
		Status: "OK",
		Data:   data,
	}

	b, err := json.Marshal(response)

	assert.NoError(t, err)
	assert.Equal(t, string(b), string(m.Ctx.Response.Body()))
}

// GetResponseData retrieves a response from the service.
func (m *MockAutheliaCtx) GetResponseData(t *testing.T, data interface{}) {
	okResponse := middlewares.OKResponse{}
	okResponse.Data = data
	err := json.Unmarshal(m.Ctx.Response.Body(), &okResponse)
	require.NoError(t, err)
}

// GetResponseError retrieves an error response from the service.
func (m *MockAutheliaCtx) GetResponseError(t *testing.T) (errResponse middlewares.ErrorResponse) {
	err := json.Unmarshal(m.Ctx.Response.Body(), &errResponse)
	require.NoError(t, err)

	return errResponse
}
