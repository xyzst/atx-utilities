package util

type Output interface {
	Init()
	Write([]string)
}
