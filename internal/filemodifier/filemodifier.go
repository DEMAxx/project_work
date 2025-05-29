package filemodifier

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/DEMAxx/project_work/internal/filesearch"
	"github.com/DEMAxx/project_work/internal/lrucache"
	"github.com/DEMAxx/project_work/pkg/config"
	"github.com/h2non/bimg"
	"github.com/rs/zerolog"
)

type Modifier interface {
	ResizeImage() ([]byte, error)
	GetFromCache() (interface{}, bool)
}

type fileModifier struct {
	height          int
	width           int
	imageURL        string
	UploadPath      string
	fetchedFilePath string
	cacheKey        lrucache.Key
	cache           lrucache.Cache
	logger          *zerolog.Logger
}

func (fileModifier *fileModifier) ResizeImage() ([]byte, error) {
	fetchedFilePath := fmt.Sprintf(
		"%s/%d_%d.jpg",
		fileModifier.UploadPath,
		fileModifier.width,
		fileModifier.height,
	)

	resp, err := filesearch.FetchFileFromURL(fileModifier.imageURL, fetchedFilePath, fileModifier.logger) //nolint
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch image")
	}

	image, err := bimg.Read(fileModifier.fetchedFilePath)
	if err != nil {
		return nil, err
	}

	resizedImage, err := bimg.NewImage(image).Process(bimg.Options{
		Width:  fileModifier.width,
		Height: fileModifier.height,
		Crop:   true,
		Type:   bimg.JPEG,
	})
	if err != nil {
		return nil, err
	}

	if err := fileModifier.cache.Set(fileModifier.cacheKey, resizedImage); err {
		return nil, errors.New("failed to store image in cache")
	}

	return resizedImage, nil
}

func (fileModifier *fileModifier) GetFromCache() (cachedImage interface{}, found bool) {
	cacheKey := lrucache.Key(
		fmt.Sprintf(
			"%d_%d_%s",
			fileModifier.width,
			fileModifier.height,
			fileModifier.imageURL,
		),
	)

	return fileModifier.cache.Get(cacheKey)
}

func New(parts []string, logger *zerolog.Logger, cnf *config.Config, cache lrucache.Cache) (Modifier, error) {
	if len(parts) < 3 {
		return nil, errors.New("not enough parts")
	}

	height, width, imageURL := parts[0], parts[1], strings.Join(parts[2:], "/")

	logger.Info().Msg(
		fmt.Sprintf(
			"Extracted vars - height: %s, width: %s, image url: %s", height, width, imageURL,
		),
	)

	if !strings.HasSuffix(imageURL, ".jpg") {
		return nil, errors.New("invalid image URL format. Only .jpg files are supported")
	}

	// Resize the image
	widthInt, err := strconv.Atoi(width)
	if err != nil {
		return nil, errors.New("invalid width value")
	}

	heightInt, err := strconv.Atoi(height)
	if err != nil {
		return nil, errors.New("invalid height value")
	}

	if widthInt <= 0 || heightInt <= 0 {
		return nil, errors.New("width or height must be positive")
	}

	fetchedFilePath := fmt.Sprintf("%s/%s_%s.jpg", cnf.UploadPath, width, height)

	cacheKey := lrucache.Key(
		fmt.Sprintf(
			"%d_%d_%s",
			widthInt,
			heightInt,
			imageURL,
		),
	)

	return &fileModifier{
		height:          heightInt,
		width:           widthInt,
		imageURL:        imageURL,
		UploadPath:      cnf.UploadPath,
		fetchedFilePath: fetchedFilePath,
		cacheKey:        cacheKey,
		cache:           cache,
		logger:          logger,
	}, nil
}
