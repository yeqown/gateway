package mongostore

import (
	"github.com/yeqown/gateway/config/presistence"
	"github.com/yeqown/gateway/config/rule"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const plgProxyPathCollName = "plg_proxy_path"

// NewPathRule func
func (s *Store) NewPathRule(r rule.PathRuler) error {
	r.SetID(bson.NewObjectId().Hex())
	for _, c := range r.CombineReqCfgs() {
		c.SetID(bson.NewObjectId().Hex())
	}
	doc := loadPathRulerModelFromPathRuler(r)
	// check dulicated
	var ex = new(pathRulerModel)
	if err := s.C(plgProxyPathCollName).
		Find(bson.M{"path": doc.Path(), "method": doc.Method()}).
		One(ex); err != nil {
		return errRuleExists
	} else if err != mgo.ErrNotFound {
		return err
	}
	if err := s.C(plgProxyPathCollName).Insert(doc); err != nil {
		return err
	}
	s.notify(presistence.PlgCodeProxyPath)
	return nil
}

// DelPathRule func
func (s *Store) DelPathRule(id string) error {
	if err := s.C(plgProxyPathCollName).
		RemoveId(bson.ObjectIdHex(id)); err != nil {
		return err
	}
	s.notify(presistence.PlgCodeProxyPath)
	return nil
}

// UpdatePathRule func ...
func (s *Store) UpdatePathRule(id string, r rule.PathRuler) error {
	r.SetID(id)
	model := loadPathRulerModelFromPathRuler(r)
	if err := s.C(plgProxyPathCollName).
		UpdateId(bson.ObjectIdHex(id), model); err != nil {
		return err
	}
	s.notify(presistence.PlgCodeProxyPath)
	return nil
}

// PathRuleByID ...
func (s *Store) PathRuleByID(id string) rule.PathRuler {
	model := new(pathRulerModel)
	if err := s.C(plgProxyPathCollName).FindId(bson.ObjectIdHex(id)).One(model); err != nil {
		panic(err)
	}
	return model
}

// PathRulesPage func
func (s *Store) PathRulesPage(offset, limit int) ([]rule.PathRuler, int) {
	models := make([]*pathRulerModel, 0, limit)
	q := s.C(plgProxyPathCollName).Find(nil)
	total, err := q.Count()
	if err != nil {
		panic(err)
	}
	if err := q.Skip(offset).Limit(limit).All(&models); err != nil {
		panic(err)
	}

	rules := make([]rule.PathRuler, len(models))
	for idx, r := range models {
		rules[idx] = r
	}

	return rules, total
}

// PathRulesCount func
// func (s *Store) PathRulesCount() int {
// 	cnt, err := s.C(plgProxyPathCollName).Count()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return cnt
// }
