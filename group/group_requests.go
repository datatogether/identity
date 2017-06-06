package group

import (
	"github.com/archivers-space/errors"
	"github.com/archivers-space/identity/user"
	"github.com/archivers-space/sqlutil"
)

// Groups holds all types of requests for groups
// it's based on an int b/c it's stateless and Go lets us
// do this sort of thing
type GroupRequests struct {
	Store sqlutil.Transactable
}

// GroupsRequest defines a request for users, outlining all possible
// options for scoping & shaping the desired response
type GroupsListParams struct {
	// the user performing the request
	User *user.User `required:"true"`
	// users requests embeds pagination info
	Limit  int
	Offset int
}

func (r GroupRequests) List(p *GroupsListParams, res *[]*Group) error {
	groups, err := ListGroups(r.Store, p.Limit, p.Offset)
	if err != nil {
		return err
	}

	*res = groups
	return nil
}

type GroupsGetParams struct {
	User  *user.User
	Group *Group
}

func (r GroupRequests) Get(p *GroupsGetParams, res *Group) error {
	if err := p.Group.Read(r.Store); err != nil {
		return err
	}
	*res = *p.Group
	return nil
}

type GroupsCreateParams struct {
	User  *user.User
	Group *Group
}

func (r GroupRequests) Create(p *GroupsCreateParams, res *Group) error {
	p.Group.Creator = p.User
	if err := p.Group.Save(r.Store); err != nil {
		return err
	}
	*res = *p.Group
	return nil
}

type GroupsSaveParams struct {
	User  *user.User
	Group *Group
}

func (r GroupRequests) Save(p *GroupsSaveParams, res *Group) error {
	// if !r.User.isAdmin && r.User.Id != r.Subject.Id {
	// 	return nil, ErrAccessDenied
	// }

	// log.Info(r.Group.Id, r.Subject.Id, r.Group.Id == r.Subject.Id)
	if err := p.Group.Save(r.Store); err != nil {
		return err
	}

	*res = *p.Group
	return nil
}

type GroupsDeleteParams struct {
	Interface string
	User      *user.User
	Group     *Group
}

func (r GroupRequests) Delete(p *GroupsDeleteParams, res *bool) error {
	if err := p.Group.Read(r.Store); err != nil {
		return err
	}

	// TODO - make isAdmin public? permissions package?
	// if !p.User.isAdmin && p.Group.Creator.Id != p.User.Id {
	// 	return errors.ErrAccessDenied
	// }

	if p.Group.Creator.Id != p.User.Id {
		return errors.ErrAccessDenied
	}

	if err := p.Group.Delete(r.Store); err != nil {
		return err
	}

	*res = true
	return nil
}
