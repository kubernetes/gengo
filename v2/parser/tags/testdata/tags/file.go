package tags

type T1 struct {
	A string `json:"a"`
	B string `json:"b,omitempty"`
	C string `json:",inline"`
	D string `json:"-"`
	E string `json:""`
}
