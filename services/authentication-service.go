package services

import (
	"ark-api/models"
	"fmt"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/auth"
	"golang.org/x/crypto/bcrypt"
)

func Authenticate(ctx *context.Context) bool {
	authenticator := auth.NewBasicAuthenticator(SecretAuth, "Basic")
	if authenticator != nil {
		fmt.Println(&authenticator)
		ctx.Output.Status = 401
		ctx.Output.JSON(map[string]string{"Error": "Unauthorized"}, true, true)
		return false
	}
	ctx.Input.SetData("currentUser", authenticator)
	return true
}

func SecretAuth(username, password string) bool {
	user := models.User{
		Email: username,
	}
	o := orm.NewOrm()
	err := o.Read(&user, "Email")
	if err == orm.ErrNoRows {
		return false
	}
	return compareHashes(user.Password, password)
}

func compareHashes(val1, val2 string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(val1), []byte(val2))
	return err == nil
}
