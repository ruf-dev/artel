package trello

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/ruf-dev/artel/internal/domain"
	"github.com/ruf-dev/artel/internal/utils"
	"go.redsock.ru/rerrors"
	"google.golang.org/grpc/codes"
)

const baseURL = "https://api.trello.com/1"

var ErrUnauthorized = rerrors.New("invalid trello credentials", codes.Unauthenticated)

type Client interface {
	GetMember(ctx context.Context) (domain.TrelloMember, error)
	ListBoards(ctx context.Context) ([]domain.TrelloBoard, error)
}

type TaskTrackerCredentials interface {
	GetClient() Client
}

type TrelloCredentials struct {
	ApiKey   string
	ApiToken string
}

func (c TrelloCredentials) GetClient() Client {
	return New(c.ApiKey, c.ApiToken)
}

type client struct {
	apiKey   string
	apiToken string
}

func New(apiKey, apiToken string) Client {
	return &client{apiKey: apiKey, apiToken: apiToken}
}

func (c *client) GetMember(ctx context.Context) (domain.TrelloMember, error) {
	u := fmt.Sprintf("%s/members/me?key=%s&token=%s", baseURL, url.QueryEscape(c.apiKey), url.QueryEscape(c.apiToken))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return domain.TrelloMember{}, rerrors.Wrap(err, "error building trello request")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return domain.TrelloMember{}, rerrors.Wrap(err, "error calling trello")
	}
	defer utils.CloseWithLog(resp.Body, "trello response")

	if resp.StatusCode == http.StatusUnauthorized {
		return domain.TrelloMember{}, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return domain.TrelloMember{}, rerrors.New(
			fmt.Sprintf("trello returned %d: %s", resp.StatusCode, string(body)),
			codes.Internal,
		)
	}

	var m member

	err = json.NewDecoder(resp.Body).Decode(&m)
	if err != nil {
		return domain.TrelloMember{}, rerrors.Wrap(err, "error decoding trello member")
	}

	return domain.TrelloMember{FullName: m.FullName}, nil
}

func (c *client) ListBoards(ctx context.Context) ([]domain.TrelloBoard, error) {
	u := fmt.Sprintf(
		"%s/members/me/boards?key=%s&token=%s&fields=id,name",
		baseURL,
		url.QueryEscape(c.apiKey),
		url.QueryEscape(c.apiToken),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, rerrors.Wrap(err, "error building trello request")
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, rerrors.Wrap(err, "error calling trello")
	}
	defer utils.CloseWithLog(resp.Body, "trello response")

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, ErrUnauthorized
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)

		return nil, rerrors.New(fmt.Sprintf("trello returned %d: %s", resp.StatusCode, string(body)), codes.Internal)
	}

	var boards []board

	err = json.NewDecoder(resp.Body).Decode(&boards)
	if err != nil {
		return nil, rerrors.Wrap(err, "error decoding trello boards")
	}

	out := make([]domain.TrelloBoard, len(boards))
	for i, b := range boards {
		out[i] = domain.TrelloBoard{Id: b.Id, Name: b.Name}
	}

	return out, nil
}
