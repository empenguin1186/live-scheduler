package presentation

import (
	"github.com/labstack/echo/v4"
	"live-scheduler/domain"
	"net/http"
	"strconv"
	"time"
)

const LAYOUT = "2006-01-02"

type LiveHandler struct {
	liveService       domain.LiveService
	liveDescService   domain.LiveDescService
	bandService       domain.BandService
	bandMemberService domain.BandMemberService
	playerService     domain.PlayerService
}

func NewLiveHandler(
	liveService domain.LiveService,
	liveDescService domain.LiveDescService,
	bandService domain.BandService,
	bandMemberService domain.BandMemberService,
	playerService domain.PlayerService) *LiveHandler {
	return &LiveHandler{
		liveService:       liveService,
		liveDescService:   liveDescService,
		bandService:       bandService,
		bandMemberService: bandMemberService,
		playerService:     playerService,
	}
}

func (h *LiveHandler) GetLives(context echo.Context) error {
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

	lives, err := h.liveService.GetByPeriod(&start, &end)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	var liveResponse []*LiveResponse
	for _, e := range lives {
		liveResponse = append(liveResponse, NewLiveResponse(e))
	}
	return context.JSON(http.StatusOK, liveResponse)
}

func (h *LiveHandler) GetLive(context echo.Context) error {
	liveId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	liveModel, err := h.liveDescService.GetById(int(liveId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return context.JSON(http.StatusOK, NewLiveDescResponse(liveModel))
}

func (h *LiveHandler) PostLive(context echo.Context) error {
	live := new(LiveCreateRequest)
	if err := context.Bind(live); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := context.Validate(live); err != nil {
		return err
	}
	err := h.liveService.Register(live.ToModel())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return context.JSON(http.StatusOK, live)
}

func (h *LiveHandler) PatchLive(context echo.Context) error {
	live := new(LivePatchRequest)
	if err := context.Bind(live); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := context.Validate(live); err != nil {
		return err
	}
	err := h.liveService.Update(live.ToModel())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return context.JSON(http.StatusOK, live)
}

func (h *LiveHandler) DeleteLive(context echo.Context) error {
	liveId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = h.liveService.Delete(int(liveId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return context.NoContent(http.StatusOK)
}

type BandHandler struct {
	bandService       domain.BandService
	bandMemberService domain.BandMemberService
}

func NewBandHandler(
	bandService domain.BandService,
	bandMemberService domain.BandMemberService) *BandHandler {
	return &BandHandler{
		bandService:       bandService,
		bandMemberService: bandMemberService,
	}
}

func (h *LiveHandler) GetBand(context echo.Context) error {
	liveId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	bands, err := h.bandService.GetByLiveId(int(liveId))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	var response []*BandResponsePart
	for _, e := range bands {
		response = append(response, NewBandResponsePart(e))
	}
	return context.JSON(http.StatusOK, response)
}

func (h *LiveHandler) PostBand(context echo.Context) error {
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

	err = h.bandService.Register(band.ToModel())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return context.JSON(http.StatusOK, band)
}

func (h *LiveHandler) PatchBand(context echo.Context) error {
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
	err = h.bandService.Update(int(liveId), int(turn), band.ToModel())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return context.JSON(http.StatusOK, band)
}

func (h *LiveHandler) DeleteBand(context echo.Context) error {
	liveId, err := strconv.ParseInt(context.Param("live_id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	turn, err := strconv.ParseInt(context.Param("turn"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	err = h.bandService.Delete(int(liveId), int(turn))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return context.NoContent(http.StatusOK)
}

func (h *LiveHandler) GetPart(context echo.Context) error {
	part := domain.Part(context.QueryParam("part"))
	players, err := h.playerService.GetByPart(&part)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	var response []*MemberResponsePart
	for _, p := range players {
		response = append(response, NewPlayerResponse(p))
	}
	return context.JSON(http.StatusOK, response)
}

func (h *LiveHandler) PostPart(context echo.Context) error {
	player := new(PlayerRequest)
	if err := context.Bind(player); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := context.Validate(player); err != nil {
		return err
	}
	err := h.playerService.Register(player.ToModel())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return context.JSON(http.StatusOK, player)
}

func (h *LiveHandler) DeletePart(context echo.Context) error {
	player := new(PlayerRequest)
	if err := context.Bind(player); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := context.Validate(player); err != nil {
		return err
	}
	err := h.playerService.Delete(player.ToModel())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return context.JSON(http.StatusOK, player)
}
