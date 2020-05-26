package resizer

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImageMagicResizer(t *testing.T) {
	ctx := context.Background()
	r := NewImageMagick()
	os.RemoveAll("testdata/out.jpg")
	requireImageDimensions(t, "testdata/in.jpg", 80, 80)
	require.NoError(t, r.Resize(ctx, "testdata/out.jpg", "testdata/in.jpg", 40, 40))
	require.FileExists(t, "testdata/out.jpg")
	requireImageDimensions(t, "testdata/out.jpg", 40, 40)
}
