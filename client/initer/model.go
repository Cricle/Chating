package initer

type Config struct {
	Address string
	Name    string `json:"-"`
	Pwd     string `json:"-"`
}
