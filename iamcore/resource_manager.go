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
	// Returns ErrUnknown error in case of unexpected response from iamcore server.
	CreateResource(ctx context.Context, authorizationHeader http.Header, application, tenantID, resourceType, resourcePath, resourceID string) error

	// DeleteResource deletes resource on iamcore.
	//
	// Returns ErrSDKDisabled error in case SDK is disabled.
	// Returns ErrUnauthenticated error in case of unauthenticated access.
	// Returns ErrForbidden error in case authenticated principal does not have sufficient permissions to delete the resource.
	// Returns ErrUnknown error in case of unexpected response from iamcore server.
	DeleteResource(ctx context.Context, authorizationHeader http.Header, application, tenantID, resourceType, resourcePath, resourceID string) error
}

func (c *сlient) CreateResource(ctx context.Context, authorizationHeader http.Header, application, tenantID, resourceType, resourcePath, resourceID string,
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

func (c *сlient) DeleteResource(ctx context.Context, authorizationHeader http.Header, application, tenantID, resourceType, resourcePath, resourceID string,
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
