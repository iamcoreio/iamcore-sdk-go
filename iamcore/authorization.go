package iamcore

import (
	"context"
	"errors"
	"net/http"

	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

var ErrSDKDisabled = errors.New("SDK disabled")

type AuthorizationClient interface {
	// Authorize returns resources to which user has ALL the requested actions granted.
	//
	// If the requested resources is empty, the function will return resources having specified resource type to which user has ALL the requested actions granted.
	// If the requested resources is not empty, the function will return requested resources if user has ALL the requested actions granted on ALL resources.
	//
	// Neither passed resources nor actions can contain wildcards.
	// All the resources must have the same type.
	//
	// Returns ErrSDKDisabled error in case SDK is disabled.
	// Returns ErrUnauthenticated error in case of unauthorized access.
	// Returns ErrForbidden error in case authenticated principal does not have sufficient permissions to requested resources.
	Authorize(ctx context.Context, authorizationHeader http.Header, accountID, application, tenantID, resourceType, resourcePath string,
		resourceIDs []string, action string) ([]string, error)
}

func (c *сlient) Authorize(ctx context.Context, authorizationHeader http.Header, accountID, application, tenantID, resourceType,
	resourcePath string, resourceIDs []string, action string) (
	[]string, error,
) {
	if c.disabled {
		return nil, ErrSDKDisabled
	}

	if len(resourceIDs) != 0 {
		resourceIRNs, err := buildResourceIRNs(accountID, application, tenantID, resourceType, resourcePath, resourceIDs)
		if err != nil {
			return nil, err
		}

		if err = c.iamcoreClient.AuthorizeOnResources(ctx, authorizationHeader, action, resourceIRNs); err != nil {
			return nil, err
		}

		return resourceIDs, nil
	}

	resourceIRNs, err := c.iamcoreClient.AuthorizedOnResourceType(ctx, authorizationHeader, tenantID, application, resourceType, action)
	if err != nil {
		return nil, err
	}

	resourceIDs = make([]string, len(resourceIRNs))

	for i := range resourceIRNs {
		resourceIDs[i] = resourceIRNs[i].GetResourceID()
	}

	return resourceIDs, nil
}

func buildResourceIRNs(accountID, application, tenantID, resourceType, resourcePath string, resourceIDs []string) ([]*irn.IRN, error) {
	resourceIRNs := make([]*irn.IRN, len(resourceIDs))

	for i := range resourceIDs {
		resourceIRN, err := irn.NewIRN(accountID, application, tenantID, nil, resourceType, irn.SplitPath(resourcePath), resourceIDs[i])
		if err != nil {
			return nil, err
		}

		resourceIRNs[i] = resourceIRN
	}

	return resourceIRNs, nil
}
