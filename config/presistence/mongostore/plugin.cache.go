package mongostore

import (
	"github.com/yeqown/gateway/config/presistence"
	"github.com/yeqown/gateway/config/rule"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const plgCacheCollName = "plg_cache"

// NewNocacheRule func
func (s *Store) NewNocacheRule(r rule.Nocacher) error {
	r.SetID(bson.NewObjectId().Hex())
	doc := loadNocacherModelFromNocacher(r)
	// check duplicate
	var ex = new(nocacherModel)
	if err := s.C(plgCacheCollName).Find(bson.M{"regular": doc.Regular()}).
		One(ex); err == nil {
		return errRuleExists
	} else if err != mgo.ErrNotFound {
		return err
	}
	if err := s.C(plgCacheCollName).Insert(doc); err != nil {
		return err
	}
	s.notify(presistence.PlgCodeCache)
	return nil
}

// DelNocacheRule func
func (s *Store) DelNocacheRule(id string) error {
	if err := s.C(plgCacheCollName).
		RemoveId(bson.ObjectIdHex(id)); err != nil {
		return err
	}
	s.notify(presistence.PlgCodeCache)
	return nil
}

// UpdateNocacheRule func ...
func (s *Store) UpdateNocacheRule(id string, r rule.Nocacher) error {
	r.SetID(id)
	model := loadNocacherModelFromNocacher(r)
	if err := s.C(plgCacheCollName).
		UpdateId(bson.ObjectIdHex(id), model); err != nil {
		return err
	}
	s.notify(presistence.PlgCodeCache)
	return nil
}

// NocacheRules ...
func (s *Store) NocacheRules(offset, limit int) ([]rule.Nocacher, int) {
	models := make([]*nocacherModel, 0, limit)
	q := s.C(plgCacheCollName).Find(nil)
	total, err := q.Count()
	if err != nil {
		panic(err)
	}
	if err := q.Skip(offset).Limit(limit).All(&models); err != nil {
		panic(err)
	}

	rules := make([]rule.Nocacher, len(models))
	for idx, r := range models {
		rules[idx] = r
	}

	return rules, total
}

// NocacheRuleByID ...
func (s *Store) NocacheRuleByID(id string) rule.Nocacher {
	model := new(nocacherModel)
	if err := s.C(plgCacheCollName).FindId(
		bson.ObjectIdHex(id)).One(model); err != nil {
		panic(err)
	}
	return model
}

// NocacheRulesCount ...
// func (s *Store) NocacheRulesCount() int {
// 	cnt, err := s.C(plgCacheCollName).Count()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return cnt
// }
