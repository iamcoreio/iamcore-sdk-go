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
	// Returns ErrBadRequest error in case of invalid request.
	Authorize(ctx context.Context, authorizationHeader http.Header, accountID, application, tenantID, resourceType, resourcePath string,
		resourceIDs []string, action string) ([]string, error)

	// AuthorizeResources returns existing iamcore resources to which user has ALL the requested actions granted.
	//
	// If the requested resources is empty, the function will return resources having specified resource type to which user has ALL the requested actions granted.
	// If the requested resources is not empty, the function will return requested resources if user has ALL the requested actions granted on ALL resources.
	//
	// Neither passed resources nor actions can contain wildcards.
	// All the resources must have the same type.
	//
	// Returns ErrSDKDisabled error in case SDK is disabled.
	// Returns ErrUnauthenticated error in case of unauthorized access.
	// Returns ErrNotFound error in case of requested resources do not exist in iamcore.
	// Returns ErrForbidden error in case authenticated principal does not have sufficient permissions to requested resources.
	// Returns ErrBadRequest error in case of invalid request.
	AuthorizeResources(ctx context.Context, authorizationHeader http.Header, accountID, application, tenantID, resourceType, resourcePath string,
		resourceIDs []string, action string) ([]string, error)

	// AuthorizationDBQueryFilter retrieves the authorization query filter by database engine.
	//
	// Returns ErrSDKDisabled error in case SDK is disabled.
	// Returns ErrUnauthenticated error in case of unauthorized access.
	// Returns ErrForbidden error in case authenticated principal does not have sufficient permissions to any resources.
	// Returns ErrBadRequest error in case of invalid request.
	AuthorizationDBQueryFilter(ctx context.Context, authorizationHeader http.Header, action, database string) (string, error)

	// EvaluateActionsOnIRNs evaluates a list of actions against a list of IRNs and returns a map
	// associating each action with its corresponding permitted and prohibited IRNs matching requested IRNs.
	//
	// Returns ErrSDKDisabled error in case SDK is disabled.
	// Returns ErrUnauthenticated error in case of unauthorized access.
	// Returns ErrBadRequest error in case of invalid request.
	EvaluateActionsOnIRNs(ctx context.Context, authorizationHeader http.Header, actions []string, irns []*irn.IRN) (map[string]*AllowedAndDeniedIRNs, error)

	// FilterAuthorizedResources filters the list of resources and returns a subset, to which user has the requested action granted within the specified tenant.
	//
	// Neither passed resources nor action can contain wildcards.
	// All the resources must have the same type.
	//
	// Returns ErrSDKDisabled error in case SDK is disabled.
	// Returns ErrUnauthenticated error in case of unauthorized access.
	// Returns ErrBadRequest error in case of invalid request.
	FilterAuthorizedResources(ctx context.Context, authorizationHeader http.Header, accountID, application, tenantID, resourceType, resourcePath string,
		resourceIDs []string, action string) ([]string, error)
}

type AuthorizationFunction func(ctx context.Context, authorizationHeader http.Header, action string, resources []*irn.IRN) error

func (c *сlient) authorize(ctx context.Context, authorizationHeader http.Header, accountID, application, tenantID, resourceType,
	resourcePath string, resourceIDs []string, action string, function AuthorizationFunction) (
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

		if err = function(ctx, authorizationHeader, action, resourceIRNs); err != nil {
			return nil, err
		}

		return resourceIDs, nil
	}

	resourceIRNs, err := c.iamcoreClient.AuthorizedOnResourceType(ctx, authorizationHeader, application, tenantID, resourceType, action)
	if err != nil {
		return nil, err
	}

	return getResourceIDs(resourceIRNs), nil
}

func (c *сlient) Authorize(ctx context.Context, authorizationHeader http.Header, accountID, application,
	tenantID, resourceType, resourcePath string, resourceIDs []string, action string) (
	[]string, error,
) {
	return c.authorize(ctx, authorizationHeader, accountID, application, tenantID,
		resourceType, resourcePath, resourceIDs, action, c.iamcoreClient.AuthorizeOnResources)
}

func (c *сlient) AuthorizeResources(ctx context.Context, authorizationHeader http.Header, accountID, application,
	tenantID, resourceType, resourcePath string, resourceIDs []string, action string) (
	[]string, error,
) {
	return c.authorize(ctx, authorizationHeader, accountID, application, tenantID,
		resourceType, resourcePath, resourceIDs, action, c.iamcoreClient.AuthorizeResources)
}

func (c *сlient) FilterAuthorizedResources(ctx context.Context, authorizationHeader http.Header, accountID, application, tenantID, resourceType,
	resourcePath string, resourceIDs []string, action string) (
	[]string, error,
) {
	if c.disabled {
		return nil, ErrSDKDisabled
	}

	resourceIRNs, err := buildResourceIRNs(accountID, application, tenantID, resourceType, resourcePath, resourceIDs)
	if err != nil {
		return nil, err
	}

	authorizedResources, err := c.iamcoreClient.FilterAuthorizedResources(ctx, authorizationHeader, action, resourceIRNs)
	if err != nil {
		return nil, err
	}

	return getResourceIDs(authorizedResources), nil
}

func (c *сlient) AuthorizationDBQueryFilter(ctx context.Context, authorizationHeader http.Header, action, database string) (string, error) {
	if c.disabled {
		return "", ErrSDKDisabled
	}

	return c.iamcoreClient.AuthorizationDBQueryFilter(ctx, authorizationHeader, action, database)
}

func (c *сlient) EvaluateActionsOnIRNs(ctx context.Context, authorizationHeader http.Header, actions []string, irns []*irn.IRN) (
	map[string]*AllowedAndDeniedIRNs, error,
) {
	if c.disabled {
		return nil, ErrSDKDisabled
	}

	return c.iamcoreClient.EvaluateActionsOnIRNs(ctx, authorizationHeader, actions, irns)
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

func getResourceIDs(resourceIRNs []*irn.IRN) []string {
	resourceIDs := make([]string, len(resourceIRNs))

	for i := range resourceIRNs {
		resourceIDs[i] = resourceIRNs[i].GetResourceID()
	}

	return resourceIDs
}
