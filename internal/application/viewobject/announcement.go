package viewobject

import (
	"jcourse_go/internal/domain/announcement"
)

type AnnouncementVO struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	PublishedAt int64  `json:"published_at"`
}

func NewAnnouncementVO(announcement *announcement.Announcement) AnnouncementVO {
	return AnnouncementVO{
		ID:          int64(announcement.ID),
		Title:       announcement.Title,
		Content:     announcement.Content,
		PublishedAt: announcement.PublishedAt.Unix(),
	}
}
