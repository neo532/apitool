package entity

const (
	EmptyPb      = "google.protobuf.Empty"
	EmptyVarName = "Empty"
	EmptyType    = "emptypb.Empty"

	AnyPb      = "google.protobuf.Any"
	AnyVarName = "Any"
	AnyType    = "anypb.Any"

	Wrapper = "Wrapper"
)

var (
	SpecialMap = map[string]string{
		EmptyPb: EmptyType,
		AnyPb:   AnyType,
	}
)

type MethodType uint8

const (
	UnaryType          MethodType = 1
	TwoWayStreamsType  MethodType = 2
	RequestStreamsType MethodType = 3
	ReturnsStreamsType MethodType = 4
)

func GetMethodType(streamsRequest, streamsReturns bool) MethodType {
	if !streamsRequest && !streamsReturns {
		return UnaryType
	} else if streamsRequest && streamsReturns {
		return TwoWayStreamsType
	} else if streamsRequest {
		return RequestStreamsType
	} else if streamsReturns {
		return ReturnsStreamsType
	}
	return UnaryType
}
