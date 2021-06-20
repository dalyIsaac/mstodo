/*
Copyright © 2021 Isaac Daly <isaac.daly@outlook.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

// Based on https://gist.github.com/guumaster/2c7f48ac3567ae6c456f4020c857c375

package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var (
	errTokenNotFound = errors.New("token not found")
	errTokenWeb      = errors.New("error getting token from web")
	errTokenOpen     = errors.New("error opening token file")
	errTokenSave     = errors.New("error saving token file")
)

type TokenManager struct {
	conf                *oauth2.Config
	token               *oauth2.Token
	originalAccessToken string
	filepath            string
}

func GetToken() (*oauth2.Token, error) {
	tm, err := GetTokenManager()
	if err != nil {
		return nil, err
	}

	return tm.token, nil
}

func GetTokenManager() (*TokenManager, error) {
	tm := &TokenManager{
		conf: &oauth2.Config{
			ClientID:     viper.GetString("client-id"),
			ClientSecret: viper.GetString("client-secret"),
			Scopes:       getScopes(),
			Endpoint:     microsoft.AzureADEndpoint(""),
		},
		filepath: path.Join(viper.GetString("config-dir"), "token.json"),
	}

	isNewToken := false
	token, err := tm.getFromFile()
	if errors.Is(err, errTokenOpen) || errors.Is(err, errTokenNotFound) {
		token, err = tm.getFromWeb()
		if err != nil {
			return nil, fmt.Errorf("error getting token from web: %w", err)
		}
		isNewToken = true
	}
	if err != nil {
		return nil, fmt.Errorf("error getting token: %w", err)
	}

	tm.token = token
	tm.originalAccessToken = token.AccessToken

	// This will refresh the token when needed
	ts := tm.TokenSource(context.Background())
	newTok, err := ts.Token()
	if err != nil {
		return nil, err
	}

	tokenRefreshed := token.AccessToken != newTok.AccessToken

	if isNewToken || tokenRefreshed {
		tm.token = newTok
		err = tm.save()
		if err != nil {
			return nil, err
		}
	}

	return tm, nil
}

func (t *TokenManager) TokenSource(ctx context.Context) oauth2.TokenSource {
	return t.conf.TokenSource(ctx, t.token)
}

// save stores the token in json file.
func (t *TokenManager) save() error {
	fmt.Printf("Saving token to file: %s\n", t.filepath)
	f, err := os.OpenFile(t.filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("%q: %w", err, errTokenSave)
	}

	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	err = enc.Encode(t.token)
	if err != nil {
		return fmt.Errorf("%q: %w", err, errTokenSave)
	}

	return nil
}

// getFromFile retrieves a token from a local file.
func (t *TokenManager) getFromFile() (*oauth2.Token, error) {
	f, err := os.Open(t.filepath)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", err, errTokenOpen)
	}
	defer f.Close()
	tok := new(oauth2.Token)
	err = json.NewDecoder(f).Decode(tok)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", err, errTokenOpen)
	}

	t.originalAccessToken = tok.AccessToken

	return tok, err
}

// getFromWeb Starts a local server and the oauth flow
func (t *TokenManager) getFromWeb() (*oauth2.Token, error) {
	client, err := authenticateUser(t.conf, viper.GetInt("port"), viper.GetInt("auth-timeout"))
	if err != nil {
		return nil, fmt.Errorf("%q: %w", err, errTokenWeb)
	}
	return client.Token, nil
}
