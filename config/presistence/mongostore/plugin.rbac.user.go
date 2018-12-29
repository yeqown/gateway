package mongostore

import (
	"github.com/yeqown/gateway/config/rule"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const plgRbacUserCollName = "plg_rbac_user"

// NewUser (r rule.User) error
func (s *Store) NewUser(r rule.User) error {
	ex := new(userModel)
	if err := s.C(plgRbacUserCollName).
		Find(bson.M{"user_id": r.UserID()}).One(ex); err == nil {
		return errRuleExists
	} else if err != mgo.ErrNotFound {
		panic(err)
	}

	r.SetID(bson.NewObjectId().Hex())
	doc := loadUserModelFromUser(r)
	return s.C(plgRbacUserCollName).Insert(doc)
}

// DelUser (id string) error
func (s *Store) DelUser(id string) error {
	return s.C(plgRbacUserCollName).RemoveId(bson.ObjectIdHex(id))
}

// EditUser (id string, r rule.User) error
func (s *Store) EditUser(id string, r rule.User) error {
	ex := new(userModel)
	if err := s.C(plgRbacUserCollName).
		Find(bson.M{"user_id": r.UserID()}).One(ex); err == nil {
		return errRuleExists
	} else if err != mgo.ErrNotFound {
		panic(err)
	}
	r.SetID(id)
	doc := loadUserModelFromUser(r)
	return s.C(plgRbacUserCollName).
		UpdateId(bson.ObjectIdHex(id), doc)
}

// UserPage (limit, offset int) ([]rule.User, int)
func (s *Store) UserPage(limit, offset int) ([]rule.User, int) {
	models := make([]*userModel, 0, limit)
	total := 0
	if n, err := s.C(plgRbacUserCollName).Count(); err != nil {
		panic(err)
	} else {
		total = n
	}

	if err := s.C(plgRbacUserCollName).
		Find(nil).Skip(offset).Limit(limit).All(&models); err != nil {
		panic(err)
	}
	rules := make([]rule.User, len(models))
	for idx, model := range models {
		rules[idx] = model
		model.RolesMap = make(map[string]rule.Role)
		for _, roleID := range model.RoleIDs {
			id := roleID.Hex()
			role := s.GetRoleByID(id)
			model.RolesMap[id] = role
		}
		// log.Printf("user: %v", rules[idx])
	}
	return rules, total
}

// GetUserByID ...
func (s *Store) GetUserByID(id string) rule.User {
	model := new(userModel)
	if err := s.C(plgRbacUserCollName).
		FindId(bson.ObjectIdHex(id)).One(model); err != nil {
		// log.Printf("find user by id err: %v\n", err)
		return nil
	}
	return model
}

// AssignRole (id, roleid string) error
// func (s *Store) AssignRole(id, roleid string) error {
// 	doc := s.GetUserByID(id)
// 	perm := s.GetRoleByID(roleid)
// 	doc.Assign(perm)
// 	return s.EditUser(id, doc)
// }

// RevokeRole (id, roleid string) error
// func (s *Store) RevokeRole(id, roleid string) error {
// 	doc := s.GetUserByID(id)
// 	perm := s.GetRoleByID(roleid)
// 	doc.Revoke(perm)
// 	return s.EditUser(id, doc)
// }
