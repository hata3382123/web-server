package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"webook/internal/web"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
}

func (s *ArticleTestSuite) SetupSuite() {
	//在所有测试执行之前，初始化一些内容
	s.server = gin.Default()
	artHdl := web.NewArticleHandler()
	//注册好了路由
	artHdl.RegisterRoutes(s.server)
}

func (s *ArticleTestSuite) TestEdit() {
	t := s.T()
	testCases := []struct {
		name string

		//集成测试准备数据
		before func(t *testing.T)
		//集成测试验证数据
		after func(t *testing.T)
		//预期中的输入
		art Article
		//HTTP响应码
		wantCode int
		//HTTP响应体
		wantRes Result[int64]
	}{
		{},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			//构造请求
			//执行
			//验证
			tc.before(t)
			reqBody, err := json.Marshal(tc.art)
			assert.NoError(t, err)
			req, err := http.NewRequest(http.MethodPost, "/articles/edit", bytes.NewBuffer([]byte(reqBody)))
			require.NoError(t, err)
			//数据是JSON格式
			req.Header.Set("Content-Type", "application/json")
			//这里可以继续使用req
			resp := httptest.NewRecorder()
			//这里就是HTTP请求进去GIN框架的入口
			//当你这样调用的时候 GIN会处理这个请求
			//响应写回 resp
			s.server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			if resp.Code != 200 {
				return
			}
			var webRes Result[int64]
			err = json.NewDecoder(resp.Body).Decode(&webRes)
			require.NoError(t, err)
			assert.Equal(t, tc.wantRes, webRes)
			tc.after(t)
		})
	}
}

func (s *ArticleTestSuite) TestABC() {
	s.T().Log("这是一个test")
}
func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

type Article struct {
	title   string
	content string
}
type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
