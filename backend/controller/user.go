package controller

import (
	"github.com/NubeIO/nubeio-rubix-lib-helpers-go/pkg/uuid"
	"github.com/NubeIO/rubix-updater/model"
	"github.com/NubeIO/rubix-updater/model/schema"
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"net/http"
	"time"
)

func getUserBody(ctx *gin.Context) (dto *model.User, err error) {
	err = ctx.ShouldBindJSON(&dto)
	return dto, err
}

func (base *Controller) UsersSchema(ctx *gin.Context) {
	reposeHandler(schema.GetUserSchema(), err, ctx)
}

func (base *Controller) GetUser(c *gin.Context) {
	host, err := base.DB.GetUser(c.Params.ByName("id"))
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(host, err, c)
}

func (base *Controller) GetUsers(c *gin.Context) {
	hosts, err := base.DB.GetUsers()
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(hosts, err, c)
}

func (base *Controller) UpdateUser(c *gin.Context) {
	body, _ := getUserBody(c)
	host, err := base.DB.UpdateUser(c.Params.ByName("id"), body)
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(host, err, c)
}

func (base *Controller) DeleteUser(c *gin.Context) {
	q, err := base.DB.DeleteUser(c.Params.ByName("id"))
	if err != nil {
		reposeHandler(nil, err, c)
	} else {
		reposeHandler(q, err, c)
	}
}

func (base *Controller) DropUsers(c *gin.Context) {
	host, err := base.DB.DropUsers()
	if err != nil {
		reposeHandler(nil, err, c)
		return
	}
	reposeHandler(host, err, c)
}

func (base *Controller) Login(c *gin.Context) (interface{}, error) {
	var loginVals model.LoginUser
	var user model.User
	if err := c.ShouldBindJSON(&loginVals); err != nil {
		return "", jwt.ErrMissingLoginValues
	}
	email := loginVals.Email
	if result := base.DB.DB.Where("email = ?", email).First(&user); result.Error != nil {
		return "", jwt.ErrFailedAuthentication
	} else {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(loginVals.Password)); err != nil {
			return "", jwt.ErrFailedAuthentication
		}
		return user, nil
	}
}

const charset string = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const length int = 8

func GenerateUID() string {
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (base *Controller) AddUser(c *gin.Context) {
	var user model.User
	var newUser model.NewUser

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		panic(err.Error())
	}

	if result := base.DB.DB.Where("email = ?", newUser.Email).First(&user); result.Error != nil {
		// TODO: Differentiate between server error and user user not found error
		hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		user = model.User{Username: newUser.Username, Email: newUser.Email, Hash: string(hash), UID: GenerateUID()}
		user.ID, _ = uuid.MakeUUID()
		if err := base.DB.DB.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, &user)
		return
	} else {
		c.JSON(http.StatusConflict, gin.H{"Error": "User already registered"})
		return
	}
}
