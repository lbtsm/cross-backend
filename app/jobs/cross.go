package jobs

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-admin/app/jobs/models"
	"io"
	"net/http"
	"strings"
	"time"

	ormModels "go-admin/app/admin/models"

	"github.com/go-admin-team/go-admin-core/sdk"
	"gorm.io/gorm"
)

func InitJob() {
	jobList = map[string]JobExec{
		"CrossDataSyncJob": CrossDataSyncJob{},
	}
}

type CrossDataSyncJob struct {
}

type CrossListItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	// Add other fields based on actual API response structure
}

// CrossListResponse represents the API response structure
type CrossListResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    []CrossListItem `json:"data"`
	NextKey string          `json:"next_key"` // This will be used as the key for next request
}

func (j CrossDataSyncJob) Exec(arg interface{}) error {
	db := sdk.Runtime.GetDbByKey("*")
	// Check if database connection is valid
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	// Start synchronization process
	startKey := ""
	limit := 50
	host := fmt.Sprintf("http://%s:6041", arg)
	lastOne := &ormModels.CrossInfo{}
	err := db.Order("id desc").First(lastOne).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if lastOne.OrderId != "" {
		startKey = "cross:" + lastOne.OrderId
	}

	fmt.Println("startKey --------------- ", startKey)
	// Construct API URL with parameters
	apiURL := fmt.Sprintf("%s/cross/list?limit=%d", host, limit)
	if startKey != "" {
		apiURL = fmt.Sprintf("%s&key=%s", apiURL, startKey)
	}

	// Fetch data from API
	err = j.fetchAndStoreData(apiURL, db)
	if err != nil {
		fmt.Printf("%s [ERROR] CrossDataSyncJob failed: %v\n", time.Now().Format(timeFormat), err)
		return err
	}

	return nil
}

// fetchAndStoreData fetches data from API and stores it in database
func (j CrossDataSyncJob) fetchAndStoreData(url string, db *gorm.DB) error {
	// Make HTTP request
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	var response models.CrossListResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}
	if len(response.Data) <= 0 { // list end
		return nil
	}

	for _, item := range response.Data {
		if item.CrossSet.Src == nil {
			continue
		}
		srcInfoBytes, _ := json.Marshal(item.CrossSet.Src)
		dstInfoBytes, _ := json.Marshal(item.CrossSet.Dest)
		relayInfoBytes, _ := json.Marshal(item.CrossSet.Relay)
		mapDstInfoBytes, _ := json.Marshal(item.CrossSet.MapDst)

		relayHash := ""
		if item.CrossSet.Relay != nil {
			relayHash = item.CrossSet.Relay.TxHash
		}
		dstChain := ""
		dstTxHash := ""
		usedTime := int64(0)
		if item.CrossSet.Dest != nil {
			dstChain = item.CrossSet.Dest.Chain
			dstTxHash = item.CrossSet.Dest.TxHash
			if item.CrossSet.Src != nil {
				usedTime = item.CrossSet.Dest.Timestamp - item.CrossSet.Src.Timestamp
			}
		}
		if item.CrossSet.Dest == nil && item.CrossSet.MapDst != nil {
			dstChain = item.CrossSet.MapDst.Chain
			dstTxHash = item.CrossSet.MapDst.TxHash
		}
		mapDstTxHash := ""
		if item.CrossSet.MapDst != nil {
			mapDstTxHash = item.CrossSet.MapDst.TxHash
		}
		createdAt := int64(0)
		if item.CrossSet.Src != nil {
			createdAt = item.CrossSet.Src.Timestamp
		}
		fmt.Println("item --------- ", item)

		crossData := ormModels.CrossInfo{
			Project:      "",
			Status:       int64(item.CrossSet.Status),
			OrderId:      strings.Split(item.Key, ":")[1],
			SrcChain:     item.CrossSet.Src.Chain,
			SrcTxHash:    item.CrossSet.Src.TxHash,
			RelayTxHash:  relayHash,
			DstChain:     dstChain,
			SrcInfo:      string(srcInfoBytes),
			RelayInfo:    string(relayInfoBytes),
			DstTxHash:    dstTxHash,
			DstInfo:      string(dstInfoBytes),
			MapDstTxHash: mapDstTxHash,
			MapDstInfo:   string(mapDstInfoBytes),
			CostTime:     usedTime,
			CreatedAt:    time.Unix(createdAt, 0),
			UpdatedAt:    time.Now(),
		}

		// Try to create record, or update if exists
		err := db.Create(&crossData).Error
		if err != nil {
			if strings.Contains(err.Error(), "Duplicate") {
				continue
			}
			fmt.Printf("%s [WARN] Failed to save cross data with key %s: %v\n",
				time.Now().Format(timeFormat), item.Key, err)
			return nil
		}
		time.Sleep(time.Millisecond * 50)
	}

	return nil
}
