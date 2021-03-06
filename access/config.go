package access

import (
	"errors"
	"strconv"
	"strings"

	"github.com/lovego/config/conf"
	"github.com/lovego/strmap"
	"github.com/lovego/xiaomei/release"
)

type Config struct {
	*conf.Conf
	App, Web *service
	Data     strmap.StrMap
}

func getConfig(env, downAddr string) (Config, error) {
	data := Config{
		Conf: release.AppConf(env),
		App:  newService(`app`, env, downAddr),
		Web:  newService(`web`, env, downAddr),
		Data: release.AppData(env),
	}
	if data.App == nil && data.Web == nil {
		return Config{}, errors.New(`neither app nor web service defined.`)
	}
	return data, nil
}

type service struct {
	*conf.Conf
	svcName  string
	downAddr string
	addrs    []string
}

func newService(svcName, env, downAddr string) *service {
	if release.HasService(svcName, env) {
		return &service{Conf: release.AppConf(env), svcName: svcName, downAddr: downAddr}
	} else {
		return nil
	}
}

func (s *service) Addrs() ([]string, error) {
	if s == nil {
		return nil, nil
	}
	if s.addrs == nil {
		addrs := []string{}
		ports := release.GetService(s.svcName, s.Env).Ports
		for _, node := range s.Nodes() {
			for _, port := range ports {
				upstreamAddr := node.GetServiceAddr() + `:` + strconv.FormatInt(int64(port), 10)
				if s.downAddr != "" && s.downAddr == node.Addr {
					upstreamAddr += " down"
				}
				addrs = append(addrs, upstreamAddr)
			}
		}
		s.addrs = addrs
		if len(addrs) == 0 {
			return nil, errors.New(`no instance defined for: ` + s.svcName)
		}
	}
	return s.addrs, nil
}

func (s *service) Nodes() (nodes []release.Node) {
	if s == nil {
		return nil
	}
	labels := release.GetService(s.svcName, s.Env).Nodes
	for _, node := range release.GetCluster(s.Env).GetNodes(``) {
		if node.Match(labels) {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (s *service) DeployName() string {
	if s == nil {
		return ``
	}
	return release.AppConf(s.Env).DeployName()
}

func (s *service) Domain() string {
	if s == nil {
		return ``
	}
	domain := release.AppConf(s.Env).Domain
	parts := strings.SplitN(domain, `.`, 2)
	if len(parts) == 2 {
		return parts[0] + `-` + s.svcName + `.` + parts[1]
	} else {
		return domain + `-` + s.svcName
	}
}
