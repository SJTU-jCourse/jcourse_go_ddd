package query

import (
	"jcourse_go/internal/application/viewobject"
	"jcourse_go/internal/domain/announcement"
	"jcourse_go/internal/domain/common"
	"jcourse_go/pkg/apperror"
)

type AnnouncementQueryService interface {
	GetAnnouncements(commonCtx *common.CommonContext) ([]viewobject.AnnouncementVO, error)
}

type announcementQueryService struct {
	announcementRepo announcement.AnnouncementRepository
}

func NewAnnouncementQueryService(announcementRepo announcement.AnnouncementRepository) AnnouncementQueryService {
	return &announcementQueryService{
		announcementRepo: announcementRepo,
	}
}

func (s *announcementQueryService) GetAnnouncements(commonCtx *common.CommonContext) ([]viewobject.AnnouncementVO, error) {
	announcements, err := s.announcementRepo.FindPublished(commonCtx.Ctx)
	if err != nil {
		return nil, apperror.ErrDB.Wrap(err)
	}

	announcementList := make([]viewobject.AnnouncementVO, len(announcements))
	for i, a := range announcements {
		announcementList[i] = viewobject.NewAnnouncementVO(&a)
	}
	return announcementList, nil
}
