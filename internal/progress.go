package internal

import (
	"fmt"
	"time"
)

type progress struct {
	Total     int
	Success   int
	Error     int
	StartTime time.Time
}

type Progress interface {
	String() string
	AddError()
	AddSuccess()
	Duration() time.Duration
	EstimatedTimeLeft() time.Duration
}

var _ Progress = (*progress)(nil)

func NewProgress(total int) Progress {
	progress := &progress{Total: total, StartTime: time.Now()}
	return Progress(progress)
}

func (p *progress) AddError() {
	p.Error += 1
}

func (p *progress) AddSuccess() {
	p.Success += 1
}

func (p progress) Duration() time.Duration {
	return time.Now().Sub(p.StartTime)
}

func (p progress) EstimatedTimeLeft() time.Duration {
	itemsProcessed := p.Success + p.Error
	timePerItemMs := float32(p.Duration().Milliseconds()) / float32(itemsProcessed)
	return time.Millisecond * time.Duration(float32(p.Total-itemsProcessed)*timePerItemMs)
}

func (p progress) String() string {
	timeLeft := p.EstimatedTimeLeft()
	return fmt.Sprintf("%v / %v - Skipped %v      ETA: %v",
		p.Success+p.Error,
		p.Total,
		p.Error,
		timeLeft.String(),
	)
}
