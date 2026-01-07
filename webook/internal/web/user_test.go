package web

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"webook/internal/domain"
	svcmocks "webook/internal/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name string
	}{}
	server := gin.Default()
	h := NewUserHandler(nil, nil)
	h.RegisterRoutes(server)
	req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewBuffer([]byte(`
{
	"email":"123@qq.com",
	"password":"123456"

}`)))
	require.NoError(t, err)
	//继续使用req
	t.Log(req)
	resp := httptest.NewRecorder()
	t.Log(resp)
	//HTTP请求进入gin的路口
	//当你这样调用 GIN会处理这个请求
	//响应写回resp里
	server.ServeHTTP(resp, req)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := NewUserHandler(nil, nil)
			ctx := &gin.Context{}
			handler.SignUp(ctx)
		})
	}
}
func TestMock(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	usersvc := svcmocks.NewMockUserService(ctrl)
	usersvc.EXPECT().SignUp(gomock.Any(), gomock.Any()).Return(errors.New("mock err"))
	err := usersvc.SignUp(context.Background(), domain.User{
		Email: "123@qq.com",
	})
	t.Log(err)
}
