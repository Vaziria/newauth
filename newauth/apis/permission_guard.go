package apis

import (
	"log"
	"net/http"
	"strconv"

	"github.com/PDC-Repository/newauth/newauth/models"
	"gorm.io/gorm"
)

type Guard struct {
	perm *models.Permission
	DB   *gorm.DB
}

type ErrorResponse struct{}

func (guard *Guard) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("USER_ID")

		errresp := &ApiResponse{
			Code:    "err_permission",
			Message: "dont have access",
		}

		if err != nil {
			SetError(w, errresp)

			return
		} else {
			userId, err := strconv.Atoi(cookie.Value)
			if err != nil {
				SetError(w, errresp)
				return
			}

			var userPermission models.UserPermission

			err = guard.DB.First(&userPermission, userId).Error
			if err != nil {
				SetError(w, errresp)
			}

			next.ServeHTTP(w, r)
		}

	})
}

func NewGuard(db *gorm.DB, key string) *Guard {
	log.Println("initiating guard ", key)

	var permission models.Permission

	err := db.FirstOrCreate(&permission, models.Permission{Key: key}).Error

	if err != nil {
		panic("cannot find or create permission")
	}

	return &Guard{
		DB:   db,
		perm: &permission,
	}
}
