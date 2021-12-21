package cmd

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/afero"
)

// RotItem is modelling one file or folder to rot
type RotItem struct {
	Path             string    `json:"path"`
	IsFolder         bool      `json:"is_folder,omitempty"`
	Hash             string    `json:"hash"` // currently SHA-256
	DeleteIfModified bool      `json:"delete_if_modified,omitempty"`
	AddedAt          time.Time `json:"added_at"`
	DeleteAt         time.Time `json:"delete_at,omitempty"`
}

// NewRotItem returns a new RotItem if the referenced path exists
func NewRotItem(path string, deleteIfModified bool) (RotItem, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}

	ri := RotItem{}
	fi, err := AppFs.Stat(absPath)
	if err != nil {
		// check existence
		if os.IsNotExist(err) {
			return ri, errors.New("File/Folder not existing: " + absPath)
		}
		return ri, err
	}
	mode := fi.Mode()
	if !mode.IsDir() && !mode.IsRegular() {
		return ri, errors.New("Neither file nor folder: " + absPath)
	}

	isFolder := mode.IsDir()
	hash := getHash(path, isFolder)
	if hash == "" {
		return ri, errors.New("Hashing failed for file or folder: " + absPath)
	}

	ri = RotItem{
		Path:             absPath,
		IsFolder:         isFolder,
		AddedAt:          time.Now(),
		Hash:             hash,
		DeleteIfModified: deleteIfModified,
	}
	return ri, nil
}

// SetDeletionDate sets the date for deletion of the referenced item
func (r *RotItem) SetDeletionDate(deletionDate time.Time) {
	r.DeleteAt = deletionDate
}

// SetDeletionDuration sets the date for deletion of the referenced item by supplying an ISO 8601 duration string
func (r *RotItem) SetDeletionDuration(deletionDuration string) {
	duration := parseDuration(deletionDuration)
	r.DeleteAt = r.AddedAt.Add(duration)
}

// IsRotten determines if the referenced item could be deleted at given time.
func (r RotItem) IsRotten(referenceDate time.Time) bool {
	return !r.DeleteAt.IsZero() && referenceDate.After(r.DeleteAt)
}

// HasChanged determines if the referenced item has changed since being staged for rotting.
func (r RotItem) HasChanged() bool {
	currHash := getHash(r.Path, r.IsFolder)
	return currHash != r.Hash
}

// Clean removes this item if it is rotten at the given date.
func (r RotItem) Clean(referenceDate time.Time, dryRun bool) (deleted bool, err error) {
	if !r.IsRotten(referenceDate) {
		return false, nil
	}
	if !r.DeleteIfModified && r.HasChanged() {
		return false, nil
	}

	if !dryRun {
		if r.IsFolder {
			err = AppFs.RemoveAll(r.Path)
		} else {
			err = AppFs.Remove(r.Path)
		}
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func parseDuration(str string) time.Duration {
	durationRegex := regexp.MustCompile(`^P?(?P<years>\d+Y)?(?P<months>\d+M)?(?P<days>\d+D)?T?(?P<hours>\d+H)?(?P<minutes>\d+M)?(?P<seconds>\d+S)?$`)
	matches := durationRegex.FindStringSubmatch(strings.ToUpper(str))

	years := parseInt64(matches[1])
	months := parseInt64(matches[2])
	days := parseInt64(matches[3])
	hours := parseInt64(matches[4])
	minutes := parseInt64(matches[5])
	seconds := parseInt64(matches[6])

	hour := int64(time.Hour)
	minute := int64(time.Minute)
	second := int64(time.Second)
	return time.Duration(years*24*365*hour + months*30*24*hour + days*24*hour + hours*hour + minutes*minute + seconds*second)
}

func parseInt64(value string) int64 {
	if len(value) == 0 {
		return 0
	}
	parsed, err := strconv.Atoi(value[:len(value)-1])
	if err != nil {
		return 0
	}
	return int64(parsed)
}

func getHash(path string, isFolder bool) string {
	var hash string
	if isFolder {
		hash = hashFolder(path)
	} else {
		hash = hashFile(path)
	}
	return hash
}

func hashFolder(path string) string {
	var hash = ""
	h := sha256.New()
	var printFile = func(innerPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		hash = hash + hashFile(innerPath)
		h.Reset()
		h.Write([]byte(hash))
		hash = fmt.Sprintf("%x", h.Sum(nil))
		return nil
	}
	err := afero.Walk(AppFs, path, printFile)
	if err != nil {
		return ""
	}
	return hash
}

func hashFile(filename string) string {
	f, err := AppFs.Open(filename)
	if err != nil {
		return ""
	}
	defer func(f afero.File) {
		err := f.Close()
		if err != nil {
			log.Fatalf("could not hash file: %s", filename)
		}
	}(f)

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return ""
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
