package meeting

import (
	"context"
	"github.com/openimsdk/openmeeting-server/pkg/common/config"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/cache/redis"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/controller"
	"github.com/openimsdk/openmeeting-server/pkg/common/storage/database/mgo"
	"github.com/openimsdk/openmeeting-server/pkg/rpcclient"
	"github.com/openimsdk/openmeeting-server/pkg/rtc"
	"github.com/openimsdk/openmeeting-server/pkg/rtc/livekit"
	userfind "github.com/openimsdk/openmeeting-server/pkg/user"
	pbmeeting "github.com/openimsdk/protocol/openmeeting/meeting"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"
	registry "github.com/openimsdk/tools/discovery"
	"google.golang.org/grpc"
)

type meetingServer struct {
	meetingStorageHandler controller.Meeting
	RegisterCenter        registry.SvcDiscoveryRegistry
	meetingRtc            rtc.MeetingRtc
	config                *Config
	userRpc               *rpcclient.User
}

type Config struct {
	Rpc       config.Meeting
	Redis     config.Redis
	Mongo     config.Mongo
	Discovery config.Discovery
	Share     config.Share
	Rtc       config.RTC
}

func Start(ctx context.Context, config *Config, client registry.SvcDiscoveryRegistry, server *grpc.Server) error {
	mgoCli, err := mongoutil.NewMongoDB(ctx, config.Mongo.Build())
	if err != nil {
		return err
	}
	rdb, err := redisutil.NewRedisClient(ctx, config.Redis.Build())
	if err != nil {
		return err
	}

	meetingDB, err := mgo.NewMeetingMongo(mgoCli.GetDB())
	if err != nil {
		return err
	}
	meetingCache := redis.NewMeeting(rdb, meetingDB, redis.GetDefaultOpt())
	database := controller.NewMeeting(meetingDB, meetingCache, mgoCli.GetTx())
	meetingRtc := livekit.NewLiveKit(&config.Rtc)

	user := userfind.NewMeeting(client, config.Share.RpcRegisterName.User)

	// init rpc client here
	userRpc := rpcclient.NewUser(user)

	u := &meetingServer{
		meetingStorageHandler: database,
		RegisterCenter:        client,
		config:                config,
		meetingRtc:            meetingRtc,
		userRpc:               userRpc,
	}
	pbmeeting.RegisterMeetingServiceServer(server, u)
	return nil
}
