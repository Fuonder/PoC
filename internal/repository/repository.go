package repository

type MyService interface {
	Run() error
	Stop() error
}
