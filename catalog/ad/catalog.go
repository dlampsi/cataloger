package ad

import (
	"fmt"

	"cataloger/client"
)

type Catalog struct {
	cl         *client.Client
	Attributes *Attributes
}

type Attributes struct {
	SearchBase      string
	SearchAttribute string
}

func NewCatalog(clientConfig *client.Config, attr *Attributes) (*Catalog, error) {
	catalog := Catalog{}
	catalog.Attributes = attr
	cl, err := client.New(clientConfig)
	if err != nil {
		return nil, fmt.Errorf("can't init client: %s", err.Error())
	}
	catalog.cl = cl
	return &catalog, nil
}

func CheckConnection(clientConfig *client.Config) error {
	_, err := client.New(clientConfig)
	return err
}
