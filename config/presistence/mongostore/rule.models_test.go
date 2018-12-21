package mongostore

import (
	"testing"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func TestModels(t *testing.T) {
	var (
		modelNocacher = &nocacherModel{
			Idx:      bson.NewObjectId(),
			Regexp:   "askdjalks",
			EEnabled: true,
		}
		modelPathRuler = &pathRulerModel{
			Idx:          bson.NewObjectId(),
			PPath:        "/path",
			MMethod:      "GET",
			SServerName:  "srv1",
			RRewritePath: "/path_2",
			NNeedCombine: true,
			CCombineReqCfgs: []*combinerModel{
				&combinerModel{
					Idx:         bson.NewObjectId(),
					SServerName: "srv1",
					PPath:       "/pp/1",
					FField:      "field1",
					MMethod:     "GET",
				},
				&combinerModel{
					Idx:         bson.NewObjectId(),
					SServerName: "srv1",
					PPath:       "/pp/2",
					FField:      "field1",
					MMethod:     "GET",
				},
			},
		}
		modelServerRuler = &serverRulerModel{
			Idx:              bson.NewObjectId(),
			PPrefix:          "/prefix",
			SServerName:      "srv1",
			NNeedStripPrefix: true,
		}
		modelReverseServer = &reverseServerModel{
			Idx:   bson.NewObjectId(),
			NName: "srv1",
			AAddr: "http://localhost:8089",
			WW:    5,
		}
	)

	session, err := mgo.Dial("mongodb://localhost:27017")
	defer session.Close()
	if err != nil {
		t.Error(err)
	}

	db := session.DB("gateway")
	if err := db.C("plugin_cache").Insert(modelNocacher); err != nil {
		t.Error(err)
	}
	if err := db.C("plugin_proxy_path").Insert(modelPathRuler); err != nil {
		t.Error(err)
	}
	if err := db.C("plugin_proxy_server").Insert(modelServerRuler); err != nil {
		t.Error(err)
	}
	if err := db.C("plugin_proxy_reverse").Insert(modelReverseServer); err != nil {
		t.Error(err)
	}
}
