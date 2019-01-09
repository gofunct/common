package runtime

import "encoding"

type Encoder struct{}

func (Encoder) Marshal() encoding.BinaryMarshaler {
	panic("implement me")
}

func (Encoder) UnMarshal() encoding.BinaryUnmarshaler {
	panic("implement me")
}
