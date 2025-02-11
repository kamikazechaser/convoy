package datastore

import (
	"context"
)

type APIKeyRepository interface {
	CreateAPIKey(context.Context, *APIKey) error
	UpdateAPIKey(context.Context, *APIKey) error
	FindAPIKeyByID(context.Context, string) (*APIKey, error)
	FindAPIKeyByMaskID(context.Context, string) (*APIKey, error)
	FindAPIKeyByHash(context.Context, string) (*APIKey, error)
	RevokeAPIKeys(context.Context, []string) error
	LoadAPIKeysPaged(context.Context, *ApiKeyFilter, *Pageable) ([]APIKey, PaginationData, error)
}

type EventDeliveryRepository interface {
	CreateEventDelivery(context.Context, *EventDelivery) error
	FindEventDeliveryByID(context.Context, string) (*EventDelivery, error)
	FindEventDeliveriesByIDs(context.Context, []string) ([]EventDelivery, error)
	FindEventDeliveriesByEventID(context.Context, string) ([]EventDelivery, error)
	CountDeliveriesByStatus(context.Context, EventDeliveryStatus, SearchParams) (int64, error)
	UpdateStatusOfEventDelivery(context.Context, EventDelivery, EventDeliveryStatus) error
	UpdateStatusOfEventDeliveries(context.Context, []string, EventDeliveryStatus) error
	FindDiscardedEventDeliveries(ctx context.Context, appId, deviceId string, searchParams SearchParams) ([]EventDelivery, error)

	UpdateEventDeliveryWithAttempt(context.Context, EventDelivery, DeliveryAttempt) error
	CountEventDeliveries(context.Context, string, string, string, []EventDeliveryStatus, SearchParams) (int64, error)
	DeleteGroupEventDeliveries(ctx context.Context, filter *EventDeliveryFilter, hardDelete bool) error
	LoadEventDeliveriesPaged(context.Context, string, string, string, []EventDeliveryStatus, SearchParams, Pageable) ([]EventDelivery, PaginationData, error)
}

type EventRepository interface {
	CreateEvent(context.Context, *Event) error
	LoadEventIntervals(context.Context, string, SearchParams, Period, int) ([]EventInterval, error)
	FindEventByID(ctx context.Context, id string) (*Event, error)
	FindEventsByIDs(context.Context, []string) ([]Event, error)
	CountGroupMessages(ctx context.Context, groupID string) (int64, error)
	LoadEventsPaged(context.Context, *Filter) ([]Event, PaginationData, error)
	DeleteGroupEvents(context.Context, *EventFilter, bool) error
}

type GroupRepository interface {
	LoadGroups(context.Context, *GroupFilter) ([]*Group, error)
	CreateGroup(context.Context, *Group) error
	UpdateGroup(context.Context, *Group) error
	DeleteGroup(ctx context.Context, uid string) error
	FetchGroupByID(context.Context, string) (*Group, error)
	FetchGroupsByIDs(context.Context, []string) ([]Group, error)
	FillGroupsStatistics(ctx context.Context, groups []*Group) error
}

type OrganisationRepository interface {
	LoadOrganisationsPaged(context.Context, Pageable) ([]Organisation, PaginationData, error)
	CreateOrganisation(context.Context, *Organisation) error
	UpdateOrganisation(context.Context, *Organisation) error
	DeleteOrganisation(context.Context, string) error
	FetchOrganisationByID(context.Context, string) (*Organisation, error)
}

type OrganisationInviteRepository interface {
	LoadOrganisationsInvitesPaged(ctx context.Context, orgID string, inviteStatus InviteStatus, pageable Pageable) ([]OrganisationInvite, PaginationData, error)
	CreateOrganisationInvite(ctx context.Context, iv *OrganisationInvite) error
	UpdateOrganisationInvite(ctx context.Context, iv *OrganisationInvite) error
	DeleteOrganisationInvite(ctx context.Context, uid string) error
	FetchOrganisationInviteByID(ctx context.Context, uid string) (*OrganisationInvite, error)
	FetchOrganisationInviteByToken(ctx context.Context, token string) (*OrganisationInvite, error)
}

type OrganisationMemberRepository interface {
	LoadOrganisationMembersPaged(ctx context.Context, organisationID string, pageable Pageable) ([]*OrganisationMember, PaginationData, error)
	LoadUserOrganisationsPaged(ctx context.Context, userID string, pageable Pageable) ([]Organisation, PaginationData, error)
	CreateOrganisationMember(ctx context.Context, member *OrganisationMember) error
	UpdateOrganisationMember(ctx context.Context, member *OrganisationMember) error
	DeleteOrganisationMember(ctx context.Context, memberID string, orgID string) error
	FetchOrganisationMemberByID(ctx context.Context, memberID string, organisationID string) (*OrganisationMember, error)
	FetchOrganisationMemberByUserID(ctx context.Context, userID string, organisationID string) (*OrganisationMember, error)
}

type ApplicationRepository interface {
	CreateApplication(context.Context, *Application, string) error
	LoadApplicationsPaged(context.Context, string, string, Pageable) ([]Application, PaginationData, error)
	FindApplicationByID(context.Context, string) (*Application, error)
	UpdateApplication(context.Context, *Application, string) error
	DeleteApplication(context.Context, *Application) error
	CountGroupApplications(ctx context.Context, groupID string) (int64, error)
	DeleteGroupApps(context.Context, string) error
	LoadApplicationsPagedByGroupId(context.Context, string, Pageable) ([]Application, PaginationData, error)
	SearchApplicationsByGroupId(context.Context, string, SearchParams) ([]Application, error)
	FindApplicationEndpointByID(context.Context, string, string) (*Endpoint, error)
	CreateApplicationEndpoint(context.Context, string, string, *Endpoint) error
	ExpireSecret(ctx context.Context, appID string, endpointID string, secrets []Secret) error
}

type SubscriptionRepository interface {
	CreateSubscription(context.Context, string, *Subscription) error
	UpdateSubscription(context.Context, string, *Subscription) error
	LoadSubscriptionsPaged(context.Context, string, *FilterBy, Pageable) ([]Subscription, PaginationData, error)
	DeleteSubscription(context.Context, string, *Subscription) error
	FindSubscriptionByID(context.Context, string, string) (*Subscription, error)
	FindSubscriptionsByEventType(context.Context, string, string, EventType) ([]Subscription, error)
	FindSubscriptionsBySourceIDs(context.Context, string, string) ([]Subscription, error)
	FindSubscriptionsByAppID(ctx context.Context, groupId string, appID string) ([]Subscription, error)
	FindSubscriptionByDeviceID(ctx context.Context, groupId string, deviceID string) (*Subscription, error)
	UpdateSubscriptionStatus(context.Context, string, string, SubscriptionStatus) error
}

type SourceRepository interface {
	CreateSource(context.Context, *Source) error
	UpdateSource(ctx context.Context, groupID string, source *Source) error
	FindSourceByID(ctx context.Context, groupID string, id string) (*Source, error)
	FindSourceByMaskID(ctx context.Context, maskID string) (*Source, error)
	DeleteSourceByID(ctx context.Context, groupID string, id string) error
	LoadSourcesPaged(ctx context.Context, groupID string, filter *SourceFilter, pageable Pageable) ([]Source, PaginationData, error)
}

type DeviceRepository interface {
	CreateDevice(ctx context.Context, device *Device) error
	UpdateDevice(ctx context.Context, device *Device, appID, groupID string) error
	UpdateDeviceLastSeen(ctx context.Context, device *Device, appID, groupID string, status DeviceStatus) error
	DeleteDevice(ctx context.Context, uid string, appID, groupID string) error
	FetchDeviceByID(ctx context.Context, uid string, appID, groupID string) (*Device, error)
	FetchDeviceByHostName(ctx context.Context, hostName string, appID, groupID string) (*Device, error)
	LoadDevicesPaged(ctx context.Context, groupID string, filter *ApiKeyFilter, pageable Pageable) ([]Device, PaginationData, error)
}

type UserRepository interface {
	CreateUser(context.Context, *User) error
	UpdateUser(ctx context.Context, user *User) error
	FindUserByEmail(context.Context, string) (*User, error)
	FindUserByID(context.Context, string) (*User, error)
	FindUserByToken(context.Context, string) (*User, error)
	LoadUsersPaged(context.Context, Pageable) ([]User, PaginationData, error)
}

type ConfigurationRepository interface {
	CreateConfiguration(context.Context, *Configuration) error
	LoadConfiguration(context.Context) (*Configuration, error)
	UpdateConfiguration(context.Context, *Configuration) error
}
