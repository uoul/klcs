package domain

import "github.com/uoul/go-common/auth"

type OidcUser struct {
	JwtId                      string   `json:"jti,omitempty"`
	Issuer                     string   `json:"iss,omitempty"`
	Subject                    string   `json:"sub,omitempty"`
	Type                       string   `json:"typ,omitempty"`
	AuthorizedParty            string   `json:"azp,omitempty"`
	SessionState               string   `json:"session_state,omitempty"`
	AuthenticationContextClass string   `acr:"jti,omitempty"`
	AllowedOrigins             []string `json:"allowed-origins,omitempty"`
	RealmAccess                struct {
		Roles []string `json:"roles,omitempty"`
	} `json:"realm_access,omitempty"`
	SessionId string `json:"sid,omitempty"`
	Name      string `json:"name,omitempty"`
	UserName  string `json:"preferred_username,omitempty"`
	FirstName string `json:"given_name,omitempty"`
	LastName  string `json:"family_name,omitempty"`
	roles     map[string]int
}

// HasRole implements iface.IUserIdentity.
func (i *OidcUser) HasRole(role string) bool {
	if i.roles == nil {
		i.roles = map[string]int{}
		for index, r := range i.RealmAccess.Roles {
			i.roles[r] = index
		}
	}
	_, exists := i.roles[role]
	return exists
}

// GetRoles implements IUserIdentity.
func (i *OidcUser) GetRoles() []string {
	return i.RealmAccess.Roles
}

// GetUsername implements IUserIdentity.
func (i *OidcUser) GetUsername() string {
	return i.UserName
}

func NewUserIdentiy() auth.IUserIdentity {
	return &OidcUser{}
}
