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
	Name         string   `json:"name"`
	Application  string   `json:"application"`
	Path         string   `json:"path"`
	ResourceType string   `json:"resourceType"`
	Enabled      bool     `json:"enabled"`
	TenantID     string   `json:"tenantID"`
	PoolIDs      []string `json:"poolIDs,omitempty"`
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

type AttachPolicyRequestDTO struct {
	PolicyIDs []string `json:"policyIDs"`
}

type AllowedAndDeniedIRNs struct {
	Allowed []*irn.IRN
	Denied  []*irn.IRN
}

type EvaluateDebugResourcesRequestDTO struct {
	Application string       `json:"application"`
	Actions     []string     `json:"actions"`
	Resources   []*irn.IRN64 `json:"resources"`
}

type DebugEvaluationActionDetail struct {
	Action        string   `json:"action"`
	Decision      string   `json:"decision"`
	AllowPolicies []string `json:"allowPolicies,omitempty"`
	DenyPolicies  []string `json:"denyPolicies,omitempty"`
}

type DebugEvaluationResourceItem struct {
	ID       string                         `json:"id"`
	IRN      *irn.IRN                       `json:"irn"`
	Decision string                         `json:"decision"`
	Actions  []*DebugEvaluationActionDetail `json:"actions"`
}

type EvaluateDebugResourcesResponseDTO struct {
	Data                 []*DebugEvaluationResourceItem `json:"data"`
	EvaluationTimeMillis int                            `json:"evaluationTimeMillis"`
}

type PoolResponseDTO struct {
	ID          string   `json:"id"`
	IRN         *irn.IRN `json:"irn"`
	Name        string   `json:"name"`
	ResourceIDs []string `json:"resourceIDs"`
}

type PoolsResponseDTO struct {
	Data     []*PoolResponseDTO `json:"data"`
	Count    int                `json:"count"`
	Page     int                `json:"page"`
	PageSize int                `json:"pageSize"`
}
