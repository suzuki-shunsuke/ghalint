package act

import (
	"io"

	"github.com/spf13/afero"
)

type Controller struct {
	fs     afero.Fs
	stderr io.Writer
}

type InputNew struct {
	Stderr io.Writer
}

func New(fs afero.Fs, input *InputNew) *Controller {
	return &Controller{
		fs:     fs,
		stderr: input.Stderr,
	}
}
