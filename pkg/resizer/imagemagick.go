package resizer

import (
	"context"
	"fmt"
	"os/exec"
)

// ImageMagickResizer implements Resizer using the `magick` CLI.
type ImageMagickResizer struct{}

func (r *ImageMagickResizer) Resize(ctx context.Context, out string, in string, width int64, height int64) error {
	return exec.CommandContext(ctx, "magick", in, "-resize", fmt.Sprintf("%dx%d", width, height), out).Run()
}

func NewImageMagick() Resizer {
	return &ImageMagickResizer{}
}
