package iamcore

import (
	"context"

	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

type ResourceManager interface {
	// CreateResource creates resource on iamcore using the API key.
	//
	// Returns ErrSDKDisabled error in case SDK is disabled.
	// Returns ErrUnauthenticated error in case of unauthenticated access.
	// Returns ErrConflict error in case duplicated resource found.
	// Returns ErrForbidden error in case authenticated principal does not have sufficient permissions to create the resource.
	// Returns ErrUnknown error in case of unexpected response from iamcore server.
	CreateResource(ctx context.Context, resourceType, resourcePath, resourceID string) error

	// DeleteResource deletes resource on iamcore using the API key.
	//
	// Returns ErrSDKDisabled error in case SDK is disabled.
	// Returns ErrUnauthenticated error in case of unauthenticated access.
	// Returns ErrForbidden error in case authenticated principal does not have sufficient permissions to delete the resource.
	// Returns ErrUnknown error in case of unexpected response from iamcore server.
	DeleteResource(ctx context.Context, resourceType, resourcePath, resourceID string) error
}

func (c *сlient) CreateResource(ctx context.Context, resourceType, resourcePath, resourceID string) error {
	if c.disabled {
		return ErrSDKDisabled
	}

	if resourcePath == "" {
		resourcePath = "/"
	}

	tenantID, err := TenantID(ctx)
	if err != nil {
		return err
	}

	createResourceRequestDTO := CreateResourceRequestDTO{
		Name:         resourceID,
		ResourceType: resourceType,
		Path:         resourcePath,
		Enabled:      true,
		TenantID:     tenantID,
	}

	return c.iamcoreClient.CreateResource(ctx, c.GetAPIKeyAuthorizationHeader(), createResourceRequestDTO)
}

func (c *сlient) DeleteResource(ctx context.Context, resourceType, resourcePath, resourceID string) error {
	if c.disabled {
		return ErrSDKDisabled
	}

	tenantID, err := TenantID(ctx)
	if err != nil {
		return err
	}

	principalIRN, err := c.iamcoreClient.GetPrincipalIRN(ctx, c.GetAPIKeyAuthorizationHeader())
	if err != nil {
		return err
	}

	resourceIRN, err := irn.NewIRN(principalIRN.GetAccountID(), "iamcore", tenantID, nil, resourceType, irn.SplitPath(resourcePath), resourceID)
	if err != nil {
		return err
	}

	return c.iamcoreClient.DeleteResource(ctx, c.GetAPIKeyAuthorizationHeader(), resourceIRN)
}