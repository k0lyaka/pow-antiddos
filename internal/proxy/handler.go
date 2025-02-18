package proxy

import (
	"context"
	"html/template"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"

	"github.com/go-redis/redis_rate/v10"
	"github.com/k0lyaka/pow-antiddos/internal/config"
	"github.com/k0lyaka/pow-antiddos/internal/redis"
	"github.com/k0lyaka/pow-antiddos/internal/session"
	"github.com/k0lyaka/pow-antiddos/internal/utils"
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

	templates.ExecuteTemplate(w, "challenge.html", map[string]string{"Prefix": ses.Prefix, "Difficulty": strconv.Itoa(config.Config.Difficulty), "PrefixShort": ses.Prefix[:16]})
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

	// apply rate limiter
	ip, err := utils.ExtractIP(r)

	if err != nil {
		w.Write([]byte("500: Internal Server Error"))
		return
	}

	if config.RateLimitEnabled {
		ctx := context.Background()
		res, err := redis.Limiter.Allow(ctx, "rate-limiter:"+ip, redis_rate.PerSecond(config.RateLimit))

		if err != nil {
			w.Write([]byte("500: Internal Server Error"))
			return
		}

		if res.Remaining == 0 {
			templates.ExecuteTemplate(w, "409.html", nil)
			return
		}
	}

	backendURL, _ := url.Parse(config.BackendURL)
	proxy := httputil.NewSingleHostReverseProxy(backendURL)
	proxy.ServeHTTP(w, r)
}
