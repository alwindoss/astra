package handlers

import (
	"net/http"

	"github.com/alwindoss/astra"
	"github.com/justinas/nosurf"
)

type NoSurf struct {
	Cfg *astra.Config
}

func (n NoSurf) NoSurfMW(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   n.Cfg.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}
