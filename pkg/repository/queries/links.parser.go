package queries

import (
	"tikkin/pkg/model"
	"tikkin/pkg/utils"
)

func (l *Link) ToModel() *model.Link {
	return &model.Link{
		ID:          l.ID,
		UserId:      l.UserID,
		Slug:        l.Slug,
		Description: utils.NullString(l.Description),
		Banned:      l.Banned,
		ExpireAt:    l.ExpireAt,
		TargetUrl:   l.TargetUrl,
		CreatedAt:   l.CreatedAt,
		UpdatedAt:   l.UpdatedAt,
	}
}
