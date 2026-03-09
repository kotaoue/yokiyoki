package locale

import (
	"embed"

	"github.com/BurntSushi/toml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.toml
var localeFS embed.FS

var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
	if _, err := bundle.LoadMessageFileFS(localeFS, "locales/active.en.toml"); err != nil {
		panic(err)
	}
	if _, err := bundle.LoadMessageFileFS(localeFS, "locales/active.ja.toml"); err != nil {
		panic(err)
	}
}

// NewLocalizer creates a new i18n.Localizer for the given language tag.
func NewLocalizer(lang string) *i18n.Localizer {
	return i18n.NewLocalizer(bundle, lang)
}
