package iamcore

import "gitlab.kaaiot.net/core/lib/iamcore/irn.git"

type PrincipalIRNResponseDTO struct {
	Data *irn.IRN `json:"data"`
}

type ErrorResponseDTO struct {
	Message string `json:"message"`
}

type AuthorizedOnResourceTypeRequestDTO struct {
	Action       string `json:"action"`
	ResourceType string `json:"resourceType"`
}

type AuthorizedOnResourceTypeResponseDTO struct {
	Data []*irn.IRN `json:"data"`
}

type AuthorizedOnResourceListRequestDTO struct {
	Resources []*irn.IRN `json:"resources"`
	Action    string     `json:"action"`
}
