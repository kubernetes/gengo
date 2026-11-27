package p

// Original implements Tracker.
type Original struct{}

// +custom:generate=true
type Alias = Original

// Original implements Tracker.

type DetachedOriginal struct{}

// +custom:generate=true

type DetachedAlias = DetachedOriginal
