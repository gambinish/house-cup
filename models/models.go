package models

type Tournament struct {
	ID              int64   `json:"id"`
	Tournament_Name string  `json:"tournament_name"`
	Created_At      string  `json:"created_at"`
	Ended_At        *string `json:"ended_at"`
}

type House struct {
	ID            int64  `json:"id"`
	House_Name    string `json:"house_name"`
	House_Points  int64  `json:"house_points"`
	Tournament_ID int64  `json:"tournament_id"`
}

type Student struct {
	ID           int64  `json:"id"`
	Student_Name string `json:"student_name"`
	Points       int64  `json:"points"`
	House_ID     int64  `json:"house_id"`
}

type Point struct {
	ID         int64  `json:"id"`
	Points     int64  `json:"points"`
	Notes      string `json:"notes"`
	Student_ID *int64 `json:"student_id"`
	House_ID   int64  `json:"house_id"`
}
