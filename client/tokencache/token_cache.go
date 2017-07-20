package tokencache

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/coreos/go-oidc"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"

	tbnflag "github.com/turbinelabs/nonstdlib/flag"
	"github.com/turbinelabs/nonstdlib/ptr"
	tbntime "github.com/turbinelabs/nonstdlib/time"
)

// nowOffest represents how much we pad time.Now for network travel time
// when computing token expiry.
const nowOffset = -2 * time.Second

// PathFromFlags allow passing an accessor to a cache path around that hasn't
// been resolved yet and may be set by command line flags.
type PathFromFlags interface {
	// CachePath returns either the configured cache path or Default() if it
	// hasn't been set.
	CachePath() string

	// Default returns some default that should be used if the path isn't set.
	Default() string
}

type pathFromFlags struct {
	path   string
	prefix string
}

func (cp *pathFromFlags) Default() string {
	home := os.Getenv("HOME")
	pf := "tbn"
	if cp.prefix != "" {
		pf = cp.prefix
	}
	return filepath.Join(home, fmt.Sprintf(".%s-auth-cache", pf))
}

func (cp *pathFromFlags) CachePath() string {
	if cp.path == "" {
		return cp.Default()
	}
	return cp.path
}

// NewStaticPath returns a PathFromFlags that always returns the path p.
func NewStaticPath(p string) PathFromFlags {
	return &pathFromFlags{p, ""}
}

// NewPathFromFlags updates flagSet with the flags necessary to configure a
// TokenCache file path and constructs a default file with the specified pfix.
func NewPathFromFlags(pfix string, flagSet tbnflag.FlagSet) PathFromFlags {
	cp := &pathFromFlags{"", pfix}
	flagSet.StringVar(
		&cp.path,
		"token-cache",
		cp.Default(),
		"Turbine Labs API auth tokens are cached between logins in this file",
	)

	return cp
}

// TokenCache stores cached data about a login endpoint and an auth token previously
// retrieved through user/pass authentication. It is sufficient to take subsequent
// actions without needing to reauthenticate with the Turbine Labs API.
type TokenCache struct {
	ClientID    string
	ClientKey   string
	ProviderURL string
	Username    string

	ExpiresAt        *time.Time
	RefreshExpiresAt *time.Time

	Token *oauth2.Token

	timeSource tbntime.Source
}

type OAuth2Token interface {
	Token() *oauth2.Token
	ExpiresIn() (int, bool)
	RefreshExpiresIn() (int, bool)
}

type otw struct {
	tkn *oauth2.Token
}

type extra interface {
	Extra(string) interface{}
}

func intExtra(e extra, k string) (int, bool) {
	if e == nil {
		return 0, false
	}

	v := e.Extra(k)
	if v == nil {
		return 0, false
	}

	processUint := func(ui uint64) (int, bool) {
		i := int(ui)
		if i < 0 {
			return 0, false
		}
		return i, true
	}

	switch t := v.(type) {
	case float64:
		return int(t), true
	case float32:
		return int(t), true
	case int:
		return t, true
	case int8:
		return int(t), true
	case int16:
		return int(t), true
	case int32:
		return int(t), true
	case int64:
		return int(t), true
	case uint:
		return processUint(uint64(t))
	case uint8:
		return processUint(uint64(t))
	case uint16:
		return processUint(uint64(t))
	case uint32:
		return processUint(uint64(t))
	case uint64:
		return processUint(t)
	default:
		return 0, false
	}
}

func (t otw) Token() *oauth2.Token { return t.tkn }

func (t otw) ExpiresIn() (int, bool) {
	return intExtra(t.tkn, "expires_in")
}

func (t otw) RefreshExpiresIn() (int, bool) {
	return intExtra(t.tkn, "refresh_expires_in")
}

func WrapOAuth2Token(ot *oauth2.Token) OAuth2Token {
	return otw{ot}
}

// SetToken processes an auth token from the Turbine Labs API and extracts
// expiration information saving all this into the TokenCache struct.
func (tc *TokenCache) SetToken(tkn OAuth2Token) {
	if tkn == nil || tkn.Token() == nil {
		tc.ExpiresAt = nil
		tc.RefreshExpiresAt = nil
		tc.Token = nil
		return
	}

	now := tc.timeSource.Now().UTC().Add(nowOffset)
	tokenExpireSec, ok := tkn.ExpiresIn()
	if ok {
		tc.ExpiresAt = ptr.Time(now.Add(time.Duration(tokenExpireSec) * time.Second))
	} else {
		tc.ExpiresAt = &tkn.Token().Expiry
	}

	refreshExpireSec, ok := tkn.RefreshExpiresIn()
	if ok {
		tc.RefreshExpiresAt = ptr.Time(now.Add(time.Duration(refreshExpireSec) * time.Second))
	} else {
		tc.RefreshExpiresAt = nil
	}
}

// ToOAuthConfig constructs an OAuth2.0 config with the stored login endpoint
// and client info.
func ToOAuthConfig(tc TokenCache) (oauth2.Config, error) {
	provider, err := oidc.NewProvider(context.Background(), tc.ProviderURL)
	if err != nil {
		return oauth2.Config{}, fmt.Errorf("unable to construct OIDC provider: %v", err)
	}

	return oauth2.Config{
		ClientID:     tc.ClientID,
		ClientSecret: tc.ClientKey,
		Endpoint:     provider.Endpoint(),
	}, nil
}

// Expired checks the expiration time of the access token and the refresh
// token to see if the cached auth info can be used or if the user will
// need to login again.
func (tc TokenCache) Expired() bool {
	if tc.Token == nil {
		return true
	}

	now := tc.timeSource.Now().Add(2 * time.Second)

	// maybe the initial token is cool
	if tc.ExpiresAt == nil {
		return false
	}

	if tc.ExpiresAt.After(now) {
		return false
	}

	// the initial token is expired but what about the refresh token
	if tc.RefreshExpiresAt == nil {
		return false
	}

	// both the initial and refresh token are expired
	return tc.RefreshExpiresAt.Before(now)
}

// NewFromFile constructs a new TokenCache populated with information that was
// saved in the file at path.
func NewFromFile(path string) (TokenCache, error) {
	tc := TokenCache{timeSource: tbntime.NewSource()}

	b, err := ioutil.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return TokenCache{}, err
	}

	if len(b) != 0 {
		err = json.Unmarshal(b, &tc)
		if err != nil {
			return TokenCache{}, err
		}
	}

	return tc, nil
}

// Save attempts to write the TokenCache information to a file at path. If the
// file does not exist it will be created with 0600 file permissions (user only
// read/write).
func (tc TokenCache) Save(path string) error {
	c, err := json.MarshalIndent(tc, "", "  ")
	if err != nil {
		return fmt.Errorf("unable to marshal TokenCache: %v", err)
	}

	err = ioutil.WriteFile(path, c, 0600)
	if err != nil {
		return fmt.Errorf("failed to save TokenCache to %v: %v", path, err)
	}

	return nil
}
