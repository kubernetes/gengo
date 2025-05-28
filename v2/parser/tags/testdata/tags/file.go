package tags

type T1 struct {
	A string `json:"a"`
	B string `json:"b,omitempty"`
	C string `json:",inline"`
	D string `json:"-"`
	E string `json:""`

	T2
	*T3
}

type T2 struct {
	Z string `json:"z"`
}

type T3 struct {
	Y string `json:"y"`
}
