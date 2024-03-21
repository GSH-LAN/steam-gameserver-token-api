package app

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gsh-lan/steam-gameserver-token-api/src/logger"
	"github.com/gsh-lan/steam-gameserver-token-api/src/steam"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"
)

const headerAuthorization = "authorization"

var log *zap.SugaredLogger

func init() {
	log = logger.GetSugaredLogger()
}

// App contains references to global necessities
type App struct {
	Router    *mux.Router
	authToken string
	steam     *steam.Steam
	cache     map[string]string
	cacheLock sync.RWMutex
}

// Run server on specific interface
func (a *App) Run(addr, apiKey, authToken string, backgroundProcessingInterval time.Duration) {
	a.cache = make(map[string]string)
	a.authToken = authToken
	a.steam = steam.New(apiKey)
	a.cacheLock = sync.RWMutex{}

	a.registerRoutes()

	go a.backgroundProcessor(backgroundProcessingInterval)

	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) registerRoutes() {
	a.Router = mux.NewRouter().StrictSlash(true)
	a.Router.HandleFunc("/token/{appID}/{memo}", a.pullToken).Methods("GET")
}

// RespondWithJSON uses a struct, for a JSON response.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := jsoniter.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// RespondWithText returns text/plain.
func respondWithText(w http.ResponseWriter, code int, payload string) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	w.Write([]byte(payload))
}

// RespondWithError standardizes error messages, through the use of RespondWithJSON.
func respondWithError(w http.ResponseWriter, code int, message string) {
	log.Info(message)
	respondWithJSON(w, code, map[string]string{"error": message})
}

func (a *App) pullToken(w http.ResponseWriter, r *http.Request) {
	if a.authToken != "" {
		authHeader := r.Header.Get(headerAuthorization)
		if authHeader != fmt.Sprintf("Bearer %s", a.authToken) {
			log.Infof("Invalid auth token provided: %s", authHeader)
			respondWithError(w, http.StatusForbidden, "Invalid authorization header")
			return
		}
	}

	vars := mux.Vars(r)
	if _, ok := vars["appID"]; !ok {
		respondWithError(w, http.StatusBadRequest, "Missing appID")
		return
	}
	if _, ok := vars["memo"]; !ok {
		respondWithError(w, http.StatusBadRequest, "Missing memo")
		return
	}

	appID, err := strconv.Atoi(vars["appID"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("bad appID: %s", err))
		return
	}

	// check if token is cached
	a.cacheLock.RLock()
	cachedToken, ok := a.cache[buildCacheKey(appID, vars["memo"])]
	a.cacheLock.RUnlock()
	if ok {
		log.Infof("Successfully fetched cached logintoken for appid %d with memo %s", appID, vars["memo"])
		respondWithText(w, http.StatusOK, cachedToken)
		return
	}

	accounts, err := a.steam.GetAccountList()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Unable to list existing tokens: %s", err))
		return
	}

	// Check for existing account
	var account steam.Account
	for _, acct := range accounts {
		if acct.Memo == vars["memo"] && int(acct.AppID) == appID {
			account = acct
			break
		}
	}

	// Create new if not found
	if account.SteamID == "" {
		account, err = a.steam.CreateAccount(appID, vars["memo"])
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Refresh token if found and expired
	if account.IsExpired == true {
		account, err = a.steam.ResetLoginToken(account.SteamID)
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	a.cacheToken(appID, vars["memo"], account.LoginToken)

	log.Infof("Successfully fetched logintoken for appid %d with memo %s", appID, vars["memo"])
	respondWithText(w, http.StatusOK, account.LoginToken)
}

func (a *App) backgroundProcessor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			func() {
				accounts, err := a.steam.GetAccountList()
				if err != nil {
					log.Error(err)
					return
				}

				for _, acct := range accounts {
					if acct.IsExpired {
						_, err := a.steam.ResetLoginToken(acct.SteamID)
						if err != nil {
							log.Error(err)
						} else {
							a.cacheToken(int(acct.AppID), acct.Memo, acct.LoginToken)
						}
					} else {
						a.cacheToken(int(acct.AppID), acct.Memo, acct.LoginToken)
					}
				}
			}()
		}
	}
}

// cacheToken caches a new Token for a specific appID and memo
func (a *App) cacheToken(appID int, memo, token string) {
	a.cacheLock.Lock()
	defer a.cacheLock.Unlock()
	a.cache[buildCacheKey(appID, memo)] = token
}

func buildCacheKey(appID int, memo string) string {
	return fmt.Sprintf("%d-%s", appID, memo)
}
