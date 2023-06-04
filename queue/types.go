package queue

type Queue interface {
	In(val any) error
	Out() (val any, err error)
}
