package tests

import (
	"net/url"

	"github.com/go-openapi/strfmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

type TestedData struct {
	Id   strfmt.UUID4 `json:"id"`
	Name string       `json:"name"`
}

func TestApiMsg(t *T) {
	ctx := testInit(t)
	getPathName := "pop"
	createPathName := "push"
	removePathName := "remove"
	query := url.Values{}
	createdMsg := TestedData{
		Id:   strfmt.UUID4(uuid.NewString()),
		Name: RandomShortString(),
	}
	var msg TestedData
	t.Run("push and get one msg", func(t *T) {
		require.NoError(t, tRequestPostByName(t, ctx, MockServer, query, createPathName, createdMsg, nil))
		require.NoError(t, tRequestGetByName(t, ctx, MockServer, query, getPathName, &msg))
		require.NoError(t, tRequestDeleteByName(t, ctx, MockServer, query, removePathName, nil, nil))

		require.Equal(t, createdMsg.Id, msg.Id)
		require.Equal(t, createdMsg.Name, msg.Name)
	})

	t.Run("push and get multiple msg", func(t *T) {
		createdMsg1 := TestedData{
			Id:   strfmt.UUID4(uuid.NewString()),
			Name: RandomShortString(),
		}
		createdMsg2 := TestedData{
			Id:   strfmt.UUID4(uuid.NewString()),
			Name: RandomShortString(),
		}
		createdMsg3 := TestedData{
			Id:   strfmt.UUID4(uuid.NewString()),
			Name: RandomShortString(),
		}
		require.NoError(t, tRequestPostByName(t, ctx, MockServer, query, createPathName, createdMsg1, nil))
		require.NoError(t, tRequestPostByName(t, ctx, MockServer, query, createPathName, createdMsg2, nil))
		require.NoError(t, tRequestPostByName(t, ctx, MockServer, query, createPathName, createdMsg3, nil))

		require.NoError(t, tRequestGetByName(t, ctx, MockServer, query, getPathName, &msg))
		require.Equal(t, createdMsg1.Id, msg.Id)
		require.Equal(t, createdMsg1.Name, msg.Name)
		require.NoError(t, tRequestDeleteByName(t, ctx, MockServer, query, removePathName, nil, nil))

		require.NoError(t, tRequestGetByName(t, ctx, MockServer, query, getPathName, &msg))
		require.Equal(t, createdMsg2.Id, msg.Id)
		require.Equal(t, createdMsg2.Name, msg.Name)
		require.NoError(t, tRequestDeleteByName(t, ctx, MockServer, query, removePathName, nil, nil))

		require.NoError(t, tRequestGetByName(t, ctx, MockServer, query, getPathName, &msg))
		require.Equal(t, createdMsg3.Id, msg.Id)
		require.Equal(t, createdMsg3.Name, msg.Name)
		require.NoError(t, tRequestDeleteByName(t, ctx, MockServer, query, removePathName, nil, nil))
	})
}
