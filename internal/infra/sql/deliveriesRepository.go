package sql

import (
	"context"

	"github.com/maitesin/marvin/internal/app"
	"github.com/maitesin/marvin/internal/domain"
	"github.com/upper/db/v4"
)

const (
	deliveryTable = "deliveries"
)

type Event struct {
	Timestamp   string `json:"timestamp"`
	Information string `json:"information"`
}

type Delivery struct {
	Identifier string  `db:"identifier"`
	Events     []Event `db:"events"`
	Delivered  bool    `db:"delivered"`
}

type DeliveriesRepository struct {
	sess db.Session
}

func NewDeliveriesRepository(sess db.Session) *DeliveriesRepository {
	return &DeliveriesRepository{sess: sess}
}

func (dr *DeliveriesRepository) Insert(ctx context.Context, delivery domain.Delivery) error {
	sDelivery := domain2SQLDelivery(delivery)

	_, err := dr.sess.WithContext(ctx).
		Collection(deliveryTable).
		Insert(sDelivery)

	return err
}

func (dr *DeliveriesRepository) FindByIdentifier(ctx context.Context, identifier string) (domain.Delivery, error) {
	var delivery Delivery
	err := dr.sess.WithContext(ctx).
		Collection(deliveryTable).
		Find(db.Cond{"identifier": identifier}).
		One(&delivery)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return domain.Delivery{}, app.NewDeliveryNotFound(identifier)
		}
		return domain.Delivery{}, err
	}
	return domain.NewDelivery(delivery.Identifier, sql2DomainEvents(delivery.Events...)), nil
}

func (dr *DeliveriesRepository) FindAllNotDelivered(ctx context.Context) ([]domain.Delivery, error) {
	var deliveries []Delivery
	err := dr.sess.WithContext(ctx).
		Collection(deliveryTable).
		Find(db.Cond{"delivered": false}).
		All(&deliveries)
	if err != nil {
		return nil, err
	}
	return sql2DomainDeliveries(deliveries...), nil
}

func sql2DomainEvents(sEvents ...Event) []domain.DeliveryEvent {
	events := make([]domain.DeliveryEvent, len(sEvents))

	for i := range sEvents {
		events[i] = domain.DeliveryEvent{
			Timestamp:   sEvents[i].Timestamp,
			Information: sEvents[i].Information,
		}
	}

	return events
}

func sql2DomainDeliveries(sDeliveries ...Delivery) []domain.Delivery {
	deliveries := make([]domain.Delivery, len(sDeliveries))

	for i := range sDeliveries {
		deliveries[i] = domain.NewDelivery(sDeliveries[i].Identifier, sql2DomainEvents(sDeliveries[i].Events...))
	}

	return deliveries
}

func domain2SQLEvents(dEvents ...domain.DeliveryEvent) []Event {
	events := make([]Event, len(dEvents))

	for i := range dEvents {
		events[i] = Event{
			Timestamp:   dEvents[i].Timestamp,
			Information: dEvents[i].Information,
		}
	}

	return events
}

func domain2SQLDelivery(dDelivery domain.Delivery) Delivery {
	return Delivery{
		Identifier: dDelivery.Identifier,
		Events:     domain2SQLEvents(dDelivery.Events...),
	}
}
