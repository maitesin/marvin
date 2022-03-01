package tracking

import "github.com/maitesin/marvin/internal/domain"

type Tracker interface {
	Track(id string) ([]domain.DeliveryEvent, error)
}
