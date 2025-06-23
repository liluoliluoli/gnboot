package region_util

import (
	"context"
	"fmt"
	"github.com/biter777/countries"
	"github.com/liluoliluoli/gnboot/internal/common/constant"
	"github.com/liluoliluoli/gnboot/internal/common/utils/httpclient_util"
	"github.com/liluoliluoli/gnboot/internal/integration/dto/regiondto"
	"github.com/samber/lo"
)

func GetCnNameByName(ctx context.Context, name string) string {
	matche := constant.RegularChinese.MatchString(name)
	if matche {
		return name
	}
	if constant.RegionMap[name] != "" {
		return constant.RegionMap[name]
	}
	country := countries.ByName(name)
	regions, err := httpclient_util.DoGet[[]regiondto.Region](ctx, fmt.Sprintf(constant.GetCountryDetail, country.Alpha2()), nil)
	if err != nil {
		return ""
	}
	if len(lo.FromPtr(regions)) == 0 || lo.FromPtr(regions)[0].Translations == nil {
		return ""
	}
	if lo.FromPtr(regions)[0].Translations["zho"] != nil && lo.FromPtr(regions)[0].Translations["zho"]["common"] != "" {
		region := lo.FromPtr(regions)[0].Translations["zho"]["common"]
		constant.RegionMap[name] = region
		return region
	}
	return ""
}
