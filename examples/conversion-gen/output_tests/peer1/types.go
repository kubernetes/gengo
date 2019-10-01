package peer1

type Foo struct {
	Str       string
	Int64Ptr  *int64
	SubStruct SubFooStruct
}

type SubFooStruct struct {
	BoolSlice []bool
}
