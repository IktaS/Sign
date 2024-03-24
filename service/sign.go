package service

import (
	"bytes"
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"github.com/IktaS/sign/auth"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
	"github.com/unidoc/unipdf/v3/creator"
	"github.com/unidoc/unipdf/v3/model"
)

type SignRequest struct {
	Username            string
	Password            string
	LocationXPercentage int
	LocationYPercentage int
	QRSize              int
	QRPage              *int
	IsAllPage           *bool
	File                multipart.File
	Filename            string
}

type SignatureInfo struct {
	Filename  string
	CreatedAt time.Time
	Fullname  string
}

type SignRepo interface {
	CreateUser(ctx context.Context, username, passwordHash, fullName string) (int, error)
	ValidateUser(ctx context.Context, username string, password string) (int, bool, error)
	GetSignatureInfo(ctx context.Context, key string) (SignatureInfo, error)
	SaveSignature(ctx context.Context, uuid string, filename string, fileHash string, createdAt time.Time, createdBy int) error
}

func (s *SignService) SignFile(ctx context.Context, req SignRequest) ([]byte, error) {
	userId, isUserVerified, err := s.signRepo.ValidateUser(ctx, req.Username, req.Password)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if !isUserVerified {
		log.Println(err)
		return nil, errors.New("user not verified")
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	uuidStr := uuid.String()

	qrString := fmt.Sprintf("%s/%s", s.verifyPath, uuidStr)

	var png []byte
	png, err = qrcode.Encode(qrString, qrcode.Highest, req.QRSize)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	c := creator.New()

	img, err := c.NewImageFromData(png)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	pdfReader, err := model.NewPdfReader(req.File)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	for i := 0; i < numPages; i++ {
		page, err := pdfReader.GetPage(i + 1)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		// Add the page.
		err = c.AddPage(page)
		if err != nil {
			log.Println(err)
			return nil, err
		}

		if req.IsAllPage != nil && *req.IsAllPage {
			err := c.Draw(img)
			if err != nil {
				log.Println(err)
				return nil, err
			}
		}

		if req.IsAllPage == nil && req.QRPage != nil && i+1 == *req.QRPage {
			err := c.Draw(img)
			if err != nil {
				log.Println(err)
				return nil, err
			}
		}
	}

	var b bytes.Buffer
	err = c.Write(&b)
	if err != nil {
		return nil, err
	}
	updatedPdfBytes := b.Bytes()

	h := sha512.New()
	h.Write(updatedPdfBytes)

	err = s.signRepo.SaveSignature(ctx, uuidStr, req.Filename, string(h.Sum(nil)), time.Now(), userId)
	if err != nil {
		return nil, err
	}

	return updatedPdfBytes, nil
}

func (s *SignService) CreateUser(ctx context.Context, username, password, fullName string) (int, error) {
	passwordHash, err := auth.GenerateEncodedHash([]byte(password), &auth.DefaultParams)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return s.signRepo.CreateUser(ctx, username, passwordHash, fullName)
}
