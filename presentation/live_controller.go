package presentation

import (
	"github.com/labstack/echo/v4"
	"live-scheduler/domain"
	"net/http"
	"strconv"
	"time"
)

const LAYOUT = "2006-01-02"

type Server struct {
	e               *echo.Echo
	liveService     domain.LiveService
	liveDescService domain.LiveDescService
}

func NewLiveController(e *echo.Echo, liveService domain.LiveService, liveDescService domain.LiveDescService) *Server {
	return &Server{
		e:               e,
		liveService:     liveService,
		liveDescService: liveDescService,
	}
}

func (s *Server) Start() {
	s.e.Validator = NewCustomValidator()
	s.e.GET("/live", func(context echo.Context) error {
		var start time.Time
		err := echo.QueryParamsBinder(context).Time("start", &start, LAYOUT).BindError()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		var end time.Time
		err = echo.QueryParamsBinder(context).Time("end", &end, LAYOUT).BindError()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		lives, err := s.liveService.GetByPeriod(&start, &end)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		var liveResponse []*LiveResponse
		for _, e := range lives {
			liveResponse = append(liveResponse, NewLiveResponse(e))
		}
		return context.JSON(http.StatusOK, liveResponse)
	})

	s.e.GET("/live/:id", func(context echo.Context) error {
		liveId, err := strconv.ParseInt(context.Param("id"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		liveModel, err := s.liveDescService.GetById(int(liveId))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return context.JSON(http.StatusOK, NewLiveDescResponse(liveModel))
	})

	s.e.POST("/live", func(context echo.Context) error {
		live := new(LiveCreateRequest)
		if err := context.Bind(live); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := context.Validate(live); err != nil {
			return err
		}
		err := s.liveService.Register(live.ToModel())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return context.JSON(http.StatusOK, live)
	})

	s.e.PATCH("/live", func(context echo.Context) error {
		live := new(LivePatchRequest)
		if err := context.Bind(live); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := context.Validate(live); err != nil {
			return err
		}
		err := s.liveService.Update(live.ToModel())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return context.JSON(http.StatusOK, live)
	})

	s.e.DELETE("/live/:id", func(context echo.Context) error {
		liveId, err := strconv.ParseInt(context.Param("id"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err = s.liveService.Delete(int(liveId))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return context.NoContent(http.StatusOK)
	})

	s.e.Logger.Fatal(s.e.Start(":1323"))
}
