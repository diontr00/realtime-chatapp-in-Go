package controller

import (
	"github.com/go-playground/validator"
	"realtime-chat/translator"
)

type maincontroller struct {
	validator  *validator.Validate
	translator *translator.UTtrans
}

func NewMainController(v *validator.Validate, t *translator.UTtrans) *maincontroller {

	return &maincontroller{
		validator:  v,
		translator: t,
	}

}
