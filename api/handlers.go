package api

import (
	"database/sql"
	// "encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"

	"github.com/Bnei-Baruch/feed-api/core"
	"github.com/Bnei-Baruch/feed-api/data_models"
	"github.com/Bnei-Baruch/feed-api/recommendations"
	"github.com/Bnei-Baruch/feed-api/utils"
)

// Responds with JSON of given response or aborts the request with the given error.
func concludeRequest(c *gin.Context, resp interface{}, err *HttpError) {
	if err == nil {
		c.JSON(http.StatusOK, resp)
	} else {
		err.Abort(c)
	}
}

// More
type MoreResponse struct {
	Feed []core.ContentItem `json:"feed"`
}

func MoreHandler(c *gin.Context) {
	r := core.MoreRequest{}
	if c.Bind(&r) != nil {
		return
	}

	suggesterContext := core.SuggesterContext{
		c.MustGet("MDB_DB").(*sql.DB),
		c.MustGet("CHRONICLES_DB").(*sql.DB),
		c.MustGet("DATA_MODELS").(*data_models.DataModels),
		make(map[string]interface{}),
	}
	resp, err := handleMore(suggesterContext, r)
	concludeRequest(c, resp, err)
}

func handleMore(suggesterContext core.SuggesterContext, r core.MoreRequest) (*MoreResponse, *HttpError) {
	log.Infof("r: %+v", r)
	feed := core.MakeFeed(suggesterContext)
	if cis, err := feed.More(r); err != nil {
		return nil, NewInternalError(err)
	} else {
		return &MoreResponse{Feed: cis}, nil
	}
}

// Recommend
func RecommendHandler(c *gin.Context) {
	r := core.MoreRequest{}
	if c.Bind(&r) != nil {
		return
	}

	suggesterContext := core.SuggesterContext{
		c.MustGet("MDB_DB").(*sql.DB),
		c.MustGet("CHRONICLES_DB").(*sql.DB),
		c.MustGet("DATA_MODELS").(*data_models.DataModels),
		make(map[string]interface{}),
	}
	resp, err := handleRecommend(suggesterContext, r)
	log.Infof("Err: %+v", err)
	concludeRequest(c, resp, err)
}

func handleRecommend(suggesterContext core.SuggesterContext, r core.MoreRequest) (*MoreResponse, *HttpError) {
	log.Infof("r: %+v", r)
	log.Infof("Spec: %+v", r.Options.Spec)
	var recommend *recommendations.Recommender
	if r.Options.Spec == nil {
		var err error
		recommend, err = recommendations.MakeRecommender(suggesterContext)
		if err != nil {
			return nil, NewInternalError(err)
		}
	} else {
		if s, err := core.MakeSuggesterFromName(suggesterContext, r.Options.Spec.Name); err != nil {
			return nil, NewInternalError(err)
		} else {
			if err := s.UnmarshalSpec(suggesterContext, *r.Options.Spec); err != nil {
				return nil, NewInternalError(err)
			} else {
				recommend = &recommendations.Recommender{s}
			}
		}
	}

	// Uncomment to debug marshaling and unmarshling of specs.
	// log.Infof("S: %+v", recommend.Suggester)
	// if spec, err := recommend.Suggester.MarshalSpec(); err != nil {
	// 	return nil, NewInternalError(err)
	// } else {
	// 	if marshaledBytes, err := json.Marshal(spec); err != nil {
	// 		return nil, NewInternalError(err)
	// 	} else {
	// 		log.Infof("Spec as JSON: %s", string(marshaledBytes))
	// 	}
	// }
	//
	//		if s, err := core.MakeSuggesterFromName(db, spec.Name); err != nil {
	//			return nil, NewInternalError(err)
	//		} else {
	//			if err := s.UnmarshalSpec(db, spec); err != nil {
	//				return nil, NewInternalError(err)
	//			} else {
	//				if sSpec, err := s.MarshalSpec(); err != nil {
	//					return nil, NewInternalError(err)
	//				} else {
	//					if sMarshaledBytes, err := json.MarshalIndent(sSpec, "", "  "); err != nil {
	//						return nil, NewInternalError(err)
	//					} else {
	//						log.Infof("Spec as JSON: %s", string(sMarshaledBytes))
	//					}
	//				}
	//			}
	//		}
	//	}

	start := time.Now()
	if cis, err := recommend.Recommend(r); err != nil {
		return nil, NewInternalError(err)
	} else {
		log.Infof("Recommend: %+v", time.Now().Sub(start))
		utils.PrintProfile(true)
		return &MoreResponse{Feed: cis}, nil
	}
}
