package simple

type ExportedString string

type ExportedInt int

type ExportedStruct struct {
	I  int
	PI *int
	S  string
	PS *string
}

type privateString string

type privateInt int

type privateStruct struct {
	i  int
	pi *int
	s  string
	ps *string
}
