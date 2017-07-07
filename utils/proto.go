package utils

func ProtoString(s string) *string {
	return &s
}

func ProtoFloat64(f float64) *float64 {
	return &f
}

func ProtoInt64(i int64) *int64 {
	return &i
}

func ProtoBool(i bool) *bool {
	return &i
}