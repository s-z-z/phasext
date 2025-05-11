package v1

import (
	"github.com/suzi1037/pcmd/app/demo/apis"
	"github.com/suzi1037/pcmd/pkg/builder"
)

const Version = "v1"

var (
	_builder      = builder.NewBuilder(apis.GroupName, Version)
	AddToScheme   = _builder.AddToScheme
	RegisterTypes = _builder.RegisterTypes
)

func init() {
	RegisterTypes(
		&Foo{},
	)
}
