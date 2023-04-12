package sdk

import (
	"context"
)

type AuthorizationClient interface {
	// AuthorizedResources returns resources having specified type to which user has ALL the requested actions granted.
	// Passed actions cannot contain wildcards.
	AuthorizedResources(ctx context.Context, resourceType string, actions []string) ([]string, error)

	// Authorized checks if user has ALL the requested actions granted on ALL specified resources.
	// Neither passed resources nor actions can contain wildcards.
	Authorized(ctx context.Context, resources, actions []string) (bool, error)
}

func (c *client) Authorized(ctx context.Context, resources, actions []string) (bool, error) {
	return false, nil
}

func (c *client) AuthorizedResources(ctx context.Context, resourceType string, actions []string) ([]string, error) {
	return nil, nil
}
