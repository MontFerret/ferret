package network

import (
	"context"
	"encoding/json"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/protocol/network"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"

	"github.com/MontFerret/ferret/pkg/drivers"
)

type Manager struct {
	client  *cdp.Client
	logger  *zerolog.Logger
	headers drivers.HTTPHeaders
}

func New(
	logger *zerolog.Logger,
	client *cdp.Client,
) (*Manager, error) {
	evt, err := client.Network.RequestWillBeSent()

	if err != nil {
		return nil, err
	}

	m, err := evt.Recv()

	if err != nil {
		return nil, err
	}

	return &Manager{client, logger, make(drivers.HTTPHeaders)}, nil
}

func (m *Manager) GetCookies(ctx context.Context) (drivers.HTTPCookies, error) {
	repl, err := m.client.Network.GetAllCookies(ctx)

	if err != nil {
		return nil, errors.Wrap(err, "failed to get cookies")
	}

	cookies := make(drivers.HTTPCookies)

	if repl.Cookies == nil {
		return cookies, nil
	}

	for _, c := range repl.Cookies {
		cookies[c.Name] = toDriverCookie(c)
	}

	return cookies, nil
}

func (m *Manager) SetCookies(ctx context.Context, url string, cookies drivers.HTTPCookies) error {
	if len(cookies) == 0 {
		return nil
	}

	params := make([]network.CookieParam, 0, len(cookies))

	for _, c := range cookies {
		params = append(params, fromDriverCookie(url, c))
	}

	return m.client.Network.SetCookies(ctx, network.NewSetCookiesArgs(params))
}

func (m *Manager) DeleteCookies(ctx context.Context, url string, cookies drivers.HTTPCookies) error {
	if len(cookies) == 0 {
		return nil
	}

	var err error

	for _, c := range cookies {
		err = m.client.Network.DeleteCookies(ctx, fromDriverCookieDelete(url, c))

		if err != nil {
			break
		}
	}

	return err
}

func (m *Manager) GetHeaders(_ context.Context) (drivers.HTTPHeaders, error) {
	copied := make(drivers.HTTPHeaders)

	for k, v := range m.headers {
		copied[k] = v
	}

	return copied, nil
}

func (m *Manager) SetHeaders(ctx context.Context, headers drivers.HTTPHeaders) error {
	if len(headers) == 0 {
		return nil
	}

	m.headers = headers

	j, err := json.Marshal(headers)

	if err != nil {
		return errors.Wrap(err, "failed to marshal headers")
	}

	err = m.client.Network.SetExtraHTTPHeaders(
		ctx,
		network.NewSetExtraHTTPHeadersArgs(j),
	)

	if err != nil {
		return errors.Wrap(err, "failed to set headers")
	}

	return nil
}

func (m *Manager) WaitForPageLoad(ctx context.Context, timeout int) error {}
