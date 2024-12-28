package installer

import (
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestInstaller_InsertFrontMatter_Error(t *testing.T) {
	tpmDir := t.TempDir()
	templatePath := filepath.Join(tpmDir, "template.md")
	err := os.WriteFile(templatePath, []byte("Hello, World!"), 0)
	require.Error(t, err)

	srcFile := filepath.Join(tpmDir, "file.md")
	err = os.WriteFile(srcFile, []byte("Hello, World!"), 0)
	require.NoError(t, err)

	in := New(srcFile, templatePath, nil, 1)

	err = in.Run()
	require.Error(t, err)
}
