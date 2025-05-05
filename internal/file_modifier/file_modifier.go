package file_modifier

import "github.com/h2non/bimg"

// ResizeImage resizes an image to the specified width and height.
func ResizeImage(inputPath string, width int, height int) ([]byte, error) {
	image, err := bimg.Read(inputPath)
	if err != nil {
		return nil, err
	}

	resizedImage, err := bimg.NewImage(image).Process(bimg.Options{
		Width:  width,
		Height: height,
		Crop:   true,
		Type:   bimg.JPEG,
	})

	if err != nil {
		return nil, err
	}

	return resizedImage, nil
}
