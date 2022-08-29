package socket

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/frain-dev/convoy"
	"github.com/frain-dev/convoy/config"
	"github.com/frain-dev/convoy/datastore"
	m "github.com/frain-dev/convoy/datastore/mongo"
	"github.com/frain-dev/convoy/pkg/httpheader"
	"github.com/frain-dev/convoy/util"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
)

// Register new clients.
var register chan *Client

// Unregister clients.
var unregister chan *Client

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	lock sync.RWMutex // prevent data race on deviceClients

	// Registered clients.
	deviceClients map[string]*Client

	watchCollectionFn WatchCollectionFn

	close chan struct{}

	// events from the change stream are written to thi channel and are sent to the respective devices
	events chan *CLIEvent

	// this ticker is used to periodically set inactive (or incorrectly disconnected) devices to offline
	ticker *time.Ticker
}

type AckEventDelivery struct {
	UID string `json:"uid"`
}

type CLIEvent struct {
	UID     string                `json:"uid"`
	Headers httpheader.HTTPHeader `json:"headers" bson:"headers"`
	Data    json.RawMessage       `json:"data"`

	// for filtering this event delivery
	EventType string `json:"-"`
	DeviceID  string `json:"-"`
	AppID     string `json:"-"`
	GroupID   string `json:"-"`
}

type WatchCollectionFn func(fn func(doc map[string]interface{}) error, pipeline mongo.Pipeline, collection string, stop chan struct{}) error

func NewHub() *Hub {
	register = make(chan *Client, 1)
	unregister = make(chan *Client, 1)

	return &Hub{
		watchCollectionFn: watchCollection,
		deviceClients:     map[string]*Client{},
		events:            make(chan *CLIEvent, 10),
		close:             make(chan struct{}),
	}
}

func (h *Hub) StartEventSender() {
	for {
		select {
		case ev := <-h.events:
			h.lock.RLock()
			client := h.deviceClients[ev.DeviceID]
			h.lock.RUnlock()

			// there is no valid client for this event delivery, so skip it
			if client == nil {
				continue
			}

			if !client.IsOnline() {
				client.GoOffline()
				client.Close()
				continue
			}

			if client.Device.GroupID != ev.GroupID {
				continue
			}

			if !util.IsStringEmpty(client.Device.AppID) {
				if client.Device.AppID != ev.AppID {
					continue
				}
			}

			if !client.HasEventType(ev.EventType) {
				continue
			}

			j, err := json.Marshal(ev)
			if err != nil {
				log.WithError(err).Error("failed to marshal cli event")
				continue
			}

			err = client.conn.WriteMessage(websocket.BinaryMessage, j)
			if err != nil {
				log.WithError(err).Error("failed to write event to socket")
			}
		case <-h.close:
			return
		}
	}
}

func (h *Hub) StartEventWatcher() {
	fn := h.watchEventDeliveriesCollection()
	err := h.watchCollectionFn(fn, mongo.Pipeline{}, m.EventDeliveryCollection, h.close)
	if err != nil {
		log.WithError(err).Error("database collection watcher exited unexpectedly")
	}
}

func (h *Hub) watchEventDeliveriesCollection() func(doc map[string]interface{}) error {
	return func(doc map[string]interface{}) error {
		var ed *datastore.EventDelivery
		b, err := json.Marshal(doc)
		if err != nil {
			return err
		}

		err = json.Unmarshal(b, &ed)
		if err != nil {
			return err
		}

		if ed.CLIMetadata == nil {
			return nil
		}

		h.events <- &CLIEvent{
			UID:       ed.UID,
			Data:      ed.Metadata.Data,
			Headers:   ed.Headers,
			EventType: ed.CLIMetadata.EventType,
			AppID:     ed.AppID,
			DeviceID:  ed.DeviceID,
			GroupID:   ed.GroupID,
		}

		return nil
	}
}

func (h *Hub) StartRegister() {
	for {
		select {
		case client := <-register:
			h.lock.Lock()
			h.deviceClients[client.deviceID] = client
			h.lock.Unlock()
		case <-h.close:
			return
		}
	}
}

func (h *Hub) StartUnregister() {
	for {
		select {
		case client := <-unregister:
			h.lock.Lock()
			delete(h.deviceClients, client.deviceID)
			h.lock.Unlock()
		case <-h.close:
			return
		}
	}
}

func (h *Hub) StartClientStatusWatcher() {
	h.ticker = time.NewTicker(time.Second * 5)
	defer h.ticker.Stop()

	for {
		select {
		case <-h.ticker.C:
			// TODO(raymond): it might be better to read directly from the devicr collection then disable them
			for k, v := range h.deviceClients {
				h.lock.Lock()
				if !h.deviceClients[k].IsOnline() {
					h.deviceClients[k].GoOffline()
					h.deviceClients[k].Close()
					log.Printf("Device %s has be set to offline after inactivity for 2 munites", v.Device.HostName)
				}
				h.lock.Unlock()
			}
		case <-h.close:
			return
		}
	}
}

func (h *Hub) Stop() {
	close(h.close)
}

type WatcherFn func(map[string]interface{}) error

func watchCollection(fn func(map[string]interface{}) error, pipeline mongo.Pipeline, collection string, stop chan struct{}) error {
	cfg, err := config.Get()
	if err != nil {
		return err
	}

	if cfg.Database.Type != "mongodb" {
		return convoy.ErrUnsupportedDatebase
	}

	client, err := m.New(cfg)
	if err != nil {
		return err
	}

	db := client.Client().(*mongo.Database)
	coll := db.Collection(collection)
	ctx := context.Background()

	cs, err := coll.Watch(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cs.Close(ctx)

	for {
		select {
		case <-stop:
			logrus.Println("Exiting Database watcher")
			return nil
		default:
			ok := cs.Next(ctx)
			if ok {
				var document *convoy.GenericMap
				err := cs.Decode(&document)
				if err != nil {
					return err
				}

				if (*document)["operationType"].(string) == "insert" {
					doc := (*document)["fullDocument"].(convoy.GenericMap)
					err := fn(doc)
					if err != nil {
						return err
					}
				}
			}
		}
	}
}
