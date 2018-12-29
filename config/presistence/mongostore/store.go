package mongostore

import (
	"errors"

	"github.com/yeqown/gateway/config/presistence"
	"github.com/yeqown/gateway/config/rule"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	_             presistence.Store = &Store{}
	errRuleExists                   = errors.New("rule exists")
)

// New ...
func New(url string, databaseName string) (presistence.Store, error) {
	// url like: mongodb://myuser:mypass@localhost:3306
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &Store{
		session:  session,
		db:       session.DB(databaseName),
		changedC: make(chan presistence.ChangedChan, 5),
	}, nil
}

// Store ...
type Store struct {
	db       *mgo.Database
	session  *mgo.Session
	changedC chan presistence.ChangedChan
}

// C means collection
func (s *Store) C(name string) *mgo.Collection {
	return s.db.C(name)
}

func (s *Store) loadNocacherRules() []rule.Nocacher {
	models := make([]*nocacherModel, 0)
	if err := s.C(plgCacheCollName).Find(nil).All(&models); err != nil {
		panic(err)
	}
	rules := make([]rule.Nocacher, len(models))
	for idx, m := range models {
		rules[idx] = m
	}
	return rules
}

func (s *Store) loadProxyServerRules() []rule.ServerRuler {
	models := make([]*serverRulerModel, 0)
	if err := s.C(plgProxyServerCollName).Find(nil).All(&models); err != nil {
		panic(err)
	}
	rules := make([]rule.ServerRuler, len(models))
	for idx, m := range models {
		rules[idx] = m
	}
	return rules
}

func (s *Store) loadProxyPathRules() []rule.PathRuler {
	models := make([]*pathRulerModel, 0)
	if err := s.C(plgProxyPathCollName).Find(nil).All(&models); err != nil {
		panic(err)
	}
	rules := make([]rule.PathRuler, len(models))
	for idx, m := range models {
		rules[idx] = m
	}
	return rules
}

func (s *Store) loadReverseServerRules() map[string][]rule.ReverseServer {
	mapNameCnt := s.ReverseServerGroups()
	result := make(map[string][]rule.ReverseServer)
	for name, cnt := range mapNameCnt {
		result[name] = make([]rule.ReverseServer, cnt)
		models := make([]*reverseServerModel, cnt)
		if err := s.C(plgPorxyReverseSrvCollName).
			Find(bson.M{"group": name}).All(&models); err != nil {
			panic(err)
		}
		for idx, m := range models {
			result[name][idx] = m
		}
	}

	return result
}

func (s *Store) loadRbacUsers() []rule.User {
	models := make([]*userModel, 0)
	if err := s.C(plgRbacUserCollName).
		Find(nil).All(&models); err != nil {
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
	}
	return rules
}

func (s *Store) loadRbacPermitURLs() []rule.PermitURL {
	models := make([]*permitURLModel, 0)
	if err := s.C(plgRbacPermURLCollName).
		Find(nil).All(&models); err != nil {
		panic(err)
	}
	rules := make([]rule.PermitURL, len(models))
	for idx, model := range models {
		perm := s.GetPermissionByID(model.PermID.Hex())
		model.Perm = perm
		rules[idx] = model
	}
	return rules
}

// Instance ...
func (s *Store) Instance() *presistence.Instance {
	return &presistence.Instance{
		ProxyServerRules:    s.loadProxyServerRules(),
		ProxyPathRules:      s.loadProxyPathRules(),
		ProxyReverseServers: s.loadReverseServerRules(),
		Nocache:             s.loadNocacherRules(),
		Users:               s.loadRbacUsers(),
		URLS:                s.loadRbacPermitURLs(),
	}
}

// Updated ...
func (s *Store) Updated() <-chan presistence.ChangedChan {
	return s.changedC
}

func (s *Store) notify(code presistence.PluginCode) {
	s.changedC <- presistence.ChangedChan{
		Code: code,
	}
}
