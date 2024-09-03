package haproxy

import (
	"fmt"
	"io"
	"os"
	"testproject/internal/m"

	"github.com/rs/zerolog/log"
)

func (h *Haproxy) GenerateConfig(reload bool) error {
	log.Info().Msg("generating config")
	tx := h.s.DB().Begin()
	// default config
	defaultsCfg := "defaults" +
		"\n  timeout client 1m" +
		"\n  timeout server 1m" +
		"\n  timeout connect 1m" +
		// require modern certificate standards
		// generated 2024-06-03, Mozilla Guideline v5.7, HAProxy 3.0, OpenSSL 1.1.1k, modern configuration
		// https://ssl-config.mozilla.org/#server=haproxy&version=3.0&config=modern&openssl=1.1.1k&guideline=5.7
		"\n\nglobal"
	if h.s.ENV().Socket {
		defaultsCfg += "\n  stats socket /var/run/haproxy.sock mode 600 expose-fd listeners level user"
	}
	defaultsCfg += "\n  ssl-default-bind-ciphersuites TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256" +
		"\n  ssl-default-bind-options prefer-client-ciphers no-sslv3 no-tlsv10 no-tlsv11 no-tlsv12 no-tls-tickets" +
		"\n  ssl-default-server-ciphersuites TLS_AES_128_GCM_SHA256:TLS_AES_256_GCM_SHA384:TLS_CHACHA20_POLY1305_SHA256" +
		"\n  ssl-default-server-options no-sslv3 no-tlsv10 no-tlsv11 no-tlsv12 no-tls-tickets" +
		"\n  ca-base /etc/ssl/certs" +
		// define default backend
		"\n\nbackend no-match" +
		"\n  mode http" +
		"\n  http-request deny deny_status 410"

	hostname, err := os.Hostname()
	if err != nil {
		log.Warn().Err(err).Msg("failed to get hostname, using localhost")
		hostname = "localhost"
	}
	peersCfg := fmt.Sprintf("\n\npeers peerscfg\n  peer %s 127.0.0.1:10000", hostname) +
		"\n  table stick_out type ipv6 size 1m expire 3600s store bytes_out_rate(1s)" +
		"\n  table stick_in type ipv6 size 1m expire 3600s store bytes_in_rate(1s)"

	frontendCfg := ``

	// prometheus monitoring
	frontendCfg += "\n\nfrontend prometheus" +
		"\n  bind 127.0.0.1:8405" +
		"\n  mode http" +
		"\n  http-request use-service prometheus-exporter if { path /metrics }" +
		"\n  no log"

	// frontend config
	var ports []int
	if err := tx.Model(&m.Frontend{}).Group("port").Pluck("port", &ports).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, port := range ports {
		var frontends []m.Frontend
		if err := tx.Model(&m.Frontend{}).
			Preload("Backends").
			Preload("Aliases").
			Where("port = ?", port).
			Find(&frontends).Error; err != nil {
			tx.Rollback()
			return err
		}
		frontendName := frontendName(port)
		// add port listener
		frontendCfg += fmt.Sprintf("\n\nfrontend %s\n  mode http\n  timeout client 1m",
			frontendName,
		)
		if frontends[0].Https {
			frontendCfg += fmt.Sprintf("\n  bind :%d ssl crt %s/certs/", port, h.s.ENV().WorkDir)
		} else {
			frontendCfg += fmt.Sprintf("\n  bind :%d", port)
		}

		// match domains with acls to backends
		for _, frontend := range frontends {
			aclFrontendName := fmt.Sprintf("ACL_%d", frontend.ID)
			matchAclDomain := fmt.Sprintf("%s:%d", frontend.Domain, frontend.Port)
			if frontend.Port == 80 || frontend.Port == 443 {
				matchAclDomain = frontend.Domain
			}
			aclRule := fmt.Sprintf("\n  acl %s hdr(host) -i %s", aclFrontendName, matchAclDomain)
			for _, alias := range frontend.Aliases {
				matchAclAliasDomain := fmt.Sprintf("%s:%d", alias.Domain, frontend.Port)
				if frontend.Port == 80 {
					matchAclAliasDomain = alias.Domain
				}
				aclRule += fmt.Sprintf(" || hdr(host) -i %s", matchAclAliasDomain)
			}
			frontendCfg += aclRule
		}
		for _, frontend := range frontends {
			aclFrontendName := fmt.Sprintf("ACL_%d", frontend.ID)
			aclFrontendRequestBodyLimitName := fmt.Sprintf("ACL_REQUEST_BODY_LIMIT_%d", frontend.ID)
			bwLimitInName := fmt.Sprintf("bwlimit_in_%d", frontend.ID)
			bwLimitOutName := fmt.Sprintf("bwlimit_out_%d", frontend.ID)
			stickTableInName := "peerscfg/stick_in"
			stickTableOutName := "peerscfg/stick_out"

			// frontendCfg += "\n  option httplog\n  log stdout format raw local0 info"

			frontendCfg += "\n  option httpclose"

			frontendCfg += "\n  capture request header Host len 64"
			frontendCfg += "\n  capture request header User-Agent len 256"
			frontendCfg += "\n  capture request header Content-Length len 10"
			frontendCfg += "\n  capture request header Referer len 256"
			frontendCfg += "\n  capture response header Content-Length len 10"

			// // block other common attacks
			// frontendCfg += "\n  acl forbidden_all path_dir changelog.txt .git .svn .hg .bzr node_modules logs settings.py database.yml dist build web.config"
			// frontendCfg += "\n  acl forbidden_all path_dir .env .env.local .env.production .env.test"
			// frontendCfg += "\n  acl forbidden_all path_dir secret private keys .ssh .htpasswd id_rsa id_dsa key.pem server.key server.pem"
			// frontendCfg += "\n  acl forbidden_all path_dir .htaccess .php_cs.cache test.php i.php info.php phpinfo.php adminer.php config.php composer.json composer.lock vb_test.php lfm.php"
			// // frontendCfg += "\n  acl forbidden_all path_dir "
			// frontendCfg += "\n  acl forbidden_all url_reg -i \\.(bak|backup|tmp|temp|swp|DS_Store|log|zip|tar|tar\\.gz|tgz|exe|dll|so|dylib)$"
			// frontendCfg += "\n  acl forbidden_all url_reg -i .*(\\.|%2e)(\\.|%2e)(%2f|%5c|/|\\\\\\\\)"
			// frontendCfg += "\n  acl forbidden_all url_sub allow_url_include auto_prepend_file php://input"
			// // frontendCfg += "\n  acl forbidden_all url_reg -i "

			// // Block requests with forbidden HTTP methods
			// frontendCfg += "\n  acl forbidden_all method TRACE TRACK"
			// // Block User-Agents known for malicious activities (example: some automated tools or outdated browsers)
			// frontendCfg += "\n  acl forbidden_agents hdr_sub(User-Agent) -i sqlmap nikto w3af acunetix netsparker zgrab jbrofuzz sqlninja sqlpowerinjector sqlsus yersinia commix andsql rips xsser openvas skipfish grabber appscan netspider gobuster"
			// // Check if the 'host' header appears more than once in the request
			// frontendCfg += "\n  acl forbidden_all hdr_cnt(host) gt 1"
			// // Check if the 'content-length' header appears more than once in the request
			// frontendCfg += "\n  acl forbidden_all hdr_cnt(content-length) gt 1"
			// // Check if the 'content-length' header value is less than 0
			// frontendCfg += "\n  acl forbidden_all hdr_val(content-length) lt 0"

			// Block all forbidden requests
			// frontendCfg += "\n  http-request deny deny_status 403 if forbidden_all"

			// match bandwith limits with acls
			if frontend.DefBwInLimit > 0 {
				frontendCfg += fmt.Sprintf(
					"\n  filter bwlim-in %s key src,ipmask(32,64) table %s limit %d min-size 2896",
					bwLimitInName,
					stickTableInName,
					uint((frontend.DefBwInLimit*frontend.DefBwInLimitUnit)/frontend.DefBwInPeriod),
				)
				frontendCfg += fmt.Sprintf("\n  http-request set-bandwidth-limit %s if %s",
					bwLimitInName,
					aclFrontendName,
				)
			}
			if frontend.DefBwOutLimit > 0 {
				frontendCfg += fmt.Sprintf(
					"\n  filter bwlim-out %s key src,ipmask(32,64) table %s limit %d min-size 2896",
					bwLimitOutName,
					stickTableOutName,
					uint((frontend.DefBwOutLimit*frontend.DefBwOutLimitUnit)/frontend.DefBwOutPeriod),
				)
				frontendCfg += fmt.Sprintf("\n  http-request set-bandwidth-limit %s if %s",
					bwLimitOutName,
					aclFrontendName,
				)
			}
			// add http rate limit
			if frontend.DefRateLimit > 0 {
				// soft limit
				peersCfg += fmt.Sprintf(
					"\n  table stick_http_%d type ipv6 size 5m expire %ds store http_req_rate(%ds)",
					frontend.ID,
					frontend.DefRatePeriod,
					frontend.DefRatePeriod,
				)
				// hard limit
				peersCfg += fmt.Sprintf(
					"\n  table stick_http_hard_%d type ipv6 size 5m expire %ds store http_req_rate(%ds)",
					frontend.ID,
					frontend.DefHardRateLimit,
					frontend.DefHardRatePeriod,
				)

				frontendCfg +=
					fmt.Sprintf("\n  http-request track-sc0 src table peerscfg/stick_http_hard_%d", frontend.ID) +
						fmt.Sprintf("\n  http-request track-sc1 src table peerscfg/stick_http_%d", frontend.ID) +
						fmt.Sprintf(
							"\n  http-request silent-drop if { sc_http_req_rate(0,peerscfg/stick_http_hard_%d) gt %d }",
							frontend.ID,
							frontend.DefHardRateLimit,
						) +
						fmt.Sprintf(
							"\n  http-request deny deny_status 429 if { sc_http_req_rate(1,peerscfg/stick_http_%d) gt %d }",
							frontend.ID,
							frontend.DefRateLimit,
						)
			}
			// add request body limit
			if frontend.DefRequestBodyLimit > 0 {
				bodyLimit := frontend.DefRequestBodyLimit * frontend.DefRequestBodyLimitUnit
				frontendCfg += fmt.Sprintf("\n  acl %s req.body_size gt %d", aclFrontendRequestBodyLimitName, bodyLimit)
				frontendCfg += fmt.Sprintf(
					"\n  http-request deny deny_status 413 if %s %s",
					aclFrontendName,
					aclFrontendRequestBodyLimitName,
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

	// backend config
	backendCfg := ``
	var frontends []m.Frontend
	if err := tx.Model(&m.Frontend{}).
		Preload("Backends").
		Preload("Aliases").
		Find(&frontends).Error; err != nil {
		tx.Rollback()
		return err
	}
	for _, frontend := range frontends {
		backendName := backendName(frontend)
		// backend base config
		backendCfg += fmt.Sprintf("\n\nbackend %s\n  mode http\n  balance roundrobin", backendName)
		// Adds X-Forwarded-For header
		// Docs: https://www.haproxy.com/documentation/haproxy-configuration-manual/latest/#option%20forwardfor
		backendCfg += "\n  option forwardfor"
		// Analyze all server responses and block responses with cacheable cookies
		// Docs: https://www.haproxy.com/documentation/haproxy-configuration-manual/latest/#4-option%20checkcache
		backendCfg += "\n  option checkcache"
		// Closes the connection between frontend and server
		// Docs: https://www.haproxy.com/documentation/haproxy-configuration-manual/latest/#option%20httpclose
		backendCfg += "\n  option httpclose"
		// timeout websocket tunnel
		// Docs: https://www.haproxy.com/documentation/haproxy-configuration-tutorials/load-balancing/websocket/
		backendCfg += "\n  timeout tunnel 1h"

		if frontend.HttpCheck != nil &&
			frontend.HttpCheckMethod != nil &&
			frontend.HttpCheckPath != nil &&
			frontend.HttpCheckExpectStatus != nil &&
			*frontend.HttpCheck {
			// backend health check
			backendCfg += fmt.Sprintf(
				"\n  option httpchk\n  http-check send meth %s  uri %s\n  http-check expect status %d",
				*frontend.HttpCheckMethod,
				*frontend.HttpCheckPath,
				*frontend.HttpCheckExpectStatus,
			)
		}
		// backend servers
		for _, backend := range frontend.Backends {
			serverName := serverName(frontend, backend)
			backendCfg += fmt.Sprintf("\n  server %s %s", serverName, backend.Address)
			// https
			if backend.Https {
				backendCfg += " ssl"
				if !backend.HttpsVerify {
					backendCfg += " verify none"
				} else {
					backendCfg += " verify required ca-file ca-certificates.crt"
				}
			}
			if frontend.HttpCheck != nil &&
				frontend.HttpCheckMethod != nil &&
				frontend.HttpCheckPath != nil &&
				frontend.HttpCheckExpectStatus != nil &&
				frontend.HttpCheckInterval != nil &&
				frontend.HttpCheckFailAfter != nil &&
				frontend.HttpCheckRecoverAfter != nil &&
				*frontend.HttpCheck {
				// health check
				backendCfg += fmt.Sprintf(" check inter %ds fall %d rise %d",
					*frontend.HttpCheckInterval,
					*frontend.HttpCheckFailAfter,
					*frontend.HttpCheckRecoverAfter,
				)
			}
			// backendCfg += fmt.Sprintf("\n  server %s %s", serverName, backend.Address)
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	// assemble config
	cfg := defaultsCfg + peersCfg + frontendCfg + backendCfg + "\n"

	// get current config
	var currentCfg []byte
	if file, err := os.Open(h.ConfigPath()); err == nil {
		defer file.Close()
		currentCfg, err = io.ReadAll(file)
		if err != nil {
			return fmt.Errorf("failed to read current config: %w", err)
		}
	}

	// write config
	if err := os.WriteFile(h.ConfigPath(), []byte(cfg), 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	// check if config is valid
	if err := h.CheckConfig(); err != nil {
		if len(currentCfg) > 0 {
			log.Info().Msg("rollback config")
			// copy new config to backup
			if err := os.WriteFile(h.ConfigPath()+".bak", []byte(cfg), 0644); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}
			// rollback config
			if err := os.WriteFile(h.ConfigPath(), []byte(currentCfg), 0644); err != nil {
				return fmt.Errorf("failed to rollback config: %w", err)
			}
		}
		return fmt.Errorf("config is invalid: %w", err)
	}

	if reload {
		if err := h.Reload(); err != nil {
			return err
		}
	}
	return nil
}

func frontendName(port int) string {
	return fmt.Sprintf("f_p%d", port)
}

func backendName(frontned m.Frontend) string {
	return fmt.Sprintf("backend_%d", frontned.ID)
}

func serverName(frontend m.Frontend, backend m.Backend) string {
	return fmt.Sprintf("srv_%d_%d", frontend.ID, backend.ID)
}
