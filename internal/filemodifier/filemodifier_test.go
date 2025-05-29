package filemodifier

import (
	"fmt"
	"strings"
	"testing"

	"github.com/DEMAxx/project_work/internal/lrucache"
	"github.com/DEMAxx/project_work/pkg/config"
	"github.com/DEMAxx/project_work/pkg/logger"
	"github.com/stretchr/testify/assert"
)

// Путь к директории с тестовыми изображениями.
const testImagesDir = "testdata"

func TestResizeImage(t *testing.T) {
	fileURL := "raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg" //nolint
	log := logger.MustSetupLogger("previewer", "Test", true, "info")
	cnf := config.Config{}
	cnf.UploadPath = testImagesDir
	cnf.Capability = 1
	cache := lrucache.NewCache(cnf.Capability, cnf.UploadPath, log)

	t.Run("success", func(t *testing.T) {
		path := fmt.Sprintf(
			"%d/%d/%s",
			100,
			100,
			fileURL,
		)

		modifier, err := New(
			strings.Split(path, "/"),
			&log,
			&cnf,
			cache,
		)

		assert.NoError(t, err)

		cachedImage, found := modifier.GetFromCache()

		assert.False(t, found)
		assert.Nil(t, cachedImage)

		resizedImage, err := modifier.ResizeImage()

		assert.NoError(t, err)
		assert.NotNil(t, resizedImage)

		cache.Clear()
	})

	t.Run("success different dimensions", func(t *testing.T) {
		path := fmt.Sprintf(
			"%d/%d/%s",
			200,
			200,
			fileURL,
		)

		modifier, err := New(
			strings.Split(path, "/"),
			&log,
			&cnf,
			cache,
		)

		assert.NoError(t, err)

		cachedImage, found := modifier.GetFromCache()

		assert.False(t, found)
		assert.Nil(t, cachedImage)

		resizedImage, err := modifier.ResizeImage()

		assert.NoError(t, err)
		assert.NotNil(t, resizedImage)

		cache.Clear()

		path = fmt.Sprintf(
			"%d/%d/%s",
			200,
			200,
			fileURL,
		)

		modifier, err = New(
			strings.Split(path, "/"),
			&log,
			&cnf,
			cache,
		)

		assert.NoError(t, err)

		cachedImage, found = modifier.GetFromCache()

		assert.False(t, found)
		assert.Nil(t, cachedImage)

		resizedImage, err = modifier.ResizeImage()

		assert.NoError(t, err)
		assert.NotNil(t, resizedImage)

		cache.Clear()
	})

	t.Run("success from cache", func(t *testing.T) {
		path := fmt.Sprintf(
			"%d/%d/%s",
			200,
			200,
			fileURL,
		)

		modifier, err := New(
			strings.Split(path, "/"),
			&log,
			&cnf,
			cache,
		)

		assert.NoError(t, err)

		cachedImage, found := modifier.GetFromCache()

		assert.False(t, found)
		assert.Nil(t, cachedImage)

		resizedImage, err := modifier.ResizeImage()

		assert.NoError(t, err)
		assert.NotNil(t, resizedImage)

		modifier, err = New(
			strings.Split(path, "/"),
			&log,
			&cnf,
			cache,
		)

		assert.NoError(t, err)

		cachedImage, found = modifier.GetFromCache()

		assert.True(t, found)
		assert.NotNil(t, cachedImage)
	})
}

func TestFailResizeImage(t *testing.T) {
	fileURL := "raw.githubusercontent.com/OtusGolang/final_project/master/examples/image-previewer/_gopher_original_1024x504.jpg" //nolint
	log := logger.MustSetupLogger("previewer", "Test", true, "info")
	cnf := config.Config{}
	cnf.UploadPath = testImagesDir
	cnf.Capability = 1
	cache := lrucache.NewCache(cnf.Capability, cnf.UploadPath, log)

	t.Run("zero dimensions", func(t *testing.T) {
		path := fmt.Sprintf(
			"%d/%d/%s",
			0,
			0,
			fileURL,
		)

		_, err := New(
			strings.Split(path, "/"),
			&log,
			&cnf,
			cache,
		)

		assert.Error(t, err)
	})

	t.Run("invalid path", func(t *testing.T) {
		path := fmt.Sprintf(
			"%d/%d/%s",
			100,
			100,
			"test",
		)

		_, err := New(
			strings.Split(path, "/"),
			&log,
			&cnf,
			cache,
		)

		assert.Error(t, err)
	})

	t.Run("negative dimensions", func(t *testing.T) {
		path := fmt.Sprintf(
			"%d/%d/%s",
			-100,
			-100,
			fileURL,
		)

		_, err := New(
			strings.Split(path, "/"),
			&log,
			&cnf,
			cache,
		)

		assert.Error(t, err)
	})
}
