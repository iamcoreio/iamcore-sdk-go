package iamcore

import (
	"time"

	"gitlab.kaaiot.net/core/lib/iamcore/irn.git"
)

type PrincipalIRNResponseDTO struct {
	Data *irn.IRN `json:"data"`
}

type ErrorResponseDTO struct {
	Message string `json:"message"`
}

type AuthorizedOnResourceTypeRequestDTO struct {
	Action       string `json:"action"`
	ResourceType string `json:"resourceType"`
	Application  string `json:"application"`
	TenantID     string `json:"tenantID"`
}

type AuthorizedOnResourceTypeResponseDTO struct {
	Data []*irn.IRN `json:"data"`
}

type AuthorizedOnResourceListRequestDTO struct {
	Resources []*irn.IRN `json:"resources"`
	Action    string     `json:"action"`
}

type CreateResourceRequestDTO struct {
	Name         string `json:"name"`
	Application  string `json:"application"`
	Path         string `json:"path"`
	ResourceType string `json:"resourceType"`
	Enabled      bool   `json:"enabled"`
	TenantID     string `json:"tenantID"`
}

type CreateResourceTypeRequestDTO struct {
	Type         string   `json:"type"`
	Description  string   `json:"description"`
	ActionPrefix string   `json:"actionPrefix"`
	Operations   []string `json:"operations"`
}

type ResourceTypeResponseDTO struct {
	ID           string    `json:"id"`
	IRN          *irn.IRN  `json:"irn"`
	Type         string    `json:"type"`
	Description  string    `json:"description"`
	ActionPrefix string    `json:"actionPrefix"`
	Operations   []string  `json:"operations"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}

type EvaluateActionsOnIRNsRequestDTO struct {
	IRNs    []*irn.IRN `json:"irns"`
	Actions []string   `json:"actions"`
}

type QueryFilterOnEvaluatedResourcesRequestDTO struct {
	Action   string `json:"action"`
	Database string `json:"database"`
}

type AllowedAndDeniedIRNs struct {
	Allowed []*irn.IRN
	Denied  []*irn.IRN
}
