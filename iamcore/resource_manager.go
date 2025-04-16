package iamcore

import (
	"context"
	"net/http"

	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

type ResourceManager interface {
	// CreateResource creates resource on iamcore.
	//
	// Returns ErrSDKDisabled error in case SDK is disabled.
	// Returns ErrUnauthenticated error in case of unauthenticated access.
	// Returns ErrConflict error in case duplicated resource found.
	// Returns ErrForbidden error in case authenticated principal does not have sufficient permissions to create the resource.
	// Returns ErrBadRequest error in case of invalid request.
	// Returns ErrUnknown error in case of unexpected response from iamcore server.
	CreateResource(ctx context.Context, authorizationHeader http.Header, application, tenantID, resourceType, resourcePath, resourceID string) error

	// DeleteResource deletes resource on iamcore.
	//
	// Returns ErrSDKDisabled error in case SDK is disabled.
	// Returns ErrUnauthenticated error in case of unauthenticated access.
	// Returns ErrForbidden error in case authenticated principal does not have sufficient permissions to delete the resource.
	// Returns ErrBadRequest error in case of invalid request.
	// Returns ErrUnknown error in case of unexpected response from iamcore server.
	DeleteResource(ctx context.Context, authorizationHeader http.Header, application, tenantID, resourceType, resourcePath, resourceID string) error

	// CreateResourceType creates a new resource type for application on iamcore.
	//
	// Returns ErrSDKDisabled error in case SDK is disabled.
	// Returns ErrUnauthenticated error in case of unauthenticated access.
	// Returns ErrForbidden error in case authenticated principal does not have sufficient permissions to create resource type.
	// Returns ErrBadRequest error in case of invalid request.
	// Returns ErrUnknown error in case of unexpected response from iamcore server.
	CreateResourceType(ctx context.Context, authorizationHeader http.Header, accountID, application, resourceType, actionPrefix string, operations []string) error

	// GetResourceTypes retrieves application`s resource types on iamcore.
	//
	// Returns ErrSDKDisabled error in case SDK is disabled.
	// Returns ErrUnauthenticated error in case of unauthenticated access.
	// Returns ErrForbidden error in case authenticated principal does not have sufficient permissions to read resource types.
	// Returns ErrBadRequest error in case of invalid request.
	// Returns ErrUnknown error in case of unexpected response from iamcore server.
	GetResourceTypes(ctx context.Context, authorizationHeader http.Header, accountID, application string) ([]*ResourceTypeResponseDTO, error)

	// AttachUserToPolicy attaches authenticated principal to policy.
	//
	// Returns ErrSDKDisabled error in case SDK is disabled.
	// Returns ErrUnauthenticated error in case of unauthenticated access.
	// Returns ErrForbidden error in case authenticated principal does not have sufficient permissions to read resource types.
	// Returns ErrBadRequest error in case of invalid request.
	// Returns ErrUnknown error in case of unexpected response from iamcore server.
	AttachUserToPolicy(ctx context.Context, authorizationHeader http.Header, application, tenantID, resourceType, policyID string, userIRN *irn.IRN) error
}

func (c *client) CreateResource(ctx context.Context, authorizationHeader http.Header, application, tenantID, resourceType, resourcePath, resourceID string,
) error {
	if c.disabled {
		return ErrSDKDisabled
	}

	if resourcePath == "" {
		resourcePath = "/"
	}

	createResourceRequestDTO := CreateResourceRequestDTO{
		Name:         resourceID,
		Application:  application,
		ResourceType: resourceType,
		Path:         resourcePath,
		Enabled:      true,
		TenantID:     tenantID,
	}

	return c.iamcoreClient.CreateResource(ctx, authorizationHeader, createResourceRequestDTO)
}

func (c *client) DeleteResource(ctx context.Context, authorizationHeader http.Header, application, tenantID, resourceType, resourcePath, resourceID string,
) error {
	if c.disabled {
		return ErrSDKDisabled
	}

	principalIRN, err := c.iamcoreClient.GetPrincipalIRN(ctx, authorizationHeader)
	if err != nil {
		return err
	}

	resourceIRN, err := irn.NewIRN(principalIRN.GetAccountID(), application, tenantID, nil, resourceType, irn.SplitPath(resourcePath), resourceID)
	if err != nil {
		return err
	}

	return c.iamcoreClient.DeleteResource(ctx, authorizationHeader, resourceIRN)
}

func (c *client) AttachUserToPolicy(ctx context.Context, authorizationHeader http.Header, application, tenantID, resourceType, policyID string, userIRN *irn.IRN) error {
	if c.disabled {
		return ErrSDKDisabled
	}

	policyIRN, err := irn.NewIRN(userIRN.GetAccountID(), application, tenantID, nil, resourceType, nil, policyID)
	if err != nil {
		return err
	}

	return c.iamcoreClient.AttachUserToPolicy(ctx, authorizationHeader, userIRN, policyIRN)
}

func (c *client) CreateResourceType(ctx context.Context, authorizationHeader http.Header, accountID, application, resourceType,
	actionPrefix string, operations []string,
) error {
	if c.disabled {
		return ErrSDKDisabled
	}

	applicationIRN, err := irn.NewIRN(accountID, "iamcore", "", nil, "application", nil, application)
	if err != nil {
		return err
	}

	requestDTO := &CreateResourceTypeRequestDTO{
		Type:         resourceType,
		ActionPrefix: actionPrefix,
		Operations:   dropEmptyStrings(operations),
	}

	return c.iamcoreClient.CreateResourceType(ctx, authorizationHeader, applicationIRN, requestDTO)
}

func (c *client) GetResourceTypes(ctx context.Context, authorizationHeader http.Header, accountID, application string) ([]*ResourceTypeResponseDTO, error) {
	if c.disabled {
		return nil, ErrSDKDisabled
	}

	applicationIRN, err := irn.NewIRN(accountID, "iamcore", "", nil, "application", nil, application)
	if err != nil {
		return nil, err
	}

	return c.iamcoreClient.GetResourceTypes(ctx, authorizationHeader, applicationIRN)
}

// dropEmptyStrings filters a slice to only non-empty strings.
func dropEmptyStrings(in []string) []string {
	out := make([]string, 0)

	if in == nil {
		return out
	}

	for i := range in {
		if in[i] != "" {
			out = append(out, in[i])
		}
	}

	return out
}
