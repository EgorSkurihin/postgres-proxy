package pgspy

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"time"
)

type QueriesCollector struct {
	data         *SessionData
	dataFilePath string
}

type SessionData struct {
	ClientIP   string     `json:"client_ip"`
	StartedAt  string     `json:"started_at"`
	EndedAt    string     `json:"ended_at"`
	SQLQueries []SQLQuery `json:"sql_queries"`
}

type SQLQuery struct {
	Query     string `json:"query"`
	IsSuccess bool   `json:"is_success"`
}

func newQueriesCollector(clientIP string, dataFilePath string) *QueriesCollector {
	return &QueriesCollector{
		data: &SessionData{
			StartedAt:  time.Now().Format(time.DateTime),
			ClientIP:   clientIP,
			SQLQueries: []SQLQuery{},
		},
		dataFilePath: dataFilePath,
	}
}

func (qc *QueriesCollector) addQuery(query *SQLQuery) {
	qc.data.SQLQueries = append(qc.data.SQLQueries, *query)
}

func (qc *QueriesCollector) saveDataToFile() {
	if len(qc.data.SQLQueries) == 0 {
		return
	}
	qc.data.EndedAt = time.Now().Format(time.DateTime)

	existingData, err := os.ReadFile(qc.dataFilePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("Unable to read data file: %v", err)
		return
	}

	var existingSessions []SessionData
	if err == nil || errors.Is(err, os.ErrNotExist) {
		json.Unmarshal(existingData, &existingSessions)
	} else {
		existingSessions = []SessionData{}
	}

	existingSessions = append([]SessionData{*qc.data}, existingSessions...)

	jsonData, err := json.Marshal(existingSessions)
	if err != nil {
		log.Printf("Failed to save queries info to file %v:", err)
		return
	}

	err = os.WriteFile(qc.dataFilePath, jsonData, 0644)
	if err != nil {
		log.Printf("Failed to write sessions in file %v:", err)
		return
	}
	log.Println("Saved sessions data in file")
}
