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

package api

import (
	"github.com/dalyisaac/mstodo/auth"
	"github.com/go-resty/resty/v2"
)

func GetClient() *resty.Client {
	client := resty.New()
	client.EnableTrace()
	client.SetHostURL("https://graph.microsoft.com/v1.0")

	return client
}

func CreateRequest() (*resty.Request, error) {
	client := GetClient()

	token, err := auth.GetToken()
	if err != nil {
		return nil, err
	}

	request := client.R().SetAuthToken(token.AccessToken)
	return request, nil
}
