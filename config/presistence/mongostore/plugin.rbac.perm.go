package mongostore

import (
	"log"

	"github.com/yeqown/gateway/config/rule"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const plgRbacPermCollName = "plg_rbac_perm"

// NewPermission (r rule.Permission) error
func (s *Store) NewPermission(r rule.Permission) error {
	ex := new(permissionModel)
	if err := s.C(plgRbacPermCollName).
		Find(bson.M{"action": r.Action(), "resource": r.Resource()}).
		One(ex); err == nil {
		return errRuleExists
	} else if err != mgo.ErrNotFound {
		panic(err)
	}

	r.SetID(bson.NewObjectId().Hex())
	doc := loadPermissionModelFromPermission(r)
	return s.C(plgRbacPermCollName).Insert(doc)
}

// DelPermission (id string) error
func (s *Store) DelPermission(id string) error {
	return s.C(plgRbacPermCollName).RemoveId(bson.ObjectIdHex(id))
}

// EditPermission (id string, r rule.Permission) error
func (s *Store) EditPermission(id string, r rule.Permission) error {
	ex := new(permissionModel)
	if err := s.C(plgRbacPermCollName).
		Find(bson.M{"action": r.Action(), "resource": r.Resource()}).
		One(ex); err == nil {
		return errRuleExists
	} else if err != mgo.ErrNotFound {
		panic(err)
	}
	r.SetID(id)
	doc := loadPermissionModelFromPermission(r)
	return s.C(plgRbacPermCollName).UpdateId(bson.ObjectIdHex(id), doc)
}

// GetPermissionByID ...
func (s *Store) GetPermissionByID(id string) rule.Permission {
	model := new(permissionModel)
	if err := s.C(plgRbacPermCollName).FindId(bson.ObjectIdHex(id)).One(model); err != nil {
		log.Printf("find permission by id err: %v\n", err)
		return nil
	}
	return model
}

// PermissionPage (limit, offset int) ([]rule.Permission, int)
func (s *Store) PermissionPage(limit, offset int) ([]rule.Permission, int) {
	models := make([]*permissionModel, 0, limit)
	total := 0
	if n, err := s.C(plgRbacPermCollName).Count(); err != nil {
		panic(err)
	} else {
		total = n
	}

	if err := s.C(plgRbacPermCollName).
		Find(nil).Skip(offset).Limit(limit).All(&models); err != nil {
		panic(err)
	}
	rules := make([]rule.Permission, len(models))
	for idx, model := range models {
		rules[idx] = model
	}
	return rules, total
}
