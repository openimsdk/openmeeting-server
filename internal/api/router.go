package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	apiMw "github.com/openimsdk/openmeeting-server/internal/api/mw"
	"github.com/openimsdk/openmeeting-server/pkg/common/token"
	"github.com/openimsdk/openmeeting-server/pkg/rpcclient"
	"github.com/openimsdk/openmeeting-server/pkg/user"
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
	userRpc := user.NewMeetingUserClient(disCov, config.Share.RpcRegisterName.User)
	meetingRpc := rpcclient.NewMeeting(disCov, config.Share.RpcRegisterName.Meeting)

	userToken := token.New(config.API.Expire, config.API.Secret)
	mwApi := apiMw.New(userRpc, userToken)
	u := NewUserApi(userRpc)
	userRouterGroup := r.Group("/user")
	{
		userRouterGroup.POST("/register", u.UserRegister)
		userRouterGroup.POST("/login", u.UserLogin)
		userRouterGroup.POST("/get_users_info", mwApi.CheckToken, u.GetUsersPublicInfo)
		userRouterGroup.POST("/update_user_password", mwApi.CheckToken, u.UpdateUserPassword)
		userRouterGroup.POST("/logout", mwApi.CheckToken, u.UserLogout)

	}

	m := NewMeetingApi(*meetingRpc)
	meetingRouterGroup := r.Group("/meeting")
	{
		meetingRouterGroup.POST("/book_meeting", mwApi.CheckToken, m.BookMeeting)
		meetingRouterGroup.POST("/create_immediate_meeting", mwApi.CheckToken, m.CreateImmediateMeeting)
		meetingRouterGroup.POST("/join_meeting", mwApi.CheckToken, m.JoinMeeting)
		meetingRouterGroup.POST("/get_meeting_token", mwApi.CheckToken, m.GetMeetingToken)

		meetingRouterGroup.POST("/update_meeting", mwApi.CheckToken, m.UpdateMeeting)
		meetingRouterGroup.POST("/get_meeting", mwApi.CheckToken, m.GetMeeting)
		meetingRouterGroup.POST("/get_meetings", mwApi.CheckToken, m.GetMeetings)
		meetingRouterGroup.POST("/leave_meeting", mwApi.CheckToken, m.LeaveMeeting)
		meetingRouterGroup.POST("/end_meeting", mwApi.CheckToken, m.EndMeeting)
		meetingRouterGroup.POST("/set_personal_setting", mwApi.CheckToken, m.SetPersonalMeetingSettings)
		meetingRouterGroup.POST("/get_personal_setting", mwApi.CheckToken, m.GetPersonalMeetingSettings)
		meetingRouterGroup.POST("/operate_meeting_all_stream", mwApi.CheckToken, m.OperateMeetingAllStream)
		meetingRouterGroup.POST("/modify_meeting_participant_name", mwApi.CheckToken, m.ModifyMeetingParticipantNickName)
		meetingRouterGroup.POST("/remove_participants", mwApi.CheckToken, m.RemoveMeetingParticipants)
		meetingRouterGroup.POST("/set_meeting_host_info", mwApi.CheckToken, m.SetMeetingHostInfo)
		meetingRouterGroup.POST("/toggle_record_meeting", mwApi.CheckToken, m.ToggleRecordMeeting)
	}
	return r
}
