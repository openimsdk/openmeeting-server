package admin

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/openmeeting-server/pkg/apistruct"
	"github.com/openimsdk/openmeeting-server/pkg/common/securetools"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/controller"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/model"
	"github.com/openimsdk/openmeeting-server/pkg/common/token"
	"github.com/openimsdk/openmeeting-server/pkg/common/xlsx"
	"github.com/openimsdk/openmeeting-server/pkg/common/xlsx/definition"
	"github.com/openimsdk/openmeeting-server/pkg/protocol/admin"
	"github.com/openimsdk/openmeeting-server/pkg/rpcclient"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
	"github.com/xuri/excelize/v2"
)

type ApiAdmin struct {
	client             rpcclient.User
	userStorageHandler controller.User
	config             *Config
	tokenVerify        *token.Token
}

func NewAdminApi(userStorage controller.User, client rpcclient.User, t *token.Token) *ApiAdmin {
	return &ApiAdmin{
		client:             client,
		userStorageHandler: userStorage,
		tokenVerify:        t,
	}
}

func (a *ApiAdmin) AdminLogin(c *gin.Context) {
	req, err := a2r.ParseRequest[admin.UserLoginReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	user, err := a.userStorageHandler.GetByAccount(c, req.Account)
	if err != nil {
		apiresp.GinError(c, errs.WrapMsg(err, "login failed, not found account, please check"))
		return
	}
	saltPasswd := securetools.VerifyPassword(req.Password, user.SaltValue)
	if saltPasswd != user.Password {
		apiresp.GinError(c, errs.WrapMsg(err, "wrong password or user account"))
		return
	}
	userToken, err := a.tokenVerify.CreateToken(user.UserID)
	if err != nil {
		apiresp.GinError(c, errs.WrapMsg(err, "create token failed, please check"))
		return
	}
	if err := a.userStorageHandler.StoreToken(c, user.UserID, userToken); err != nil {
		apiresp.GinError(c, errs.WrapMsg(err, "set token failed, please check"))
		return
	}
	apiresp.GinSuccess(c, &apistruct.AdminLoginResp{
		AdminAccount: user.Account,
		AdminToken:   userToken,
		Nickname:     user.Nickname,
	})
}

func (a *ApiAdmin) ImportUserByJson(c *gin.Context) {
	formFile, err := c.FormFile("data")
	if err != nil {
		apiresp.GinError(c, errs.WrapMsg(err, "get form file failed"))
		return
	}

	// open file
	f, err := formFile.Open()
	if err != nil {
		apiresp.GinError(c, errs.WrapMsg(err, "open file failed"))
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
		}
	}()

	var users []definition.User
	if err := json.NewDecoder(f).Decode(&users); err != nil {
		apiresp.GinError(c, errs.WrapMsg(err, "decode json failed"))
		return
	}

	dbUsers := make([]*model.User, 0, len(users))
	for _, user := range users {
		passwd, salt := securetools.HashPassword(user.Password)
		dbUsers = append(dbUsers, &model.User{
			UserID:    user.UserID,
			Nickname:  user.Nickname,
			Account:   user.Account,
			Password:  passwd,
			SaltValue: salt,
		})
	}
	if len(dbUsers) > 0 {
		if err := a.userStorageHandler.Create(c, dbUsers); err != nil {
			apiresp.GinError(c, errs.ErrArgs.WrapMsg("create users failed "+err.Error()))
			return
		}
	}
	apiresp.GinSuccess(c, nil)
}

func (a *ApiAdmin) ImportUserByXlsx(c *gin.Context) {
	formFile, _, err := c.Request.FormFile("data")
	if err != nil {
		apiresp.GinError(c, errs.WrapMsg(err, "get form file failed"))
		return
	}

	file, err := excelize.OpenReader(formFile)
	if err != nil {
		apiresp.GinError(c, errs.WrapMsg(err, "open failed"))
		return
	}

	// read sheet
	rows, err := file.GetRows("Sheet1")
	if err != nil {
		apiresp.GinError(c, errs.WrapMsg(err, "get rows failed"))
		return
	}

	// get header
	headers := rows[0]
	colIndex := xlsx.GetColumnIndex(headers)

	// get data
	var users []definition.User
	for _, row := range rows[1:] { // skip header
		user := definition.User{}
		if err := xlsx.SetStructValues(&user, row, colIndex); err != nil {
			apiresp.GinError(c, errs.WrapMsg(err, "please check xlsx structure"))
			return
		}
		users = append(users, user)
	}
	dbUsers := make([]*model.User, 0, len(users))
	for _, user := range users {
		passwd, salt := securetools.HashPassword(user.Password)
		dbUsers = append(dbUsers, &model.User{
			UserID:    user.UserID,
			Nickname:  user.Nickname,
			Account:   user.Account,
			Password:  passwd,
			SaltValue: salt,
		})
	}
	if err := a.userStorageHandler.Create(c, dbUsers); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("create users failed "+err.Error()))
		return
	}
	apiresp.GinSuccess(c, nil)
}

func (a *ApiAdmin) RegisterUser(c *gin.Context) {
	req, err := a2r.ParseRequest[admin.UserRegisterReq](c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	passwd, salt := securetools.HashPassword(req.Password)
	dbUser := &model.User{
		UserID:    req.UserID,
		Nickname:  req.Nickname,
		Account:   req.Account,
		Password:  passwd,
		SaltValue: salt,
	}
	if err := a.userStorageHandler.Create(c, []*model.User{dbUser}); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WrapMsg("create users failed "+err.Error()))
		return
	}
	apiresp.GinSuccess(c, nil)
}
