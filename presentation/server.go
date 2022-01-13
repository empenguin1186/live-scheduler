package presentation

import (
	"github.com/labstack/echo/v4"
	"live-scheduler/domain"
	"net/http"
	"strconv"
	"time"
)

const LAYOUT = "2006-01-02"

type LiveServer struct {
	e                 *echo.Echo
	liveService       domain.LiveService
	liveDescService   domain.LiveDescService
	bandService       domain.BandService
	bandMemberService domain.BandMemberService
}

func NewLiveServer(
	e *echo.Echo,
	liveService domain.LiveService,
	liveDescService domain.LiveDescService,
	bandService domain.BandService,
	bandMemberService domain.BandMemberService) *LiveServer {
	return &LiveServer{
		e:                 e,
		liveService:       liveService,
		liveDescService:   liveDescService,
		bandService:       bandService,
		bandMemberService: bandMemberService,
	}
}

func (s *LiveServer) Start() {
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

	s.e.GET("/live/:id/band", func(context echo.Context) error {
		liveId, err := strconv.ParseInt(context.Param("id"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		bands, err := s.bandService.GetByLiveId(int(liveId))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		var response []*BandResponsePart
		for _, e := range bands {
			response = append(response, NewBandResponsePart(e))
		}
		return context.JSON(http.StatusOK, response)
	})

	s.e.POST("/live/:id/band", func(context echo.Context) error {
		liveId, err := strconv.ParseInt(context.Param("id"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		band := new(BandCreateRequest)
		if err := context.Bind(band); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := context.Validate(band); err != nil {
			return err
		}
		if band.LiveId != int(liveId) {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		err = s.bandService.Register(band.ToModel())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}

		return context.JSON(http.StatusOK, band)
	})

	s.e.PATCH("/live/:live_id/band/:turn", func(context echo.Context) error {
		liveId, err := strconv.ParseInt(context.Param("live_id"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		turn, err := strconv.ParseInt(context.Param("turn"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		band := new(BandPatchRequest)
		if err := context.Bind(band); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		if err := context.Validate(band); err != nil {
			return err
		}
		err = s.bandService.Update(int(liveId), int(turn), band.ToModel())
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return context.JSON(http.StatusOK, band)
	})

	s.e.DELETE("/live/:live_id/band/:turn", func(context echo.Context) error {
		liveId, err := strconv.ParseInt(context.Param("live_id"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		turn, err := strconv.ParseInt(context.Param("turn"), 10, 64)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
		err = s.bandService.Delete(int(liveId), int(turn))
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
		return context.NoContent(http.StatusOK)
	})

	s.e.Logger.Fatal(s.e.Start(":1323"))
}
