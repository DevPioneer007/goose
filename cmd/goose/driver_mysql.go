// +build !no_mysql

package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/go-sql-driver/mysql"
	_ "github.com/ziutek/mymysql/godrv"
)

// normalizeMySQLDSN parses the dsn used with the mysql driver to always have
// the parameter `parseTime` set to true. This allows internal goose logic
// to assume that DATETIME/DATE/TIMESTAMP can be scanned into the time.Time
// type.
func normalizeDBString(driver string, str string, tls bool) string {
	if driver == "mysql" {
		var err error
		str, err = normalizeMySQLDSN(str, tls)
		if err != nil {
			log.Fatalf("failed to normalize MySQL connection string: %v", err)
		}
	}
	return str
}

const tlsConfigKey = "custom"

func normalizeMySQLDSN(dsn string, tls bool) (string, error) {
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return "", err
	}
	config.ParseTime = true
	if tls {
		config.TLSConfig = tlsConfigKey
	}
	return config.FormatDSN(), nil
}

func registerTLSConfig(pemfile string) error {
	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(pemfile)
	if err != nil {
		return err
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return fmt.Errorf("failed to append PEM: %q", pemfile)
	}
	return mysql.RegisterTLSConfig(tlsConfigKey, &tls.Config{
		RootCAs: rootCertPool,
	})
}
