//go:build integration
// +build integration

package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/frain-dev/convoy/internal/pkg/metrics"

	"github.com/frain-dev/convoy/server/testdb"
	"github.com/stretchr/testify/suite"

	"github.com/frain-dev/convoy/config"
	"github.com/frain-dev/convoy/datastore"
	cm "github.com/frain-dev/convoy/datastore/mongo"
	"github.com/frain-dev/convoy/server/models"
	"github.com/google/uuid"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type DashboardIntegrationTestSuite struct {
	suite.Suite
	DB              cm.Client
	Router          http.Handler
	ConvoyApp       *ApplicationHandler
	AuthenticatorFn AuthenticatorFn
	DefaultUser     *datastore.User
	DefaultOrg      *datastore.Organisation
	DefaultGroup    *datastore.Group
}

func (s *DashboardIntegrationTestSuite) SetupSuite() {
	s.DB = getDB()
	s.ConvoyApp = buildServer()
	s.Router = s.ConvoyApp.BuildRoutes()
}

func (s *DashboardIntegrationTestSuite) SetupTest() {
	testdb.PurgeDB(s.T(), s.DB)

	// Setup Default User
	user, err := testdb.SeedDefaultUser(s.ConvoyApp.A.Store)
	require.NoError(s.T(), err)
	s.DefaultUser = user

	// Setup Default Organisation
	org, err := testdb.SeedDefaultOrganisation(s.ConvoyApp.A.Store, user)
	require.NoError(s.T(), err)
	s.DefaultOrg = org

	// Setup Default Group.
	s.DefaultGroup, _ = testdb.SeedDefaultGroup(s.ConvoyApp.A.Store, s.DefaultOrg.UID)

	s.AuthenticatorFn = authenticateRequest(&models.LoginUser{
		Username: user.Email,
		Password: testdb.DefaultUserPassword,
	})

	// Setup Config.
	err = config.LoadConfig("./testdata/Auth_Config/full-convoy-with-jwt-realm.json")
	require.NoError(s.T(), err)

	apiRepo := cm.NewApiKeyRepo(s.ConvoyApp.A.Store)
	userRepo := cm.NewUserRepo(s.ConvoyApp.A.Store)
	initRealmChain(s.T(), apiRepo, userRepo, s.ConvoyApp.A.Cache)
}

func (s *DashboardIntegrationTestSuite) TearDownTest() {
	testdb.PurgeDB(s.T(), s.DB)
	metrics.Reset()
}

func (s *DashboardIntegrationTestSuite) TestGetDashboardSummary() {
	ctx := context.Background()
	application := &datastore.Application{
		UID:            "abc",
		GroupID:        s.DefaultGroup.UID,
		Title:          "test-app",
		SupportEmail:   "test@suport.com",
		Endpoints:      []datastore.Endpoint{},
		CreatedAt:      primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:      primitive.NewDateTimeFromTime(time.Now()),
		DocumentStatus: datastore.ActiveDocumentStatus,
	}

	appRepo := cm.NewApplicationRepo(s.ConvoyApp.A.Store)
	err := appRepo.CreateApplication(ctx, application, application.GroupID)
	require.NoError(s.T(), err)

	events := []datastore.Event{
		{
			UID:              uuid.New().String(),
			EventType:        "*",
			MatchedEndpoints: 1,
			ProviderID:       "provider_id",
			Data:             json.RawMessage(`{"data":"12345"}`),
			GroupID:          s.DefaultGroup.UID,
			AppID:            application.UID,
			CreatedAt:        primitive.NewDateTimeFromTime(time.Date(2021, time.January, 1, 1, 1, 1, 0, time.UTC)),
			UpdatedAt:        primitive.NewDateTimeFromTime(time.Date(2021, time.January, 1, 1, 1, 1, 0, time.UTC)),
			DocumentStatus:   datastore.ActiveDocumentStatus,
		},
		{
			UID:              uuid.New().String(),
			EventType:        "*",
			MatchedEndpoints: 1,
			ProviderID:       "provider_id",
			Data:             json.RawMessage(`{"data":"12345"}`),
			GroupID:          s.DefaultGroup.UID,
			AppID:            application.UID,
			CreatedAt:        primitive.NewDateTimeFromTime(time.Date(2021, time.January, 10, 1, 1, 1, 0, time.UTC)),
			UpdatedAt:        primitive.NewDateTimeFromTime(time.Date(2021, time.January, 10, 1, 1, 1, 0, time.UTC)),
			DocumentStatus:   datastore.ActiveDocumentStatus,
		},
		{
			UID:              uuid.New().String(),
			EventType:        "*",
			MatchedEndpoints: 1,
			ProviderID:       "provider_id",
			Data:             json.RawMessage(`{"data":"12345"}`),
			GroupID:          s.DefaultGroup.UID,
			AppID:            application.UID,
			CreatedAt:        primitive.NewDateTimeFromTime(time.Date(2022, time.March, 20, 1, 1, 1, 0, time.UTC)),
			UpdatedAt:        primitive.NewDateTimeFromTime(time.Date(2022, time.March, 20, 1, 1, 1, 0, time.UTC)),
			DocumentStatus:   datastore.ActiveDocumentStatus,
		},
		{
			UID:              uuid.New().String(),
			EventType:        "*",
			MatchedEndpoints: 1,
			ProviderID:       "provider_id",
			Data:             json.RawMessage(`{"data":"12345"}`),
			GroupID:          s.DefaultGroup.UID,
			AppID:            application.UID,
			CreatedAt:        primitive.NewDateTimeFromTime(time.Date(2022, time.March, 20, 1, 1, 1, 0, time.UTC)),
			UpdatedAt:        primitive.NewDateTimeFromTime(time.Date(2022, time.March, 20, 1, 1, 1, 0, time.UTC)),
			DocumentStatus:   datastore.ActiveDocumentStatus,
		},
		{
			UID:              uuid.New().String(),
			EventType:        "*",
			MatchedEndpoints: 1,
			ProviderID:       "provider_id",
			Data:             json.RawMessage(`{"data":"12345"}`),
			GroupID:          s.DefaultGroup.UID,
			AppID:            application.UID,
			CreatedAt:        primitive.NewDateTimeFromTime(time.Date(2022, time.March, 20, 1, 1, 1, 0, time.UTC)),
			UpdatedAt:        primitive.NewDateTimeFromTime(time.Date(2022, time.March, 20, 1, 1, 1, 0, time.UTC)),
			DocumentStatus:   datastore.ActiveDocumentStatus,
		},
		{
			UID:              uuid.New().String(),
			EventType:        "*",
			MatchedEndpoints: 1,
			ProviderID:       "provider_id",
			Data:             json.RawMessage(`{"data":"12345"}`),
			GroupID:          s.DefaultGroup.UID,
			AppID:            application.UID,
			CreatedAt:        primitive.NewDateTimeFromTime(time.Date(2022, time.March, 20, 1, 1, 1, 0, time.UTC)),
			UpdatedAt:        primitive.NewDateTimeFromTime(time.Date(2022, time.March, 20, 1, 1, 1, 0, time.UTC)),
			DocumentStatus:   datastore.ActiveDocumentStatus,
		},
	}

	eventRepo := cm.NewEventRepository(s.ConvoyApp.A.Store)
	for i := range events {
		err = eventRepo.CreateEvent(ctx, &events[i])
		require.NoError(s.T(), err)
	}

	type urlQuery struct {
		groupID   string
		startDate string
		endDate   string
		Type      string
	}

	tests := []struct {
		name       string
		method     string
		urlQuery   urlQuery
		statusCode int
	}{
		{
			name:       "should_fetch_yearly_dashboard_summary",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			urlQuery: urlQuery{
				groupID:   s.DefaultGroup.UID,
				startDate: "2021-01-01T00:00:00",
				endDate:   "2021-01-30T00:00:00",
				Type:      "yearly",
			},
		},
		{
			name:       "should_fetch_monthly_dashboard_summary",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			urlQuery: urlQuery{
				groupID:   s.DefaultGroup.UID,
				startDate: "2021-01-01T00:00:00",
				endDate:   "2022-12-27T00:00:00",
				Type:      "monthly",
			},
		},
		{
			name:       "should_fetch_weekly_dashboard_summary",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			urlQuery: urlQuery{
				groupID:   s.DefaultGroup.UID,
				startDate: "2021-01-01T00:00:00",
				endDate:   "2022-12-27T00:00:00",
				Type:      "weekly",
			},
		},
		{
			name:       "should_fetch_daily_dashboard_summary",
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			urlQuery: urlQuery{
				groupID:   s.DefaultGroup.UID,
				startDate: "2021-01-01T00:00:00",
				endDate:   "2022-12-27T00:00:00",
				Type:      "daily",
			},
		},
		{
			name:       "should_error_for_empty_startDate",
			method:     http.MethodGet,
			statusCode: http.StatusBadRequest,
			urlQuery: urlQuery{
				groupID: s.DefaultGroup.UID,
				endDate: "2022-12-27T00:00:00",
				Type:    "daily",
			},
		},
		{
			name:       "should_error_for_invalid_startDate",
			method:     http.MethodGet,
			statusCode: http.StatusBadRequest,
			urlQuery: urlQuery{
				groupID:   s.DefaultGroup.UID,
				startDate: "2021-01-01",
				endDate:   "2022-12-27T00:00:00",
				Type:      "daily",
			},
		},
		{
			name:       "should_error_for_invalid_type",
			method:     http.MethodGet,
			statusCode: http.StatusBadRequest,
			urlQuery: urlQuery{
				groupID:   s.DefaultGroup.UID,
				startDate: "2021-01-01T00:00:00",
				endDate:   "2022-12-27T00:00:00",
				Type:      "abc",
			},
		},
		{
			name:       "should_error_for_startDate_greater_than_endDate",
			method:     http.MethodGet,
			statusCode: http.StatusBadRequest,
			urlQuery: urlQuery{
				groupID:   s.DefaultGroup.UID,
				startDate: "2021-01-01T00:00:00",
				endDate:   "2020-12-27T00:00:00",
				Type:      "daily",
			},
		},
	}
	for _, tc := range tests {
		s.T().Run(tc.name, func(t *testing.T) {
			err := config.LoadConfig("./testdata/Auth_Config/full-convoy-with-jwt-realm.json")
			if err != nil {
				t.Errorf("Failed to load config file: %v", err)
			}
			apiRepo := cm.NewApiKeyRepo(s.ConvoyApp.A.Store)
			userRepo := cm.NewUserRepo(s.ConvoyApp.A.Store)
			initRealmChain(t, apiRepo, userRepo, s.ConvoyApp.A.Cache)

			req := httptest.NewRequest(tc.method, fmt.Sprintf("/ui/organisations/%s/projects/%s/dashboard/summary?startDate=%s&endDate=%s&type=%s", s.DefaultOrg.UID, tc.urlQuery.groupID, tc.urlQuery.startDate, tc.urlQuery.endDate, tc.urlQuery.Type), nil)
			err = s.AuthenticatorFn(req, s.Router)
			require.NoError(s.T(), err)

			w := httptest.NewRecorder()

			s.Router.ServeHTTP(w, req)

			if w.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, w.Code)
			}

			verifyMatch(t, *w)
		})
	}
}

func TestDashboardIntegrationTestSuiteTest(t *testing.T) {
	suite.Run(t, new(DashboardIntegrationTestSuite))
}

func verifyMatch(t *testing.T, w httptest.ResponseRecorder) {
	g := goldie.New(
		t,
		goldie.WithDiffEngine(goldie.ColoredDiff),
	)
	g.Assert(t, t.Name(), w.Body.Bytes())
}
