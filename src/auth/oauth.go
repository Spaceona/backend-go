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

type GoogleToken struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type SpaceonaUserToken struct {
	UserInfo    GoogleUserInfo `json:"email"`
	GoogleToken GoogleToken
}

func Callback(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	code := query.Get("code")

	googleToken, tokenErr := GetGoogleToken(code)
	if tokenErr != nil {
		slog.Error(tokenErr.Error())
		http.Error(w, "failed to login", http.StatusInternalServerError)
		return
	}
	userInfo, userInfoErr := GetGoogleUserInfo(googleToken)
	if userInfoErr != nil {
		slog.Error("user info err", "error", userInfoErr.Error())
		http.Error(w, "failed to login", http.StatusInternalServerError)
	}
	//TODO check if the user is an admin

	token := SpaceonaUserToken{}
	token.GoogleToken = googleToken
	token.UserInfo = userInfo
	spaceonaToken, spaceonaTokenError := GenToken(token, time.Duration(googleToken.ExpiresIn)*time.Second)

	slog.Info(spaceonaToken)
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
		Expires:  time.Now().Add(time.Duration(googleToken.ExpiresIn) * time.Second),
		MaxAge:   googleToken.ExpiresIn,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "https://spaceona.com", http.StatusPermanentRedirect)
}

func AuthWithRefreshToken(w http.ResponseWriter, r *http.Request) {

}

func GetGoogleToken(code string) (GoogleToken, error) {
	googleTokenRequest := GoogleTokenRequest{
		AccessCode:   code,
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		ClientId:     os.Getenv("GOOGLE_CLIENT_ID"),
		RedirectUri:  os.Getenv("REDIRECT_URL"),
		GrantType:    "authorization_code",
	}

	requestString, marshalError := json.Marshal(&googleTokenRequest)
	if marshalError != nil {
		return GoogleToken{}, marshalError
	}

	req, requestCreationError := http.NewRequest("POST", "https://oauth2.googleapis.com/token", bytes.NewBuffer(requestString))
	if requestCreationError != nil {
		return GoogleToken{}, requestCreationError
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, tokenError := client.Do(req)
	if tokenError != nil {
		return GoogleToken{}, tokenError
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
	var resJson GoogleToken
	decodeError := json.Unmarshal(body, &resJson)
	if decodeError != nil {
		return GoogleToken{}, decodeError
	}
	return resJson, nil
}

type GoogleUserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

func GetGoogleUserInfo(token GoogleToken) (GoogleUserInfo, error) {
	req, requestCreationError := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v1/userinfo?alt=json", bytes.NewBuffer([]byte("")))
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	if requestCreationError != nil {
		return GoogleUserInfo{}, requestCreationError
	}
	client := &http.Client{}
	response, infoErr := client.Do(req)
	if infoErr != nil {
		return GoogleUserInfo{}, infoErr
	}
	body, resReadError := io.ReadAll(response.Body)
	if resReadError != nil {
		return GoogleUserInfo{}, resReadError
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			slog.Error(err.Error())
		}
	}()
	var userInfo GoogleUserInfo
	decodeError := json.Unmarshal(body, &userInfo)
	if decodeError != nil {
		return GoogleUserInfo{}, decodeError
	}
	return userInfo, nil
}
