package db

type Model interface {
	GetById(ID int32)
	GetAll() []Model
}

type User struct {
	ID       int32
	Name     string
	Email    string
	Password string
}

func (u *User) GetById(ID int32) {

}

type Privilege struct {
	ID   int32
	Name string
}

type Workspace struct {
	ID          int32
	Name        string
	DateCreated string
}

type UserPrivilege struct {
	PrivilegeId int32
	UserId      int32
	WorkspaceId int32
}

type Desk struct {
	ID          int32
	Name        string
	DateCreated string
	WorkspaceId int32
}

type Column struct {
	ID     int32
	Name   string
	DeskId int32
}

type Card struct {
	ID             int32
	Name           string
	Description    string
	DateCreated    string
	DateExpiration string
	IsDone         bool
	ColumnId       int32
	Assigned       int32
	Creator        int32
}

type Label struct {
	ID    int32
	Name  string
	Color string
}

type Attachment struct {
	ID     int32
	Path   string
	CardId int32
}

type CardsLabel struct {
	CardId  int32
	LabelId int32
}
