// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package user

import (
	"errors"
)

func (x *GetDesignateUsersReq) Check() error {
	if x.UserIDs == nil {
		return errors.New("UserIDs is empty")
	}
	return nil
}

func (x *UserRegisterReq) Check() error {
	if x.Users == nil {
		return errors.New("users are empty")
	}
	for _, u := range x.Users {
		if u.Nickname == "" {
			return errors.New("nickname is empty")
		}
	}

	return nil
}
