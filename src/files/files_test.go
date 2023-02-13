package files

import (
	"context"
	"fmt"
	mocks "github.com/ahmed.sukhera/dls3obj/mocks/src/files"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"strings"
	"testing"
)

func TestDownloadFile(t *testing.T) {
	type testCase struct {
		name        string
		fileOpener  func() FileOpener
		bucketName  string
		bucketPath  string
		expectError bool
	}

	testCases := []testCase{
		{
			name: "invalid fileOpener returns error",
			fileOpener: func() FileOpener {
				return nil
			},
			expectError: true,
		},
		{
			name: "empty bucket name returns error",
			fileOpener: func() FileOpener {
				return mocks.NewFileOpener(t)
			},
			expectError: true,
		},
		{
			name: "empty bucket path returns error",
			fileOpener: func() FileOpener {
				return mocks.NewFileOpener(t)
			},
			bucketName:  "some-bucket",
			expectError: true,
		},
		{
			name: "file opener fails and returns error",
			fileOpener: func() FileOpener {
				m := mocks.NewFileOpener(t)
				m.On("OpenFile", mock.Anything, "some-bucket", "some-path").
					Return(nil, fmt.Errorf("some-error"))
				return m
			},
			bucketName:  "some-bucket",
			bucketPath:  "some-path",
			expectError: true,
		},
		{
			name: "file opener opens file successfully",
			fileOpener: func() FileOpener {
				m := mocks.NewFileOpener(t)
				m.On("OpenFile", mock.Anything, "some-bucket", "some-path").
					Return(io.NopCloser(strings.NewReader("some data")), nil)
				return m
			},
			bucketName:  "some-bucket",
			bucketPath:  "some-path",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			data, err := DownloadFile(ctx, tc.fileOpener(), tc.bucketName, tc.bucketPath)
			if tc.expectError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.NotNil(t, data)
		})
	}
}
