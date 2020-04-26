package utils

import (
	"fmt"
	"io"
	"os"

	"github.com/google/renameio"
)

// CreateAtomic creates or overwrites a file at dest atomically using
// data from r. Atomic means that even in case of a power outage,
// dest will never be a zero-length file. It will always either contain
// the previous data (or not exist) or the new data but never anything
// in between.
func CreateAtomic(dest string, fileMode os.FileMode, r io.Reader) error {
	tmpFile, err := renameio.TempFile("", dest)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Cleanup() //nolint:errcheck

	if err := tmpFile.Chmod(fileMode); err != nil {
		return fmt.Errorf("failed to update mode bits of temp file: %w", err)
	}

	if _, err := io.Copy(tmpFile, r); err != nil {
		return fmt.Errorf("failed to copy source file: %w", err)
	}

	if err := tmpFile.CloseAtomicallyReplace(); err != nil {
		return fmt.Errorf("failed to rename temp file to %q", dest)
	}

	return nil
}

// CopyAtomicMode is like CreateAtomic but copies data from src.
func CopyAtomicMode(src, dst string, mode os.FileMode) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	return CreateAtomic(dst, mode, f)
}

// CopyAtomicKeepMode is like CopyAtomicMode by tries to keep the
// mode bits of dst if it exists. If dst does not yet exist the
// mode bits are set to defaultMode.
func CopyAtomicKeepMode(src, dst string, defaultMode os.FileMode) error {
	mode := defaultMode
	dstStat, err := os.Lstat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		mode = dstStat.Mode()
	}

	return CopyAtomicMode(src, dst, mode)
}
