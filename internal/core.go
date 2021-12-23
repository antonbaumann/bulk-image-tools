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

// listAllFiles returns a list of all files in the given directory and its subdirectories
func listAllFiles(root string) ([]string, error) {
	fileList := make([]string, 0)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.Type().IsRegular() {
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return fmt.Errorf("failed to list files: %v", err)
			}
			fileList = append(fileList, relPath)
		}
		return err
	})
	return fileList, err
}

// createTasks creates transformation tasks from images in the given directory
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
		return nil, fmt.Errorf("failed to create tasks: %v", err)
	}

	tasks := make([]Task, 0, len(files))
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

// worker processes tasks in the `tasks` channel and reports the results
func worker(tasks <-chan Task, results chan<- Result) {
	for task := range tasks {
		err := task.Run()
		results <- Result{
			Success:  err == nil,
			FilePath: filepath.Join(task.RootFrom, task.RelPath),
			Error:    err,
		}
	}
}

// RunTransformations runs the transformations in parallel
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
		return fmt.Errorf("failed to run transformations: %v", err)
	}

	// Create channels to send tasks to workers and receive results from workers
	tasks := make(chan Task, len(taskList))
	results := make(chan Result, len(taskList))

	// start workers
	for i := 0; i < nrWorkers; i++ {
		go worker(tasks, results)
	}

	// send tasks
	for _, task := range taskList {
		tasks <- task
	}
	close(tasks)

	// receive results and track progress
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
