package engine

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/jacoblai/tinygm/models"
	"net/http"
	"strings"

	"github.com/jacoblai/httprouter"
)

// TokenAuth 校验中间件
func (d *DbEngine) TokenAuth(next httprouter.Handle) httprouter.Handle {
	//权限验证
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(auth, "Bearer ") {
			auth = strings.Replace(auth, "Bearer ", "", -1)
		}
		if auth != "" {
			if strings.Contains(auth, ".") {
				rtoken, err := jwt.ParseWithClaims(auth, &models.JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
					return jwtSign, nil
				})
				if err == nil {
					if tk, ok := rtoken.Claims.(*models.JwtClaims); ok && rtoken.Valid && tk.Issuer == "newworld_sys" {
						r.Header.Set("user_id", tk.Id)
						r.Header.Set("sys_role", tk.SysRole)
						next(w, r, ps)
						return
					}
				}
			}
		}
		// // Request Basic Authentication otherwise
		w.Header().Set("WWW-Authenticate", "Bearer realm=Restricted")
		//d.Logger.Debug("[Authenticate unknow]", http.StatusText(http.StatusUnauthorized))
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}
