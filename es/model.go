package es

const (
	DocType = "_doc"
)

type Model interface {
	GetID() string
}
