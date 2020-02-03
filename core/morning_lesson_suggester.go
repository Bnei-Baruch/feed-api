package core

import (
	"database/sql"
	"fmt"

	"github.com/Bnei-Baruch/sqlboiler/queries/qm"

	"github.com/Bnei-Baruch/feed-api/consts"
	"github.com/Bnei-Baruch/feed-api/mdb"
	mdbmodels "github.com/Bnei-Baruch/feed-api/mdb/models"
	"github.com/Bnei-Baruch/feed-api/utils"
)

var SECURE_PUBLISHED_MOD = qm.Where(fmt.Sprintf("secure=%d AND published IS TRUE", consts.SEC_PUBLIC))

type MorningLessonSuggester struct {
	db *sql.DB
}

func MakeMorningLessonSuggester(db *sql.DB) *MorningLessonSuggester {
	return &MorningLessonSuggester{db: db}
}

func (suggester *MorningLessonSuggester) More(request MoreRequest) ([]ContentItem, error) {
	currentLessonUIDs := []string(nil)
	for _, ci := range request.CurrentFeed {
		if ci.ContentType == consts.CT_DAILY_LESSON {
			currentLessonUIDs = append(currentLessonUIDs, ci.UID)
		}
	}
	return suggester.fetchMorningLesson(currentLessonUIDs, request.MoreItems)
}

func (suggester *MorningLessonSuggester) fetchMorningLesson(currentLessonUIDs []string, moreItems int) ([]ContentItem, error) {
	mods := []qm.QueryMod{SECURE_PUBLISHED_MOD}
	mods = append(mods, qm.Where(fmt.Sprintf("type_id = %d", mdb.CONTENT_TYPE_REGISTRY.ByName[consts.CT_DAILY_LESSON].ID)))
	if len(currentLessonUIDs) > 0 {
		mods = append(mods, qm.WhereIn("uid NOT IN ?", utils.ToInterfaceSlice(currentLessonUIDs)...))
	}
	mods = append(mods, qm.OrderBy("created_at desc"))
	mods = append(mods, qm.Limit(moreItems))
	if lessons, err := mdbmodels.Collections(suggester.db, mods...).All(); err != nil {
		return nil, err
	} else {
		ret := []ContentItem(nil)
		for _, lesson := range lessons {
			ret = append(ret, ContentItem{UID: lesson.UID, CreatedAt: lesson.CreatedAt, ContentType: consts.CT_DAILY_LESSON})
		}
		return ret, nil
	}
}
