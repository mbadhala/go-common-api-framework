package core

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type Handler func(ctx context.Context, req Request) Response

func Wrap[I any, O any](fn func(context.Context, Request, I) (O, error)) Handler {
	return func(ctx context.Context, req Request) Response {
		var input I
		if err := json.Unmarshal(req.Body, &input); err != nil {
			return Error(400, "invalid json")
		}

		out, err := fn(ctx, req, input)
		if err != nil {
			return Error(500, err.Error())
		}

		return JSON(200, out)
	}
}

func WrapNoBody[O any](fn func(context.Context, Request) (O, error)) Handler {
	return func(ctx context.Context, req Request) Response {
		out, err := fn(ctx, req)
		if err != nil {
			return Error(500, err.Error())
		}
		return JSON(200, out)
	}
}

func WrapUpload1[O any](h func(ctx context.Context, req Request) (O, error)) Handler {
	return func(ctx context.Context, req Request) Response {

		if len(req.Files) == 0 {
			return Error(400, "no file uploaded")
		}

		out, err := h(ctx, req)
		if err != nil {
			return Error(500, err.Error())
		}

		return JSON(200, out)
	}
}

func WrapUpload(
	registry *UploadRegistry,
	basePath string,
) Handler {

	return func(ctx context.Context, req Request) Response {

		file, ok := req.Files["file"]
		if !ok {
			return Error(400, "file required")
		}

		component := req.Form["component"]
		if component == "" {
			return Error(400, "component required")
		}

		processor, ok := registry.Get(component)
		if !ok {
			return Error(400, "invalid component")
		}

		// checksum
		hash := sha256.Sum256(file.Content)
		checksum := hex.EncodeToString(hash[:])

		now := time.Now()
		ext := filepath.Ext(file.Filename)

		finalPath := filepath.Join(
			basePath,
			component,
			now.Format("2006"),
			now.Format("01"),
			now.Format("02"),
			checksum+ext,
		)

		if err := os.MkdirAll(filepath.Dir(finalPath), 0755); err != nil {
			return Error(500, err.Error())
		}

		if err := os.WriteFile(finalPath, file.Content, 0644); err != nil {
			return Error(500, err.Error())
		}

		record := &FileRecord{
			ID:           uuid.New().String(),
			Component:    component,
			OriginalName: file.Filename,
			StoragePath:  finalPath,
			Checksum:     checksum,
			Size:         file.Size,
			UploadedAt:   now,
		}

		// Inject into context
		ctx = context.WithValue(ctx, FileRecordKey, record)

		// Call processor
		if err := processor.Process(ctx, *record); err != nil {
			return Error(500, err.Error())
		}

		return JSON(200, record)
	}
}
