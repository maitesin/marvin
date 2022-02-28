// Correos is the Spanish national mail service https://www.correos.es/

package correos

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/maitesin/marvin/pkg/tracking"
)

const urlRegex = "https://api1.correos.es/digital-services/searchengines/api/v1/?text=%s&language=ES&searchType=envio"

// Tracker for the Correos delivery service
type Tracker struct {
	client *http.Client
}

// NewTracker constructor for the Correos tracker
func NewTracker(client *http.Client) (*Tracker, error) {
	return &Tracker{
		client: client,
	}, nil
}

func (t *Tracker) Track(id string) ([]tracking.Event, error) {
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

	if len(body.Shipments) != 1 {
		return nil, fmt.Errorf("expected information from a single shipment, found %d", len(body.Shipments))
	}

	events := make([]tracking.Event, len(body.Shipments[0].Events))
	for i, event := range body.Shipments[0].Events {
		events[i] = tracking.Event{
			Timestamp:   fmt.Sprintf("%s %s", event.Date, event.Time),
			Information: fmt.Sprintf("%s (%s)", event.SummaryText, event.ExtendedText),
		}
	}

	return events, nil
}
