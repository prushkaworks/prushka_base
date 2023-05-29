package db

import (
	"time"
)

type User struct {
	ID            int32           `json:"id" gorm:"primaryKey"`
	Name          string          `json:"name"`
	Email         string          `json:"email"`
	Password      string          `json:"-"`
	IsAuthorized  bool            `json:"is_authorized"`
	UserPrivilege []UserPrivilege `json:"-"`
	CreatorCard   []Card          `json:"-" gorm:"foreignKey:CreatorID;"`
	AssignedCard  []Card          `json:"-" gorm:"foreignKey:AssignedID;"`
}

type Privilege struct {
	ID            int32           `json:"id" gorm:"primaryKey"`
	Name          string          `json:"name"`
	UserPrivilege []UserPrivilege `json:"-"`
}

type Workspace struct {
	ID            int32           `json:"id" gorm:"primaryKey"`
	Name          string          `json:"name"`
	DateCreated   time.Time       `json:"date_created" gorm:"autoCreateTime"`
	UserPrivilege []UserPrivilege `json:"-"`
	Desk          []Desk          `json:"-"`
}

type UserPrivilege struct {
	ID          int32 `json:"id" gorm:"primaryKey"`
	PrivilegeID uint  `json:"privilege_id"`
	UserID      uint  `json:"user_id"`
	WorkspaceID uint  `json:"workspace_id"`
}

type Desk struct {
	ID          int32     `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	DateCreated time.Time `json:"date_created" gorm:"autoCreateTime"`
	WorkspaceID uint      `json:"workspace_id"`
	Column      []Column  `json:"-"`
}

type Column struct {
	ID     int32  `json:"id" gorm:"primaryKey"`
	Name   string `json:"name"`
	DeskID uint   `json:"desk_id"`
	Card   []Card `json:"-"`
}

type Card struct {
	ID             int32        `json:"id" gorm:"primaryKey"`
	Name           string       `json:"name"`
	Description    string       `json:"description" gorm:"text"`
	DateCreated    time.Time    `json:"date_created" gorm:"autoCreateTime"`
	DateExpiration time.Time    `json:"date_expiration"`
	IsDone         bool         `json:"is_done"`
	ColumnID       uint         `json:"column_id"`
	AssignedID     uint         `json:"assigned_id"`
	CreatorID      uint         `json:"creator_id"`
	Attachment     []Attachment `json:"-"`
	CardsLabel     []CardsLabel `json:"-"`
}

type Label struct {
	ID         int32        `json:"id" gorm:"primaryKey"`
	Name       string       `json:"name"`
	Color      string       `json:"color"`
	CardsLabel []CardsLabel `json:"-"`
}

type Attachment struct {
	ID     int32  `json:"id" gorm:"primaryKey"`
	Path   string `json:"path"`
	CardID uint   `json:"card_id"`
}

type CardsLabel struct {
	ID      int32 `json:"id" gorm:"primaryKey"`
	CardID  uint  `json:"card_id"`
	LabelID uint  `json:"label_id"`
}
