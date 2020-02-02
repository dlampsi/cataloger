package client

import "testing"

func Test_Connect(t *testing.T) {
	f := func(name string, cfg *Config, ok bool) {
		t.Helper()
		t.Run(name, func(t *testing.T) {
			cl := Client{
				config: cfg,
			}
			_, err := cl.Connect()
			if err != nil && ok {
				t.Fatalf("unexpected error from Connect: %s", err.Error())
			}
			if err == nil && !ok {
				t.Fatal("expecting error from Connect")
			}
		})
	}
	f("BadURI", &Config{
		Host:         "fale.host",
		Port:         636,
		SSL:          true,
		Insecure:     true,
		BindDN:       "",
		BindPassword: "",
	}, false)
}
