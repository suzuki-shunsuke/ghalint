package token

import (
	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/ghalint/pkg/github"
	"github.com/suzuki-shunsuke/slog-logrus/slogrus"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/keyring/ghtoken"
	"github.com/urfave/cli/v3"
)

func New(logE *logrus.Entry) *cli.Command {
	return ghtoken.Command(ghtoken.NewActor(slogrus.Convert(logE), github.KeyService))
}
