package iamcore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

const (
	userIRNPath                = "api/v1/users/me/irn"
	evaluateOnResources        = "api/v1/evaluate"
	evaluateOnResourceTypePath = "api/v1/evaluate/resources"
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrForbidden       = errors.New("forbidden")
	ErrUnknown         = errors.New("unknown error")
)

type ServerClient struct {
	host string

	httpClient *http.Client
}

func NewServerClient(host string, httpClient *http.Client) *ServerClient {
	return &ServerClient{
		host: host,

		httpClient: httpClient,
	}
}

func (c *ServerClient) GetPrincipalIRN(ctx context.Context, authorizationHeader http.Header) (*irn.IRN, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, c.getURL(userIRNPath), nil)
	if err != nil {
		return nil, err
	}

	request.Header = authorizationHeader

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		responseDTO := &PrincipalIRNResponseDTO{}
		if err = json.NewDecoder(response.Body).Decode(&responseDTO); err != nil {
			return nil, err
		}

		return responseDTO.Data, nil
	}

	return nil, handleServerErrorResponse(response)
}

func (c *ServerClient) AuthorizeOnResources(ctx context.Context, authorizationHeader http.Header, action string, resources []*irn.IRN) error {
	requestDTO, err := json.Marshal(&AuthorizedOnResourceListRequestDTO{
		Action:    action,
		Resources: resources,
	})
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.getURL(evaluateOnResources), bytes.NewReader(requestDTO))
	if err != nil {
		return err
	}

	request.Header = authorizationHeader

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		return nil
	}

	return handleServerErrorResponse(response)
}

func (c *ServerClient) AuthorizedOnResourceType(ctx context.Context, authorizationHeader http.Header, action, resourceType string) ([]*irn.IRN, error) {
	requestDTO, err := json.Marshal(&AuthorizedOnResourceTypeRequestDTO{
		Action:       action,
		ResourceType: resourceType,
	})
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.getURL(evaluateOnResourceTypePath), bytes.NewReader(requestDTO))
	if err != nil {
		return nil, err
	}

	request.Header = authorizationHeader

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		responseDTO := &AuthorizedOnResourceTypeResponseDTO{}
		if err = json.NewDecoder(response.Body).Decode(&responseDTO); err != nil {
			return nil, err
		}

		return responseDTO.Data, nil
	}

	return nil, handleServerErrorResponse(response)
}

func (c *ServerClient) getURL(path string) string {
	return fmt.Sprintf("https://%s/%s", c.host, path)
}

func handleServerErrorResponse(response *http.Response) error {
	responseDTO := &ErrorResponseDTO{}
	if err := json.NewDecoder(response.Body).Decode(&responseDTO); err != nil {
		return err
	}

	switch response.StatusCode {
	case http.StatusUnauthorized:
		return fmt.Errorf("%s: %w", responseDTO.Message, ErrUnauthenticated)
	case http.StatusForbidden:
		return fmt.Errorf("%s: %w", responseDTO.Message, ErrForbidden)
	default:
		return fmt.Errorf("%s: %w", responseDTO.Message, ErrUnknown)
	}
}
