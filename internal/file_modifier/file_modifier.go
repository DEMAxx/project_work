package file_modifier

import "github.com/h2non/bimg"

// ResizeImage resizes an image to the specified width and height.
func ResizeImage(inputPath string, outputPath string, width int, height int) error {
	// Read the image from the input path
	image, err := bimg.Read(inputPath)
	if err != nil {
		return err
	}

	// Resize the image
	resizedImage, err := bimg.NewImage(image).Resize(width, height)
	if err != nil {
		return err
	}

	// Write the resized image to the output path
	err = bimg.Write(outputPath, resizedImage)
	if err != nil {
		return err
	}

	return nil
}
