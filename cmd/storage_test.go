package cmd

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestLoadSave(t *testing.T) {
	AppFs = afero.NewMemMapFs()
	data, err := load()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(data))

	ri := RotItem{
		Path:     "folder one",
		IsFolder: true,
		AddedAt:  makeTestDate(2017, 11, 22),
		DeleteAt: makeTestDate(2118, 1, 2),
	}
	data = append(data, ri)

	err = save(data)
	assert.NoError(t, err)

	dataLoaded, err := load()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(dataLoaded))
	riLoaded := dataLoaded[0]
	assert.Equal(t, ri.Path, riLoaded.Path)
	assert.Equal(t, ri.IsFolder, riLoaded.IsFolder)
	assert.Equal(t, ri.AddedAt, riLoaded.AddedAt)
	assert.Equal(t, ri.DeleteAt, riLoaded.DeleteAt)
}
