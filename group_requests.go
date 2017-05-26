package main

// GroupsRequest defines a request for users, outlining all possible
// options for scoping & shaping the desired response
type GroupsRequest struct {
	Interface string
	// the user performing the request
	User *User `required:"true"`
	// users requests embeds pagination info
	Page
}

func (r *GroupsRequest) Exec() (interface{}, error) {
	return Groups(appDB, r.Page.Size, r.Page.Offset())
}

type GroupRequest struct {
	Interface string
	User      *User
	Group     *Group
}

func (r *GroupRequest) Exec() (interface{}, error) {
	if err := r.Group.Read(appDB); err != nil {
		return nil, err
	}
	return r.Group, nil
}

type CreateGroupRequest struct {
	Interface string
	User      *User
	Group     *Group
}

func (r *CreateGroupRequest) Exec() (interface{}, error) {
	r.Group.Creator = r.User
	if err := r.Group.Save(appDB); err != nil {
		return nil, err
	}
	return r.Group, nil
}

type SaveGroupRequest struct {
	Interface string
	User      *User
	Group     *Group
}

func (r *SaveGroupRequest) Exec() (interface{}, error) {
	// if !r.User.isAdmin && r.User.Id != r.Subject.Id {
	// 	return nil, ErrAccessDenied
	// }

	// log.Info(r.Group.Id, r.Subject.Id, r.Group.Id == r.Subject.Id)
	if err := r.Group.Save(appDB); err != nil {
		return nil, err
	}

	return r.Group, nil
}

type DeleteGroupRequest struct {
	Interface string
	User      *User
	Group     *Group
}

func (r *DeleteGroupRequest) Exec() (interface{}, error) {
	if err := r.Group.Read(appDB); err != nil {
		return nil, err
	}

	if !r.User.isAdmin && r.Group.Creator.Id != r.User.Id {
		return nil, ErrAccessDenied
	}

	if err := r.Group.Delete(appDB); err != nil {
		return nil, err
	}

	return "ok", nil
}
