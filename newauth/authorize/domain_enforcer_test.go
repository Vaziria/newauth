package authorize_test

import (
	"testing"

	"github.com/PDC-Repository/newauth/newauth"
	"github.com/PDC-Repository/newauth/newauth/authorize"
)

func TestDomainEnforcer(t *testing.T) {
	db := newauth.NewDatabase()
	authorize.NewDomainEnforcer(db)

}
