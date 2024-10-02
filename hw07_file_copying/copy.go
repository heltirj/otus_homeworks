package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile          = errors.New("unsupported file")
	ErrOffsetExceedsFileSize    = errors.New("offset exceeds file size")
	ErrUndefinedSourceFile      = errors.New("undefined source file")
	ErrUndefinedDestinationFile = errors.New("undefined destination file")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" {
		return ErrUndefinedSourceFile
	}

	if toPath == "" {
		return ErrUndefinedDestinationFile
	}

	src, err := os.OpenFile(fromPath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer src.Close()

	stats, err := src.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}

	srcSize := stats.Size()
	if srcSize < offset {
		return ErrOffsetExceedsFileSize
	}

	dest, err := os.Create(toPath)
	if err != nil {
		return err
	}

	defer dest.Close()

	if offset > 0 {
		src.Seek(offset, io.SeekStart)
	}

	if limit > (srcSize-offset) || limit == 0 {
		limit = srcSize - offset
	}

	bar := pb.Full.Start64(limit)
	_, err = io.CopyN(dest, bar.NewProxyReader(src), limit)

	bar.Finish()

	return err
}
