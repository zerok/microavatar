package resizer

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog"
)

// ImageMagickResizer implements Resizer using the `magick` CLI.
type ImageMagickResizer struct{}

func (r *ImageMagickResizer) Resize(ctx context.Context, out string, in string, width int64, height int64) error {
	logger := zerolog.Ctx(ctx)
	bin, err := exec.LookPath("magick")
	if err != nil {
		logger.Debug().Msg("magick command not found. Falling back to convert.")
		// Let's see if the im 6.x commands are available
		bin, err = exec.LookPath("convert")
		if err != nil {
			return err
		}
	}
	logger.Debug().Msgf("%s -> %s", in, out)
	cmd := exec.CommandContext(ctx, bin, in, "-resize", fmt.Sprintf("%dx%d", width, height), out)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func NewImageMagick() Resizer {
	return &ImageMagickResizer{}
}
