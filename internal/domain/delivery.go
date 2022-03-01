package domain

type DeliveryEvent struct {
	Timestamp   string
	Information string
}

type Delivery struct {
	Identifier string
	Events     []DeliveryEvent
}

func NewDelivery(code string, events []DeliveryEvent) Delivery {
	return Delivery{
		Identifier: code,
		Events:     events,
	}
}

func (d *Delivery) AddEvents(events ...DeliveryEvent) {
	d.Events = append(d.Events, events...)
}
