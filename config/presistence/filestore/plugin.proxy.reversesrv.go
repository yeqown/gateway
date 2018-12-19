package filestore

import (
	"fmt"
	"log"
	"strings"

	"github.com/yeqown/gateway/config/rule"
	"github.com/yeqown/gateway/utils"
)

// NewReverseServer func
func (f *JSONFileStore) NewReverseServer(group string, s rule.ReverseServer) error {
	_, ok := f.cfg.ProxyReverseServers[group]
	if !ok {
		return ErrGroupNotExist
	}

	cnt := f.ReverseServerGroupPageCount(group)
	idx := fmt.Sprintf("%s#%d", group, cnt)
	s.SetID(idx)
	// update cfg
	f.cfg.ProxyReverseServers[group] = append(f.cfg.ProxyReverseServers[group], s)
	// update f.JSONFile
	cfg := &ReverseServerCfg{
		NName:  s.Name(),
		AAddr:  s.Addr(),
		Weight: s.W(),
		Idx:    idx,
	}
	f.JSONFile.ProxyReverseServers[group] = append(f.JSONFile.ProxyReverseServers[group], cfg)
	f.updatedC <- true
	return nil
}

// DelReverseServer func
func (f *JSONFileStore) DelReverseServer(id string) error {
	slices := strings.Split(id, "#")
	if len(slices) != 2 {
		panic("incorrect id, shoule be like: group#2")
	}
	group, idxs := slices[0], slices[1]
	idx := utils.Atoi(idxs)
	g, ok := f.ProxyReverseServers[group]
	if !ok {
		return ErrGroupNotExist
	}
	length := len(g)

	// del from JSONFile
	f.ProxyReverseServers[group][idx] = f.ProxyReverseServers[group][length-1]
	f.ProxyReverseServers[group][length-1] = nil
	f.ProxyReverseServers[group] = f.ProxyReverseServers[group][:length-1]
	// del cfg
	f.cfg.ProxyReverseServers[group][idx] = f.cfg.ProxyReverseServers[group][length-1]
	f.cfg.ProxyReverseServers[group][length-1] = nil
	f.cfg.ProxyReverseServers[group] = f.cfg.ProxyReverseServers[group][:length-1]
	f.updatedC <- true

	// 重新设置ID修复ID不对应的情况
	// for idx, r := range f.cfg.ProxyReverseServers[group] {
	// 	r.SetID(fmt.Sprintf("%d", idx))
	// }
	return nil
}

// UpdateReverseServer func ...
func (f *JSONFileStore) UpdateReverseServer(
	id string, s rule.ReverseServer) error {
	slices := strings.Split(id, "#")
	if len(slices) != 2 {
		panic("incorrect id, shoule be like: group#2")
	}
	group, idxs := slices[0], slices[1]
	idx := utils.Atoi(idxs)
	_, ok := f.ProxyReverseServers[group]
	if !ok {
		panic(ErrGroupNotExist)
	}
	// update JSONFile
	f.ProxyReverseServers[group][idx].AAddr = s.Addr()
	f.ProxyReverseServers[group][idx].NName = s.Name()
	f.ProxyReverseServers[group][idx].Weight = s.W()
	// update cfg
	f.cfg.ProxyReverseServers[group][idx] = s
	f.updatedC <- true
	return nil
}

// DelReverseServerGroup func
func (f *JSONFileStore) DelReverseServerGroup(group string) error {
	delete(f.ProxyReverseServers, group)
	delete(f.cfg.ProxyReverseServers, group)
	f.updatedC <- true
	return nil
}

// ReverseServerByID ...
func (f *JSONFileStore) ReverseServerByID(group string, id string) rule.ReverseServer {
	idx := utils.Atoi(id)
	return f.cfg.ProxyReverseServers[group][idx]
}

// ReverseServerGroup func
func (f *JSONFileStore) ReverseServerGroup(group string, start, end int) []rule.ReverseServer {
	g, ok := f.cfg.ProxyReverseServers[group]
	if !ok {
		log.Printf("no such group: %s\n", group)
		return nil
	}

	cnt := len(g)
	if start > cnt {
		return nil
	}
	if end > cnt {
		end = cnt
	}
	return g[start:end]
}

// ReverseServerGroupCount func
func (f *JSONFileStore) ReverseServerGroupCount() int {
	return len(f.cfg.ProxyReverseServers)
}

// ReverseServerGroupPageCount func
func (f *JSONFileStore) ReverseServerGroupPageCount(group string) int {
	g, ok := f.cfg.ProxyReverseServers[group]
	if !ok {
		log.Printf("no such group: %s\n", group)
		return 0
	}
	return len(g)
}
