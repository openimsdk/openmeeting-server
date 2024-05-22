package mw

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/openmeeting-server/pkg/common/token"
	pbuser "github.com/openimsdk/openmeeting-server/pkg/protocol/user"
	"github.com/openimsdk/openmeeting-server/pkg/rpcclient"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
)

type MW struct {
	client      *rpcclient.User
	tokenVerify *token.Token
}

func New(c *rpcclient.User, t *token.Token) *MW {
	return &MW{
		client:      c,
		tokenVerify: t,
	}
}

func (o *MW) parseToken(c *gin.Context) (string, string, error) {
	userToken := c.GetHeader("token")
	if userToken == "" {
		return "", "", errs.ErrArgs.WrapMsg("token is empty")
	}
	userID, err := o.tokenVerify.GetToken(userToken)
	if err != nil {
		return "", "", err
	}
	return userID, userToken, nil
}

func (o *MW) CheckToken(c *gin.Context) {
	userID, userToken, err := o.parseToken(c)
	if err != nil {
		c.Abort()
		apiresp.GinError(c, err)
		return
	}
	if err := o.isValidToken(c, userID, userToken); err != nil {
		c.Abort()
		apiresp.GinError(c, err)
		return
	}
	o.setToken(c, userID)
}

func (o *MW) isValidToken(c *gin.Context, userID, userToken string) error {
	resp, err := o.client.Client.GetUserToken(c, &pbuser.GetUserTokenReq{UserID: userID})
	if err != nil {
		return err
	}
	if resp.Token == "" {
		return errs.ErrTokenExpired.Wrap()
	}
	return nil
}

func (o *MW) setToken(c *gin.Context, userID string) {
	SetToken(c, userID)
}

func SetToken(c *gin.Context, userID string) {
	c.Set(constant.OpUserID, userID)
}
