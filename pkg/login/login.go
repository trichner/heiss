package login

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/gob"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/trichner/heiss/pkg/login/session"
	"github.com/trichner/heiss/pkg/login/signer"
)

const tokenCookie = "token"

//go:embed login.html
var loginPageHtml string

type contextKey string

var sessionContextKey = contextKey("session")

type Signer interface {
	Sign(data []byte) (signature []byte, err error)
	Verify(data []byte, signature []byte) error
}

func WithContextSession(ctx context.Context, session *session.Session) context.Context {
	return context.WithValue(ctx, sessionContextKey, session)
}

func GetContextSession(ctx context.Context) *session.Session {
	v := ctx.Value(sessionContextKey)
	s, ok := v.(*session.Session)
	if !ok {
		return nil
	}
	return s
}

func NewLoginHandler(next http.Handler, sessionManager session.Manager, key []byte) http.Handler {
	signer := signer.NewHmacSigner(key)
	return &pwHandler{
		next:     next,
		userRepo: sessionManager,
		signer:   signer,
	}
}

type pwHandler struct {
	next     http.Handler
	userRepo session.Manager
	signer   Signer
}

func (p *pwHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	// let's parse the URL
	reqUrl, err := url.ParseRequestURI(req.RequestURI)
	if err != nil {
		http.Error(res, "", 400)
		return
	}

	var session *session.Session

	// find the associated session or create one
	cookie, err := req.Cookie(tokenCookie)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		http.Error(res, "", 500)
		return
	}

	// we have a session! let's add it to the context
	if err == nil {
		session, err = p.verify(cookie)
		if err != nil {
			// session expired/invalid, let's send the user to login again
			p.redirectToLogin(res, req)
			return
		}
		newContext := WithContextSession(req.Context(), session)
		req = req.WithContext(newContext)
	}

	// should we just log out?
	if reqUrl.Path == "/logout" && req.Method == http.MethodGet {
		p.redirectToLogin(res, req)
		return
	}

	// are we on the login page?
	if reqUrl.Path == "/login" {

		if session != nil {
			http.Redirect(res, req, "/", http.StatusFound)
			return
		}

		// just showing?
		if req.Method == http.MethodGet {
			password := reqUrl.Query().Get("p")
			if password != "" {
				p.handlePasswordEntered(res, req, password)
				return
			}
			err = p.showLoginPage(res, req)
			if err != nil {
				http.Error(res, "", 500)
			}
			return
		}

		// are we sending credentials? check them!
		if req.Method == http.MethodPost {
			err := req.ParseForm()
			if err != nil {
				http.Error(res, "", 400)
				return
			}

			password := req.FormValue("password")

			p.handlePasswordEntered(res, req, password)
			return
		}
	}

	if session != nil {
		p.next.ServeHTTP(res, req)
		return
	}

	// neither logged-in nor on the login page, let's guide the user there
	p.redirectToLogin(res, req)
	return
}

func (p *pwHandler) redirectToLogin(res http.ResponseWriter, req *http.Request) {
	http.SetCookie(res, &http.Cookie{Name: tokenCookie, Expires: time.UnixMilli(0)})
	http.Redirect(res, req, "/login", http.StatusFound)
}

func (p *pwHandler) handlePasswordEntered(res http.ResponseWriter, req *http.Request, password string) {
	token, err := p.createSession(password)
	if err != nil {
		// bad password!
		err = p.showBadLoginPage(res, req)
		return
	}
	http.SetCookie(res, &http.Cookie{
		Name:       tokenCookie,
		Value:      token,
		Expires:    time.Now().Add(time.Hour * 24 * 30),
		RawExpires: "",
		MaxAge:     0,
		Secure:     false, // TODO
		HttpOnly:   true,  // TODO
		SameSite:   http.SameSiteStrictMode,
	})
	http.Redirect(res, req, "/", http.StatusFound)
}

// createSession tries to log in a user based on the password and returns
// the signed session token or an error
func (p *pwHandler) createSession(password string) (string, error) {
	session, err := p.userRepo.Login(password)
	if err != nil {
		return "", fmt.Errorf("cannot login user: %w", err)
	}
	signedToken, err := signSession(p.signer, session)
	if err != nil {
		return "", fmt.Errorf("failed to sign session: %w", err)
	}
	token := CreateCookie(signedToken)
	return token, nil
}

func (p *pwHandler) showLoginPage(res http.ResponseWriter, req *http.Request) error {
	res.Header().Set("content-type", "text/html;charset=utf8")
	_, err := res.Write([]byte(loginPageHtml))
	return err
}

func (p *pwHandler) showBadLoginPage(res http.ResponseWriter, req *http.Request) error {
	res.Header().Set("content-type", "text/html;charset=utf8")
	res.WriteHeader(http.StatusForbidden)
	_, err := res.Write([]byte(loginPageHtml))
	return err
}

func (p *pwHandler) verify(cookie *http.Cookie) (*session.Session, error) {
	// verify signature

	signedToken, err := ParseCookie(cookie.Value)
	if err != nil {
		return nil, fmt.Errorf("cannot parse token: %w", err)
	}

	err = p.signer.Verify(signedToken.Payload, signedToken.Signature)
	if err != nil {
		return nil, fmt.Errorf("cannot verify token: %w", err)
	}

	return decodeSession(signedToken.Payload)
}

func signSession(signer Signer, session *session.Session) (*SignedToken, error) {
	data, err := encodeSession(session)
	if err != nil {
		return nil, err
	}
	signature, err := signer.Sign(data)
	if err != nil {
		return nil, err
	}
	return &SignedToken{
		Payload:   data,
		Signature: signature,
	}, nil
}

func encodeSession(session *session.Session) ([]byte, error) {
	var buf bytes.Buffer
	e := gob.NewEncoder(&buf)

	err := e.Encode(*session)
	if err != nil {
		return nil, fmt.Errorf("cannot encode session: %w", err)
	}
	return buf.Bytes(), nil
}

func decodeSession(data []byte) (*session.Session, error) {
	d := gob.NewDecoder(bytes.NewReader(data))

	var session session.Session
	err := d.Decode(&session)
	if err != nil {
		return nil, fmt.Errorf("cannot decode session: %w", err)
	}
	return &session, nil
}
