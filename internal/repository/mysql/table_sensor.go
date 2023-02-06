package mysql

import (
	"context"
	"database/sql"
	"time"
)

type SensorDataDB struct {
	Value     *int    `json:"value"`
	ID1       *int    `json:"id1"`
	ID2       *string `json:"id2"`
	Timestamp *int64  `json:"timestamp"`
}

type SensorDataDTO struct {
	Value     int    `json:"value"`
	ID1       int    `json:"id1"`
	ID2       string `json:"id2"`
	Timestamp string `json:"timestamp"`
}

type TableSensor interface {
	GetDataByID(ctx context.Context, id1 int, id2 string) (data []SensorDataDTO, err error)
	GetDataByTimestamp(ctx context.Context, start, end int) (data []SensorDataDTO, err error)
	SaveData(ctx context.Context, data SensorDataDB) (tx *sql.Tx, err error)
}

func (x *instance) GetDataByID(ctx context.Context, id1 int, id2 string) (data []SensorDataDTO, err error) {
	return x.SelectDataSensor(ctx, "Select value, id1, id2, timestamp from sensor where id1 = ? and id2 = ?", id1, id2)
}

func (x *instance) GetDataByTimestamp(ctx context.Context, start, end int) (data []SensorDataDTO, err error) {
	return x.SelectDataSensor(ctx, "Select value, id1, id2, timestamp from sensor where timestamp between ? and ?", start, end)
}

func (x *instance) SelectDataSensor(ctx context.Context, query string, param ...any) (data []SensorDataDTO, err error) {
	stmt, err := x.QueryContext(ctx, query, param...)
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	for stmt.Next() {
		var temp SensorDataDB
		err = stmt.Scan(&temp.Value, &temp.ID1, &temp.ID2, &temp.Timestamp)
		if err != nil {
			return nil, err
		}

		data = append(data, SensorDataDTO{
			Value:     *temp.Value,
			ID1:       *temp.ID1,
			ID2:       *temp.ID2,
			Timestamp: time.Unix(*temp.Timestamp, 0).Format("Mon 01/02/2006-15:04:05"),
		})
	}

	if len(data) == 0 {
		return nil, sql.ErrNoRows
	}

	return
}

func (x *instance) SaveData(ctx context.Context, data SensorDataDB) (tx *sql.Tx, err error) {
	query := "insert sensor (value, id1, id2, timestamp) values (?, ?, ?, ?)"
	return x.ExecContext(ctx, query, data.Value, data.ID1, data.ID2, data.Timestamp)
}
