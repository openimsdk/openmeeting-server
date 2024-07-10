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

// UnknownCode represents the error code when code is not parsed or parsed code equals 0.
const UnknownCode = 1000

// Error codes for various error scenarios.
const (
	HasRegistered        = 10001 // user has already registered
	PasswordErr          = 10002 // Password error
	NotFoundAccountErr   = 10003 // not found user account
	NotFoundUserTokenErr = 10004 // not found user token
	KickOffMeetingError  = 10010

	MeetingUserLimitError = 20001 // one user joins more than one meeting
	MeetingPasswordError  = 20002 // password not match error
	MeetingAuthCheckError = 20003 // meeting auth check permission error
	MeetingCompleteError  = 20004 // meeting update check error
)

// General error codes.
const (
	NoError       = 0     // No error
	DatabaseError = 90002 // Database error (redis/mysql, etc.)
	NetworkError  = 90004 // Network error
	DataError     = 90007 // Data error

	// General error codes.
	ServerInternalError = 500  // Server internal error
	ArgsError           = 1001 // Input parameter error

	// Account error codes.
	UserIDNotFoundError    = 1101 // UserID does not exist or is not registered
	RegisteredAlreadyError = 1102 // user is already registered

	// Group error codes.
	GroupIDNotFoundError = 1201 // GroupID does not exist
	GroupIDExisted       = 1202 // GroupID already exists

)
