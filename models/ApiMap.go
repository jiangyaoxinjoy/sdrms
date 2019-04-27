package models

func (c *Map) TableName() string {
	return MapTBName()
}
func (c *MapStatus) TableName() string {
	return MapStatusTBName()
}

type Map struct {
	Id        int `orm:"pk;column(id);`
	Longitude float32
	Latitude  float32
	Name      string
	Address   string
	MapStatus *MapStatus `orm:"rel(fk)"` //设置一对多关系
}

type MapStatus struct {
	Id    int `orm:"pk;column(id);`
	Name  string
	Lever int
	Maps  []*Map `orm:"pk;reverse(many)"` // 设置一对多的反向关系
}
