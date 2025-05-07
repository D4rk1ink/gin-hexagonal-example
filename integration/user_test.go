package integration_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	http_apigen "github.com/D4rk1ink/gin-hexagonal-example/internal/application/handler/http/apigen"
	"github.com/oapi-codegen/runtime/types"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("User Integration", Label("Integration"), func() {
	BeforeEach(func() {
		clearDb()
	})

	registerUser := func(prefix string) {
		payload := http_apigen.RegisterReq{
			Name:            prefix + "mock",
			Email:           types.Email(prefix + "mock@email.com"),
			Password:        prefix + "password",
			ConfirmPassword: "password",
		}
		b, _ := json.Marshal(payload)
		req := httptest.NewRequest("POST", server.URL+"/api/auth/register", bytes.NewReader(b))
		res := httptest.NewRecorder()
		router.ServeHTTP(res, req)
	}

	loginUser := func(prefix string) string {
		payload := http_apigen.LoginReq{
			Email:    types.Email(prefix + "mock@email.com"),
			Password: prefix + "password",
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

	Context("GET /api/users", func() {
		It("should return list of users", func() {
			registerUser("")
			accessToken := loginUser("")

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

			Expect(body).ToNot(BeNil())
			Expect(body.Data).To(HaveLen(1))
		})
	})
})
