package p

// Original implements Tracker.
// +k8s:deepcopy-gen=true
type Original struct{}

// Alias is a compatibility name.
// +k8s:deepcopy-gen=true
type Alias = Original

// Original implements Tracker.
// +k8s:deepcopy-gen=true

type DetachedOriginal struct{}

// Alias is a compatibility name.
// +k8s:deepcopy-gen=true

type DetachedAlias = DetachedOriginal
