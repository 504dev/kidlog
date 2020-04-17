package types

type Dashboard struct {
	Id         int         `db:"id"          json:"id"`
	OwnerId    int         `db:"owner_id"    json:"owner_id"`
	Name       string      `db:"name"        json:"name"`
	PublicKey  string      `db:"public_key"  json:"public_key"`
	PrivateKey string      `db:"private_key" json:"private_key"`
	Members    DashMembers `json:"members"`
}
type Dashboards []*Dashboard

func (ds Dashboards) Ids() []int {
	ids := make([]int, len(ds))
	for k, v := range ds {
		ids[k] = v.Id
	}
	return ids
}
func (ds Dashboards) ByPrimary() map[int]*Dashboard {
	res := make(map[int]*Dashboard, len(ds))
	for _, v := range ds {
		res[v.Id] = v
	}
	return res
}

type DashKey struct {
	Id         int    `db:"id"          json:"id"`
	DashId     int    `db:"dash_id"     json:"dash_id"`
	PublicKey  string `db:"public_key"  json:"public_key"`
	PrivateKey string `db:"private_key" json:"private_key"`
}
type DashKeys []*DashKey

type DashMember struct {
	Id     int `db:"id"      json:"id"`
	DashId int `db:"dash_id" json:"dash_id"`
	UserId int `db:"user_id" json:"user_id"`
}
type DashMembers []*DashMember

func (dm DashMembers) DashIds() []int {
	ids := make([]int, len(dm))
	for k, v := range dm {
		ids[k] = v.DashId
	}
	return ids
}

type DashStatRow struct {
	DashId   int    `db:"dash_id" json:"dash_id"`
	Hostname string `db:"hostname" json:"hostname"`
	Logname  string `db:"logname"  json:"logname"`
	Level    string `db:"level"    json:"level"`
	Version  string `db:"version"  json:"version"`
	Cnt      int    `db:"cnt"      json:"cnt"`
	Updated  string `db:"updated"  json:"updated"`
}

type DashStatRows []*DashStatRow
