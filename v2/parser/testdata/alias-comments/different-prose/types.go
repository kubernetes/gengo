package p

// Original implements Tracker.
type Original struct{}

// Alias implements Tracker.
//
// Alias is preserved for backwards compatibility.
type Alias = Original

// Original implements Tracker.

type DetachedOriginal struct{}

// Alias implements Tracker.
//
// Alias is preserved for backwards compatibility.

type DetachedAlias = DetachedOriginal
