package service

import (
	"context"
	"strconv"
	"strings"
	"testing"

	"github.com/go-cinch/common/proto/params"
	"github.com/google/uuid"
	"gnboot/api/gnboot"
	"gnboot/internal/tests/mock"
)

func TestGnbootService_CreateGnboot(t *testing.T) {
	s := mock.GnbootService()
	ctx := context.Background()
	userID := uuid.NewString()
	ctx = mock.NewContextWithUserId(ctx, userID)

	_, err := s.CreateGnboot(ctx, &gnboot.CreateGnbootRequest{
		Name: "gnboot1",
	})
	if err != nil {
		t.Error(err)
		return
	}
	_, err = s.CreateGnboot(ctx, &gnboot.CreateGnbootRequest{
		Name: "gnboot2",
	})
	if err != nil {
		t.Error(err)
		return
	}
	res1, _ := s.FindGnboot(ctx, &gnboot.FindGnbootRequest{
		Page: &params.Page{
			Disable: true,
		},
	})
	if res1 == nil || len(res1.List) != 2 {
		t.Error("res1 len must be 2")
		return
	}
	res2, err := s.GetGnboot(ctx, &gnboot.GetGnbootRequest{
		Id: res1.List[0].Id,
	})
	if err != nil {
		t.Error(err)
		return
	}
	if res2.Name != res1.List[0].Name {
		t.Errorf("res2.Name must be %s", res1.List[0].Name)
		return
	}
	_, err = s.DeleteGnboot(ctx, &params.IdsRequest{
		Ids: strings.Join([]string{
			strconv.FormatUint(res1.List[0].Id, 10),
			strconv.FormatUint(res1.List[1].Id, 10),
		}, ","),
	})
	if err != nil {
		t.Error(err)
		return
	}
	res3, _ := s.FindGnboot(ctx, &gnboot.FindGnbootRequest{
		Page: &params.Page{
			Disable: true,
		},
	})
	if res3 == nil || len(res3.List) != 0 {
		t.Error("res3 len must be 0")
		return
	}
}
