package mongostore

import (
	"github.com/yeqown/gateway/config/presistence"
	"github.com/yeqown/gateway/config/rule"
	"gopkg.in/mgo.v2/bson"
)

const plgProxyServerCollName = "plg_proxy_srv"

// NewServerRule func
func (s *Store) NewServerRule(r rule.ServerRuler) error {
	r.SetID(bson.NewObjectId().Hex())
	doc := loadServerRulerModelFromServerRuler(r)
	if err := s.C(plgProxyServerCollName).Insert(doc); err != nil {
		return err
	}
	s.notify(presistence.PlgCodeProxyServer)
	return nil
}

// DelServerRule func
func (s *Store) DelServerRule(id string) error {
	if err := s.C(plgProxyServerCollName).RemoveId(
		bson.ObjectIdHex(id)); err != nil {
		return err
	}
	s.notify(presistence.PlgCodeProxyServer)
	return nil
}

// UpdateServerRule func ...
func (s *Store) UpdateServerRule(id string, r rule.ServerRuler) error {
	r.SetID(id)
	doc := loadServerRulerModelFromServerRuler(r)
	if err := s.C(plgProxyServerCollName).UpdateId(
		bson.ObjectIdHex(id), doc); err != nil {
		return err
	}
	s.notify(presistence.PlgCodeProxyServer)
	return nil
}

// ServerRuleByID ...
func (s *Store) ServerRuleByID(id string) rule.ServerRuler {
	model := new(serverRulerModel)
	if err := s.C(plgProxyServerCollName).
		FindId(bson.ObjectIdHex(id)).One(model); err != nil {
		panic(err)
	}
	return model
}

// ServerRulesPage func
func (s *Store) ServerRulesPage(offset, limit int) ([]rule.ServerRuler, int) {
	models := make([]*serverRulerModel, 0, limit)
	q := s.C(plgProxyServerCollName).Find(nil)
	total, err := q.Count()
	if err != nil {
		panic(err)
	}
	if err := q.Skip(offset).Limit(limit).All(&models); err != nil {
		panic(err)
	}

	rules := make([]rule.ServerRuler, len(models))
	for idx, r := range models {
		rules[idx] = r
	}

	return rules, total
}

// ServerRulesCount func
// func (s *Store) ServerRulesCount() int {
// 	cnt, err := s.C(plgProxyServerCollName).Count()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return cnt
// }
