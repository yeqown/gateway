package mongostore

import (
	"github.com/yeqown/gateway/config/rule"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const plgRbacPermURLCollName = "plg_rbac_permurl"

// NewPermitURL (r rule.PermitURL) error
func (s *Store) NewPermitURL(r rule.PermitURL) error {
	ex := new(permitURLModel)
	if err := s.C(plgRbacPermURLCollName).
		Find(bson.M{"url": r.URL(), "perm_id": bson.ObjectIdHex(r.Permission().ID())}).
		One(ex); err == nil {
		return errRuleExists
	} else if err != mgo.ErrNotFound {
		panic(err)
	}
	r.SetID(bson.NewObjectId().Hex())
	doc := loadPermitURLModelFormPermitURL(r)
	return s.C(plgRbacPermURLCollName).Insert(doc)
}

// DelPermitURL (id string) error
func (s *Store) DelPermitURL(id string) error {
	return s.C(plgRbacPermURLCollName).RemoveId(bson.ObjectIdHex(id))
}

// EditPermitURL (id string, r rule.PermitURL) error
func (s *Store) EditPermitURL(id string, r rule.PermitURL) error {
	ex := new(permitURLModel)
	if err := s.C(plgRbacPermURLCollName).
		Find(bson.M{"url": r.URL(), "perm_id": bson.ObjectIdHex(r.Permission().ID())}).
		One(ex); err == nil {
		return errRuleExists
	} else if err != mgo.ErrNotFound {
		panic(err)
	}
	r.SetID(id)
	doc := loadPermitURLModelFormPermitURL(r)
	return s.C(plgRbacPermURLCollName).
		UpdateId(bson.ObjectIdHex(id), doc)
}

// PermitURLPage (limit, offset int) ([]rule.PermitURL, int)
func (s *Store) PermitURLPage(limit, offset int) ([]rule.PermitURL, int) {
	models := make([]*permitURLModel, 0, limit)
	total := 0
	if n, err := s.C(plgRbacPermURLCollName).Count(); err != nil {
		panic(err)
	} else {
		total = n
	}

	if err := s.C(plgRbacPermURLCollName).
		Find(nil).Skip(offset).Limit(limit).All(&models); err != nil {
		panic(err)
	}
	rules := make([]rule.PermitURL, len(models))
	for idx, model := range models {
		perm := s.GetPermissionByID(model.PermID.Hex())
		model.Perm = perm
		rules[idx] = model
	}
	return rules, total
}
