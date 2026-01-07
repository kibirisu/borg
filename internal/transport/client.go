package transport

import (
	"encoding/json/v2"
	"net/http"

	"github.com/kibirisu/borg/internal/domain"
)

type Client interface {
	Get(string) (*domain.ObjectOrLink, error)
}

type client struct {
	client http.Client
}

var _ Client = (*client)(nil)

// Get implements Client.
func (c *client) Get(uri string) (*domain.ObjectOrLink, error) {
	var object domain.ObjectOrLink
	resp, err := c.client.Get(uri)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if err = json.UnmarshalRead(resp.Body, &object); err != nil {
		return nil, err
	}
	return &object, nil
}
