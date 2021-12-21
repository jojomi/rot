package cmd

import (
	"os"
	"testing"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var now = makeTestDate(2017, 11, 22)
var yesterday = makeTestDate(2017, 11, 21)
var tomorrow = makeTestDate(2017, 11, 23)

func TestIsRotten(t *testing.T) {
	added := makeTestDate(2016, 10, 22)
	riYesterday := RotItem{
		AddedAt:  added,
		DeleteAt: yesterday,
	}
	assert.True(t, riYesterday.IsRotten(now))

	riTomorrow := RotItem{
		AddedAt:  added,
		DeleteAt: tomorrow,
	}
	assert.False(t, riTomorrow.IsRotten(now))

	riUndefined := RotItem{
		AddedAt:  added,
		DeleteAt: time.Time{},
	}
	assert.False(t, riUndefined.IsRotten(now))
}

func TestHasChanged(t *testing.T) {
	AppFs = afero.NewMemMapFs()
	filename := "/tmp/testfile"
	makeTestFile(filename, []byte("abc"))
	ri, err := NewRotItem(filename, false)
	assert.NoError(t, err)

	assert.False(t, ri.HasChanged())

	makeTestFile(filename, []byte("abcdef"))
	assert.True(t, ri.HasChanged())

	AppFs.Remove(filename)
	assert.True(t, ri.HasChanged())
}

func TestClean(t *testing.T) {
	var (
		exists, deleted bool
		err             error
		ri              RotItem
	)
	AppFs = afero.NewMemMapFs()
	filename := "/tmp/testfile"
	makeTestFile(filename, []byte(""))
	exists, err = afero.Exists(AppFs, filename)
	assert.NoError(t, err)
	assert.True(t, exists)
	ri, err = NewRotItem(filename, false)
	assert.NoError(t, err)

	ri.DeleteAt = tomorrow
	deleted, err = ri.Clean(now, false)
	assert.NoError(t, err)
	assert.False(t, deleted)
	exists, err = afero.Exists(AppFs, filename)
	assert.NoError(t, err)
	assert.True(t, exists)

	ri.DeleteAt = yesterday

	deleted, err = ri.Clean(now, true)
	assert.NoError(t, err)
	assert.True(t, deleted)
	exists, err = afero.Exists(AppFs, filename)
	assert.NoError(t, err)
	assert.True(t, exists)

	deleted, err = ri.Clean(now, false)
	assert.NoError(t, err)
	assert.True(t, deleted)
	exists, err = afero.Exists(AppFs, filename)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestSetDeletionDuration(t *testing.T) {
	ri := RotItem{
		AddedAt: now,
	}

	ri.SetDeletionDuration("P4D")
	assert.Equal(t, makeTestDate(2017, 11, 26), ri.DeleteAt)

	ri.SetDeletionDuration("P24H")
	assert.Equal(t, tomorrow, ri.DeleteAt)

	ri.SetDeletionDuration("P1M")
	assert.Equal(t, makeTestDate(2017, 12, 22), ri.DeleteAt)

	ri.SetDeletionDuration("P3Y")
	assert.Equal(t, makeTestDate(2020, 11, 21), ri.DeleteAt)
}

func makeTestFile(filename string, content []byte) {
	afero.WriteFile(AppFs, filename, content, os.FileMode(0777))
}

func makeTestDate(year int, month time.Month, day int) time.Time {
	location, err := time.LoadLocation("UTC")
	if err != nil {
		return time.Time{}
	}
	return time.Date(year, month, day, 0, 0, 0, 0, location)
}
