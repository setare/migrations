//go:generate go run github.com/golang/mock/mockgen -package sql -destination fs_mock_test.go io/fs DirEntry

package sql

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_parseSQLFile(t *testing.T) {
	t.Run("should match a file with type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		dirEntry := NewMockDirEntry(ctrl)

		dirEntry.EXPECT().IsDir().Return(false)
		wantID := "123"
		wantDescription := "description with new data"
		wantType := "down"
		givenName := "123_description_with_new_data.down.sql"
		dirEntry.EXPECT().Name().Return(givenName)

		gotID, gotDescription, gotType := parseSQLFile(dirEntry)
		assert.Equal(t, gotID, wantID)
		assert.Equal(t, gotDescription, wantDescription)
		assert.Equal(t, gotType, wantType)
	})

	t.Run("should match a file without type", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		dirEntry := NewMockDirEntry(ctrl)

		dirEntry.EXPECT().IsDir().Return(false)
		wantID := "123"
		wantDescription := "description with new data"
		wantType := ""
		givenName := "123_description_with_new_data.sql"
		dirEntry.EXPECT().Name().Return(givenName)

		gotID, gotDescription, gotType := parseSQLFile(dirEntry)
		assert.Equal(t, gotID, wantID)
		assert.Equal(t, gotDescription, wantDescription)
		assert.Equal(t, gotType, wantType)
	})

	t.Run("should not match when the entry is a directory", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		dirEntry := NewMockDirEntry(ctrl)

		dirEntry.EXPECT().IsDir().Return(true)

		gotID, _, _ := parseSQLFile(dirEntry)
		assert.Empty(t, gotID)
	})

	t.Run("should not match file name", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		dirEntry := NewMockDirEntry(ctrl)

		dirEntry.EXPECT().IsDir().Return(false)
		dirEntry.EXPECT().Name().Return("crazy_file_name.txt")

		gotID, _, _ := parseSQLFile(dirEntry)
		assert.Empty(t, gotID)
	})
}
