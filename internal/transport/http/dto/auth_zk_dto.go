package dto

type AuthZkResponse struct {
	UserID  string                 `json:"userId"`
	Message string                 `json:"message"`
	Proofs  map[string]interface{} `json:"proofs,omitempty"`
}
