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
	FormattingError      = 10001 // Error in formatting
	HasRegistered        = 10002 // user has already registered
	NotRegistered        = 10003 // user is not registered
	PasswordErr          = 10004 // Password error
	GetIMTokenErr        = 10005 // Error in getting IM token
	RepeatSendCode       = 10006 // Repeat sending code
	MailSendCodeErr      = 10007 // Error in sending code via email
	SmsSendCodeErr       = 10008 // Error in sending code via SMS
	CodeInvalidOrExpired = 10009 // Code is invalid or expired
	RegisterFailed       = 10010 // Registration failed
	ResetPasswordFailed  = 10011 // Resetting password failed
	RegisterLimit        = 10012 // Registration limit exceeded
	LoginLimit           = 10013 // Login limit exceeded
	InvitationError      = 10014 // Error in invitation
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
