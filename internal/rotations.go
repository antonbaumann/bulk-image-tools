package internal

import "github.com/h2non/bimg"

type ImageRotation int

const (
	Rotate0   ImageRotation = 0
	Rotate90  ImageRotation = 90
	Rotate180 ImageRotation = 180
	Rotate270 ImageRotation = 270
)

func (r ImageRotation) ToBimgAngle() bimg.Angle {
	return bimg.Angle(r)
}
