package zabbix_test

import (
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	zapi "github.com/lavrenko/go-zabbix-api"
)

var (
	_host string
	_api  *zapi.API
)

func init() {
	rand.Seed(time.Now().UnixNano())

	var err error
	_host, err = os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	_host += "-testing"

	if os.Getenv("TEST_ZABBIX_URL") == "" {
		log.Fatal("Set environment variables TEST_ZABBIX_URL (and optionally TEST_ZABBIX_USER and TEST_ZABBIX_PASSWORD)")
	}
}

func getHost() string {
	return _host
}

func getAPI(t *testing.T) *zapi.API {
	if _api != nil {
		return _api
	}

	url, user, password := os.Getenv("TEST_ZABBIX_URL"), os.Getenv("TEST_ZABBIX_USER"), os.Getenv("TEST_ZABBIX_PASSWORD")

	// Zabbix client connection configuration
	var c zapi.Config
	c.Url = url

	_api = zapi.NewAPI(c)
	_api.SetClient(http.DefaultClient)
	v := os.Getenv("TEST_ZABBIX_VERBOSE")
	if v != "" && v != "0" {
		_api.Logger = log.New(os.Stderr, "[zabbix] ", 0)
	}

	if user != "" {
		auth, err := _api.Login(user, password)
		if err != nil {
			t.Fatal(err)
		}
		if auth == "" {
			t.Fatal("Login failed")
		}
	}

	return _api
}

func TestBadCalls(t *testing.T) {
	api := getAPI(t)
	res, err := api.Call("", nil)
	if err != nil {
		t.Fatal(err)
	}
	if res.Error.Code != -32602 {
		t.Errorf("Expected code -32602, got %s", res.Error)
	}
}

func TestVersion(t *testing.T) {
	api := getAPI(t)
	v, err := api.Version()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Zabbix version %s", v)
	if !regexp.MustCompile(`^\d\.\d\.\d+$`).MatchString(v) {
		t.Errorf("Unexpected version: %s", v)
	}
}
