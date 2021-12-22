package internal

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

type Task struct {
	RootFrom string
	RootTo   string
	RelPath  string
	Format   ImageFormat
	Width    int
	Height   int
	Rotation ImageRotation
}

type Result struct {
	Success  bool
	FilePath string
	Error    error
}

func createPathIfNotExist(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(dir, 0755)
		} else {
			return fmt.Errorf("createPathIfNotExist: %v", err)
		}
	}
	return nil
}

func trimExtension(path string) string {
	suffix := filepath.Ext(path)
	return strings.TrimSuffix(path, suffix)
}

func resize(img image.Image, width, height int) image.Image {
	if width > 0 || height > 0 {
		return imaging.Resize(img, width, height, imaging.CatmullRom)
	}
	return img
}

func rotate(img image.Image, rotation ImageRotation) image.Image {
	switch rotation {
	case Rotate0:
		return img
	case Rotate90:
		return imaging.Rotate90(img)
	case Rotate180:
		return imaging.Rotate180(img)
	case Rotate270:
		return imaging.Rotate270(img)
	default:
		return img
	}
}

func save(img image.Image, rootTo, relPath string, format ImageFormat) error {
	imagePathNoExt := trimExtension(filepath.Join(rootTo, relPath))
	if err := createPathIfNotExist(imagePathNoExt); err != nil {
		return err
	}

	switch format {
	case JPEG:
		err := imaging.Save(
			img,
			fmt.Sprintf("%v.jpeg", imagePathNoExt),
			imaging.JPEGQuality(90),
		)
		if err != nil {
			return err
		}
	case PNG:
		err := imaging.Save(
			img,
			fmt.Sprintf("%v.png", imagePathNoExt),
			imaging.PNGCompressionLevel(png.DefaultCompression),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (t Task) Run() error {
	absPath := filepath.Join(t.RootFrom, t.RelPath)
	img, err := imaging.Open(absPath, imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("opening image [%v]: %v", absPath, err)
	}

	result := resize(img, t.Width, t.Height)
	result = rotate(result, t.Rotation)
	return save(result, t.RootTo, t.RelPath, t.Format)
}
