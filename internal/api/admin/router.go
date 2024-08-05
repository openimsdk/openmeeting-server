package admin

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/cache/redis"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/controller"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/database/mgo"
	"github.com/openimsdk/openmeeting-server/pkg/common/token"
	"github.com/openimsdk/openmeeting-server/pkg/rpcclient"
	userfind "github.com/openimsdk/openmeeting-server/pkg/user"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/mw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Whitelist api not parse token
var whitelist = []string{
	"/admin/login",
	"/admin/user/register",
}

func secretKey(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}
}

func newAdminGinRouter(ctx context.Context, disCov discovery.SvcDiscoveryRegistry, config *Config) *gin.Engine {
	disCov.AddOption(mw.GrpcClient(), grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery(), mw.CorsHandler(), mw.GinParseOperationID(), mw.GinParseToken(secretKey(config.AdminAPI.Secret), whitelist))

	// init storage
	mgoCli, err := mongoutil.NewMongoDB(ctx, config.Mongo.Build())
	if err != nil {
		return nil
	}
	rdb, err := redisutil.NewRedisClient(ctx, config.Redis.Build())
	if err != nil {
		return nil
	}

	userDB, err := mgo.NewUserMongo(mgoCli.GetDB())
	if err != nil {
		return nil
	}
	userCache := redis.NewUser(rdb, userDB, redis.GetDefaultOpt())
	database := controller.NewUser(userDB, userCache, mgoCli.GetTx())

	user := userfind.NewMeeting(disCov, config.Share.RpcRegisterName.User)
	// init rpc client here
	userRpc := rpcclient.NewUser(user)
	userToken := token.New(config.AdminAPI.Expire, config.AdminAPI.Secret)
	u := NewAdminApi(database, *userRpc, userToken)
	adminRouterGroup := r.Group("/admin")
	{
		adminRouterGroup.POST("/login", u.AdminLogin)
		adminRouterGroup.POST("/user/register", u.RegisterUser)
		adminRouterGroup.POST("/user/import/json", u.ImportUserByJson)
		adminRouterGroup.POST("/user/import/xlsx", u.ImportUserByXlsx)
	}
	return r
}
