package db

type ControlModel struct {
	Topic string `gorm:"primaryKey"`
	Value string
}

func (ControlModel) TableName() string {
	return "virtual_controls"
}
