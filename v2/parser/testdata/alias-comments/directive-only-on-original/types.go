package p

// +custom:generate=true
type Original struct{}

// Alias is a compatibility name.
type Alias = Original

// +custom:generate=true

type DetachedOriginal struct{}

// Alias is a compatibility name.

type DetachedAlias = DetachedOriginal
