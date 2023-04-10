package iamcore

import (
	"encoding/json"
	"net/http"

	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

const (
	authorizationHeaderName = "Authorization"
	iamcoreGetUserIrnPath   = "api/v1/users/me/irn"
)

type Client struct {
	host string

	httpClient *http.Client
}

func NewClient(opts *Options, httpClient *http.Client) *Client {
	return &Client{
		host: opts.IamcoreURL,

		httpClient: httpClient,
	}
}

func (c *Client) GetUserIRN(authorizationHeader string) (*irn.IRN, error) {
	req, err := http.NewRequest(http.MethodGet, c.host+iamcoreGetUserIrnPath, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set(authorizationHeaderName, authorizationHeader)

	gotResponse, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer gotResponse.Body.Close()

	userIrnResponseDTO := &UserIrnResponseDTO{}
	if err := json.NewDecoder(gotResponse.Body).Decode(&userIrnResponseDTO); err != nil {
		return nil, err
	}

	return irn.NewIRNFromString(userIrnResponseDTO.Data)
}
