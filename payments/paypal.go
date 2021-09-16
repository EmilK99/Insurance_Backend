package payments

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flight_app/app/api"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type expirationTime int64

// RequestNewTokenBeforeExpiresIn is used by SendWithAuth and try to get new Token when it's about to expire
const RequestNewTokenBeforeExpiresIn = time.Duration(60) * time.Second

// Client represents a Paypal REST API Client
type Client struct {
	sync.Mutex
	Client               *http.Client
	ClientID             string
	Secret               string
	APIBase              string
	Log                  io.Writer // If user set log file name all requests will be logged there
	Token                *TokenResponse
	tokenExpiresAt       time.Time
	returnRepresentation bool
}

// Authorization struct
type Authorization struct {
	ID               string                `json:"id,omitempty"`
	CustomID         string                `json:"custom_id,omitempty"`
	InvoiceID        string                `json:"invoice_id,omitempty"`
	Status           string                `json:"status,omitempty"`
	StatusDetails    *CaptureStatusDetails `json:"status_details,omitempty"`
	Amount           *PurchaseUnitAmount   `json:"amount,omitempty"`
	SellerProtection *SellerProtection     `json:"seller_protection,omitempty"`
	CreateTime       *time.Time            `json:"create_time,omitempty"`
	UpdateTime       *time.Time            `json:"update_time,omitempty"`
	ExpirationTime   *time.Time            `json:"expiration_time,omitempty"`
	Links            []api.Link            `json:"links,omitempty"`
}

// TokenResponse is for API response for the /oauth2/token endpoint
type TokenResponse struct {
	RefreshToken string         `json:"refresh_token"`
	Token        string         `json:"access_token"`
	Type         string         `json:"token_type"`
	ExpiresIn    expirationTime `json:"expires_in"`
}

// NewClient returns new Client struct
// APIBase is a base API URL, for testing you can use paypal.APIBaseSandBox
func NewClient(clientID string, secret string, APIBase string) (*Client, error) {
	if clientID == "" || secret == "" || APIBase == "" {
		return nil, errors.New("ClientID, Secret and APIBase are required to create a Client")
	}

	return &Client{
		Client:   &http.Client{},
		ClientID: clientID,
		Secret:   secret,
		APIBase:  APIBase,
	}, nil
}

// GetAuthorization returns an authorization by ID
// Endpoint: GET /v2/payments/authorizations/ID
func (c *Client) GetAuthorization(ctx context.Context, authID string) (*Authorization, error) {
	buf := bytes.NewBuffer([]byte(""))
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s%s%s", c.APIBase, "/v2/payments/authorizations/", authID), buf)
	auth := &Authorization{}

	if err != nil {
		return auth, err
	}

	err = c.SendWithAuth(req, auth)
	return auth, err
}

// SendWithAuth makes a request to the API and apply OAuth2 header automatically.
// If the access token soon to be expired or already expired, it will try to get a new one before
// making the main request
// client.Token will be updated when changed
func (c *Client) SendWithAuth(req *http.Request, v interface{}) error {
	c.Lock()
	// Note: Here we do not want to `defer c.Unlock()` because we need `c.Send(...)`
	// to happen outside of the locked section.

	if c.Token != nil {
		if !c.tokenExpiresAt.IsZero() && c.tokenExpiresAt.Sub(time.Now()) < RequestNewTokenBeforeExpiresIn {
			// c.Token will be updated in GetAccessToken call
			if _, err := c.GetAccessToken(req.Context()); err != nil {
				c.Unlock()
				return err
			}
		}

		req.Header.Set("Authorization", "Bearer "+c.Token.Token)
	}

	// Unlock the client mutex before sending the request, this allows multiple requests
	// to be in progress at the same time.
	c.Unlock()
	return c.Send(req, v)
}

// GetAccessToken returns struct of TokenResponse
// No need to call SetAccessToken to apply new access token for current Client
// Endpoint: POST /v1/oauth2/token
func (c *Client) GetAccessToken(ctx context.Context) (*TokenResponse, error) {
	buf := bytes.NewBuffer([]byte("grant_type=client_credentials"))
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s%s", c.APIBase, "/v1/oauth2/token"), buf)
	if err != nil {
		return &TokenResponse{}, err
	}

	req.Header.Set("Content-type", "application/x-www-form-urlencoded")

	response := &TokenResponse{}
	err = c.SendWithBasicAuth(req, response)

	// Set Token fur current Client
	if response.Token != "" {
		c.Token = response
		c.tokenExpiresAt = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)
	}

	return response, err
}

// SendWithBasicAuth makes a request to the API using clientID:secret basic auth
func (c *Client) SendWithBasicAuth(req *http.Request, v interface{}) error {
	req.SetBasicAuth(c.ClientID, c.Secret)

	return c.Send(req, v)
}

// Send makes a request to the API, the response body will be
// unmarshaled into v, or if v is an io.Writer, the response will
// be written to it without decoding
func (c *Client) Send(req *http.Request, v interface{}) error {
	var (
		err  error
		resp *http.Response
		data []byte
	)

	// Set default headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en_US")

	// Default values for headers
	if req.Header.Get("Content-type") == "" {
		req.Header.Set("Content-type", "application/json")
	}
	if c.returnRepresentation {
		req.Header.Set("Prefer", "return=representation")
	}

	resp, err = c.Client.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		errResp := &api.ErrorResponse{Response: resp}
		data, err = ioutil.ReadAll(resp.Body)

		if err == nil && len(data) > 0 {
			json.Unmarshal(data, errResp)
		}

		return err
	}
	if v == nil {
		return nil
	}

	if w, ok := v.(io.Writer); ok {
		io.Copy(w, resp.Body)
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(v)
}
