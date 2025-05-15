package itmoschedule

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/pkg/errors"

	"github.com/hexarchy/itmo-calendar/internal/entities"
)

const (
	_basePath = "/schedule/schedule/personal"
)

// Client is schedule API adapter.
type Client struct {
	client  *http.Client
	baseURL string
}

// New creates new Client.
func New(baseURL string) *Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	httpClient := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	return &Client{
		client:  httpClient,
		baseURL: baseURL,
	}
}

// Get fetches schedule using provided token.
func (c *Client) Get(ctx context.Context, token string, from, to time.Time) ([]entities.DaySchedule, error) {
	req, err := c.buildRequest(ctx, token, from, to)
	if err != nil {
		return nil, errors.Wrap(err, "build request")
	}

	respData, err := c.executeRequest(req)
	if err != nil {
		return nil, errors.Wrap(err, "execute request")
	}

	result, err := c.transformResponse(respData)
	if err != nil {
		return nil, errors.Wrap(err, "transform response")
	}

	return result, nil
}

// buildRequest creates an HTTP request for the schedule API.
func (c *Client) buildRequest(ctx context.Context, token string, from, to time.Time) (*http.Request, error) {
	fromStr := from.Format("2006-01-02")
	toStr := to.Format("2006-01-02")

	u, err := url.Parse(c.baseURL + _basePath)
	if err != nil {
		return nil, errors.Wrap(err, "parse URL")
	}

	q := u.Query()
	q.Add("date_start", fromStr)
	q.Add("date_end", toStr)
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, errors.Wrap(err, "create request")
	}

	req.Header.Set("Authorization", "Bearer "+token)

	return req, nil
}

// executeRequest sends the request and parses the response.
func (c *Client) executeRequest(req *http.Request) (*scheduleResponse, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "send request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var response scheduleResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, errors.Wrap(err, "decode response")
	}

	return &response, nil
}

// transformResponse converts DTO to domain entities.
func (c *Client) transformResponse(response *scheduleResponse) ([]entities.DaySchedule, error) {
	result := make([]entities.DaySchedule, 0, len(response.Data))

	for _, day := range response.Data {
		daySchedule, err := c.transformDay(day)
		if err != nil {
			return nil, errors.Wrapf(err, "transform day %s", day.Date)
		}

		result = append(result, daySchedule)
	}

	return result, nil
}
