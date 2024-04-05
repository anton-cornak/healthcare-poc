package main

import "errors"

func validateConfig(cfg *config) error {
	if cfg.port == "" {
		return errors.New("missing port configuration")
	}

	if cfg.dbConn.port == "" {
		return errors.New("missing DB_PORT configuration")
	}

	if cfg.dbConn.host == "" {
		return errors.New("missing DB_HOST configuration")
	}

	if cfg.dbConn.user == "" {
		return errors.New("missing DB_USER configuration")
	}

	if cfg.dbConn.password == "" {
		return errors.New("missing DB_PASS configuration")
	}

	if cfg.dbConn.dbname == "" {
		return errors.New("missing DB_NAME configuration")
	}

	if cfg.dbConn.sslmode == "" {
		return errors.New("missing SSL_MODE configuration")
	}

	return nil
}
