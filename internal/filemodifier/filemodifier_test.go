package filemodifier

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Путь к директории с тестовыми изображениями.
const testImagesDir = "testdata"

func TestResizeImage_Success(t *testing.T) {
	fmt.Println(testImagesDir)
	inputPath := fmt.Sprintf("%s/valid_image.jpg", testImagesDir)

	// Выполнение
	resizedImage, err := ResizeImage(inputPath, 100, 100)

	// Проверка
	assert.NoError(t, err)
	assert.NotEmpty(t, resizedImage)
}

func TestResizeImage_InvalidPath(t *testing.T) {
	// Выполнение
	resizedImage, err := ResizeImage("non_existent_file.jpg", 100, 100)

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, resizedImage)
}

func TestResizeImage_ZeroDimensions(t *testing.T) {
	inputPath := fmt.Sprintf("%s/valid_image.jpg", testImagesDir)

	resizedImage, err := ResizeImage(inputPath, 0, 0)

	assert.Error(t, err)
	assert.Nil(t, resizedImage)
}

func TestResizeImage_NegativeDimensions(t *testing.T) {
	inputPath := fmt.Sprintf("%s/valid_image.jpg", testImagesDir)

	resizedImage, err := ResizeImage(inputPath, -100, -100)

	assert.Error(t, err)
	assert.Nil(t, resizedImage)
}

func TestResizeImage_DifferentDimensions(t *testing.T) {
	// Подготовка
	inputPath := fmt.Sprintf("%s/valid_image.jpg", testImagesDir)

	// Выполнение
	resizedImage100x100, err := ResizeImage(inputPath, 100, 100)
	assert.NoError(t, err)

	resizedImage200x200, err := ResizeImage(inputPath, 200, 200)
	assert.NoError(t, err)

	// Проверка
	assert.NotEqual(t, bytes.Equal(resizedImage100x100, resizedImage200x200), true)
}

func TestResizeImage_InvalidImageFormat(t *testing.T) {
	// Подготовка
	inputPath := fmt.Sprintf("%s/invalid_format.txt", testImagesDir)

	// Выполнение
	resizedImage, err := ResizeImage(inputPath, 100, 100)

	// Проверка
	assert.Error(t, err)
	assert.Nil(t, resizedImage)
}
