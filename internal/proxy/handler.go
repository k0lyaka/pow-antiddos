package proxy

import (
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/k0lyaka/pow-antiddos/internal/config"
	"github.com/k0lyaka/pow-antiddos/internal/session"
)

type ProxyHandlerWithConfig struct {
	Config    config.ConfigModel
	Templates *template.Template
}

func (p *ProxyHandlerWithConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ProxyHandler(w, r, p.Config, p.Templates)
}

func challengeHandler(w http.ResponseWriter, r *http.Request, ses *session.Session, session_id string, templates *template.Template) {
	if r.Method == http.MethodPost && !ses.Authorized {
		r.ParseForm()
		nonce := r.Form.Get("nonce")

		if Validate(ValidationRequest{Nonce: nonce, Prefix: ses.Prefix, Difficulty: config.Config.Difficulty}) {
			session.AuthorizeSession(session_id)

			http.Redirect(w, r, r.URL.Path, http.StatusFound)
			return
		}
	}

	templates.ExecuteTemplate(w, "challenge.html", map[string]string{"Prefix": ses.Prefix, "Difficulty": strconv.Itoa(config.Config.Difficulty)})
}

func handleNewSession(w http.ResponseWriter, r *http.Request, templates *template.Template) {
	sid, ses := session.NewSession()

	http.SetCookie(w, &http.Cookie{
		Name:  "__pow_session_id",
		Value: sid,
		Path:  "/",
	})

	challengeHandler(w, r, ses, sid, templates)
}

func ProxyHandler(w http.ResponseWriter, r *http.Request, config config.ConfigModel, templates *template.Template) {
	cookie, err := r.Cookie("__pow_session_id")

	if err != nil {
		handleNewSession(w, r, templates)
		return
	}

	ses, err := session.GetSession(cookie.Value)
	if err != nil {
		handleNewSession(w, r, templates)
		return
	}

	if !ses.Authorized {
		challengeHandler(w, r, ses, cookie.Value, templates)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   config.BackendURL,
	})
	proxy.ServeHTTP(w, r)
}
