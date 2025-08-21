package httpclient

type Seriazable[T any] interface {
	Codec() Codec[T]
}

type Codec[T any] interface {
	Encode(input T) ([]byte, error)
	Decode(data []byte) (T, error)
}
