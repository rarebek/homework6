package tests

import (
	"EXAM3/api-gateway/api_test/handlers"
	"EXAM3/api-gateway/api_test/storage"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApi(t *testing.T) {
	gin.SetMode(gin.TestMode)
	require.NoError(t, SetupMinimumInstance(""))
	buffer, err := OpenFile("user.json")

	require.NoError(t, err)
	// User Create
	req := NewRequest(http.MethodPost, "/user/create", buffer)
	res := httptest.NewRecorder()
	r := gin.Default()
	r.POST("/user/create", handlers.CreateUser)
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code)

	var user storage.User
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &user))

	require.Equal(t, "nomonovn2@gmail.com", user.Email)
	require.Equal(t, int64(17), user.Age)
	require.Equal(t, "Nodirbek", user.FirstName)
	require.Equal(t, "rareuz", user.Username)
	require.Equal(t, "Nodirbek2006", user.Password)
	require.NotNil(t, user.Id)

	// GetUser
	getReq := NewRequest(http.MethodGet, "/users/get", nil)
	q := getReq.URL.Query()
	q.Add("id", user.Id)
	getReq.URL.RawQuery = q.Encode()
	getRes := httptest.NewRecorder()
	r = gin.Default()
	r.GET("/users/get", handlers.GetUser)
	r.ServeHTTP(getRes, getReq)
	assert.Equal(t, http.StatusOK, getRes.Code)
	var getUserResp storage.User
	bodyBytes, err := io.ReadAll(getRes.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(bodyBytes, &getUserResp))
	assert.Equal(t, user.Id, getUserResp.Id)
	assert.Equal(t, user.FirstName, getUserResp.FirstName)
	assert.Equal(t, user.Username, getUserResp.Username)
	assert.Equal(t, user.Age, getUserResp.Age)
	assert.Equal(t, user.Email, getUserResp.Email)

	// User List
	listReq := NewRequest(http.MethodGet, "/users", nil)
	listRes := httptest.NewRecorder()
	r = gin.Default()
	r.GET("/users", handlers.ListUsers)
	r.ServeHTTP(listRes, listReq)
	assert.Equal(t, http.StatusOK, listRes.Code)
	bodyBytes, err = io.ReadAll(listRes.Body)
	assert.NoError(t, err)
	assert.NotNil(t, bodyBytes)

	// User Delete
	delReq := NewRequest(http.MethodDelete, "/user/delete?id="+user.Id, nil)
	delRes := httptest.NewRecorder()
	r.DELETE("/user/delete", handlers.DeleteUser)
	r.ServeHTTP(delRes, delReq)
	assert.Equal(t, http.StatusOK, delRes.Code)
	var respm storage.Message
	bodyBytes, _ = io.ReadAll(delRes.Body)
	require.NoError(t, json.Unmarshal(bodyBytes, &respm))
	require.Equal(t, "user was deleted successfully", respm.Message)

	// User Register
	regReq := NewRequest(http.MethodPost, "/user/register", buffer)
	regRes := httptest.NewRecorder()
	r.POST("/user/register", handlers.RegisterUser)
	r.ServeHTTP(regRes, regReq)
	assert.Equal(t, http.StatusOK, regRes.Code)
	var resp storage.Message
	bodyBytes, err = io.ReadAll(regRes.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(bodyBytes, &resp))
	require.NotNil(t, resp.Message)
	require.Equal(t, "One time password sent to your email", resp.Message)

	// User Verify with correct code
	verURLCorrect := "/user/verify/12345"
	verReqCorrect := NewRequest(http.MethodGet, verURLCorrect, buffer)
	verResCorrect := httptest.NewRecorder()
	r = gin.Default()
	r.GET("/user/verify/:code", handlers.Verify)
	r.ServeHTTP(verResCorrect, verReqCorrect)

	assert.Equal(t, http.StatusOK, verResCorrect.Code)
	var responseCorrect storage.Message
	bodyBytesCorrect, err := io.ReadAll(verResCorrect.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(bodyBytesCorrect, &responseCorrect))
	require.Equal(t, "Success", responseCorrect.Message)

	// User Verify with incorrect code
	verURLIncorrect := "/user/verify/54321" // 54321 ni xato kod deb oldim va xatoliklar uchun ham tekshirib ketdim
	verReqIncorrect := NewRequest(http.MethodGet, verURLIncorrect, buffer)
	verResIncorrect := httptest.NewRecorder()
	r.ServeHTTP(verResIncorrect, verReqIncorrect)

	assert.Equal(t, http.StatusBadRequest, verResIncorrect.Code)
	var responseIncorrect storage.Message
	bodyBytesIncorrect, err := io.ReadAll(verResIncorrect.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(bodyBytesIncorrect, &responseIncorrect))
	require.Equal(t, "Incorrect code", responseIncorrect.Message)

	//PRODUCT TEST

	gin.SetMode(gin.TestMode)
	require.NoError(t, SetupMinimumInstance(""))
	buffer, err = OpenFile("product.json")

	require.NoError(t, err)

	// Product create
	req = NewRequest(http.MethodPost, "/product/create", buffer)
	res = httptest.NewRecorder()
	r = gin.Default()
	r.POST("/product/create", handlers.CreateProduct)
	r.ServeHTTP(res, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.Code)
	var product storage.Product
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &product))
	require.Equal(t, product.Amount, int64(99))
	require.Equal(t, product.Description, "Nodirbek's Description")
	require.Equal(t, product.Name, "Nodirbek's Product")
	require.Equal(t, product.Price, float32(99.9))

	// Get Product
	getReq = NewRequest(http.MethodGet, "/product/get", buffer)
	q = getReq.URL.Query()
	q.Add("id", product.Id)
	getReq.URL.RawQuery = q.Encode()
	getRes = httptest.NewRecorder()
	r = gin.Default()
	r.GET("/product/get", handlers.GetProduct)
	r.ServeHTTP(getRes, getReq)
	assert.Equal(t, http.StatusOK, getRes.Code)
	var getProduct storage.Product
	bodyBytes, err = io.ReadAll(getRes.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(bodyBytes, &getProduct))
	require.Equal(t, product.Id, getProduct.Id)
	require.Equal(t, product.Amount, getProduct.Amount)
	require.Equal(t, product.Description, getProduct.Description)
	require.Equal(t, product.Name, getProduct.Name)
	require.Equal(t, product.Price, getProduct.Price)

	// List Products
	listReqq := NewRequest(http.MethodGet, "/products", buffer)
	listRess := httptest.NewRecorder()
	r = gin.Default()
	r.GET("/products", handlers.ListProducts)
	r.ServeHTTP(listRess, listReqq)
	assert.Equal(t, http.StatusOK, listRess.Code)
	bodyBytes, err = io.ReadAll(listRess.Body)
	assert.NoError(t, err)
	assert.NotNil(t, bodyBytes)

	// Delete Product
	delReq = NewRequest(http.MethodDelete, "/product/delete?id="+product.Id, buffer)
	delRes = httptest.NewRecorder()
	r.DELETE("/product/delete", handlers.DeleteProduct)
	r.ServeHTTP(delRes, delReq)
	assert.Equal(t, http.StatusOK, delRes.Code)
	var respmessage storage.Message
	bodyBytes, _ = io.ReadAll(delRes.Body)
	require.NoError(t, json.Unmarshal(bodyBytes, &respmessage))
	require.Equal(t, "product was deleted successfully", respmessage.Message)

}
