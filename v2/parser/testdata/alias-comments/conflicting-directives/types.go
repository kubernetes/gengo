package p

// +k8s:deepcopy-gen=true
type Original struct{}

// +k8s:deepcopy-gen=false
type Alias = Original

// +k8s:deepcopy-gen=true

type DetachedOriginal struct{}

// +k8s:deepcopy-gen=false

type DetachedAlias = DetachedOriginal
