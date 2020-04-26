package change

import (
	"encoding/hex"
	"io"
	"os"

	"github.com/twmb/murmur3"
)

// FileChecksum computes the a non-cryptographic hash
// (currently a Murmur3) suitable to compare files.
func FileChecksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := murmur3.New128()

	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}

// FileUpdateNeeded checks if target needs to be updated
// to match ref. Target may not yet exist in which case
// true is returned without an error. If target exists,
// both files are compared using their checksum (see
// FileChecksum). FileUpdateNeeded returns an error
// if ref does not exist.
func FileUpdateNeeded(ref, target string) (bool, error) {
	targetSum, err := FileChecksum(target)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	refSum, err := FileChecksum(ref)
	if err != nil {
		return false, err
	}

	return refSum != targetSum, nil
}

// EqualFileMode checks if f1 and f2 have the same
// mode bits set. If either f1 or f2 does not exist
// or failed to LStat, an error is returned.
func EqualFileMode(f1, f2 string) (bool, error) {
	f1Stat, err := os.Lstat(f1)
	if err != nil {
		return false, err
	}

	f2Stat, err := os.Lstat(f2)
	if err != nil {
		return false, err
	}

	return f1Stat.Mode() == f2Stat.Mode(), nil
}

// CheckFileMode checks if path has it's mode bits set to
// expectedMode.
func CheckFileMode(path string, expectedMode os.FileMode) (bool, error) {
	stat, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return expectedMode == stat.Mode(), nil
}

// EnsureFileMode ensures that path has mode. It returns true
// if the mode bits were updated, false otherwise.
func EnsureFileMode(path string, mode os.FileMode) (bool, error) {
	sameMode, err := CheckFileMode(path, mode)
	if err != nil || sameMode {
		return !sameMode, err
	}

	if err := os.Chmod(path, mode); err != nil {
		return false, err
	}
	return true, nil
}
