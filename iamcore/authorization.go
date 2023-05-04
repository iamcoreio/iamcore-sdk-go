package iamcore

import (
	"context"
	"errors"

	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

var ErrSDKDisabled = errors.New("SDK disabled")

type AuthorizationClient interface {
	// Authorize returns resources to which user has ALL the requested actions granted.
	//
	// The ctx must contain principal's authorization header previously extracted with the Client.WithAuth middleware.
	//
	// If the requested resources is empty, the function will return resources having specified resource type to which user has ALL the requested actions granted.
	// If the requested resources is not empty, the function will return requested resources if user has ALL the requested actions granted on ALL resources.
	//
	// Neither passed resources nor actions can contain wildcards.
	//
	// Returns ErrSDKDisabled error in case SDK is disabled.
	// Returns ErrUnauthenticated error in case of unauthorized access.
	// Returns ErrForbidden error in case authenticated principal does not have sufficient permissions to requested resources.
	Authorize(ctx context.Context, resourceType, resourcePath string, resourceIDs, actions []string) ([]string, error)
}

func (c *сlient) Authorize(ctx context.Context, resourceType, resourcePath string, resourceIDs, actions []string) ([]string, error) {
	if c.disabled {
		return resourceIDs, ErrSDKDisabled
	}

	authorizationHeader, err := PrincipalAuthorizationHeader(ctx)
	if err != nil {
		return nil, err
	}

	if len(resourceIDs) != 0 {
		resourceIRNs, err := buildResourceIRNs(ctx, resourceType, resourcePath, resourceIDs)
		if err != nil {
			return nil, err
		}

		for i := range actions {
			if err = c.iamcoreClient.AuthorizeOnResources(ctx, authorizationHeader, actions[i], resourceIRNs); err != nil {
				return nil, err
			}
		}

		return resourceIDs, nil
	}

	uniqueResourceIDs := make(map[string]bool)
	resourceIDs = make([]string, 0)

	for i := range actions {
		resourceIRNs, err := c.iamcoreClient.AuthorizedOnResourceType(ctx, authorizationHeader, actions[i], resourceType)
		if err != nil {
			return nil, err
		}

		for j := range resourceIRNs {
			resourceID := resourceIRNs[j].GetResourceID()

			if _, contains := uniqueResourceIDs[resourceID]; !contains {
				uniqueResourceIDs[resourceID] = true

				resourceIDs = append(resourceIDs, resourceID)
			}
		}
	}

	return resourceIDs, nil
}

func buildResourceIRNs(ctx context.Context, resourceType, resourcePath string, resourceIDs []string) ([]*irn.IRN, error) {
	resourceIRNs := make([]*irn.IRN, len(resourceIDs))

	accountID, err := AccountID(ctx)
	if err != nil {
		return nil, err
	}

	tenantID, err := TenantID(ctx)
	if err != nil {
		return nil, err
	}

	for i := range resourceIDs {
		resourceIRN, err := irn.NewIRN(accountID, "iamcore", tenantID, nil, resourceType, irn.SplitPath(resourcePath), resourceIDs[i])
		if err != nil {
			return nil, err
		}

		resourceIRNs[i] = resourceIRN
	}

	return resourceIRNs, nil
}
