package singlefile

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vmware-tanzu/graph-framework-for-microservices/src/gqlgen/client"
	"github.com/vmware-tanzu/graph-framework-for-microservices/src/gqlgen/graphql/handler"
)

func TestPtrToSlice(t *testing.T) {
	resolvers := &Stub{}

	c := client.New(handler.NewDefaultServer(NewExecutableSchema(Config{Resolvers: resolvers})))

	resolvers.QueryResolver.PtrToSliceContainer = func(ctx context.Context) (wrappedStruct *PtrToSliceContainer, e error) {
		ptrToSliceContainer := PtrToSliceContainer{
			PtrToSlice: &[]string{"hello"},
		}
		return &ptrToSliceContainer, nil
	}

	t.Run("pointer to slice", func(t *testing.T) {
		var resp struct {
			PtrToSliceContainer struct {
				PtrToSlice []string
			}
		}

		err := c.Post(`query { ptrToSliceContainer {  ptrToSlice }}`, &resp)
		require.NoError(t, err)

		require.Equal(t, []string{"hello"}, resp.PtrToSliceContainer.PtrToSlice)
	})
}
