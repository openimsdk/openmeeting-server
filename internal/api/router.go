package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/openimsdk/openmeeting-server/pkg/rpcclient"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/mw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Whitelist api not parse token
var whitelist = []string{
	"",
}

func secretKey(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}
}

func newGinRouter(disCov discovery.SvcDiscoveryRegistry, config *Config) *gin.Engine {
	disCov.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), mw.CorsHandler(), mw.GinParseOperationID(), mw.GinParseToken(secretKey(config.API.Secret), whitelist))
	// init rpc client here
	userRpc := rpcclient.NewUser(disCov, config.Share.RpcRegisterName.User)
	meetingRpc := rpcclient.NewMeeting(disCov, config.Share.RpcRegisterName.Meeting)

	u := NewUserApi(*userRpc)
	userRouterGroup := r.Group("/user")
	{
		userRouterGroup.POST("/user_register", u.UserRegister)
		userRouterGroup.POST("/get_users_info", u.GetUsersPublicInfo)
	}

	m := NewMeetingApi(*meetingRpc)
	meetingRouterGroup := r.Group("/meeting")
	{
		meetingRouterGroup.POST("/book_meeting", m.BookMeeting)
	}
	return r
}
