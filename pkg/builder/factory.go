package builder

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type Builder struct {
	schemeBuilder runtime.SchemeBuilder
	types         []runtime.Object
}

func NewBuilder(group, version string) *Builder {
	_builder := &Builder{
		schemeBuilder: runtime.NewSchemeBuilder(),
		types:         make([]runtime.Object, 0),
	}

	_builder.schemeBuilder.Register(func(s *runtime.Scheme) error {
		s.AddKnownTypes(
			schema.GroupVersion{Group: group, Version: version},
			_builder.types...,
		)
		return nil
	})

	return _builder
}

func (b *Builder) AddToScheme(s *runtime.Scheme) error {
	return b.schemeBuilder.AddToScheme(s)
}

func (b *Builder) RegisterTypes(o ...runtime.Object) {
	b.types = append(b.types, o...)
}
