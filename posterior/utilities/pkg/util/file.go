package util

type File interface {
	New(name string) interface{}
}
