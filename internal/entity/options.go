package entity

import entsql "entgo.io/ent/dialect/sql"

type Option struct {
	modifiers []func(s *entsql.Selector)
}

func (o *Option) Modify(modifiers ...func(s *entsql.Selector)) *Option {
	o.modifiers = append(o.modifiers, modifiers...)
	return o
}
