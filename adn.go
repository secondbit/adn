package adn

import (
	"errors"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	SCOPE_BASIC           = "basic"
	SCOPE_STREAM          = "stream"
	SCOPE_WRITE_POST      = "write_post"
	SCOPE_FOLLOW          = "follow"
	SCOPE_PUBLIC_MESSAGES = "public_messages"
	SCOPE_MESSAGES        = "messages"
	SCOPE_UPDATE_PROFILE  = "update_profile"
	SCOPE_FILES           = "files"
	SCOPE_EXPORT          = "export"
)

type ADN struct {
	Token        string
	Scopes       []string
	ClientID     string
	ClientSecret string
	RedirectURI  string
}

var TokenNotSet = errors.New("Access token not set.")
var NoScopesSet = errors.New("Scopes not set.")
var NoClientIDSet = errors.New("Client ID not set.")
var NoRedirectURISet = errors.New("Redirect URI not set.")

func NewClient(client_id, client_secret, redirect_uri string, scopes []string) *ADN {
	return &ADN{
		Scopes:       scopes,
		ClientID:     client_id,
		ClientSecret: client_secret,
		RedirectURI:  redirect_uri,
	}
}

func (a *ADN) makeURL(endpoint string) string {
	endpoint = strings.TrimLeft(endpoint, "/")
	return "https://alpha-api.app.net/" + endpoint
}

func (a *ADN) makeRequest(method, urlStr string, body io.Reader) (*http.Request, error) {
	if a.Token == "" {
		return nil, TokenNotSet
	}
	req, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		return req, err
	}
	req.Header.Add("Authorization", "Bearer "+a.Token)
	req.Header.Add("Content-Type", "application/json")
	return req, nil
}

func (a *ADN) GetClientSideAuthURL() (string, error) {
	if len(a.Scopes) == 0 {
		return "", NoScopesSet
	}
	if a.ClientID == "" {
		return "", NoClientIDSet
	}
	if a.RedirectURI == "" {
		return "", NoRedirectURISet
	}
	params := make(url.Values)
	params.Add("client_id", a.ClientID)
	params.Add("response_type", "token")
	params.Add("redirect_uri", a.RedirectURI)
	params.Add("scope", strings.Join(a.Scopes, " "))
	return "https://account.app.net/oauth/authenticate?" + params.Encode(), nil
}

type clientSideAuthListener struct {
	listener chan string
}

func (c *clientSideAuthListener) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	token := req.FormValue("access_token")
	if token == "" {
		io.WriteString(w, "Access token not found in the request fragment.")
		return
	}
	c.listener <- token
	io.WriteString(w, "Successfully obtained access token. You can close this window.")
}

var redirectPage = `<html>
<head>
<script>
window.location = "/auth/token?" + window.location.hash.substr(1, window.location.hash.length - 1);
</script>
</head>
<body>
You should be redirected momentarily.
</body>
</html>`

func ServeRedirect(w http.ResponseWriter, req *http.Request) {
	tmp, err := template.New("redirect").Parse(redirectPage)
	if err != nil {
		io.WriteString(w, "Error parsing redirect template.\n"+err.Error()+"\n")
	}
	err = tmp.Execute(w, nil)
	if err != nil {
		io.WriteString(w, "Error executing template.\n"+err.Error()+"\n")
	}
}

func (a *ADN) ListenForClientSideAuth() (string, error) {
	if a.RedirectURI == "" {
		return "", NoRedirectURISet
	}
	red_url, err := url.Parse(a.RedirectURI)
	if err != nil {
		return "", err
	}
	port := ":80"
	colIndex := strings.Index(red_url.Host, ":")
	if colIndex > -1 {
		port = red_url.Host[colIndex:]
	}
	authListener := new(clientSideAuthListener)
	authListener.listener = make(chan string)
	go func() {
		http.Handle("/auth/token", authListener)
		http.HandleFunc("/", ServeRedirect)
		err := http.ListenAndServe(port, nil)
		if err != nil {
			panic(err)
		}
	}()
	token := <-authListener.listener
	return token, nil
}
