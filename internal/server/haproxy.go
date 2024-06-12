package server

import "testproject/internal/t"

func (s *Server) HaStart() {
	s.Haproxy.Start()
}

func (s *Server) HaStop() {
	s.Haproxy.Stop()
}

func (s *Server) HaIsRunning() bool {
	return s.Haproxy.IsRunning()
}

func (s *Server) HaKeepAlive() {
	s.Haproxy.KeepAlive()
}

func (s *Server) HaGenerateConfig() error {
	return s.Haproxy.GenerateConfig()
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

func (s *Server) HaReload() error {
	return s.Haproxy.Reload()
}
