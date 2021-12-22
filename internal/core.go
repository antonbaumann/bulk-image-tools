package internal

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

type ImageFormat string

const (
	JPEG ImageFormat = "jpeg"
	PNG  ImageFormat = "png"
)

type ImageRotation int

const (
	Rotate0   ImageRotation = 0
	Rotate90  ImageRotation = 90
	Rotate180 ImageRotation = 180
	Rotate270 ImageRotation = 270
)

func listAllFiles(root string) ([]string, error) {
	fileList := make([]string, 0)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.Type().IsRegular() {
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			fileList = append(fileList, relPath)
		}
		return err
	})
	return fileList, err
}

func createTasks(
	pathFrom string,
	pathTo string,
	outputFormat ImageFormat,
	width int,
	height int,
	rotation ImageRotation,
) ([]Task, error) {
	files, err := listAllFiles(pathFrom)
	if err != nil {
		return nil, err
	}

	tasks := make([]Task, 0)
	for _, file := range files {
		tasks = append(tasks, Task{
			RootFrom: pathFrom,
			RootTo:   pathTo,
			RelPath:  file,
			Format:   outputFormat,
			Width:    width,
			Height:   height,
			Rotation: rotation,
		})
	}
	return tasks, nil
}

func worker(tasks <-chan Task, results chan<- Result) {
	for task := range tasks {
		if err := task.Run(); err != nil {
			results <- Result{
				Success:  false,
				FilePath: filepath.Join(task.RootFrom, task.RelPath),
				Error:    err,
			}
		} else {
			results <- Result{
				Success:  true,
				FilePath: filepath.Join(task.RootFrom, task.RelPath),
				Error:    nil,
			}
		}
	}
}

func RunTransformations(
	pathFrom string,
	pathTo string,
	outputFormat ImageFormat,
	width int,
	height int,
	rotation ImageRotation,
	nrWorkers int,
) error {
	taskList, err := createTasks(
		pathFrom,
		pathTo,
		outputFormat,
		width,
		height,
		rotation,
	)
	if err != nil {
		return err
	}

	tasks := make(chan Task, len(taskList))
	results := make(chan Result, len(taskList))

	for i := 0; i < nrWorkers; i++ {
		go worker(tasks, results)
	}

	for _, task := range taskList {
		tasks <- task
	}
	close(tasks)

	progress := NewProgress(len(taskList))
	for i := 0; i < len(taskList); i++ {
		result := <-results
		if result.Error != nil {
			fmt.Println(result.Error)
			progress.AddError()
		} else {
			progress.AddSuccess()
		}
		fmt.Printf("\r%v  ", progress.String())
	}

	return nil
}
