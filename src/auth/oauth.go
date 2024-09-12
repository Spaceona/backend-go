package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

// https://accounts.google.com/o/oauth2/auth?client_id={clientid}&redirect_uri={redirectURI}&scope={scope}&response_type=code
// https://www.googleapis.com/auth/userinfo.email
//
//	https://www.googleapis.com/auth/userinfo.profile
//	openid

// https://oauth2.googleapis.com/tokeninfo?id_token=XYZ123 verify token

func GoogleConsent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("wahtttttt")
	scopesArray := []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile", "openid"}
	scopes := strings.Join(scopesArray, " ")
	clientId := os.Getenv("GOOGLE_CLIENT_ID")
	redirectUrl := os.Getenv("REDIRECT_URL")
	fmt.Println("redirectUrl:", redirectUrl)
	url := fmt.Sprintf("https://accounts.google.com/o/oauth2/auth?client_id=%s&access_type=offline&&redirect_uri=%s&scope=%s&response_type=code", clientId, redirectUrl, scopes)
	fmt.Println(url)
	http.Redirect(w, r, url, http.StatusPermanentRedirect)
}

type GoogleTokenRequest struct {
	AccessCode   string `json:"code"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	RedirectUri  string `json:"redirect_uri"`
	GrantType    string `json:"grant_type"`
}

type GoogleTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type SpaceonaUserToken struct {
	Email           string `json:"email"`
	GoogleTokenInfo GoogleTokenResponse
}

func Callback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	code := query.Get("code")

	googleTokenRequest := GoogleTokenRequest{
		AccessCode:   code,
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		ClientId:     os.Getenv("GOOGLE_CLIENT_ID"),
		RedirectUri:  os.Getenv("REDIRECT_URL"),
		GrantType:    "authorization_code",
	}

	requestString, marshalError := json.Marshal(&googleTokenRequest)
	if marshalError != nil {
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}
	req, requestCreationError := http.NewRequest("POST", "https://oauth2.googleapis.com/token", bytes.NewBuffer(requestString))
	if requestCreationError != nil {
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, tokenError := client.Do(req)
	if tokenError != nil {
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}(response.Body)
	body, resReadError := io.ReadAll(response.Body)
	if resReadError != nil {
		slog.Error(resReadError.Error())
	}
	var resJson GoogleTokenResponse
	decodeError := json.Unmarshal(body, &resJson)
	if decodeError != nil {
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}
	token := SpaceonaUserToken{}
	spaceonaToken, spaceonaTokenError := GenToken(token, time.Duration(resJson.ExpiresIn)*time.Second)
	if spaceonaTokenError != nil {
		slog.Error(spaceonaTokenError.Error())
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:     "AuthToken",
		Value:    spaceonaToken,
		Domain:   "localhost",
		Path:     "/",
		Expires:  time.Now().Add(time.Duration(resJson.ExpiresIn) * time.Second),
		MaxAge:   resJson.ExpiresIn,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	fmt.Println("cookie:", cookie.String())
	http.SetCookie(w, &cookie)
	fmt.Println(w.Header().Values("Set-Cookie"))
	http.Redirect(w, r, "https://spaceona.com", http.StatusPermanentRedirect)
}
