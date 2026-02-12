package models

import (
	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type CreateShortUrlRequest struct {
	OriginalUrl    string `json:"original_url" validate:"required,url"`
	CustomShortUrl string `json:"custom_short_url,omitempty" validate:"omitempty,alphanum,min=4,max=8"`
}

type CreateShortUrlResponse struct {
	ShortUrl string `json:"short_url"`
}

func (r *CreateShortUrlRequest) Validate() error {
	return v.ValidateStruct(r,
		v.Field(&r.OriginalUrl, v.Required, is.URL),
		v.Field(&r.CustomShortUrl, v.NilOrNotEmpty, v.Length(4, 8)),
	)
}
