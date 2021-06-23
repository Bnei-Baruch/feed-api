package api

import (
	"database/sql"
	"errors"

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
	Feed  []core.ContentItem   `json:"feed"`
	Feeds [][]core.ContentItem `json:"feeds"`
}

func MoreHandler(c *gin.Context) {
	r := core.MoreRequest{}
	if c.Bind(&r) != nil {
		return
	}

	suggesterContext := core.SuggesterContext{
		c.MustGet("MDB_DB").(*sql.DB),
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

// Views
type ViewsRequest struct {
	Uids []string `json:"uids,omitempty" form:"uids,omitempty"`
}

type ViewsResponse struct {
	Views []int64 `json:"views,omitempty" form:"views,omitempty"`
}

func ViewsHandler(c *gin.Context) {
	r := ViewsRequest{}
	if c.Bind(&r) != nil {
		return
	}

	resp := ViewsResponse{}
	concludeRequest(c, resp, NewInternalError(errors.New("Not Implemented")))
}

// Recommend
func RecommendHandler(c *gin.Context) {
	r := core.MoreRequest{}
	if c.Bind(&r) != nil {
		return
	}

	suggesterContext := core.SuggesterContext{
		c.MustGet("MDB_DB").(*sql.DB),
		c.MustGet("DATA_MODELS").(*data_models.DataModels),
		make(map[string]interface{}),
	}
	resp, err := handleRecommend(suggesterContext, r)
	log.Infof("Err: %+v", err)
	concludeRequest(c, resp, err)
}

func MakeAndUnmarshal(suggesterContext core.SuggesterContext, spec *core.SuggesterSpec) (*recommendations.Recommender, error) {
	if spec.Name == "Default" {
		return recommendations.MakeRecommender(suggesterContext)
	}
	if s, err := core.MakeSuggesterFromName(suggesterContext, spec.Name); err != nil {
		return nil, err
	} else {
		if err := s.UnmarshalSpec(suggesterContext, *spec); err != nil {
			return nil, err
		} else {
			return &recommendations.Recommender{s}, nil
		}
	}
}

func handleRecommend(suggesterContext core.SuggesterContext, r core.MoreRequest) (*MoreResponse, *HttpError) {
	log.Infof("r: %+v", r)
	log.Infof("Spec: %+v", r.Options.Spec)
	log.Infof("Specs: %+v", r.Options.Specs)
	var recommend *recommendations.Recommender
	var recommends []*recommendations.Recommender
	if r.Options.Spec == nil && r.Options.Specs == nil {
		var err error
		recommend, err = recommendations.MakeRecommender(suggesterContext)
		if err != nil {
			return nil, NewInternalError(err)
		}
	} else if r.Options.Specs == nil {
		if rec, err := MakeAndUnmarshal(suggesterContext, r.Options.Spec); err != nil {
			return nil, NewInternalError(err)
		} else {
			recommend = rec
		}
	} else {
		for i, spec := range r.Options.Specs {
			log.Infof("Spec %d: %+v", i, spec)
			if rec, err := MakeAndUnmarshal(suggesterContext, spec); err != nil {
				return nil, NewInternalError(err)
			} else {
				recommends = append(recommends, rec)
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

	if len(recommends) > 0 {
		res := &MoreResponse{}
		for i, rec := range recommends {
			start := time.Now()
			if cis, err := rec.Recommend(r); err != nil {
				return nil, NewInternalError(err)
			} else {
				log.Infof("cis: %+v", cis)
				log.Infof("Recommend[%d]: %+v", i, time.Now().Sub(start))
				utils.PrintProfile(true)
				res.Feeds = append(res.Feeds, cis)
			}
		}
		log.Infof("res: %+v", res)
		return res, nil
	} else {
		start := time.Now()
		if cis, err := recommend.Recommend(r); err != nil {
			return nil, NewInternalError(err)
		} else {
			log.Infof("Recommend: %+v", time.Now().Sub(start))
			utils.PrintProfile(true)
			return &MoreResponse{Feed: cis}, nil
		}
	}
}
