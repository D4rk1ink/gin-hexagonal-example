package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	http_apigen "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/apigen"
	custom_error "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
	repository_model "github.com/D4rk1ink/gin-hexagonal-example/internal/infrastructure/repository/model"
	"github.com/guregu/null"
	"github.com/oapi-codegen/runtime/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/v2/bson"
)

var _ = Describe("User Integration", Label("Integration"), func() {
	BeforeEach(func() {
		clearDb()
	})

	registerUser := func(name string) {
		payload := http_apigen.RegisterReq{
			Name:            name,
			Email:           types.Email(name + "@email.com"),
			Password:        name + "password",
			ConfirmPassword: name + "password",
		}
		b, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(b))
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
	}

	loginUser := func(name string) string {
		payload := http_apigen.LoginReq{
			Email:    types.Email(name + "@email.com"),
			Password: name + "password",
		}
		b, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", server.URL+"/api/auth/login", bytes.NewReader(b))
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)

		var body http_apigen.LoginRes
		resBody, _ := io.ReadAll(res.Body)
		json.Unmarshal(resBody, &body)

		return body.AccessToken
	}

	findById := func(id string) repository_model.UserModel {
		var userModel repository_model.UserModel
		objId, _ := bson.ObjectIDFromHex(id)
		db.GetDb().Collection("users").FindOne(ctx, bson.M{"_id": objId}).Decode(&userModel)

		return userModel
	}

	findUserByAccessToken := func(accessToken string) repository_model.UserModel {
		token, _ := jwtInfra.ValidateAccessToken(accessToken)

		return findById(token.ID)
	}

	Context("GET /api/users", func() {
		It("should return list of users", func() {
			registerUser("mock")
			accessToken := loginUser("mock")

			req := httptest.NewRequest("GET", server.URL+"/api/users", nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			Expect(res).ToNot(BeNil())
			Expect(res.Code).To(Equal(http.StatusOK))

			var body http_apigen.UsersRes
			resBody, err := io.ReadAll(res.Body)
			Expect(err).To(BeNil())

			err = json.Unmarshal(resBody, &body)
			Expect(err).To(BeNil())

			userModel := findById(body.Data[0].Id)

			Expect(body).ToNot(BeNil())
			Expect(body.Data).ToNot(BeNil())
			Expect(body.Data[0].Id).To(Equal(userModel.ID.Hex()))
			Expect(body.Data[0].Name).To(Equal(userModel.Name))
			Expect(body.Data[0].Email).To(Equal(userModel.Email))
			Expect(body.Data[0].CreatedAt.String()).To(Equal(userModel.CreatedAt.String()))
			Expect(body.Data[0].UpdatedAt.String()).To(Equal(userModel.UpdatedAt.String()))
		})
		It("should return 401 when not authenticated", func() {
			req := httptest.NewRequest("GET", server.URL+"/api/users", nil)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			Expect(res).ToNot(BeNil())
			Expect(res.Code).To(Equal(http.StatusUnauthorized))

			var body http_apigen.ErrorRes
			resBody, err := io.ReadAll(res.Body)
			Expect(err).To(BeNil())

			err = json.Unmarshal(resBody, &body)
			Expect(err).To(BeNil())

			Expect(body).ToNot(BeNil())
			Expect(body.Error).ToNot(BeNil())
			Expect(body.Error.Code).To(Equal(custom_error.ErrUnauthorized))
		})
	})

	Context("GET /api/users/{id}", func() {
		It("should return user by id", func() {
			registerUser("mock")
			accessToken := loginUser("mock")
			userModel := findUserByAccessToken(accessToken)

			req := httptest.NewRequest("GET", server.URL+"/api/users/"+userModel.ID.Hex(), nil)
			req.Header.Set("Authorization", "Bearer "+accessToken)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			Expect(res).ToNot(BeNil())
			Expect(res.Code).To(Equal(http.StatusOK))

			var body http_apigen.UserRes
			resBody, err := io.ReadAll(res.Body)
			Expect(err).To(BeNil())

			err = json.Unmarshal(resBody, &body)
			Expect(err).To(BeNil())

			Expect(body).ToNot(BeNil())
			Expect(body.Data).ToNot(BeNil())
			Expect(body.Data.Id).To(Equal(userModel.ID.Hex()))
			Expect(body.Data.Name).To(Equal(userModel.Name))
			Expect(body.Data.Email).To(Equal(userModel.Email))
			Expect(body.Data.CreatedAt.String()).To(Equal(userModel.CreatedAt.String()))
			Expect(body.Data.UpdatedAt.String()).To(Equal(userModel.UpdatedAt.String()))
		})
		It("should return 401 when not authenticated", func() {
			req := httptest.NewRequest("GET", server.URL+"/api/users/1", nil)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			Expect(res).ToNot(BeNil())
			Expect(res.Code).To(Equal(http.StatusUnauthorized))

			var body http_apigen.ErrorRes
			resBody, err := io.ReadAll(res.Body)
			Expect(err).To(BeNil())

			err = json.Unmarshal(resBody, &body)
			Expect(err).To(BeNil())

			Expect(body).ToNot(BeNil())
			Expect(body.Error).ToNot(BeNil())
			Expect(body.Error.Code).To(Equal(custom_error.ErrUnauthorized))
		})
	})

	Context("PATCH /api/users/{id}", func() {
		It("should return user by id", func() {
			registerUser("mock")
			accessToken := loginUser("mock")
			userModel := findUserByAccessToken(accessToken)

			payload := http_apigen.UserUpdateReq{
				Name:  null.StringFrom("new_mock").Ptr(),
				Email: (*types.Email)(null.StringFrom("new_mock@email.com").Ptr()),
			}
			b, _ := json.Marshal(payload)
			req := httptest.NewRequest("PATCH", server.URL+"/api/users/"+userModel.ID.Hex(), bytes.NewReader(b))
			req.Header.Set("Authorization", "Bearer "+accessToken)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			Expect(res).ToNot(BeNil())
			Expect(res.Code).To(Equal(http.StatusOK))

			var body http_apigen.UserRes
			resBody, err := io.ReadAll(res.Body)
			Expect(err).To(BeNil())

			err = json.Unmarshal(resBody, &body)
			Expect(err).To(BeNil())

			Expect(body).ToNot(BeNil())
			Expect(body.Data).ToNot(BeNil())
			Expect(body.Data.Id).ToNot(BeEmpty())
			Expect(body.Data.Name).To(Equal("new_mock"))
			Expect(body.Data.Email).To(Equal("new_mock@email.com"))
			Expect(body.Data.CreatedAt.String()).ToNot(BeEmpty())
			Expect(body.Data.UpdatedAt.String()).ToNot(BeEmpty())

			userModel = findById(body.Data.Id)

			Expect(body).ToNot(BeNil())
			Expect(body.Data).ToNot(BeNil())
			Expect(body.Data.Id).To(Equal(userModel.ID.Hex()))
			Expect(body.Data.Name).To(Equal(userModel.Name))
			Expect(body.Data.Email).To(Equal(userModel.Email))
			Expect(body.Data.CreatedAt.UnixMilli()).To(Equal(userModel.CreatedAt.UnixMilli()))
			Expect(body.Data.UpdatedAt.UnixMilli()).To(Equal(userModel.UpdatedAt.UnixMilli()))
		})
		It("should return 400 Bad Request if invalid email format", func() {
			registerUser("mock")
			accessToken := loginUser("mock")
			userModel := findUserByAccessToken(accessToken)

			payload := http_apigen.UserUpdateReq{
				Name:  null.StringFrom("mock").Ptr(),
				Email: (*types.Email)(null.StringFrom("mock_email.com").Ptr()),
			}
			b, _ := json.Marshal(payload)
			req := httptest.NewRequest("PATCH", server.URL+"/api/users/"+userModel.ID.Hex(), bytes.NewReader(b))
			req.Header.Set("Authorization", "Bearer "+accessToken)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			Expect(res).ToNot(BeNil())
			Expect(res.Code).To(Equal(http.StatusBadRequest))

			var body http_apigen.ErrorRes
			resBody, err := io.ReadAll(res.Body)
			Expect(err).To(BeNil())

			err = json.Unmarshal(resBody, &body)
			Expect(err).To(BeNil())

			Expect(body).ToNot(BeNil())
			Expect(body.Error).ToNot(BeNil())
			Expect(body.Error.Code).To(Equal(custom_error.ErrBadRequest))
		})
		It("should return 404 Not Found if user not found", func() {
			registerUser("mock")
			accessToken := loginUser("mock")

			invalidId := bson.NewObjectID().Hex()
			payload := http_apigen.UserUpdateReq{
				Name:  null.StringFrom("new_mock").Ptr(),
				Email: (*types.Email)(null.StringFrom("new_mock@email.com").Ptr()),
			}
			b, _ := json.Marshal(payload)
			req := httptest.NewRequest("PATCH", server.URL+"/api/users/"+invalidId, bytes.NewReader(b))
			req.Header.Set("Authorization", "Bearer "+accessToken)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			Expect(res).ToNot(BeNil())
			Expect(res.Code).To(Equal(http.StatusNotFound))

			var body http_apigen.ErrorRes
			resBody, err := io.ReadAll(res.Body)
			Expect(err).To(BeNil())

			err = json.Unmarshal(resBody, &body)
			Expect(err).To(BeNil())

			Expect(body).ToNot(BeNil())
			Expect(body.Error).ToNot(BeNil())
			Expect(body.Error.Code).To(Equal(custom_error.ErrUserNotFound))
		})
		It("should return 409 Bad Request if update email with existing email", func() {
			registerUser("mock")
			registerUser("mock2")
			accessToken := loginUser("mock")
			userModel := findUserByAccessToken(accessToken)

			payload := http_apigen.UserUpdateReq{
				Email: (*types.Email)(null.StringFrom("mock2@email.com").Ptr()),
			}
			b, _ := json.Marshal(payload)
			req := httptest.NewRequest("PATCH", server.URL+"/api/users/"+userModel.ID.Hex(), bytes.NewReader(b))
			req.Header.Set("Authorization", "Bearer "+accessToken)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			Expect(res).ToNot(BeNil())
			Expect(res.Code).To(Equal(http.StatusConflict))

			var body http_apigen.ErrorRes
			resBody, err := io.ReadAll(res.Body)
			Expect(err).To(BeNil())

			err = json.Unmarshal(resBody, &body)
			Expect(err).To(BeNil())

			Expect(body).ToNot(BeNil())
			Expect(body.Error).ToNot(BeNil())
			Expect(body.Error.Code).To(Equal(custom_error.ErrUserEmailAlreadyExists))
		})
		It("should return 401 when not authenticated", func() {
			req := httptest.NewRequest("GET", server.URL+"/api/users/1", nil)
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			Expect(res).ToNot(BeNil())
			Expect(res.Code).To(Equal(http.StatusUnauthorized))

			var body http_apigen.ErrorRes
			resBody, err := io.ReadAll(res.Body)
			Expect(err).To(BeNil())

			err = json.Unmarshal(resBody, &body)
			Expect(err).To(BeNil())

			Expect(body).ToNot(BeNil())
			Expect(body.Error).ToNot(BeNil())
			Expect(body.Error.Code).To(Equal(custom_error.ErrUnauthorized))
		})
	})
})
