package auth

import (
	"context"
	"regexp"
	"time"

	"github.com/volatiletech/authboss/v3"
	_ "github.com/volatiletech/authboss/v3/auth"
	"github.com/volatiletech/authboss/v3/defaults"
	_ "github.com/volatiletech/authboss/v3/register"
	"github.com/waas-app/WaaS/config"
	"github.com/waas-app/WaaS/util"
	"go.uber.org/zap"
)

var (
	ab = authboss.New()
)

const (
	sessionCookieName = "waas_cookie"
)

func GetAuthBoss() *authboss.Authboss {
	return ab
}

func InitializeAuthBoss() error {
	authStore := NewAuthStore()

	cookieStore := NewCookieStorer([]byte(config.Spec.SessionSecret), nil)
	sessionStore := NewSessionStorer(sessionCookieName, []byte(config.Spec.SessionSecret))

	ab.Config.Paths.RootURL = config.Spec.RootURL
	ab.Config.Storage.Server = authStore
	ab.Config.Storage.SessionState = sessionStore
	ab.Config.Storage.CookieState = cookieStore
	ab.Config.Paths.Mount = "/auth"
	ab.Config.Core.ViewRenderer = defaults.JSONRenderer{}
	ab.Config.Core.Logger = NewLogger()
	ab.Config.Modules.ExpireAfter = time.Duration(30 * 24 * time.Hour)
	emailRule := defaults.Rules{
		FieldName: "email", Required: true,
		MatchError: "Must be a valid e-mail address",
		MustMatch:  regexp.MustCompile(`.*@.*\.[a-z]+`),
	}
	passwordRule := defaults.Rules{
		FieldName: "password", Required: true,
		MinLength: 4,
	}
	nameRule := defaults.Rules{
		FieldName: "name", Required: true,
		MinLength: 2,
	}

	ab.Config.Core.BodyReader = defaults.HTTPBodyReader{
		ReadJSON: true,
		Rulesets: map[string][]defaults.Rules{
			"register":    {emailRule, passwordRule, nameRule},
			"recover_end": {passwordRule},
		},
		Whitelist: map[string][]string{
			"register": {"email", "name", "password"},
		},
	}
	defaults.SetCore(&ab.Config, true, false)

	if err := ab.Init(); err != nil {
		util.Logger(context.Background()).Fatal("Failed to initialize Auth boss", zap.Error(err))
		return err
	}

	return nil
}
