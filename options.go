package session

type Options struct {
	Name        string
	IDLength    int
	UserCookie  bool
	Provider    *ProviderOptions
	AutoStart   bool
	Cookie      *CookieOptions
	GCInterval  int64
	MaxLifeTime int64
}

type CookieOptions struct {
	Domain   string
	Path     string
	Secure   bool
	LifeTime int64
	HttpOnly bool
}

type ProviderOptions struct {
	Adapter string
	Config  interface{}
}
