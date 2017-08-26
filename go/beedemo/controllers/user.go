package controllers

import (
	"encoding/json"
	"time"

	"github.com/golang/glog"

	"practice/go/beedemo/controllers/basecontroller"
	"practice/go/beedemo/models"
	"practice/go/beedemo/module/user"
	"practice/go/beedemo/util"
)

// Operations about Users
type UserController struct {
	basecontroller.Controller
}

// Search search users by name, if filter is emtpy, return all
// @router / [get]
func (c *UserController) Search() {
	// TODO: build query struct
	// just list all user for now
	all, err := user.AllUsers()
	if err != nil {
		c.HandleErr(err)
		return
	}
	c.Success(util.M{
		"users": all,
	})
}

// CreateUser create user
// @router / [post]
func (c *UserController) CreateUser() {
	var u models.User
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &u)
	if err != nil {
		err = util.InvalidArgumentErr(err, "reurest body is not valid")
		c.HandleErr(err)
		return
	}
	now := time.Now()
	u.CreateTime = now
	u.UpdateTime = now

	if !u.IsValid() {
		err = util.InvalidArgumentErr(nil, "user %#v is not valid", u)
		c.HandleErr(err)
		return
	}
	uid, err := user.CreateUser(u)
	if err != nil {
		err = util.InternalError(err, "blablsss%s", "sdf")
		e, ok := err.(*util.BaseErr)
		glog.Infof("e: %v\tok: %t\n%#v", e, ok, err)
		c.HandleErr(err)
		return
	}
	c.Success(util.M{
		"uid": uid,
	})

}

// CheckExistence check whether the given name already exist
// @router /:uname/existance [get]
func (c *UserController) CheckExistence() {
	uname := c.GetString(":uname")
	exist, err := user.CheckExistance(uname)
	if err != nil {
		c.HandleErr(err)
		return
	}
	c.Success(util.M{
		"exist": exist,
	})
}

// // @Title CreateUser
// // @Description create users
// // @Param	body		body 	models.User	true		"body for user content"
// // @Success 200 {int} models.User.Id
// // @Failure 403 body is empty
// // @router / [post]
// func (u *UserController) Post() {
// 	var user models.User
// 	json.Unmarshal(u.Ctx.Input.RequestBody, &user)
// 	uid := models.AddUser(user)
// 	u.Data["json"] = map[string]string{"uid": uid}
// 	u.ServeJSON()
// }

// // @Title GetAll
// // @Description get all Users
// // @Success 200 {object} models.User
// // @router / [get]
// func (u *UserController) GetAll() {
// 	users := models.GetAllUsers()
// 	u.Data["json"] = users
// 	u.ServeJSON()
// }

// // @Title Get
// // @Description get user by uid
// // @Param	uid		path 	string	true		"The key for staticblock"
// // @Success 200 {object} models.User
// // @Failure 403 :uid is empty
// // @router /:uid [get]
// func (u *UserController) Get() {
// 	uid := u.GetString(":uid")
// 	if uid != "" {
// 		user, err := models.GetUser(uid)
// 		if err != nil {
// 			u.Data["json"] = err.Error()
// 		} else {
// 			u.Data["json"] = user
// 		}
// 	}
// 	u.ServeJSON()
// }

// // @Title Update
// // @Description update the user
// // @Param	uid		path 	string	true		"The uid you want to update"
// // @Param	body		body 	models.User	true		"body for user content"
// // @Success 200 {object} models.User
// // @Failure 403 :uid is not int
// // @router /:uid [put]
// func (u *UserController) Put() {
// 	uid := u.GetString(":uid")
// 	if uid != "" {
// 		var user models.User
// 		json.Unmarshal(u.Ctx.Input.RequestBody, &user)
// 		uu, err := models.UpdateUser(uid, &user)
// 		if err != nil {
// 			u.Data["json"] = err.Error()
// 		} else {
// 			u.Data["json"] = uu
// 		}
// 	}
// 	u.ServeJSON()
// }

// // @Title Delete
// // @Description delete the user
// // @Param	uid		path 	string	true		"The uid you want to delete"
// // @Success 200 {string} delete success!
// // @Failure 403 uid is empty
// // @router /:uid [delete]
// func (u *UserController) Delete() {
// 	uid := u.GetString(":uid")
// 	models.DeleteUser(uid)
// 	u.Data["json"] = "delete success!"
// 	u.ServeJSON()
// }

// // @Title Login
// // @Description Logs user into the system
// // @Param	username		query 	string	true		"The username for login"
// // @Param	password		query 	string	true		"The password for login"
// // @Success 200 {string} login success
// // @Failure 403 user not exist
// // @router /login [get]
// func (u *UserController) Login() {
// 	username := u.GetString("username")
// 	password := u.GetString("password")
// 	if models.Login(username, password) {
// 		u.Data["json"] = "login success"
// 	} else {
// 		u.Data["json"] = "user not exist"
// 	}
// 	u.ServeJSON()
// }

// // @Title logout
// // @Description Logs out current logged in user session
// // @Success 200 {string} logout success
// // @router /logout [get]
// func (u *UserController) Logout() {
// 	u.Data["json"] = "logout success"
// 	u.ServeJSON()
// }