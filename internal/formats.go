package internal

import "github.com/h2non/bimg"

type ImageFormat string

const (
	JPEG ImageFormat = "jpeg"
	PNG  ImageFormat = "png"
)

func (f ImageFormat) ToBimgImageType() bimg.ImageType {
	switch f {
	case JPEG:
		return bimg.JPEG
	case PNG:
		return bimg.PNG
	default:
		return bimg.JPEG
	}
}
