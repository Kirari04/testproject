package util

import (
	"fmt"
	"os"
	"testproject/internal/app"
	"testproject/internal/m"
	"testproject/internal/t"
)

func GenerateProxyConfig(s t.Server) error {
	tx := s.DB().Begin()
	defaultsCfg := "defaults\n  timeout client 1m\n  timeout server 1m\n  timeout connect 1m\n\nbackend no-match\n  mode http\n  http-request deny deny_status 400"
	// upRatePeersTblRef := "peerscfg/uploadrate"
	// downRatePeersTblRef := "peerscfg/downloadrate"
	// peersCfg := "\n\npeers peerscfg\n  peer hapee 127.0.0.1:10000" +
	// 	"\n  table uploadrate type ip size 1m expire 3600s store bytes_in_rate(1s)" +
	// 	"\n  table downloadrate type ip size 1m expire 3600s store bytes_out_rate(1s)"

	frontendCfg := ``
	var ports []int
	if err := tx.Model(&m.Frontend{}).Group("port").Pluck("port", &ports).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, port := range ports {
		var frontends []m.Frontend
		if err := tx.Model(&m.Frontend{}).
			Preload("Backends").
			Where("port = ?", port).
			Find(&frontends).Error; err != nil {
			tx.Rollback()
			return err
		}
		frontendName := frontendName(port)
		// add port listener
		frontendCfg += fmt.Sprintf("\n\nfrontend %s\n  bind :%d\n  timeout client 1m",
			frontendName, port,
		)

		// add default bandwith limits
		// frontendCfg += fmt.Sprintf("\n  filter bwlim-in bwlimit_in_default default-limit 10000000 default-period 1s key src table %s", upRatePeersTblRef)
		// frontendCfg += fmt.Sprintf("\n  filter bwlim-out bwlimit_out_default default-limit 10000000 default-period 1s key src table %s", downRatePeersTblRef)

		// match domains with acls to backends
		for _, frontend := range frontends {
			aclFrontendName := fmt.Sprintf("ACL_%d", frontend.ID)
			matchAclDomain := fmt.Sprintf("%s:%d", frontend.Domain, frontend.Port)
			if frontend.Port == 80 {
				matchAclDomain = frontend.Domain
			}
			frontendCfg += fmt.Sprintf("\n  acl %s hdr(host) -i %s", aclFrontendName, matchAclDomain)
		}
		for _, frontend := range frontends {
			aclFrontendName := fmt.Sprintf("ACL_%d", frontend.ID)
			bwLimitInName := fmt.Sprintf("bwlimit_in_%d", frontend.ID)
			bwLimitOutName := fmt.Sprintf("bwlimit_out_%d", frontend.ID)

			// match bandwith limits with acls
			if frontend.DefBwInLimit > 0 {
				frontendCfg += fmt.Sprintf("\n  filter bwlim-in %s default-limit %d default-period %ds",
					bwLimitInName,
					frontend.DefBwInLimit*frontend.DefBwInLimitUnit,
					frontend.DefBwInPeriod,
				)
				frontendCfg += fmt.Sprintf("\n  http-request set-bandwidth-limit %s if %s",
					bwLimitInName,
					aclFrontendName,
				)
			}
			if frontend.DefBwOutLimit > 0 {
				frontendCfg += fmt.Sprintf("\n  filter bwlim-out %s default-limit %d default-period %ds",
					bwLimitOutName,
					frontend.DefBwOutLimit*frontend.DefBwOutLimitUnit,
					frontend.DefBwOutPeriod,
				)
				frontendCfg += fmt.Sprintf("\n  http-request set-bandwidth-limit %s if %s",
					bwLimitOutName,
					aclFrontendName,
				)
			}

		}
		// add backends based on acls
		for _, frontend := range frontends {
			aclFrontendName := fmt.Sprintf("ACL_%d", frontend.ID)
			backendName := backendName(frontend)
			// match backends with acls
			frontendCfg += fmt.Sprintf("\n  use_backend %s if %s", backendName, aclFrontendName)
		}
		// add default backend no-match
		frontendCfg += "\n  default_backend no-match"

	}

	backendCfg := ``
	var frontends []m.Frontend
	if err := tx.Model(&m.Frontend{}).
		Preload("Backends").
		Find(&frontends).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, frontend := range frontends {
		backendName := backendName(frontend)
		// backend base config
		backendCfg += fmt.Sprintf("\n\nbackend %s\n  mode http\n  balance roundrobin", backendName)
		// backend health check
		backendCfg += "\n  option httpchk\n  http-check send meth GET  uri /"

		// backend servers
		for i, backend := range frontend.Backends {
			serverName := serverName(frontend, i)
			backendCfg += fmt.Sprintf("\n  server %s %s check  inter 2s  fall 5  rise 1", serverName, backend.Address)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	// assemble config
	cfg := defaultsCfg + frontendCfg + backendCfg + "\n"
	wasRunning := app.Proxy.IsRunning()
	app.Proxy.Stop()
	// write config
	if err := os.WriteFile("haproxy/haproxy.cfg", []byte(cfg), 0644); err != nil {
		return err
	}

	if wasRunning {
		app.Proxy.Start()
	}

	return nil
}

func frontendName(port int) string {
	return fmt.Sprintf("f_p%d", port)
}

func backendName(frontned m.Frontend) string {
	return fmt.Sprintf("backend_%d", frontned.ID)
}

func serverName(frontend m.Frontend, i int) string {
	return fmt.Sprintf("srv_%d_%d", frontend.ID, i)
}
