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

package auth

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/color"
	rndm "github.com/nmrshll/rndm-go"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
)

type AuthorizedClient struct {
	*http.Client
	Token *oauth2.Token
}

type oauthStateStringContextKeyType int

const (
	oauthStateStringContextKey oauthStateStringContextKeyType = 987
)

// authenticateUser starts the login process
func authenticateUser(config *oauth2.Config, port int, authTimeout int) (*AuthorizedClient, error) {
	// validate config
	if config == nil {
		return nil, errors.New("OAuth2 config was unexpectedly nil")
	}

	// http transport for self-signed certificate, to be added to the context
	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sslClient := &http.Client{Transport: httpTransport}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, sslClient)

	// redirect the user to the consent page to ask permission for the scopes
	config.RedirectURL = fmt.Sprintf("http://localhost:%s/oauth/callback", strconv.Itoa(port))

	// some random string, used for getting the AuthCodeURL
	oauthStateString := rndm.String(8)
	ctx = context.WithValue(ctx, oauthStateStringContextKey, oauthStateString)
	url := config.AuthCodeURL(oauthStateString, oauth2.AccessTypeOffline)

	clientChan, stopHTTPServerChan, cancelAuthentication := startHTTPServer(ctx, config, port)
	log.Println(color.CyanString("You will now be taken to your browser for authentication or open the url below in a browser:"))
	log.Println(color.CyanString(url))

	err := open.Run(url)
	if err != nil {
		log.Println("Failed to open URL")
		return nil, err
	}

	// shutdown the server after timeout
	go func() {
		log.Printf("Authentication will be cancelled in %s seconds", strconv.Itoa(authTimeout))
		time.Sleep(time.Duration(authTimeout) * time.Second)
		stopHTTPServerChan <- struct{}{}
	}()

	select {
	// wait for client on clientChan
	case client := <-clientChan:
		// after the callbackHandler returns a client, shutdown the server gracefully
		stopHTTPServerChan <- struct{}{}
		return client, nil

	case <-cancelAuthentication:
		// if authentication process is cancelled first, return an error
		return nil, errors.New("authentication timed out and was cancelled")
	}
}

func startHTTPServer(ctx context.Context, conf *oauth2.Config, port int) (clientChan chan *AuthorizedClient, stopHTTPServerChan chan struct{}, cancelAuthentication chan struct{}) {
	// init returns
	clientChan = make(chan *AuthorizedClient)
	stopHTTPServerChan = make(chan struct{})
	cancelAuthentication = make(chan struct{})

	http.HandleFunc("/oauth/callback", callbackHandler(ctx, conf, clientChan))
	srv := &http.Server{Addr: ":" + strconv.Itoa(port)}

	// handle server shutdown signal
	go func() {
		// wait for signal on stopHTTPServerChan
		<-stopHTTPServerChan
		log.Println("Shutting down server...")

		// give it 5 sec to shutdown gracefully, else quit program
		d := time.Now().Add(5 * time.Second)
		ctx, cancel := context.WithDeadline(context.Background(), d)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			log.Printf(color.RedString("Auth server could not shutdown gracefully: %v"), err)
		}

		// after server is shutdown, quit program
		cancelAuthentication <- struct{}{}
	}()

	// handle callback request
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
		fmt.Println("Server gracefully stopped")
	}()

	return clientChan, stopHTTPServerChan, cancelAuthentication
}

const success = `
<!DOCTYPE html>
<html lang="en">

    <head>
        <meta charset="utf-8">
        <title>Todo CLI Authentication Success</title>
        <style>
            .wrapper {
                height: 100%;
                width: 100%;
                display: flex;
                flex-direction: column;
                align-items: center;
                font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
                background-color: #2ecc71;
                padding-bottom: 10px;
            }
        </style>
    </head>

    <body>
        <div class="wrapper">
            <h1>Success</h1>
            <p>You can now close this window and return to the application.</p>
        </div>
    </body>

</html>
`

func callbackHandler(ctx context.Context, config *oauth2.Config, clientChan chan *AuthorizedClient) func(w http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		requestStateString := ctx.Value(oauthStateStringContextKey).(string)
		responseStateString := r.FormValue("state")

		if responseStateString != requestStateString {
			fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", requestStateString, responseStateString)
			http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
			return
		}

		fmt.Println(r.Form.Get("error_description"))
		code := r.FormValue("code")
		token, err := config.Exchange(ctx, code)
		if err != nil {
			fmt.Printf("oauthConfig.Exchange() failed with error '%s'\n", err)
			http.Redirect(rw, r, "/", http.StatusTemporaryRedirect)
			return
		}

		// The HTTP Client returned by oauthConfig.Client will refresh the token as necessary
		client := &AuthorizedClient{
			config.Client(ctx, token),
			token,
		}

		rw.Write([]byte(success))
		clientChan <- client
	}
}
