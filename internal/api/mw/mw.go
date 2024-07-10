package mw

import (
	"github.com/gin-gonic/gin"
	cmConstant "github.com/openimsdk/openmeeting-server/pkg/common/constant"
	"github.com/openimsdk/openmeeting-server/pkg/common/servererrs"
	"github.com/openimsdk/openmeeting-server/pkg/common/token"
	"github.com/openimsdk/openmeeting-server/pkg/rpcclient"
	"github.com/openimsdk/protocol/constant"
	pbuser "github.com/openimsdk/protocol/openmeeting/user"
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
		apiresp.GinError(c, errs.WrapMsg(err, "parse token failed, invalid token"))
		return
	}
	if err := o.isValidToken(c, userID, userToken); err != nil {
		c.Abort()
		apiresp.GinError(c, errs.WrapMsg(err, "not valid token"))
		return
	}
	o.setToken(c, userID)
}

func (o *MW) isValidToken(c *gin.Context, userID, userToken string) error {
	resp, err := o.client.Client.GetUserToken(c, &pbuser.GetUserTokenReq{UserID: userID})
	if err != nil {
		return err
	}
	if resp.Token == cmConstant.KickOffMeetingMsg {
		return servererrs.ErrKickOffMeeting.WrapMsg("kick off meeting, please login again")
	}
	if resp.Token == "" || resp.Token != userToken {
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
