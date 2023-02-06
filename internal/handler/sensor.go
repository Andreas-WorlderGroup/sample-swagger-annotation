package sensor

import (
	error_response "andreas/internal/repository"
	"andreas/internal/repository/mysql"
	"database/sql"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Sensor interface {
	GetData(ctx echo.Context) (err error)
	StoreData(ctx echo.Context) (err error)
}

type Dependency struct {
	mysql.MySQL
	zerolog.Logger
}

type instance struct {
	Dependency
}

func New(dep Dependency) Sensor {
	return &instance{dep}
}

type ReqGetDataById struct {
	ID1 int    `query:"id1"` // not provide this field mean ID1 is zero
	ID2 string `query:"id2" validate:"required"`
}

type SensorData struct {
	Value     int    `json:"value"` // not provide this field mean value is zero
	ID1       int    `json:"id1"`   // not provide this field mean ID1 is zero
	ID2       string `json:"id2" validate:"required"`
	Timestamp string `json:"timestamp" validate:"required"`
}

type ReqGetDataByTimestamp struct {
	StartTimestamp int `query:"start_timestamp" validate:"required"`
	EndTimestamp   int `query:"end_timestamp" validate:"required"`
}

var (
	Location, _ = time.LoadLocation("Asia/Jakarta")
	FORMAT_DATE = "Mon 01/02/2006-15:04:05"
)

// GetData godoc
// @Summary Get data from database.
// @Description get real time data from online database.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} []mysql.SensorDataDTO
// @Failure 500 {object} error_response.ErrorData
// @Router / [get]
func (x *instance) GetData(ctx echo.Context) (err error) {
	reqId := new(ReqGetDataById)
	reqTime := new(ReqGetDataByTimestamp)
	flag := 0

	if err = ctx.Bind(reqId); err == nil {
		if err = ctx.Validate(*reqId); err == nil {
			flag = 1
		}
	}

	if flag == 0 {
		if err = ctx.Bind(reqTime); err == nil {
			if err = ctx.Validate(*reqTime); err == nil {
				flag = 2
			}
		}
	}

	var datas []mysql.SensorDataDTO
	if flag == 0 {
		return error_response.CreateError(ctx, "bad_request")
	} else if flag == 1 {
		datas, err = x.GetDataByID(ctx.Request().Context(), reqId.ID1, reqId.ID2)
	} else if flag == 2 {
		datas, err = x.GetDataByTimestamp(ctx.Request().Context(), reqTime.StartTimestamp, reqTime.EndTimestamp)
	}

	if err != nil && err == sql.ErrNoRows {
		x.Logger.Error().Msg("failed get data: " + err.Error())
		return error_response.CreateError(ctx, "data_not_found")
	} else if err != nil {
		x.Logger.Error().Msg("failed get data: " + err.Error())
		return error_response.CreateError(ctx, "general_error")
	}

	return ctx.JSON(http.StatusOK, datas)
}

func (x *instance) StoreData(ctx echo.Context) (err error) {
	reqData := new(SensorData)

	if err = ctx.Bind(reqData); err != nil {
		return error_response.CreateError(ctx, "bad_request")
	}

	if err = ctx.Validate(*reqData); err != nil {
		return error_response.CreateError(ctx, "bad_request")
	}

	// date to unix
	tm, err := time.ParseInLocation(FORMAT_DATE, reqData.Timestamp, Location)
	if err != nil {
		x.Logger.Error().Msg("wrong date format")
	}
	dt := tm.Unix()

	// transform request
	data := mysql.SensorDataDB{
		Value:     &reqData.Value,
		ID1:       &reqData.ID1,
		ID2:       &reqData.ID2,
		Timestamp: &dt,
	}

	tx, err := x.SaveData(ctx.Request().Context(), data)
	if err != nil {
		x.Logger.Error().Msg("failed insert: " + err.Error())
		return mysql.RollbackTxSql([]*sql.Tx{tx}, err)
	}

	mysql.CommitTxSql([]*sql.Tx{tx})

	return ctx.String(http.StatusOK, "success insert")
}
