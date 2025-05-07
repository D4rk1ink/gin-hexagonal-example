package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	http_apigen "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/apigen"
	custom_error "github.com/D4rk1ink/gin-hexagonal-example/internal/core/error"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth Integration", Label("Integration"), func() {
	BeforeEach(func() {
		clearDb()
	})

	Context("POST /api/auth/register", func() {
		It("should return 201 Created if register successfully", func() {
			payload := http_apigen.RegisterReq{
				Name:            "mock",
				Email:           "mock@email.com",
				Password:        "password",
				ConfirmPassword: "password",
			}
			b, _ := json.Marshal(payload)
			req := httptest.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(b))
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			Expect(res).ToNot(BeNil())
			Expect(res.Code).To(Equal(http.StatusCreated))

			var body http_apigen.RegisterRes
			resBody, err := io.ReadAll(res.Body)
			Expect(err).To(BeNil())

			err = json.Unmarshal(resBody, &body)
			Expect(err).To(BeNil())

			Expect(body).ToNot(BeNil())
			Expect(body.Success).To(BeTrue())
		})
		It("should return 400 Bad Request if some field is missing", func() {
			payload := http_apigen.RegisterReq{
				Email:           "mock@email.com",
				Password:        "password",
				ConfirmPassword: "password",
			}
			b, _ := json.Marshal(payload)
			req := httptest.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(b))
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
			Expect(body.Error.Code).To(Equal(custom_error.ErrBadRequest))
		})
		It("should return 400 Bad Request if required field is empty string", func() {
			payload := http_apigen.RegisterReq{
				Name:            "",
				Email:           "mock@email.com",
				Password:        "password",
				ConfirmPassword: "invalid_password",
			}
			b, _ := json.Marshal(payload)
			req := httptest.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(b))
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
			Expect(body.Error.Code).To(Equal(custom_error.ErrBadRequest))
		})
		It("should return 400 Bad Request if password and confirm password not match", func() {
			payload := http_apigen.RegisterReq{
				Name:            "mock",
				Email:           "mock@email.com",
				Password:        "password",
				ConfirmPassword: "invalid_password",
			}
			b, _ := json.Marshal(payload)
			req := httptest.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(b))
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
			Expect(body.Error.Code).To(Equal(custom_error.ErrAuthInvalidConfirmPassword))
		})
		It("should return 400 Bad Request if email is invalid format", func() {
			payload := http_apigen.RegisterReq{
				Name:            "mock",
				Email:           "mockemail.com",
				Password:        "password",
				ConfirmPassword: "password",
			}
			b, _ := json.Marshal(payload)
			req := httptest.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(b))
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
			Expect(body.Error.Code).To(Equal(custom_error.ErrBadRequest))
		})
		It("should return 400 Bad Request if email already registered", func() {
			_, _ = db.GetDb().Collection("users").InsertOne(ctx, map[string]interface{}{
				"name":       "mock",
				"email":      "mock@email.com",
				"password":   "password",
				"created_at": time.Now().Format(time.RFC3339),
				"updated_at": time.Now().Format(time.RFC3339),
			})
			payload := http_apigen.RegisterReq{
				Name:            "mock",
				Email:           "mock@email.com",
				Password:        "password",
				ConfirmPassword: "password",
			}
			b, _ := json.Marshal(payload)
			req := httptest.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(b))
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
			Expect(body.Error.Code).To(Equal(custom_error.ErrAuthEmailAlreadyExists))
		})
	})

	Context("POST /api/auth/login", func() {
		It("should return 200 OK with access token if login successfully", func() {
			payloadRegister := http_apigen.RegisterReq{
				Name:            "mock",
				Email:           "mock@email.com",
				Password:        "password",
				ConfirmPassword: "password",
			}
			b, _ := json.Marshal(payloadRegister)
			req := httptest.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(b))
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			payload := http_apigen.LoginReq{
				Email:    payloadRegister.Email,
				Password: payloadRegister.Password,
			}
			b, _ = json.Marshal(payload)
			req = httptest.NewRequest("POST", server.URL+"/api/auth/login", bytes.NewReader(b))
			res = httptest.NewRecorder()
			router.ServeHTTP(res, req)

			Expect(res).ToNot(BeNil())
			Expect(res.Code).To(Equal(http.StatusOK))

			var body http_apigen.LoginRes
			resBody, err := io.ReadAll(res.Body)
			Expect(err).To(BeNil())

			err = json.Unmarshal(resBody, &body)
			Expect(err).To(BeNil())

			Expect(body).ToNot(BeNil())
			Expect(body.AccessToken).ToNot(BeNil())
			Expect(body.AccessToken).ToNot(BeEmpty())
			Expect(body.TokenType).To(Equal("Bearer"))
			Expect(body.ExpiresIn).To(BeNumerically(">", 0))
		})
		It("should return 400 Bad Request if payload is invalid", func() {
			payloadRegister := http_apigen.RegisterReq{
				Email: "mock@email.com",
			}
			b, _ := json.Marshal(payloadRegister)
			req := httptest.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(b))
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			payload := http_apigen.LoginReq{
				Email: payloadRegister.Email,
			}
			b, _ = json.Marshal(payload)
			req = httptest.NewRequest("POST", server.URL+"/api/auth/login", bytes.NewReader(b))
			res = httptest.NewRecorder()
			router.ServeHTTP(res, req)

			Expect(res).ToNot(BeNil())
			Expect(res.Code).To(Equal(http.StatusBadRequest))

			var body http_apigen.ErrorRes
			resBody, err := io.ReadAll(res.Body)
			Expect(err).To(BeNil())

			err = json.Unmarshal(resBody, &body)
			Expect(err).To(BeNil())

			Expect(body).ToNot(BeNil())
			Expect(body.Error.Code).To(Equal(custom_error.ErrBadRequest))
		})
		It("should return 401 Unauthorized if email not registered", func() {
			payload := http_apigen.LoginReq{
				Email:    "mock@email.com",
				Password: "password",
			}
			b, _ := json.Marshal(payload)
			req := httptest.NewRequest("POST", server.URL+"/api/auth/login", bytes.NewReader(b))
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
			Expect(body.Error.Code).To(Equal(custom_error.ErrAuthInvalidCredentials))
		})
		It("should return 401 Unauthorized if password is invalid", func() {
			payloadRegister := http_apigen.RegisterReq{
				Name:            "mock",
				Email:           "mock@email.com",
				Password:        "password",
				ConfirmPassword: "password",
			}
			b, _ := json.Marshal(payloadRegister)
			req := httptest.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(b))
			res := httptest.NewRecorder()
			router.ServeHTTP(res, req)

			payload := http_apigen.LoginReq{
				Email:    payloadRegister.Email,
				Password: "invalid_password",
			}
			b, _ = json.Marshal(payload)
			req = httptest.NewRequest("POST", server.URL+"/api/auth/login", bytes.NewReader(b))
			res = httptest.NewRecorder()
			router.ServeHTTP(res, req)

			Expect(res).ToNot(BeNil())
			Expect(res.Code).To(Equal(http.StatusUnauthorized))

			var body http_apigen.ErrorRes
			resBody, err := io.ReadAll(res.Body)
			Expect(err).To(BeNil())

			err = json.Unmarshal(resBody, &body)
			Expect(err).To(BeNil())

			Expect(body).ToNot(BeNil())
			Expect(body.Error.Code).To(Equal(custom_error.ErrAuthInvalidCredentials))
		})
		It("should return 400 Bad Request if email is invalid format", func() {
			payload := http_apigen.LoginReq{
				Email:    "mockemail.com",
				Password: "password",
			}
			b, _ := json.Marshal(payload)
			req := httptest.NewRequest("POST", server.URL+"/api/auth/login", bytes.NewReader(b))
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
			Expect(body.Error.Code).To(Equal(custom_error.ErrBadRequest))
		})
		It("should return 400 Bad Request if required field is empty string", func() {
			payload := http_apigen.LoginReq{
				Email:    "mockemail.com",
				Password: "",
			}
			b, _ := json.Marshal(payload)
			req := httptest.NewRequest("POST", server.URL+"/api/auth/login", bytes.NewReader(b))
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
			Expect(body.Error.Code).To(Equal(custom_error.ErrBadRequest))
		})
	})
})
