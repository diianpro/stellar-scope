package response

import "github.com/diianpro/stellar-scope/internal/domain/apod"

type GetAll struct {
	Pictures []apod.Data `json:"pictures"`
}
