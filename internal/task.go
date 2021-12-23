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

// Task describes a image transformation
type Task struct {
	RootFrom string        // Root directory to read from
	RootTo   string        // Root directory to write to
	RelPath  string        // Relative path to the image
	Format   ImageFormat   // Format the image should be converted to
	Width    int           // Width of the resulting image
	Height   int           // Height of the resulting image
	Rotation ImageRotation // Rotation of the resulting image
}

// Result describes the result of a task
type Result struct {
	Success  bool   // Whether the task was successful
	FilePath string // Path to the resulting image
	Error    error  // Error message if the task failed
}

// createPathIfNotExist creates all necessary directories for the given path
func createPathIfNotExist(path string) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(dir, 0755)
		} else {
			return fmt.Errorf("failed to create directory %s: %s", dir, err)
		}
	}
	return nil
}

// trimExtension removes the extension from the given path
func trimExtension(path string) string {
	suffix := filepath.Ext(path)
	return strings.TrimSuffix(path, suffix)
}

// resize the given image to the given width and height
func resize(img image.Image, width, height int) image.Image {
	if width > 0 || height > 0 {
		return imaging.Resize(img, width, height, imaging.CatmullRom)
	}
	return img
}

// rotate the given image by the given rotation
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

// save the given image to the given path
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

// Run executes the given task
func (t Task) Run() error {
	absPath := filepath.Join(t.RootFrom, t.RelPath)
	img, err := imaging.Open(absPath, imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("failed to run task: opening image [%v]: %v", absPath, err)
	}

	result := resize(img, t.Width, t.Height)
	result = rotate(result, t.Rotation)
	return save(result, t.RootTo, t.RelPath, t.Format)
}
