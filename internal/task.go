package internal

import (
	"fmt"
	"github.com/h2non/bimg"
	"os"
	"path/filepath"
	"strings"
)

// Task describes a image transformation
type Task struct {
	RootFrom string        // Root directory to read from
	RootTo   string        // Root directory to write to
	RelPath  string        // Relative path to the image
	Format   ImageFormat   // Format to convert to
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
			return fmt.Errorf("createPathIfNotExist: %v", err)
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
func resize(buffer []byte, width, height int) ([]byte, error) {
	if width > 0 || height > 0 {
		return bimg.NewImage(buffer).Resize(width, height)
	}
	return buffer, nil
}

// rotate the given image by the given rotation
func rotate(buffer []byte, rotation ImageRotation) ([]byte, error) {
	return bimg.NewImage(buffer).Rotate(rotation.ToBimgAngle())
}

// save the given image to the given path
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

// Run executes the given task
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
