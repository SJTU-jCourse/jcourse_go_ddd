package web

import (
	"github.com/gin-gonic/gin"

	"jcourse_go/internal/application/announcement/query"
)

type AnnouncementController struct {
	announcementQueryService query.AnnouncementQueryService
}

func NewAnnouncementController(announcementQueryService query.AnnouncementQueryService) *AnnouncementController {
	return &AnnouncementController{
		announcementQueryService: announcementQueryService,
	}
}

func (c *AnnouncementController) GetAnnouncements(ctx *gin.Context) {
	commonCtx := GetCommonContext(ctx)

	announcements, err := c.announcementQueryService.GetAnnouncements(commonCtx)
	if err != nil {
		HandleError(ctx, err)
		return
	}

	HandleSuccess(ctx, announcements)
}
