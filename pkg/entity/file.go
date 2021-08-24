package entity

import "time"

type File struct {
	ID          string
	Name        string
	Path        string
	User        int
	Size        int
	ContentType string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (e *File) IsEmpty() bool {
	return e == nil || *e == File{}
}
