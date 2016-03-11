/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This file is part of Nanocloud community.
 *
 * Nanocloud community is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Nanocloud community is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package middlewares

import (
	"net/http"

	"github.com/Nanocloud/community/nanocloud/models/users"
	"gopkg.in/labstack/echo.v1"
)

type hash map[string]interface{}

func admin(c *echo.Context, handler echo.HandlerFunc) error {
	user := c.Get("user").(*users.User)

	if !user.IsAdmin {
		return c.JSON(http.StatusForbidden, hash{
			"error": "forbidden",
		})
	}
	return handler(c)
}

func Admin(handler echo.HandlerFunc) echo.HandlerFunc {
	return func(c *echo.Context) error {
		return admin(c, handler)
	}
}
