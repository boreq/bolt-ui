package commands

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"strings"
	"time"

	"github.com/boreq/bolt-ui/internal/config"
	"github.com/boreq/bolt-ui/internal/wire"
	"github.com/boreq/bolt-ui/logging"
	"github.com/boreq/guinea"
	"github.com/pkg/errors"
)

const (
	nameAddress       = "address"
	nameInsecureCORS  = "insecure-cors"
	nameInsecureToken = "insecure-token"
	nameInsecureTLS   = "insecure-tls"
)

var MainCmd = guinea.Command{
	Run: run,
	Arguments: []guinea.Argument{
		{
			Name:        "database",
			Optional:    false,
			Multiple:    false,
			Description: "Path to the database file",
		},
	},
	Options: []guinea.Option{
		{
			Name:        nameAddress,
			Type:        guinea.String,
			Default:     ":8118",
			Description: `Specifies listening address. Default: :8118`,
		},
		{
			Name:        nameInsecureCORS,
			Type:        guinea.Bool,
			Default:     false,
			Description: "Disables CORS",
		},
		{
			Name:        nameInsecureToken,
			Type:        guinea.Bool,
			Default:     false,
			Description: "Disables token validation",
		},
		{
			Name:        nameInsecureTLS,
			Type:        guinea.Bool,
			Default:     false,
			Description: "Disables serving using TLS",
		},
	},
	ShortDescription: "a web user interface for the Bolt database",
	Description: `
Thanks to bolt-ui you are able to explore a Bolt database using a web
interface. To access the web interface access the address printed out by the
program. Make sure that the address includes the token query parameter.
`,
}

var log = logging.New("main")

func run(c guinea.Context) error {
	conf, err := newConfig(c)
	if err != nil {
		return errors.Wrap(err, "could not create the config")
	}

	if conf.InsecureCORS {
		log.Warn("insecure-cors option enabled")
	}

	if conf.InsecureToken {
		log.Warn("insecure-token option enabled")
	}

	if conf.InsecureTLS {
		log.Warn("insecure-tls option enabled")
	}

	service, err := wire.BuildService(conf)
	if err != nil {
		return errors.Wrap(err, "could not create a service")
	}

	printInfo(conf)

	return service.HTTPServer.Serve()
}

func newConfig(c guinea.Context) (*config.Config, error) {
	conf := &config.Config{
		ServeAddress:  c.Options[nameAddress].Str(),
		DatabaseFile:  c.Arguments[0],
		InsecureCORS:  c.Options[nameInsecureCORS].Bool(),
		InsecureToken: c.Options[nameInsecureToken].Bool(),
		InsecureTLS:   c.Options[nameInsecureTLS].Bool(),
	}

	if !conf.InsecureToken {
		token, err := generateSecureToken()
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate secure token")
		}
		conf.Token = token
	}

	if !conf.InsecureTLS {
		cert, err := generateCertificate()
		if err != nil {
			return nil, errors.Wrap(err, "failed to generate certificate")
		}
		conf.Certificate = cert
	}

	return conf, nil
}

func generateCertificate() (tls.Certificate, error) {
	hosts := []string{
		"localhost",
		"127.0.0.1",
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, errors.Wrap(err, "failed to generate private key")
	}

	keyUsage := x509.KeyUsageDigitalSignature
	notBefore := time.Now()
	notAfter := notBefore.Add(24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return tls.Certificate{}, errors.Wrap(err, "failed to generate serial number")
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, errors.Wrap(err, "failed to create certificate")
	}

	return tls.Certificate{
		Certificate: [][]byte{
			derBytes,
		},
		PrivateKey: priv,
	}, nil
}

func printInfo(conf *config.Config) {
	addr := conf.ServeAddress
	if strings.HasPrefix(addr, ":") {
		addr = "localhost" + addr
	}

	if conf.InsecureTLS {
		addr = "http://" + addr
	} else {
		addr = "https://" + addr
	}

	if !conf.InsecureToken {
		addr = fmt.Sprintf("%s/?token=%s", addr, conf.Token)
	}

	fmt.Printf("You can view database '%s' by clicking on this link:\n", conf.DatabaseFile)
	fmt.Println(addr)
	if !conf.InsecureTLS {
		fmt.Println()
		fmt.Println("For safety check the TLS certificate fingerprint:")
		fmt.Println(certFingerprint(conf.Certificate))
	}
}

func certFingerprint(cert tls.Certificate) string {
	sha1Fingerprint := sha1.Sum(cert.Certificate[0])
	s := fmt.Sprintf("%X", sha1Fingerprint)

	var builder strings.Builder

	for i, r := range s {
		builder.WriteRune(r)

		if i != 0 && i != len(s)-1 && i%2 != 0 {
			builder.WriteRune(':')
		}
	}

	return builder.String()
}

const tokenLength = 32

func generateSecureToken() (string, error) {
	b := make([]byte, tokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", errors.Wrap(err, "failed to read random bytes")
	}
	return hex.EncodeToString(b), nil
}
