package internal

import (
	"fmt"
	"github.com/h2non/bimg"
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

func resize(buffer []byte, width, height int) ([]byte, error) {
	if width > 0 || height > 0 {
		return bimg.NewImage(buffer).Resize(width, height)
	}
	return buffer, nil
}

func rotate(buffer []byte, rotation ImageRotation) ([]byte, error) {
	return bimg.NewImage(buffer).Rotate(rotation.ToBimgAngle())
}

func save(buffer []byte, rootTo, relPath string, format ImageFormat) error {
	imagePathNoExt := trimExtension(filepath.Join(rootTo, relPath))
	if err := createPathIfNotExist(imagePathNoExt); err != nil {
		return err
	}

	convertedBuffer, err := bimg.NewImage(buffer).Convert(format.ToBimgImageType())
	if err != nil {
		return err
	}

	imagePath := fmt.Sprintf("%v.%v", imagePathNoExt, bimg.NewImage(convertedBuffer).Type())
	return bimg.Write(imagePath, convertedBuffer)
}

func (t Task) Run() error {
	absPath := filepath.Join(t.RootFrom, t.RelPath)
	buffer, err := bimg.Read(absPath)
	if err != nil {
		return fmt.Errorf("opening image [%v]: %v", absPath, err)
	}

	buffer, err = resize(buffer, t.Width, t.Height)
	if err != nil {
		return err
	}
	buffer, err = rotate(buffer, t.Rotation)
	if err != nil {
		return err
	}
	return save(buffer, t.RootTo, t.RelPath, t.Format)
}
