package db

import (
	"time"
)

// User - users table
type User struct {
	ID            int32           `json:"id" gorm:"primaryKey" example:"1"`
	Name          string          `json:"name" example:"Dmitrii"`
	Email         string          `json:"email" gorm:"unique" example:"dmitrii@mail.su"`
	Password      string          `json:"-"`
	IsAuthorized  bool            `json:"is_authorized" example:"true"`
	UserPrivilege []UserPrivilege `json:"-"`
	CreatorCard   []Card          `json:"-" gorm:"foreignKey:CreatorID;"`
	AssignedCard  []Card          `json:"-" gorm:"foreignKey:AssignedID;"`
}

// Privilege - privileges table
type Privilege struct {
	ID            int32           `json:"id" gorm:"primaryKey" example:"1"`
	Name          string          `json:"name" example:"admin"`
	UserPrivilege []UserPrivilege `json:"-"`
}

// Workspace - workspaces table
type Workspace struct {
	ID            int32           `json:"id" gorm:"primaryKey" example:"1"`
	Name          string          `json:"name" example:"SomeProjectName"`
	DateCreated   time.Time       `json:"date_created" gorm:"autoCreateTime" example:"2019-11-09T21:21:46+00:00"`
	UserPrivilege []UserPrivilege `json:"-"`
	Desk          []Desk          `json:"-"`
}

// UserPrivilege - link between user, privilege and workspace
type UserPrivilege struct {
	ID          int32 `json:"id" gorm:"primaryKey" example:"1"`
	PrivilegeID uint  `json:"privilege_id" example:"2"`
	UserID      uint  `json:"user_id" example:"3"`
	WorkspaceID uint  `json:"workspace_id" example:"4"`
}

// Desk - desks table
type Desk struct {
	ID          int32     `json:"id" gorm:"primaryKey" example:"1"`
	Name        string    `json:"name" example:"Backend"`
	DateCreated time.Time `json:"date_created" gorm:"autoCreateTime" example:"2019-11-09T21:21:46+00:00"`
	WorkspaceID uint      `json:"workspace_id" example:"1"`
	Column      []Column  `json:"-"`
}

// Column - columns table
type Column struct {
	ID     int32  `json:"id" gorm:"primaryKey" example:"1"`
	Name   string `json:"name" example:"TO DO"`
	DeskID uint   `json:"desk_id" example:"2"`
	Card   []Card `json:"-"`
}

// Card - cards table
type Card struct {
	ID             int32        `json:"id" gorm:"primaryKey" example:"1"`
	Name           string       `json:"name" example:"Create routers"`
	Description    string       `json:"description" gorm:"text" example:"Some long description"`
	DateCreated    time.Time    `json:"date_created" gorm:"autoCreateTime" example:"2019-11-09T21:21:46+00:00"`
	DateExpiration time.Time    `json:"date_expiration" example:"2020-11-09T21:21:46+00:00"`
	IsDone         bool         `json:"is_done" example:"false"`
	ColumnID       uint         `json:"column_id" example:"1"`
	AssignedID     uint         `json:"assigned_id" example:"660"`
	CreatorID      uint         `json:"creator_id" example:"661"`
	Attachment     []Attachment `json:"-"`
	CardsLabel     []CardsLabel `json:"-"`
}

// Label - labels table
type Label struct {
	ID         int32        `json:"id" gorm:"primaryKey" example:"1"`
	Name       string       `json:"name" example:"SomeLabel"`
	Color      string       `json:"color" example:"SomeHexColor"`
	CardsLabel []CardsLabel `json:"-"`
}

// Attachment - attachments table
type Attachment struct {
	ID     int32  `json:"id" gorm:"primaryKey" example:"1"`
	Path   string `json:"path" example:"/attachments/456.jpeg"`
	CardID uint   `json:"card_id" example:"32"`
}

// CardsLabel - link between card and label
type CardsLabel struct {
	ID      int32 `json:"id" gorm:"primaryKey" example:"1"`
	CardID  uint  `json:"card_id" example:"32"`
	LabelID uint  `json:"label_id" example:"64"`
}
