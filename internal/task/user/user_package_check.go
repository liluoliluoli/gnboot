package user

import (
	"github.com/liluoliluoli/gnboot/internal/conf"
	"github.com/liluoliluoli/gnboot/internal/service"
	"github.com/liluoliluoli/gnboot/internal/service/sdomain"
)

type UserPackageCheckTask struct {
	c           *conf.Bootstrap
	userService *service.UserService
}

func NewUserPackageCheckTask(c *conf.Bootstrap, userService *service.UserService) *UserPackageCheckTask {
	return &UserPackageCheckTask{
		c:           c,
		userService: userService,
	}
}

func (t *UserPackageCheckTask) Process(task *sdomain.Task) error {
	return nil
}
