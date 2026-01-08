package transport

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"net/http"
	"time"

	"github.com/kibirisu/borg/internal/domain"
)

type Client interface {
	Get(context.Context, string) (*domain.ObjectOrLink, error)
	Post(context.Context, string, *domain.Object) error
}

type client struct {
	client http.Client
}

var _ Client = (*client)(nil)

func New() Client {
	return &client{http.Client{
		Timeout: 2 * time.Second,
	}}
}

// Get implements Client.
func (c *client) Get(ctx context.Context, uri string) (*domain.ObjectOrLink, error) {
	var object domain.ObjectOrLink
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(req)
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

// Post implements Client.
func (c *client) Post(ctx context.Context, uri string, object *domain.Object) error {
	buf, err := json.Marshal(object)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return nil
}
