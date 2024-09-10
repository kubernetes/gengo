package simple

type ExportedString string

type ExportedInt int

type ExportedStruct struct {
	I  int
	PI *int
	S  string
	PS *string
	AI AliasInt
	AS AliasStruct
}

type privateString string

type privateInt int

type privateStruct struct {
	i  int
	pi *int
	s  string
	ps *string
}

type AliasInt = int

type UnderlyingStruct struct {
}
type AliasStruct = UnderlyingStruct

type AliasUnused = int
