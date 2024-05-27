package apistruct

type AdminLoginResp struct {
	AdminAccount string `json:"adminAccount"`
	AdminToken   string `json:"adminToken"`
	Nickname     string `json:"nickname"`
}
