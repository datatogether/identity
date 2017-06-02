package main

// Groups holds all types of requests for groups
// it's based on an int b/c it's stateless and Go lets us
// do this sort of thing
type Groups int

// GroupsRequest defines a request for users, outlining all possible
// options for scoping & shaping the desired response
type GroupsListParams struct {
	// the user performing the request
	User *User `required:"true"`
	// users requests embeds pagination info
	Page
}

func (Groups) List(p *GroupsListParams, res *[]*Group) error {
	groups, err := ListGroups(appDB, p.Page.Size, p.Page.Offset())
	if err != nil {
		return err
	}

	*res = groups
	return nil
}

type GroupsGetParams struct {
	User  *User
	Group *Group
}

func (Groups) Get(p *GroupsGetParams, res *Group) error {
	if err := p.Group.Read(appDB); err != nil {
		return err
	}
	*res = *p.Group
	return nil
}

type GroupsCreateParams struct {
	User  *User
	Group *Group
}

func (Groups) Create(p *GroupsCreateParams, res *Group) error {
	p.Group.Creator = p.User
	if err := p.Group.Save(appDB); err != nil {
		return err
	}
	*res = *p.Group
	return nil
}

type GroupsSaveParams struct {
	User  *User
	Group *Group
}

func (Groups) Save(p *GroupsSaveParams, res *Group) error {
	// if !r.User.isAdmin && r.User.Id != r.Subject.Id {
	// 	return nil, ErrAccessDenied
	// }

	// log.Info(r.Group.Id, r.Subject.Id, r.Group.Id == r.Subject.Id)
	if err := p.Group.Save(appDB); err != nil {
		return err
	}

	*res = *p.Group
	return nil
}

type GroupsDeleteParams struct {
	Interface string
	User      *User
	Group     *Group
}

func (Groups) Delete(p *GroupsDeleteParams, res *bool) error {
	if err := p.Group.Read(appDB); err != nil {
		return err
	}

	if !p.User.isAdmin && p.Group.Creator.Id != p.User.Id {
		return ErrAccessDenied
	}

	if err := p.Group.Delete(appDB); err != nil {
		return err
	}

	*res = true
	return nil
}
