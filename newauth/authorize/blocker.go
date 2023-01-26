package authorize

import (
	"encoding/json"
	"errors"

	"gorm.io/gorm"
)

type Subject struct {
	Role     RoleEnum `json:"role" binding:"enum"`
	UserID   uint     `json:"user_id"`
	DeviceID uint     `json:"device_id"`
}

func (s *Subject) getValue() (string, error) {
	if s.Role != "" {
		return string(s.Role), nil
	}

	if s.UserID != 0 {
		return userString(s.UserID), nil
	}

	if s.DeviceID != 0 {
		return deviceString(s.DeviceID), nil
	}

	return "", errors.New("subject value empty")
}

type BlockCommand struct {
	Subject  Subject      `json:"subject" validate:"required"`
	Domain   Domain       `json:"domain"`
	Resource ResourceEnum `json:"resource" binding:"enum"`
	Act      ActBasicEnum `json:"action" binding:"enum"`
}

type BlockItem struct {
	gorm.Model
	Name string
	Data string
}

func (b *BlockItem) GetData() ([]*BlockCommand, error) {
	var command []*BlockCommand

	err := json.Unmarshal([]byte(b.Data), &command)
	return command, err
}

func (b *BlockItem) SetData(data []*BlockCommand) error {

	datastr, err := json.Marshal(&data)
	if err == nil {
		b.Data = string(datastr)
	}
	return err
}

type Blocker struct {
	db     *gorm.DB
	forcer *Enforcer
}

func (b *Blocker) Create(name string, commands []*BlockCommand) *BlockItem {
	commandstr := make([][]string, len(commands))

	bitem := BlockItem{
		Name: name,
	}
	bitem.SetData(commands)

	err := b.db.Save(&bitem).Error
	if err != nil {
		panic(err)
	}

	for i, cmd := range commands {
		sub, err := cmd.Subject.getValue()
		if err != nil {
			panic(err)
		}
		commandstr[i] = []string{sub, string(cmd.Domain), string(cmd.Resource), string(cmd.Act), string(DenyEffect)}
	}

	_, err = b.forcer.forcer.AddPolicies(commandstr)
	if err != nil {
		panic(err)
	}

	return &bitem
}

func (b *Blocker) Get(blockID uint) (*BlockItem, error) {

	var item BlockItem

	err := b.db.First(&item, blockID).Error
	return &item, err
}

func (b *Blocker) Delete(blockID uint) error {
	item, err := b.Get(blockID)

	if err != nil {
		return err
	}

	commands, err := item.GetData()
	if err != nil {
		return err
	}
	commandstr := make([][]string, len(commands))
	for i, cmd := range commands {
		sub, err := cmd.Subject.getValue()
		if err != nil {
			panic(err)
		}
		commandstr[i] = []string{sub, string(cmd.Domain), string(cmd.Resource), string(cmd.Act), string(DenyEffect)}
	}

	_, err = b.forcer.forcer.RemovePolicies(commandstr)
	if err != nil {
		panic(err)
	}

	return b.db.Delete(&BlockItem{}, blockID).Error
}

func NewBlocker(db *gorm.DB) *Blocker {

	return &Blocker{
		db: db,
	}
}
