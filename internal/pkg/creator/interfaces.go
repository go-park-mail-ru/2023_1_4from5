package creator

//go:generate mockgen -source=interfaces.go -destination=./mocks/creator_mock.go -package=mock

type CreatorUsecase interface {
	GetPage()
}

type CreatorRepo interface {
	GetPage()
}
