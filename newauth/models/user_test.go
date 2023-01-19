package models_test

import (
	"testing"
	"time"

	"github.com/PDC-Repository/newauth/newauth/models"
	"github.com/stretchr/testify/assert"
)

func TestAllowed(t *testing.T) {
	tgl := time.Now().AddDate(0, -1, 0)

	user := &models.User{
		LastReset: tgl,
	}

	assert.True(t, user.AllowedResetPwd())
}
