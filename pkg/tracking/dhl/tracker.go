package dhl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/maitesin/marvin/internal/domain"
)

const urlRegex = "https://clientesparcel.dhl.es/LiveTracking/api/expediciones?numeroExpedicion=%s"

// Tracker for the DHL delivery service
type Tracker struct {
	client *http.Client
}

// NewTracker constructor for the DHL tracker
func NewTracker(client *http.Client) (*Tracker, error) {
	return &Tracker{
		client: client,
	}, nil
}

func (t *Tracker) Track(id string) ([]domain.DeliveryEvent, error) {
	resp, err := t.client.Get(fmt.Sprintf(urlRegex, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var body body
	err = json.Unmarshal(respBody, &body)
	if err != nil {
		return nil, err
	}

	events := make([]domain.DeliveryEvent, len(body.Shipments))
	for i, event := range body.Shipments {
		events[i] = domain.DeliveryEvent{
			Timestamp:   fmt.Sprintf("%s %s", strings.TrimSpace(event.Date), strings.TrimSpace(event.Time)),
			Information: fmt.Sprintf("%s %s", strings.TrimSpace(event.Text), strings.TrimSpace(event.City)),
		}
	}

	return events, nil
}
