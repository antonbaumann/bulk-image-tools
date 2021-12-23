package internal

import "github.com/h2non/bimg"

type ImageFormat string

const (
	JPEG ImageFormat = "jpeg"
	WEBP ImageFormat = "webp"
	PNG  ImageFormat = "png"
	TIFF ImageFormat = "tiff"
	GIF  ImageFormat = "gif"
)

func (f ImageFormat) IsValid() bool {
	return f == JPEG || f == WEBP || f == PNG || f == TIFF || f == GIF
}

func (f ImageFormat) ToBimgImageType() bimg.ImageType {
	switch f {
	case JPEG:
		return bimg.JPEG
	case WEBP:
		return bimg.WEBP
	case PNG:
		return bimg.PNG
	case TIFF:
		return bimg.TIFF
	case GIF:
		return bimg.GIF
	default:
		return bimg.JPEG
	}
}
