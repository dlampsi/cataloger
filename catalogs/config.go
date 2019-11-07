package catalogs

type Config struct {
	Host            string
	Port            int
	SSL             bool
	Insecure        bool
	BindDn          string
	BindPass        string
	SearchBase      string
	UserSearchBase  string
	GroupSearchBase string
}
