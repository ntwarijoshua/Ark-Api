package services

import (
	"github.com/astaxie/beego/plugins/auth"
	"github.com/astaxie/beego/orm"
	"golang.org/x/crypto/bcrypt"
	"ark-api/models"
	"github.com/astaxie/beego/context"
	"fmt"
)

func Authenticate(ctx *context.Context) bool{
	authenticator := auth.NewBasicAuthenticator(SecretAuth,"Basic")
	if(authenticator != nil){
		fmt.Println(&authenticator)
		ctx.Output.Status = 401
		ctx.Output.JSON(map[string]string{"Error":"Unauthorized"},true,true)
		return false
	}
	ctx.Input.SetData("currentUser",authenticator)
	return  true
}

func SecretAuth(username,password string)bool{
	user := models.User{
		Email:username,
	}
	o := orm.NewOrm()
	err := o.Read(&user,"Email")
	if(err == orm.ErrNoRows){
		return false
	}
	return compareHashes(user.Password,password)
}

func compareHashes(val1,val2 string)bool{
	err := bcrypt.CompareHashAndPassword([]byte(val1),[]byte(val2))
	return err == nil
}
