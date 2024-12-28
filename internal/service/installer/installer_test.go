package installer

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestInstaller_InsertFrontMatter_Error(t *testing.T) {
	testCases := []struct {
		name           string
		setup          func() (srcDir, path string, cleanup func())
		expectedErrMsg string
	}{
		{
			name: "Ошибка чтения шаблона",
			setup: func() (path, srcPath string, cleanup func()) {
				tpmDir := t.TempDir()
				templatePath := filepath.Join(tpmDir, "template.md")
				err := os.WriteFile(templatePath, []byte("Hello, World!"), 0000)
				require.NoError(t, err)

				srcFile := filepath.Join(tpmDir, "file.md")
				err = os.WriteFile(srcFile, []byte("Hello, World!"), 0000)
				require.NoError(t, err)

				return templatePath, srcFile, func() {
					_ = os.RemoveAll(srcFile)
					_ = os.RemoveAll(templatePath)
				}
			},
			expectedErrMsg: "не удалось прочитать шаблон",
		},
	}
	var err error
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			templatePath, srcPath, cleanup := tc.setup()
			defer cleanup()

			in := New(templatePath, srcPath, nil, 1)

			err = in.Run()
			require.Error(t, err)
			require.ErrorContains(t, err, tc.expectedErrMsg)

		})
	}
}

func TestInstaller_InsertFrontMatter_Success(t *testing.T) {
	testCases := []struct {
		name  string
		setup func() (srcDir, path string, cleanup func())
	}{
		{
			name: "Файл остался без изменений",
			setup: func() (path, srcPath string, cleanup func()) {
				tpmDir := t.TempDir()
				templatePath := filepath.Join(tpmDir, "template.md")
				err := os.WriteFile(templatePath, []byte("Hello, World!"), 0777)
				require.NoError(t, err)

				srcFile := filepath.Join(tpmDir, "file.md")
				err = os.WriteFile(srcFile, []byte("Hello, World!"), 0777)
				require.NoError(t, err)

				return templatePath, srcFile, func() {
					_ = os.RemoveAll(srcFile)
					_ = os.RemoveAll(templatePath)
				}
			},
		},
		{
			name: "Директория пропускается по паттерну",
			setup: func() (path, srcPath string, cleanup func()) {
				tpmDir := t.TempDir()
				templatePath := filepath.Join(tpmDir, "template.md")
				err := os.WriteFile(templatePath, []byte("Hello, World!"), 0777)
				require.NoError(t, err)

				customDir := t.TempDir() + "/_gaslgda"
				err = os.Mkdir(customDir, 0777)
				require.NoError(t, err)

				srcFile := filepath.Join(tpmDir, "file.md")
				err = os.WriteFile(srcFile, []byte("Hello, World!"), 0777)
				require.NoError(t, err)

				return templatePath, srcFile, func() {
					_ = os.RemoveAll(srcFile)
					_ = os.RemoveAll(templatePath)
				}
			},
		},
	}
	var err error
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			templatePath, srcPath, cleanup := tc.setup()
			defer cleanup()

			in := New(templatePath, srcPath, []string{"_", "."}, 1)

			err = in.Run()
			require.NoError(t, err)
		})
	}
}
