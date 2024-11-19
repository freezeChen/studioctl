package util

type Tag interface {
	TagName() string
	PrimaryKey() string
	Auto() string
	CreateTime() string
	UpdateTime() string
	Split() string
}

func NewTag() Tag {
	return XormTag{}
}

type XormTag struct {
}

func (XormTag) TagName() string {
	return "xorm"
}
func (XormTag) PrimaryKey() string {
	return "pk"
}
func (XormTag) Auto() string {
	return "autoincr"
}

func (XormTag) CreateTime() string {
	return "Created"
}
func (XormTag) UpdateTime() string {
	return "Updated"
}

func (XormTag) Split() string {
	return " "
}

type GormTag struct {
}

func (GormTag) TagName() string {
	return "gorm"
}

func (g GormTag) PrimaryKey() string {
	return "primaryKey"
}

func (g GormTag) Auto() string {
	return "autoIncrement"
}

func (g GormTag) CreateTime() string {
	return "autoCreateTime"
}

func (g GormTag) UpdateTime() string {
	return "autoUpdateTime"
}

func (g GormTag) Split() string {
	return ";"
}
