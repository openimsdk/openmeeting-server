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

package servererrs

import "github.com/openimsdk/tools/errs"

var (
	ErrDatabase = errs.NewCodeError(DatabaseError, "DatabaseError")
	ErrNetwork  = errs.NewCodeError(NetworkError, "NetworkError")

	ErrInternalServer         = errs.NewCodeError(ServerInternalError, "ServerInternalError")
	ErrArgs                   = errs.NewCodeError(ArgsError, "ArgsError")
	ErrUserIDNotFound         = errs.NewCodeError(UserIDNotFoundError, "UserIDNotFoundError")
	ErrRegisteredAlready      = errs.NewCodeError(HasRegistered, "RegisteredAlreadyError")
	ErrUserPasswordError      = errs.NewCodeError(PasswordErr, "PasswordErr")
	ErrUserAccountNotFoundErr = errs.NewCodeError(NotFoundAccountErr, "NotFoundAccountErr")
	ErrUserTokenNotFoundErr   = errs.NewCodeError(NotFoundUserTokenErr, "NotFoundUserTokenErr")
	ErrKickOffMeeting         = errs.NewCodeError(KickOffMeetingError, "KickOffMeetingError")

	ErrMeetingUserLimit        = errs.NewCodeError(MeetingUserLimitError, "MeetingUserLimitError")
	ErrMeetingPasswordNotMatch = errs.NewCodeError(MeetingPasswordError, "MeetingPasswordError")
	ErrMeetingAuthCheck        = errs.NewCodeError(MeetingAuthCheckError, "MeetingAuthCheckError")
	ErrMeetingAlreadyCompleted = errs.NewCodeError(MeetingCompleteError, "MeetingCompleteError")
)
