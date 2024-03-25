package service

import (
	"bytes"
	"context"
	"crypto/sha512"
	"crypto/subtle"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
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
	GetUserHash(ctx context.Context, username string) (int, string, error)
	GetSignatureInfo(ctx context.Context, id string) (SignatureInfo, error)
	GetSignatureFileHash(ctx context.Context, id string) ([]byte, error)
	SaveSignature(ctx context.Context, uuid string, filename string, fileHash []byte, createdAt time.Time, createdBy int) error
}

func (s *SignService) SignFile(ctx context.Context, req SignRequest) ([]byte, string, error) {
	userId, passwordHash, err := s.signRepo.GetUserHash(ctx, req.Username)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}
	isUserVerified, err := auth.CompareDataAndHash([]byte(req.Password), passwordHash)
	if err != nil {
		log.Println(err)
		return nil, "", err
	}
	if !isUserVerified {
		log.Println(err)
		return nil, "", errors.New("user not verified")
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		log.Println(err)
		return nil, "", err
	}
	uuidStr := uuid.String()

	qrString := fmt.Sprintf("%s?id=%s", s.verifyPath, uuidStr)

	var png []byte
	png, err = qrcode.Encode(qrString, qrcode.Highest, req.QRSize)
	if err != nil {
		log.Println(err, png)
		return nil, "", err
	}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	img, err := gopdf.ImageHolderByBytes(png)
	if err != nil {
		log.Println(err)
		return nil, "", err
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

		xLocation := ((width - (float64(req.QRSize) / 1.75)) * (float64(req.LocationXPercentage) / 100))
		ylocation := ((height - (float64(req.QRSize) / 1.75)) * (float64(100-req.LocationYPercentage) / 100))

		if req.IsAllPage != nil && *req.IsAllPage {
			pdf.ImageByHolder(
				img,
				xLocation,
				ylocation,
				nil,
			)
		}

		if req.IsAllPage == nil && req.QRPage != nil && i+1 == *req.QRPage {
			pdf.ImageByHolder(
				img,
				xLocation,
				ylocation,
				nil,
			)
		}
	}

	var b bytes.Buffer
	_, err = pdf.WriteTo(&b)
	if err != nil {
		return nil, "", err
	}
	updatedPdfBytes := b.Bytes()

	h := sha512.New()
	_, err = h.Write(updatedPdfBytes)
	if err != nil {
		return nil, "", err
	}

	filename := fmt.Sprintf("%s_signed.pdf", strings.TrimSuffix(req.Filename, filepath.Ext(req.Filename)))

	err = s.signRepo.SaveSignature(ctx, uuidStr, filename, h.Sum(nil), time.Now(), userId)
	if err != nil {
		return nil, "", err
	}

	return updatedPdfBytes, filename, nil
}

func (s *SignService) CreateUser(ctx context.Context, username, password, fullName string) (int, error) {
	passwordHash, err := auth.GenerateEncodedHash([]byte(password), &auth.DefaultParams)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return s.signRepo.CreateUser(ctx, username, passwordHash, fullName)
}

func (s *SignService) GetSignatureInfo(ctx context.Context, id string) (SignatureInfo, error) {
	return s.signRepo.GetSignatureInfo(ctx, id)
}

func (s *SignService) VerifyFileHash(ctx context.Context, id string, fileData []byte) (bool, error) {
	originalHash, err := s.signRepo.GetSignatureFileHash(ctx, id)
	if err != nil {
		return false, err
	}

	h := sha512.New()
	_, err = h.Write(fileData)
	if err != nil {
		return false, err
	}

	// Check that the contents of the hashed data are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(originalHash, h.Sum(nil)) == 1 {
		return true, nil
	}
	return false, nil
}
