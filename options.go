package session

// Options session manager options
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

// CookieOptions session manager cookie options
type CookieOptions struct {
	Domain   string
	Path     string
	Secure   bool
	LifeTime int64
	HttpOnly bool
}

// ProviderOptions session manager provider options
type ProviderOptions struct {
	Adapter string
	Config  interface{}
}
