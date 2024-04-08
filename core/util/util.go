package util

type Tag interface {
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
