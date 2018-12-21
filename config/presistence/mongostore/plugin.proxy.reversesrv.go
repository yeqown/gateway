package mongostore

import (
	"log"

	"github.com/yeqown/gateway/config/presistence"
	"github.com/yeqown/gateway/config/rule"
	"gopkg.in/mgo.v2/bson"
)

const plgPorxyReverseSrvCollName = "plg_proxy_reverse_srv"

// // NewReverseServerGroup ...
// func (s *Store) NewReverseServerGroup(group string) error {
// 	models := make([]*reverseServerModel, 0)
// 	s.C(plgPorxyReverseSrvCollName).Find(bson.M{"group": group}).All(models)
// }

// NewReverseServer func
func (s *Store) NewReverseServer(group string, r rule.ReverseServer) error {
	r.SetID(bson.NewObjectId().Hex())
	doc := loadReverseServerModelFromReverseServer(r)
	if err := s.C(plgPorxyReverseSrvCollName).Insert(doc); err != nil {
		return err
	}
	s.notify(presistence.PlgCodeProxyReverseSrv)
	return nil
}

// DelReverseServer func
func (s *Store) DelReverseServer(id string) error {
	if err := s.C(plgPorxyReverseSrvCollName).
		RemoveId(bson.ObjectIdHex(id)); err != nil {
		return err
	}
	s.notify(presistence.PlgCodeProxyReverseSrv)
	return nil
}

// DelReverseServerGroup func
func (s *Store) DelReverseServerGroup(group string) error {
	if err := s.C(plgPorxyReverseSrvCollName).
		Remove(bson.M{"group": group}); err != nil {
		return err
	}
	s.notify(presistence.PlgCodeProxyReverseSrv)
	return nil
}

// UpdateReverseServerGroupName ...
func (s *Store) UpdateReverseServerGroupName(group string, newname string) error {
	changeInfo, err := s.C(plgPorxyReverseSrvCollName).
		UpdateAll(bson.M{"group": group}, bson.M{"group": newname})

	log.Printf("UpdateReverseServerGroupName: %v\n", changeInfo)
	if err != nil {
		return err
	}
	s.notify(presistence.PlgCodeProxyReverseSrv)
	return nil
}

// UpdateReverseServer func ...
func (s *Store) UpdateReverseServer(id string, r rule.ReverseServer) error {
	r.SetID(id)
	doc := loadReverseServerModelFromReverseServer(r)
	if err := s.C(plgPorxyReverseSrvCollName).
		UpdateId(bson.ObjectIdHex(id), doc); err != nil {
		return err
	}
	s.notify(presistence.PlgCodeProxyReverseSrv)
	return nil
}

// ReverseServerByID ...
func (s *Store) ReverseServerByID(group string, id string) rule.ReverseServer {
	model := new(reverseServerModel)
	s.C(plgPorxyReverseSrvCollName).
		FindId(bson.ObjectIdHex(id)).One(model)
	return model
}

// ReverseServerByGroup func
func (s *Store) ReverseServerByGroup(group string, offset, limit int) ([]rule.ReverseServer, int) {
	models := make([]*reverseServerModel, 0, limit)
	q := s.C(plgPorxyReverseSrvCollName).Find(bson.M{"group": group})
	total, err := q.Count()
	if err != nil {
		panic(err)
	}

	if err := q.Skip(offset).Limit(limit).
		All(&models); err != nil {
		panic(err)
	}

	rules := make([]rule.ReverseServer, len(models))
	for idx, r := range models {
		rules[idx] = r
	}

	return rules, total
}

// ReverseServerGroups ... 获取所有的分组名字和数量
func (s *Store) ReverseServerGroups() map[string]int {
	result := make(map[string]int)
	groupNames := make([]string, 0)

	if err := s.C(plgPorxyReverseSrvCollName).
		Find(nil).Distinct("group", &groupNames); err != nil {
		panic(err)
	}

	log.Println(groupNames)

	for _, name := range groupNames {
		cnt, err := s.C(plgPorxyReverseSrvCollName).
			Find(bson.M{"group": name}).Count()
		if err != nil {
			panic(err)
		}
		result[name] = cnt
	}

	return result
}

// ReverseServerGroupCount func
// func (s *Store) ReverseServerGroupCount() int {
// 	cnt, err := s.C(plgPorxyReverseSrvCollName).Count()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return cnt
// }

// ReverseServerGroupPageCount func
// func (s *Store) ReverseServerGroupPageCount(group string) int {
// 	cnt, err := s.C(plgPorxyReverseSrvCollName).
// 		Find(bson.M{"group": group}).Count()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return cnt
// }
