package models

type GachaRequest struct {
	AuthkeyVer string `json:"authkey_ver"`
	SignType   string `json:"sign_type"`
	AuthAppId  string `json:"auth_appid"`
	InitType   string `json:"init_type"`
	Lang       string `json:"lang"`
	Page       string `json:"page"`
	Size       string `json:"size"`
	EndId      string `json:"end_id"`
	GachaType  string `json:"gacha_type"`
}

type GachaResponse struct {
	RetCode int               `json:"retcode"`
	Message string            `json:"message"`
	Data    GachaResponseData `json:"data"`
}

type GachaResponseData struct {
	Page   string  `json:"page"`
	Size   string  `json:"size"`
	Total  string  `json:"total"`
	List   []Gacha `json:"list"`
	Region string  `json:"region"`
}

type Gacha struct {
	UID       string `json:"uid"`
	GachaType string `json:"gacha_type"`
	ItemId    string `json:"item_id"`
	Count     string `json:"count"`
	Time      string `json:"time"`
	Name      string `json:"name"`
	Lang      string `json:"lang"`
	ItemType  string `json:"item_type"`
	RankType  string `json:"rank_type"`
	Id        string `json:"id"`
}

type GachaDetail struct {
	RateOn               bool
	Last5Stars           string
	Last4Stars           string
	Last5StarsFlag       bool
	Last4StarsFlag       bool
	CountAfterLast5Stars int
	CountAfterLast4Stars int
}

type Inventory struct {
	ItemType  int
	ItemLists []InventoryDetail
}

type InventoryDetail struct {
	Name     string
	Quantity int
}
