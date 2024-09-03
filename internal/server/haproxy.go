package server

import (
	"testproject/internal/t"
)

func (s *Server) HaStart() error {
	return s.Haproxy.Start()
}

func (s *Server) HaStop() error {
	return s.Haproxy.Stop()
}

func (s *Server) HaReload() error {
	return s.Haproxy.Reload()
}

func (s *Server) HaIsRunning() bool {
	return s.Haproxy.IsRunning()
}

func (s *Server) HaStartKeepAlive() {
	s.Haproxy.StartKeepAlive()
}

func (s *Server) HaGenerateConfig(reload bool) error {
	return s.Haproxy.GenerateConfig(reload)
}

func (s *Server) HaCheckConfig() error {
	return s.Haproxy.CheckConfig()
}

func (s *Server) HaGetStats() (*[]t.ProxyStatus, error) {
	return s.Haproxy.GetStats()
}

func (s *Server) HaStopKeepAlive() {
	s.Haproxy.StopKeepAlive()
}

func (s *Server) HaConfigPath() string {
	return s.Haproxy.ConfigPath()
}

func (s *Server) HaGetCrashReasons() t.HaproxyCrashReasonsData {
	return s.Haproxy.GetCrashReasons()
}
