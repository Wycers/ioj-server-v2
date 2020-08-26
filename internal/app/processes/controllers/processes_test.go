package controllers

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/infinity-oj/server-v2/internal/pkg/models"
	"github.com/infinity-oj/server-v2/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var r *gin.Engine
var configFile = flag.String("f", "accounts.yml", "set config file which viper will loading.")

func setup() {
	r = gin.New()
}

func TestUsersController_Get(t *testing.T) {
	flag.Parse()
	setup()

	sto := new(mocks.UsersRepository)

	sto.On("GetAccount", mock.AnythingOfType("uint64")).Return(func(ID uint64) (p *models.Detail) {
		return &models.Detail{
			ID: ID,
		}
	}, func(ID uint64) error {
		return nil
	})

	c, err := CreateUsersController(*configFile, sto)
	if err != nil {
		t.Fatalf("create product serviceerror,%+v", err)
	}

	r.GET("/proto/:id", c.Get)

	tests := []struct {
		name     string
		id       uint64
		expected uint64
	}{
		{"1", 1, 1},
		{"2", 2, 2},
		{"3", 3, 3},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uri := fmt.Sprintf("/proto/%d", test.id)
			// 构造get请求
			req := httptest.NewRequest("GET", uri, nil)
			// 初始化响应
			w := httptest.NewRecorder()

			// 调用相应的controller接口
			r.ServeHTTP(w, req)

			// 提取响应
			rs := w.Result()
			defer func() {
				_ = rs.Body.Close()
			}()

			// 读取响应body
			body, _ := ioutil.ReadAll(rs.Body)
			p := new(models.Detail)
			err := json.Unmarshal(body, p)
			if err != nil {
				t.Errorf("unmarshal response body error:%v", err)
			}

			assert.Equal(t, test.expected, p.ID)
		})
	}

}
