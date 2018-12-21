package mongostore

import (
	"github.com/yeqown/gateway/config/rule"
	"gopkg.in/mgo.v2/bson"
)

var (
	_ rule.Nocacher      = &nocacherModel{}
	_ rule.PathRuler     = &pathRulerModel{}
	_ rule.Combiner      = &combinerModel{}
	_ rule.ServerRuler   = &serverRulerModel{}
	_ rule.ReverseServer = &reverseServerModel{}
)

type nocacherModel struct {
	Idx      bson.ObjectId `bson:"_id"`
	Regexp   string        `bson:"regular"`
	EEnabled bool          `bson:"enabled"`
}

func (c *nocacherModel) Regular() string { return c.Regexp }
func (c *nocacherModel) Enabled() bool   { return c.EEnabled }
func (c *nocacherModel) ID() string      { return c.Idx.Hex() }
func (c *nocacherModel) SetID(id string) { c.Idx = bson.ObjectIdHex(id) }
func loadNocacherModelFromNocacher(r rule.Nocacher) *nocacherModel {
	return &nocacherModel{
		Idx:      bson.ObjectIdHex(r.ID()),
		Regexp:   r.Regular(),
		EEnabled: r.Enabled(),
	}
}

type combinerModel struct {
	Idx         bson.ObjectId `bson:"_id"`
	SServerName string        `bson:"sever_name"`
	PPath       string        `bson:"path"`
	FField      string        `bson:"field"`
	MMethod     string        `bson:"method"`
}

func (c *combinerModel) ID() string         { return c.Idx.Hex() }
func (c *combinerModel) SetID(id string)    { c.Idx = bson.ObjectIdHex(id) }
func (c *combinerModel) ServerName() string { return c.SServerName }
func (c *combinerModel) Path() string       { return c.PPath }
func (c *combinerModel) Field() string      { return c.FField }
func (c *combinerModel) Method() string     { return c.MMethod }
func loadCombinerModelFromCombiner(r rule.Combiner) *combinerModel {
	return &combinerModel{
		Idx:         bson.ObjectIdHex(r.ID()),
		SServerName: r.ServerName(),
		PPath:       r.Path(),
		FField:      r.Field(),
		MMethod:     r.Method(),
	}
}

type pathRulerModel struct {
	Idx             bson.ObjectId    `bson:"_id"`
	PPath           string           `bson:"path"`
	MMethod         string           `bson:"method"`
	SServerName     string           `bson:"server_name"`
	RRewritePath    string           `bson:"rewrite_path"`
	NNeedCombine    bool             `bson:"need_combine"`
	CCombineReqCfgs []*combinerModel `bson:"combine_req_cfgs"`
}

func (c *pathRulerModel) ID() string          { return c.Idx.Hex() }
func (c *pathRulerModel) SetID(id string)     { c.Idx = bson.ObjectIdHex(id) }
func (c *pathRulerModel) Path() string        { return c.PPath }
func (c *pathRulerModel) Method() string      { return c.MMethod }
func (c *pathRulerModel) ServerName() string  { return c.SServerName }
func (c *pathRulerModel) RewritePath() string { return c.RRewritePath }
func (c *pathRulerModel) NeedCombine() bool   { return c.NNeedCombine }
func (c *pathRulerModel) CombineReqCfgs() []rule.Combiner {
	rules := make([]rule.Combiner, len(c.CCombineReqCfgs))
	for idx, r := range c.CCombineReqCfgs {
		rules[idx] = r
	}
	return rules
}
func loadPathRulerModelFromPathRuler(r rule.PathRuler) *pathRulerModel {
	combReqs := make([]*combinerModel, len(r.CombineReqCfgs()))
	for idx, r := range r.CombineReqCfgs() {
		combReqs[idx] = loadCombinerModelFromCombiner(r)
	}

	return &pathRulerModel{
		Idx:             bson.ObjectIdHex(r.ID()),
		PPath:           r.Path(),
		MMethod:         r.Method(),
		SServerName:     r.ServerName(),
		RRewritePath:    r.RewritePath(),
		NNeedCombine:    r.NeedCombine(),
		CCombineReqCfgs: combReqs,
	}
}

type serverRulerModel struct {
	Idx              bson.ObjectId `bson:"_id"`
	PPrefix          string        `bson:"prefix"`
	SServerName      string        `bson:"server_name"`
	NNeedStripPrefix bool          `bson:"need_strip_prefix"`
}

func (c *serverRulerModel) ID() string            { return c.Idx.Hex() }
func (c *serverRulerModel) SetID(id string)       { c.Idx = bson.ObjectIdHex(id) }
func (c *serverRulerModel) Prefix() string        { return c.PPrefix }
func (c *serverRulerModel) ServerName() string    { return c.SServerName }
func (c *serverRulerModel) NeedStripPrefix() bool { return c.NNeedStripPrefix }
func loadServerRulerModelFromServerRuler(r rule.ServerRuler) *serverRulerModel {
	return &serverRulerModel{
		Idx:              bson.ObjectIdHex(r.ID()),
		PPrefix:          r.Prefix(),
		SServerName:      r.ServerName(),
		NNeedStripPrefix: r.NeedStripPrefix(),
	}
}

type reverseServerModel struct {
	Idx    bson.ObjectId `bson:"_id"`
	NName  string        `bson:"name"`
	AAddr  string        `bson:"addr"`
	WW     int           `bson:"w"`
	GGroup string        `bson:"group"`
}

func (c *reverseServerModel) ID() string      { return c.Idx.Hex() }
func (c *reverseServerModel) SetID(id string) { c.Idx = bson.ObjectIdHex(id) }
func (c *reverseServerModel) Name() string    { return c.NName }
func (c *reverseServerModel) Addr() string    { return c.AAddr }
func (c *reverseServerModel) W() int          { return c.WW }
func (c *reverseServerModel) Group() string   { return c.GGroup }
func loadReverseServerModelFromReverseServer(r rule.ReverseServer) *reverseServerModel {
	return &reverseServerModel{
		Idx:    bson.ObjectIdHex(r.ID()),
		NName:  r.Name(),
		AAddr:  r.Addr(),
		WW:     r.W(),
		GGroup: r.Group(),
	}
}
