package v1

import (
	"EXAM3/api-gateway/api/handlers/v1/tokens"
	"EXAM3/api-gateway/api/model"
	pb "EXAM3/api-gateway/genproto/user_service"
	"EXAM3/api-gateway/pkg/codegen"
	"EXAM3/api-gateway/pkg/etc"
	"EXAM3/api-gateway/pkg/logger"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"google.golang.org/protobuf/encoding/protojson"
)

// REGISTER USER
// @Summary Register User
// @Description Api for Registering
// @Tags User
// @Accept json
// @Produce json
// @Param user body model.RegisterUserRequest true "user"
// @Success 200 {object} model.RegisterUserResponse
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Router /v1/user/register [post]
func (h *handlerV1) RegisterUser(c *gin.Context) {
	var (
		body        model.RegisterUserRequest
		code        string
		jsbpMarshal protojson.MarshalOptions
	)
	jsbpMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
	}

	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	err = body.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorBadRequest,
			Error: "wrong format, correct your email or password format amd try again",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	exists, err := h.serviceManager.MockUserService().CheckField(ctx, &pb.CheckFieldRequest{
		Field: "email",
		Data:  body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
		h.log.Error("failed to check uniqueness: ", logger.Error(err))
		return
	}

	if exists.Status {
		c.JSON(http.StatusConflict, model.ResponseError{
			Code:  ErrorCodeAlreadyExists,
			Error: "this email is already in use",
		})
		h.log.Error("email is already exist in database")
		return
	}

	code = codegen.GenerateCode()
	type PageData struct {
		OTP string
	}
	tpl := template.Must(template.ParseFiles("index.html"))
	data := PageData{
		OTP: code,
	}
	var buf bytes.Buffer
	tpl.Execute(&buf, data)
	htmlContent := buf.Bytes()

	auth := smtp.PlainAuth("", "nodirbekgolang@gmail.com", "ecncwhvfdyvjghux", "smtp.gmail.com")
	err = smtp.SendMail("smtp.gmail.com:587", auth, "nodirbekgolang@gmail.com", []string{body.Email}, []byte("To: "+body.Email+"\r\nSubject: Email verification\r\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"+string(htmlContent)))
	if err != nil {
		log.Fatalf("Error sending otp to email: %v", err)
	}
	log.Println("Email sent successfully")
	body.OTP = code

	byteUser, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
		h.log.Error("failed while marshalling user data")
		return
	}

	if err := h.reds.SetWithTTL(body.Email, string(byteUser), int(time.Second)*300); err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
		h.log.Error("cannot sset redis")
		return
	}

	c.JSON(http.StatusOK, model.RegisterUserResponse{
		Message: "One time verification password sent to your email. Please verify",
	})
}

// Verify User
// @Summary verify user
// @Tags User
// @Description Verify a user with code sent to their email
// @Accept json
// @Product json
// @Param email path string true "email"
// @Param code path string true "code"
// @Success 201 {object} model.UserModel
// @Failure 400 string error models.ResponseError
// @Failure 400 string error models.ResponseError
// @Router /v1/user/verify/{email}/{code} [post]
func (h *handlerV1) Verify(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	userEmail := c.Param("email")
	code := c.Param("code")

	user, err := redis.Bytes(h.reds.Get(userEmail))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorCodeUnauthorized,
			Error: "code is expired, try again.",
		})
		h.log.Error("Code is expired, TTL is over.")
		return
	}

	var respUser model.UserModel
	if err := json.Unmarshal(user, &respUser); err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
		h.log.Error("cannot unmarshal uslaer from redis", logger.Error(err))
		fmt.Println(respUser)
		return
	}

	if respUser.OTP != code {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorCodeInvalidCode,
			Error: "code is incorrect, try again.",
		})
		h.log.Error("verification failed", logger.Error(err))
		return
	}

	respUser.Password, err = etc.HashPassword(respUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
		h.log.Error("cannot hash the password", logger.Error(err))
		return
	}

	h.jwtHandler = tokens.JWTHandler{
		Sub:       respUser.Id,
		Role:      "user",
		SignInKey: h.cfg.SigningKey,
		Log:       h.log,
		Timeout:   h.cfg.AccessTokenTimeout,
	}

	access, refresh, err := h.jwtHandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
		h.log.Error("cannot create access and refresh token", logger.Error(err))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	userResp := model.UserModel{
		Id:          respUser.Id,
		Name:        respUser.Name,
		Age:         respUser.Age,
		Username:    respUser.Username,
		Email:       respUser.Email,
		Password:    respUser.Password,
		AccessToken: access,
	}
	id := uuid.New().String()

	_, err = h.serviceManager.MockUserService().CreateUser(ctx, &pb.User{
		Id:           id,
		Name:         respUser.Name,
		Age:          int64(respUser.Age),
		Username:     respUser.Username,
		Email:        respUser.Email,
		Password:     respUser.Password,
		RefreshToken: refresh,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot create user", logger.Error(err))
		return
	}
	c.JSON(http.StatusOK, userResp)
}

// LOGIN USER
// @Summary Log In User
// @Description Api for Logging in
// @Tags User
// @Accept json
// @Produce json
// @Param email path string true "Email"
// @Param password path string true "Password"
// @Success 200 {object} model.LogInResponse
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Router /v1/user/login/{email}/{password} [post]
func (h *handlerV1) Login(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	email := c.Param("email")
	password := c.Param("password")

	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	resp, err := h.serviceManager.MockUserService().GetUserByEmail(ctx, &pb.Email{
		Email: email,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to get user by email", logger.Error(err))
		return
	}

	if !etc.CompareHashPassword(resp.Password, password) {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorCodeInvalidCode,
			Error: "wrong password",
		})
		return
	}
	h.jwtHandler = tokens.JWTHandler{
		Sub:       resp.Id,
		Role:      "user",
		SignInKey: h.cfg.SigningKey,
		Log:       h.log,
		Timeout:   h.cfg.AccessTokenTimeout,
	}

	access, _, err := h.jwtHandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorBadRequest,
			Error: err.Error(),
		})
		h.log.Error("cannot create access token", logger.Error(err))
		return
	}

	res := model.LogInResponse{
		Message:     "Successfully logged in",
		AccessToken: access,
	}

	c.JSON(http.StatusOK, res)
}

// Generate access token for admin
// @Summary Access token generator for admin
// @Description Access token
// @Tags Admin
// @Accept json
// @Produce json
// @Param username path string true "username"
// @Param password path string true "password"
// @Success 200 {object} model.LogInResponse
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Router /v1/admin/{username}/{password} [get]
func (h *handlerV1) GenerateAccessTokenForAdmin(c *gin.Context) {
	username := c.Param("username")
	password := c.Param("password")
	fmt.Println(username, password)
	if username == "golang" && password == "backend" {
		h.jwtHandler = tokens.JWTHandler{
			Sub:       "golang",
			Role:      "admin",
			SignInKey: h.cfg.SigningKey,
			Log:       h.log,
			Timeout:   h.cfg.AccessTokenTimeout,
		}

		access, _, err := h.jwtHandler.GenerateAuthJWT()
		if err != nil {
			c.JSON(http.StatusBadRequest, model.ResponseError{
				Code:  ErrorBadRequest,
				Error: err.Error(),
			})
			h.log.Error("cannot create access token for admin", logger.Error(err))
			return
		}

		res := model.LogInResponse{
			Message:     "Successfully logged in",
			AccessToken: access,
		}

		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorCodeUnauthorized,
			Error: "invalid username and password",
		})
		h.log.Error("unauthorized request")
		return
	}
}

// CREATE USER
// @Summary Create User
// @Description Api for creating Users
// @Security BearerAuth
// @Tags User
// @Accept json
// @Produce json
// @Param user body model.RegisterUserRequest true "user"
// @Success 200 {object} model.UserModel
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Router /v1/user/create [post]
func (h *handlerV1) CreateUser(c *gin.Context) {
	var (
		body        model.RegisterUserRequest
		jsbpMarshal protojson.MarshalOptions
	)
	jsbpMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
		h.log.Error(err.Error())
		return
	}

	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	exists, err := h.serviceManager.MockUserService().CheckField(ctx, &pb.CheckFieldRequest{
		Field: "email",
		Data:  body.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
		h.log.Error("failed to check uniqueness: ", logger.Error(err))
		return
	}

	if exists.Status {
		c.JSON(http.StatusConflict, model.ResponseError{
			Code:  ErrorCodeAlreadyExists,
			Error: "email is already exist",
		})
		h.log.Error("email is already exist in database")
		return
	}

	id := uuid.New().String()

	h.jwtHandler = tokens.JWTHandler{
		Sub:       id,
		Role:      "user",
		SignInKey: h.cfg.SigningKey,
		Log:       h.log,
		Timeout:   h.cfg.AccessTokenTimeout,
	}
	access, refresh, err := h.jwtHandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorBadRequest,
			Error: err.Error(),
		})
		h.log.Error("cannot create access token", logger.Error(err))
		return
	}

	resp, err := h.serviceManager.MockUserService().CreateUser(ctx, &pb.User{
		Id:           id,
		Name:         body.Name,
		Age:          int64(body.Age),
		Username:     body.Username,
		Email:        body.Email,
		Password:     body.Password,
		RefreshToken: refresh,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, model.UserModel{
		AccessToken: access,
		Id:          resp.Id,
		Name:        resp.Name,
		Age:         int(resp.Age),
		Username:    resp.Username,
		Email:       resp.Email,
		Password:    resp.Password,
	})
}

// UPDATE USER
// @Summary Update User
// @Description Api for updating Users
// @Security BearerAuth
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Param user body model.RegisterUserRequest true "user"
// @Success 200 {object} model.RegisterUserResponse
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Router /v1/user/update/{id} [post]
func (h *handlerV1) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user pb.User

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
	}
	user.Id = id

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	result, err := h.serviceManager.MockUserService().UpdateUserById(ctx, &user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, result)
}

// DELETE USER
// @Summary Delete User
// @Description Api for deleting Users
// @Security BearerAuth
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "Id"
// @Success 200 {object} model.RegisterUserResponse
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Router /v1/user/delete/{id} [delete]
func (h *handlerV1) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	_, err := h.serviceManager.MockUserService().DeleteUser(ctx, &pb.UserId{
		UserId: id,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ResponseError{
			Code:  ErrorCodeInternalServerError,
			Error: err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted successfully",
	})

}

// GET ALL USERS
// @Summary Get All Users
// @Description Api to get all Users
// @Security BearerAuth
// @Tags User
// @Accept json
// @Produce json
// @Param page path string true "page"
// @Param limit path string true "limit"
// @Success 200 {object} model.GetAllUserResponse
// @Failure 400 {object} model.ResponseError
// @Failure 500 {object} model.ResponseError
// @Router /v1/user/getall/{page}/{limit} [get]
func (h *handlerV1) GetAllUsers(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()
	page := c.Param("page")
	fmt.Println(page)
	pageToInt, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorBadRequest,
			Error: err.Error(),
		})
		h.log.Error("cannot parse page query param", logger.Error(err))
		return
	}

	limit := c.Param("limit")
	fmt.Println(limit)
	LimitToInt, err := strconv.Atoi(limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ResponseError{
			Code:  ErrorBadRequest,
			Error: err.Error(),
		})
		h.log.Error("cannot parse limit query param", logger.Error(err))
		return
	}

	response, err := h.serviceManager.MockUserService().ListUser(ctx, &pb.GetAllUserRequest{
		Page:  int64(pageToInt),
		Limit: int64(LimitToInt),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("cannot get all users", logger.Error(err))
		return
	}

	c.JSON(http.StatusOK, response)
}

// func MessageProcess(c *gin.Context) {
// 	freelancer_id := c.Param("fre_id")
// 	employee_id := c.Param("emp_id")

// }
