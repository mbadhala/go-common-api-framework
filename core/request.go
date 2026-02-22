package core

import (
	"context"
	"time"
)

type contextKey string

const FileRecordKey contextKey = "fileRecord"

func GetFileRecord(ctx context.Context) (*FileRecord, bool) {
	v, ok := ctx.Value(FileRecordKey).(*FileRecord)
	return v, ok
}

type File struct {
	Filename string
	Size     int64
	Content  []byte
	Header   map[string][]string
}

type Request struct {
	Method      string
	Path        string
	Headers     map[string]string
	Body        []byte
	PathParams  map[string]string
	QueryParams map[string]string

	Files map[string]File   // <--- NEW
	Form  map[string]string // normal form fields
}

type FileRecord struct {
	ID           string
	Component    string
	OriginalName string
	StoragePath  string
	Checksum     string
	Size         int64
	UploadedBy   string
	UploadedAt   time.Time
}

type UploadProcessor interface {
	Component() string
	Process(ctx context.Context, file FileRecord) error
}

type UploadRegistry struct {
	processors map[string]UploadProcessor
}

func NewUploadRegistry() *UploadRegistry {
	return &UploadRegistry{
		processors: make(map[string]UploadProcessor),
	}
}

func (r *UploadRegistry) Register(p UploadProcessor) {
	r.processors[p.Component()] = p
}

func (r *UploadRegistry) Get(component string) (UploadProcessor, bool) {
	p, ok := r.processors[component]
	return p, ok
}
