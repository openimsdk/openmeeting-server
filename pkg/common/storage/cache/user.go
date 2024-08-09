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

package cache

import (
	"context"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
)

type User interface {
	Meta
	NewCache() User
	GetUsersInfo(ctx context.Context, userIDs []string) ([]*model.User, error)
	DelUsersInfo(userIDs ...string) User
	GetUserByAccount(ctx context.Context, account string) (*model.User, error)
	CacheUserToken(ctx context.Context, userID, userToken string) error
	GetUserToken(ctx context.Context, userID string) (string, error)
	ClearUserToken(ctx context.Context, userID string) error
	GenerateUserID(ctx context.Context) (string, error)
}
