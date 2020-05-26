package resizer

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func requireImageDimensions(t *testing.T, path string, width int64, height int64) {
	t.Helper()
	var output bytes.Buffer
	cmd := exec.Command("magick", "identify", path)
	cmd.Stdout = &output
	require.NoError(t, cmd.Run())
	elems := strings.Split(output.String(), " ")
	if len(elems) < 3 {
		t.Errorf("Failed to retrieve dimensions of %s", path)
		t.FailNow()
	}
	require.Equal(t, fmt.Sprintf("%dx%d", width, height), elems[2])
}
