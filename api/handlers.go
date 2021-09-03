package api

import (
	"database/sql"

	// "encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"

	"github.com/Bnei-Baruch/feed-api/core"
	"github.com/Bnei-Baruch/feed-api/data_models"
	"github.com/Bnei-Baruch/feed-api/recommendations"
	"github.com/Bnei-Baruch/feed-api/utils"
)

// Premetheus handler.
func PrometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

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
	if err != nil {
		log.Infof("Err: %+v", err)
	}
	concludeRequest(c, resp, err)
}

func handleMore(suggesterContext core.SuggesterContext, r core.MoreRequest) (*MoreResponse, *HttpError) {
	log.Debugf("r: %+v", r)
	var feed *core.Feed
	if r.Options.Spec == nil {
		feed = core.MakeFeed(suggesterContext)
	} else {
		if s, err := MakeAndUnmarshalSuggester(suggesterContext, r.Options.Spec); err != nil {
			return nil, NewInternalError(err)
		} else {
			feed = core.MakeFeedFromSuggester(s, suggesterContext)
		}
	}
	if cis, err := feed.More(r); err != nil {
		return nil, NewInternalError(err)
	} else {
		return &MoreResponse{Feed: cis}, nil
	}
}

// Watching Now
type WatchingNowRequest struct {
	Uids []string `json:"uids,omitempty" form:"uids,omitempty"`
}

type WatchingNowResponse struct {
	WatchingNow []int64 `json:"watching_now,omitempty" form:"watching_now,omitempty"`
}

func WatchingNowHandler(c *gin.Context) {
	r := WatchingNowRequest{}
	if c.Bind(&r) != nil {
		return
	}

	dm := c.MustGet("DATA_MODELS").(*data_models.DataModels)
	resp := WatchingNowResponse{}
	if watchingNow, err := dm.SqlDataModel.WatchingNow(r.Uids); err != nil {
		concludeRequest(c, resp, NewInternalError(err))
	} else {
		resp.WatchingNow = watchingNow
		concludeRequest(c, resp, nil)
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

	dm := c.MustGet("DATA_MODELS").(*data_models.DataModels)
	resp := ViewsResponse{}
	if views, err := dm.SqlDataModel.Views(r.Uids); err != nil {
		concludeRequest(c, resp, NewInternalError(err))
	} else {
		resp.Views = views
		concludeRequest(c, resp, nil)
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
		c.MustGet("DATA_MODELS").(*data_models.DataModels),
		make(map[string]interface{}),
	}
	resp, err := handleRecommend(suggesterContext, r)
	if err != nil {
		log.Infof("Err: %+v", err)
	}
	concludeRequest(c, resp, err)
}

func MakeAndUnmarshalSuggester(suggesterContext core.SuggesterContext, spec *core.SuggesterSpec) (core.Suggester, error) {
	if spec.Name == "Default" {
		return core.MakeDefaultSuggester(suggesterContext)
	}
	if s, err := core.MakeSuggesterFromName(suggesterContext, spec.Name); err != nil {
		return nil, err
	} else {
		if err := s.UnmarshalSpec(suggesterContext, *spec); err != nil {
			return nil, err
		} else {
			return s, nil
		}
	}
}

func MakeAndUnmarshalRecommender(suggesterContext core.SuggesterContext, spec *core.SuggesterSpec) (*recommendations.Recommender, error) {
	if s, err := MakeAndUnmarshalSuggester(suggesterContext, spec); err == nil {
		return &recommendations.Recommender{s}, nil
	} else {
		return nil, err
	}
}

func handleRecommend(suggesterContext core.SuggesterContext, r core.MoreRequest) (*MoreResponse, *HttpError) {
	log.Debugf("r: %+v", r)
	log.Debugf("Spec: %+v", r.Options.Spec)
	log.Debugf("Specs: %+v", r.Options.Specs)
	var recommend *recommendations.Recommender
	var recommends []*recommendations.Recommender
	if r.Options.Spec == nil && r.Options.Specs == nil {
		var err error
		recommend, err = recommendations.MakeRecommender(suggesterContext)
		if err != nil {
			return nil, NewInternalError(err)
		}
	} else if r.Options.Specs == nil {
		if rec, err := MakeAndUnmarshalRecommender(suggesterContext, r.Options.Spec); err != nil {
			return nil, NewInternalError(err)
		} else {
			recommend = rec
		}
	} else {
		for i, spec := range r.Options.Specs {
			log.Debugf("Spec %d: %+v", i, spec)
			if rec, err := MakeAndUnmarshalRecommender(suggesterContext, spec); err != nil {
				return nil, NewInternalError(err)
			} else {
				log.Debugf("Rec: %+v", rec)
				recommends = append(recommends, rec)
			}
		}
	}

	// Uncomment to debug marshaling and unmarshling of specs.
	// log.Debugf("S: %+v", recommend.Suggester)
	// if spec, err := recommend.Suggester.MarshalSpec(); err != nil {
	// 	return nil, NewInternalError(err)
	// } else {
	// 	if marshaledBytes, err := json.Marshal(spec); err != nil {
	// 		return nil, NewInternalError(err)
	// 	} else {
	// 		log.Debugf("Spec as JSON: %s", string(marshaledBytes))
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
	//						log.Debugf("Spec as JSON: %s", string(sMarshaledBytes))
	//					}
	//				}
	//			}
	//		}
	//	}

	skipUidsMap := make(map[string]bool)
	for _, uid := range r.Options.SkipUids {
		skipUidsMap[uid] = true
	}

	if len(recommends) > 0 {
		log.Debugf("len: %+v", len(recommends))
		res := &MoreResponse{}
		for i, rec := range recommends {
			log.Debugf("Rec: %+v", rec)
			start := time.Now()
			skipUids := []string(nil)
			for uid, _ := range skipUidsMap {
				skipUids = append(skipUids, uid)
			}
			r.Options.SkipUids = skipUids
			if cis, err := rec.Recommend(r); err != nil {
				return nil, NewInternalError(err)
			} else {
				for _, ci := range cis {
					skipUidsMap[ci.UID] = true
				}
				// log.Debugf("cis: %+v", cis)
				log.Debugf("Recommend[%d]: %+v", i, time.Now().Sub(start))
				utils.PrintProfile(true)
				res.Feeds = append(res.Feeds, cis)
			}
		}
		// log.Debugf("res: %+v", res)
		return res, nil
	} else {
		start := time.Now()
		if cis, err := recommend.Recommend(r); err != nil {
			return nil, NewInternalError(err)
		} else {
			log.Debugf("Recommend: %+v", time.Now().Sub(start))
			utils.PrintProfile(true)
			return &MoreResponse{Feed: cis}, nil
		}
	}
}
