package service

type SignService struct {
	signRepo   SignRepo
	verifyPath string
}

func NewSignService(signRepo SignRepo, verifyPath string) *SignService {
	return &SignService{
		signRepo:   signRepo,
		verifyPath: verifyPath,
	}
}
