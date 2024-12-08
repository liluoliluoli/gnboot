package adaptor

import (
	"context"
	"github.com/liluoliluoli/gnboot/api/keyword"
	"github.com/liluoliluoli/gnboot/internal/service"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
	"github.com/samber/lo"
)

type KeywordProvider struct {
	keyword.UnimplementedKeywordRemoteServiceServer
	keyword *service.KeywordService
}

func NewKeywordProvider(keyword *service.KeywordService) *KeywordProvider {
	return &KeywordProvider{keyword: keyword}
}

func (s *KeywordProvider) FindKeyword(ctx context.Context, req *keyword.FindKeywordRequest) (*keyword.FindKeywordResp, error) {
	res, err := s.keyword.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	return &keyword.FindKeywordResp{
		Keywords: lo.Map(res, func(item *sdomain.Keyword, index int) *keyword.KeywordResp {
			return &keyword.KeywordResp{
				Name: item.Name,
				Id:   int32(item.ID),
			}
		}),
	}, nil

}
