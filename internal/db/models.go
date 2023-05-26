package db

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       int32  `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Privilege struct {
	gorm.Model
	ID   int32  `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
}

type Workspace struct {
	gorm.Model
	ID          int32     `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	DateCreated time.Time `json:"date_created" gorm:"autoCreateTime"`
}

type UserPrivilege struct {
	gorm.Model
	Privilege Privilege `json:"privilege" gorm:"embedded;embeddedPrefix:privilege_"`
	User      User      `json:"user" gorm:"embedded;embeddedPrefix:user_"`
	Workspace Workspace `json:"workspace" gorm:"embedded;embeddedPrefix:workspace_"`
}

type Desk struct {
	ID          int32     `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	DateCreated time.Time `json:"date_created" gorm:"autoCreateTime"`
	Workspace   Workspace `json:"workspace" gorm:"embedded;embeddedPrefix:workspace_"`
}

type Column struct {
	ID   int32  `json:"id" gorm:"primaryKey"`
	Name string `json:"name"`
	Desk Desk   `json:"desk" gorm:"embedded;embeddedPrefix:desk_"`
}

type Card struct {
	ID             int32     `json:"id" gorm:"primaryKey"`
	Name           string    `json:"name"`
	Description    string    `json:"description" gorm:"text"`
	DateCreated    time.Time `json:"date_created" gorm:"autoCreateTime"`
	DateExpiration time.Time `json:"date_expiration"`
	IsDone         bool      `json:"is_done"`
	Column         Column    `json:"column" gorm:"embedded;embeddedPrefix:column_"`
	Assigned       User      `json:"assigned" gorm:"embedded;embeddedPrefix:assigned_"`
	Creator        User      `json:"creator" gorm:"embedded;embeddedPrefix:creator_"`
}

type Label struct {
	ID    int32  `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Attachment struct {
	ID   int32  `json:"id" gorm:"primaryKey"`
	Path string `json:"path"`
	Card Card   `json:"card" gorm:"embedded;embeddedPrefix:card_"`
}

type CardsLabel struct {
	Card  Card  `json:"card" gorm:"embedded;embeddedPrefix:card_"`
	Label Label `json:"label" gorm:"embedded;embeddedPrefix:label_"`
}
