package iamcore

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

const (
	userIRNPath = "api/v1/users/me/irn"
)

type Client struct {
	host string

	httpClient *http.Client
}

func NewClient(host string, httpClient *http.Client) *Client {
	return &Client{
		host: host,

		httpClient: httpClient,
	}
}

func (c *Client) GetPrincipalIRN(headerName, headerValue string) (*irn.IRN, error) {
	request, err := http.NewRequest(http.MethodGet, c.GetURL(userIRNPath), nil)
	if err != nil {
		return nil, err
	}

	request.Header = http.Header{
		headerName: {headerValue},
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	responseDTO := &UserIRNResponseDTO{}
	if err := json.NewDecoder(response.Body).Decode(&responseDTO); err != nil {
		return nil, err
	}

	return irn.NewIRNFromString(responseDTO.Data)
}

func (c *Client) GetURL(path string) string {
	return fmt.Sprintf("https://%s/%s", c.host, path)
}
