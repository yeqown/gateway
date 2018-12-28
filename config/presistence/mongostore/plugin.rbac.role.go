package mongostore

import (
	"github.com/yeqown/gateway/config/rule"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const plgRbacRoleCollName = "plg_rbac_role"

// NewRole (r rule.Role) error
func (s *Store) NewRole(r rule.Role) error {
	ex := new(roleModel)
	if err := s.C(plgRbacRoleCollName).
		Find(bson.M{"name": r.Name()}).One(ex); err == nil {
		return errRuleExists
	} else if err != mgo.ErrNotFound {
		panic(err)
	}

	r.SetID(bson.NewObjectId().Hex())
	doc := loadRoleModelFromRole(r)
	// log.Printf("%v \n", *doc)
	return s.C(plgRbacRoleCollName).Insert(doc)
}

// DelRole (id string) error
func (s *Store) DelRole(id string) error {
	return s.C(plgRbacRoleCollName).
		RemoveId(bson.ObjectIdHex(id))
}

// EditRole (id string, r rule.Role) error
func (s *Store) EditRole(id string, r rule.Role) error {
	ex := new(roleModel)
	if err := s.C(plgRbacRoleCollName).
		Find(bson.M{"name": r.Name()}).One(ex); err == nil {
		return errRuleExists
	} else if err != mgo.ErrNotFound {
		panic(err)
	}
	r.SetID(id)
	doc := loadRoleModelFromRole(r)
	return s.C(plgRbacRoleCollName).
		UpdateId(bson.ObjectIdHex(id), doc)
}

// RolePage (limit, offset int) ([]rule.Role, int)
func (s *Store) RolePage(limit, offset int) ([]rule.Role, int) {
	models := make([]*roleModel, 0, limit)
	total := 0
	if n, err := s.C(plgRbacRoleCollName).Count(); err != nil {
		panic(err)
	} else {
		total = n
	}

	if err := s.C(plgRbacRoleCollName).
		Find(nil).Skip(offset).Limit(limit).All(&models); err != nil {
		panic(err)
	}
	rules := make([]rule.Role, len(models))
	for idx, model := range models {
		model.permissions = make(map[string]rule.Permission)
		for _, permID := range model.PermIDs {
			id := permID.Hex()
			model.permissions[id] = s.GetPermissionByID(id)
		}
		rules[idx] = model
	}
	return rules, total
}

// GetRoleByID ...
func (s *Store) GetRoleByID(id string) rule.Role {
	model := new(roleModel)
	if err := s.C(plgRbacRoleCollName).
		FindId(bson.ObjectIdHex(id)).One(model); err != nil {
		// log.Printf("find role by id err: %v\n", err)
		return nil
	}

	model.permissions = make(map[string]rule.Permission)
	for _, permID := range model.PermIDs {
		id := permID.Hex()
		model.permissions[id] = s.GetPermissionByID(id)
	}

	return model
}

// AssignPerm (id, permid string) error
func (s *Store) AssignPerm(id string, permids ...string) error {
	doc := s.GetRoleByID(id)
	for _, permid := range permids {
		perm := s.GetPermissionByID(permid)
		doc.Assign(perm)
	}
	return s.EditRole(id, doc)
}

// RevokePerm (id, permid string) error
func (s *Store) RevokePerm(id string, permids ...string) error {
	doc := s.GetRoleByID(id)
	for _, permid := range permids {
		perm := s.GetPermissionByID(permid)
		doc.Revoke(perm)
	}
	return s.EditRole(id, doc)
}
