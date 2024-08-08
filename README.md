<p align="center">
    <a href="https://openim.io">
        <img src="./assets/openim-logo.gif" width="60%" height="30%"/>
    </a>
</p>

<div align="center">


</div>


## :busts_in_silhouette: Join our community
## Ⓜ️ About OpenMeeting
OpenMeeting is an open-source real-time audio and video conferencing system developed using Golang. OpenMeeting provides user management, meeting management, audio and video transmission, instant meetings, scheduled meetings, screen sharing, and more, aiming to offer users a convenient remote meeting experience. It is similar to Zoom and Tencent Meeting, with support for private deployment to ensure the security and privacy of enterprise and individual user data.

[//]: # (![Relations of App-OpenMeeting]&#40;./assets/open-meeting-design.png&#41;)



## 🌐 Introduction to OpenMeetingServer
+ **OpenMeetingServer**  include:
  - Instant Meetings: Users can create instant meetings at any time, invite others to join, and engage in efficient remote communication.
  - Scheduled Meetings: Supports scheduling future meetings, setting meeting times and participants, and the system will remind users before the meeting starts.
  - Screen Sharing: Users can share their screens during meetings for demonstration and collaboration.
  - High-Quality Audio and Video: Provides high-quality audio and video transmission to ensure smooth meetings.
  - Multi-Platform Support: Supports various operating systems, including Windows, macOS, Linux, and more.
  - Microservices Architecture: Supports cluster mode, including a gateway and multiple RPC services.
  - Multiple Deployment Methods: Supports source code, Kubernetes, or Docker deployment.


### Enhanced Business Features:
+ **REST API**：Provides REST API for business systems, offering client interfaces.
+ **RPC API**： Provides corresponding services through gRPC, including user and meeting, to extend more business forms.

[//]: # (![architecture]&#40;./assets/architecture-layers.png&#41;)



## :rocket: Quick Start
To facilitate user experience, we provide multiple deployment solutions. You can choose the suitable deployment method from the list below:

[//]: # (+ **[Source Code Deployment Guide]&#40;https://github.com/openimsdk/openmeeting-server/blob/main/deployments/deployment.md&#41;**)
### OpenMeeting Server Source Code Deployment

#### 1. Downloading the Source Code

```bash
git clone https://github.com/openimsdk/openmeeting-server.git && cd openmeeting-server
```


#### 2. Deploying Related Dependencies Component (Etcd, MongoDB, Redis, LiveKit)
```bash
# install dependencies component
docker compose up -d

# Checking if Related Dependencies components are Running Properly
docker ps
```

#### 3. Set external IP
```bash
Modify the `url` in `config/live.yml` to `ws://external_IP:17880` or a domain name.
```

#### 4. Initialization
Before the first compilation, execute on Linux/Mac platforms:
```bash
bash bootstrap.sh
```
On Windows execute:
```bash
bootstrap.bat
```

#### 5. compile and run
```bash
mage && mage start
```



+ **[Docker Deployment Guide]()**

### How to add user in meeting server
+ Replace your_ip_or_domain, your_userID, your_password, your_account, and your_nickname with the appropriate values. Then, run the command in your Bash terminal.
```bash
curl -X POST "http://your_ip_or_domain:11022/admin/user/register" \
-H "Content-Type: application/json" \
-H "operationID: 123456789" \
-d '{
  "userID": "your_userID",
  "password": "your_password",
  "account": "your_account",
  "nickname": "your_nickname"
}'
```

+ Then you can use this account's account and password to log in to the client.




## System Support
Supports Linux, Windows, Mac systems, as well as ARM and AMD CPU architectures.

## :link: Related Links


+ **[Developer Manual]()**
+ **[Changelog]()**

## :writing_hand: How to Contribute
We welcome contributions of any kind! Before submitting a Pull Request, please ensure you have read our Contributor Documentation.

+ **[Report a Bug](https://github.com/openimsdk/openmeeting-server/issues/new?assignees=&labels=kind%2Fbug&projects=&template=bug-report.yaml&title=%5BBUG%5D+)**
+ **[Request a Feature](https://github.com/openimsdk/openmeeting-server/issues/new?assignees=&labels=feature+request&projects=&template=feature-request.yaml&title=%5BFEATURE+REQUEST%5D+)**
+ **[Submit a Pull Request](https://github.com/openimsdk/openmeeting-server/pulls)**

Thank you for your contributions, let's build a powerful instant audio and video conferencing system together!

## :closed_book: License
OpenMeeting is available under the GNU AFFERO GENERAL PUBLIC LICENSE 3.0. See the [LICENSE file](https://github.com/openimsdk/openmeeting-server/blob/main/LICENSE) for more information.

## 🔮 Thanks to our contributors!

<a href="https://github.com/openimsdk/openmeeting-server/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=openimsdk/openmeeting-server" />
</a>