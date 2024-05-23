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
	resourcePath               = "/api/v1/resources"
	applicationPath            = "/api/v1/applications"
	evaluatePath               = "/api/v1/evaluate"
	evaluateOnResourceTypePath = evaluatePath + "/resources"
	evaluateActionsOnIRNsPath  = evaluatePath + "/irns/actions"
	evaluateDBQueryFilterPath  = evaluatePath + "/database-query-filter"
	pageSize                   = 100000
)

var (
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrForbidden       = errors.New("forbidden")
	ErrConflict        = errors.New("conflict")
	ErrBadRequest      = errors.New("bad request")
	ErrUnknown         = errors.New("unknown error")
)

type ServerClient struct {
	serverURL string

	httpClient *http.Client
}

func NewServerClient(serverURL string, httpClient *http.Client) *ServerClient {
	return &ServerClient{
		serverURL: serverURL,

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

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.getURL(evaluatePath), bytes.NewReader(requestDTO))
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

func (c *ServerClient) FilterAuthorizedResources(ctx context.Context, authorizationHeader http.Header, action string, resources []*irn.IRN) (
	authorizedResources []*irn.IRN, err error,
) {
	requestDTO, err := json.Marshal(&AuthorizedOnResourceListRequestDTO{
		Action:    action,
		Resources: resources,
	})
	if err != nil {
		return nil, err
	}

	url := c.getURL(fmt.Sprintf("%s?filterResources=%t", evaluatePath, true))

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(requestDTO))
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
		if err = json.NewDecoder(response.Body).Decode(&authorizedResources); err != nil {
			return nil, err
		}

		return authorizedResources, nil
	}

	return nil, handleServerErrorResponse(response)
}

func (c *ServerClient) AuthorizedOnResourceType(ctx context.Context, authorizationHeader http.Header, application, tenantID, resourceType, action string) (
	[]*irn.IRN, error,
) {
	authorizedOnResourceTypeRequestDTO := &AuthorizedOnResourceTypeRequestDTO{
		Action:       action,
		ResourceType: resourceType,
		Application:  application,
	}

	if tenantID != "" {
		authorizedOnResourceTypeRequestDTO.TenantID = tenantID
	}

	requestDTO, err := json.Marshal(authorizedOnResourceTypeRequestDTO)
	if err != nil {
		return nil, err
	}

	url := c.getURL(fmt.Sprintf("%s?pageSize=%d", evaluateOnResourceTypePath, pageSize))

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(requestDTO))
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

func (c *ServerClient) CreateResourceType(ctx context.Context, authorizationHeader http.Header,
	applicationIRN *irn.IRN, createDTO *CreateResourceTypeRequestDTO,
) error {
	url := c.getURL(fmt.Sprintf("%s/%s/resource-types", applicationPath, applicationIRN.Base64()))

	requestDTO, err := json.Marshal(createDTO)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(requestDTO))
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

func (c *ServerClient) GetResourceTypes(ctx context.Context, authorizationHeader http.Header, applicationIRN *irn.IRN) ([]*ResourceTypeResponseDTO, error) {
	url := c.getURL(fmt.Sprintf("%s/%s/resource-types?pageSize=%d", applicationPath, applicationIRN.Base64(), pageSize))

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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
		var responseDTO struct {
			Data []*ResourceTypeResponseDTO `json:"data"`
		}

		if err = json.NewDecoder(response.Body).Decode(&responseDTO); err != nil {
			return nil, err
		}

		return responseDTO.Data, nil
	}

	return nil, handleServerErrorResponse(response)
}

func (c *ServerClient) EvaluateActionsOnIRNs(ctx context.Context, authorizationHeader http.Header, actions []string, irns []*irn.IRN) (
	map[string]*AllowedAndDeniedIRNs, error,
) {
	evaluateActionsRequestDTO := &EvaluateActionsOnIRNsRequestDTO{
		IRNs:    irns,
		Actions: actions,
	}

	requestDTO, err := json.Marshal(evaluateActionsRequestDTO)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.getURL(evaluateActionsOnIRNsPath), bytes.NewReader(requestDTO))
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
		responseDTO := map[string]*AllowedAndDeniedIRNs{}
		if err = json.NewDecoder(response.Body).Decode(&responseDTO); err != nil {
			return nil, err
		}

		return responseDTO, nil
	}

	return nil, handleServerErrorResponse(response)
}

func (c *ServerClient) AuthorizationDBQueryFilter(ctx context.Context, authorizationHeader http.Header, action, database string) (string, error) {
	queryFilterRequestDTO := &QueryFilterOnEvaluatedResourcesRequestDTO{
		Action:   action,
		Database: database,
	}

	requestDTO, err := json.Marshal(queryFilterRequestDTO)
	if err != nil {
		return "", err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.getURL(evaluateDBQueryFilterPath), bytes.NewReader(requestDTO))
	if err != nil {
		return "", err
	}

	request.Header = authorizationHeader

	response, err := c.httpClient.Do(request)
	if err != nil {
		return "", err
	}

	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		var responseDTO struct {
			Data string `json:"data"`
		}

		if err = json.NewDecoder(response.Body).Decode(&responseDTO); err != nil {
			return "", err
		}

		return responseDTO.Data, nil
	}

	return "", handleServerErrorResponse(response)
}

func (c *ServerClient) getURL(path string) string {
	return fmt.Sprintf("%s%s", c.serverURL, path)
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
	case http.StatusBadRequest:
		return fmt.Errorf("%s: %w", responseDTO.Message, ErrBadRequest)
	default:
		return fmt.Errorf("%s: %w", responseDTO.Message, ErrUnknown)
	}
}
