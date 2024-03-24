package service

import (
	"bytes"
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"time"

	"github.com/IktaS/sign/auth"
	"github.com/google/uuid"
	"github.com/phpdave11/gofpdi"
	"github.com/signintech/gopdf"
	"github.com/skip2/go-qrcode"
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
		log.Println(err, png)
		return nil, err
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	img, err := gopdf.ImageHolderByBytes(png)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	importer := gofpdi.NewImporter()
	readSeeker := io.ReadSeeker(req.File)
	importer.SetSourceStream(&readSeeker)
	pageSizes := importer.GetPageSizes()

	for i := 1; i <= len(pageSizes); i++ {
		tmplid := pdf.ImportPageStream(&readSeeker, i, "/MediaBox")

		width := pageSizes[i]["/MediaBox"]["w"]
		height := pageSizes[i]["/MediaBox"]["h"]

		pdf.AddPageWithOption(gopdf.PageOption{PageSize: &gopdf.Rect{W: width, H: height}})
		pdf.UseImportedTemplate(tmplid, 0, 0, width, height)

		log.Println(width*(float64(req.LocationXPercentage)/100), height-(height*(float64(req.LocationYPercentage)/100))+float64(req.QRSize))

		pdf.ImageByHolder(img,
			width*(float64(req.LocationXPercentage)/100),
			height-(height*(float64(req.LocationYPercentage)/100))+float64(req.QRSize),
			nil,
		)
		// // Add the page.
		// err = c.AddPage(page)
		// if err != nil {
		// 	log.Println(err)
		// 	return nil, err
		// }

		// if req.IsAllPage != nil && *req.IsAllPage {
		// 	err := c.Draw(img)
		// 	if err != nil {
		// 		log.Println(err)
		// 		return nil, err
		// 	}
		// }

		// if req.IsAllPage == nil && req.QRPage != nil && i+1 == *req.QRPage {
		// 	err := c.Draw(img)
		// 	if err != nil {
		// 		log.Println(err)
		// 		return nil, err
		// 	}
		// }
	}

	var b bytes.Buffer
	_, err = pdf.WriteTo(&b)
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
