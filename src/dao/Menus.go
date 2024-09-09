package dao

// Menu 菜单
type Menu struct {
	Id       uint32 `json:"id"`        //编号
	ParentId uint32 `json:"parent_id"` //父级菜单编号
	Name     string `json:"name"`      //名称
	Url      string `json:"url"`       //URL
	Level    uint8  `json:"level"`     //级别
	Children []Menu `json:"children"`  //子级菜单
}
