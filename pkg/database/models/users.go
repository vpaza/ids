/*
 * Copyright Daniel Hawton
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package models

import "github.com/adh-partnership/api/pkg/database"

type User struct {
	ID        int64  `json:"id"`
	CID       uint   `json:"cid" gorm:"column:cid"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Roles     string `json:"roles"`
}

func FindUser(cid interface{}) (*User, error) {
	var user User
	if err := database.DB.Where("cid = ?", cid).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
