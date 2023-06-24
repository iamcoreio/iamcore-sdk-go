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
	userIRNPath                = "/api/v1/users/me/irn"
	evaluateOnResources        = "/api/v1/evaluate"
	evaluateOnResourceTypePath = "/api/v1/evaluate/resources"
	resourcePath               = "/api/v1/resources"
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrForbidden       = errors.New("forbidden")
	ErrConflict        = errors.New("conflict")
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

func (c *ServerClient) AuthorizedOnResourceType(ctx context.Context, authorizationHeader http.Header, application, resourceType, action string) (
	[]*irn.IRN, error,
) {
	requestDTO, err := json.Marshal(&AuthorizedOnResourceTypeRequestDTO{
		Action:       action,
		ResourceType: resourceType,
		Application:  application,
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

func (c *ServerClient) CreateResource(ctx context.Context, authorizationHeader http.Header, createResourceDTO CreateResourceRequestDTO) error {
	requestDTO, err := json.Marshal(createResourceDTO)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.getURL(resourcePath), bytes.NewReader(requestDTO))
	if err != nil {
		return err
	}

	request.Header = authorizationHeader

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusCreated {
		return nil
	}

	return handleServerErrorResponse(response)
}

func (c *ServerClient) DeleteResource(ctx context.Context, authorizationHeader http.Header, resourceIRN *irn.IRN) error {
	url := fmt.Sprintf("%s/%s", c.getURL(resourcePath), resourceIRN.Base64())

	request, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	request.Header = authorizationHeader

	response, err := c.httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusNoContent {
		return nil
	}

	return handleServerErrorResponse(response)
}

func (c *ServerClient) getURL(path string) string {
	return fmt.Sprintf("https://%s%s", c.host, path)
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
	case http.StatusConflict:
		return fmt.Errorf("%s: %w", responseDTO.Message, ErrConflict)
	default:
		return fmt.Errorf("%s: %w", responseDTO.Message, ErrUnknown)
	}
}
