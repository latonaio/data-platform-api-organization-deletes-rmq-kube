package requests

type Organization struct {
	BusinessPartner     int    `json:"BusinessPartner"`
	Organization        string `json:"Organization"`
	IsMarkedForDeletion *bool  `json:"IsMarkedForDeletion"`
}
