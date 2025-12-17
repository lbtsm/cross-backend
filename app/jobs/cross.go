package jobs

import (
	"encoding/json"
	"fmt"
	"go-admin/app/jobs/models"
	"go-admin/internal/constants"
	"io"
	"net/http"
	"strings"
	"time"

	ormModels "go-admin/app/admin/models"

	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func InitJob() {
	jobList = map[string]JobExec{
		"CrossDataSync":       CrossDataSyncJob{},
		"SignleCrossDataSync": SignleSyncJob{},
	}
}

type CrossDataSyncJob struct {
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

	fmt.Println("startKey ---", startKey)
	// Construct API URL with parameters
	apiURL := fmt.Sprintf("%s/cross/list?limit=%d", host, limit)
	if startKey != "" {
		apiURL = fmt.Sprintf("%s&key=%s", apiURL, startKey)
	}

	fmt.Println("apiURL ---", apiURL)
	// Fetch data from API
	err = j.fetchAndStoreData(apiURL, db)
	if err != nil {
		fmt.Printf("%s [ERROR] CrossDataSyncJob failed: %v\n", time.Now().Format(timeFormat), err)
		return err
	}
	fmt.Println("ignoreCount ------- ", ignoreCount)

	return nil
}

var ignoreCount int64

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
		fmt.Println("response.key ----------- ", item.Key)
		if item.CrossSet.Src == nil {
			ignoreCount++
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

type SignleSyncJob struct {
}

func (j SignleSyncJob) Exec(arg interface{}) error {
	db := sdk.Runtime.GetDbByKey("*")
	if db == nil {
		return fmt.Errorf("database connection is nil")
	}
	// Start synchronization process
	host := fmt.Sprintf("http://%s:6041", arg)
	unCompletedList := make([]ormModels.CrossInfo, 0)
	err := db.Where("status != ?", constants.StatusOfCompleted).Order("id desc").
		Limit(20).Find(&unCompletedList).Error
	if err != nil {
		return err
	}
	if len(unCompletedList) <= 0 {
		return nil
	}

	for _, item := range unCompletedList {
		fmt.Println("unCompletedList ------- ", item.Id)
		apiURL := fmt.Sprintf("%s/cross/signle?key=cross:%s", host, item.OrderId)
		err = j.fetchAndStoreData(item.Id, apiURL, db)
		if err != nil {
			fmt.Printf("%s [ERROR] SignleSyncJob failed: %v\n", time.Now().Format(timeFormat), err)
			return err
		}

		time.Sleep(time.Millisecond * 10)
	}

	return nil
}

func (j SignleSyncJob) fetchAndStoreData(id int, apiURL string, db *gorm.DB) error {
	// Make HTTP request
	resp, err := http.Get(apiURL)
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

	type signleResponse struct {
		Data models.CrossSet `json:"data"`
	}
	var response signleResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return fmt.Errorf("failed to parse JSON response: %v", err)
	}

	if response.Data.Src == nil {
		return nil
	}
	srcInfoBytes, _ := json.Marshal(response.Data.Src)
	dstInfoBytes, _ := json.Marshal(response.Data.Dest)
	relayInfoBytes, _ := json.Marshal(response.Data.Relay)
	mapDstInfoBytes, _ := json.Marshal(response.Data.MapDst)

	relayHash := ""
	if response.Data.Relay != nil {
		relayHash = response.Data.Relay.TxHash
	}
	dstChain := ""
	dstTxHash := ""
	usedTime := int64(0)
	if response.Data.Dest != nil {
		dstChain = response.Data.Dest.Chain
		dstTxHash = response.Data.Dest.TxHash
		if response.Data.Src != nil {
			usedTime = response.Data.Dest.Timestamp - response.Data.Src.Timestamp
		}
	}
	if response.Data.Dest == nil && response.Data.MapDst != nil {
		dstChain = response.Data.MapDst.Chain
		dstTxHash = response.Data.MapDst.TxHash
	}
	mapDstTxHash := ""
	if response.Data.MapDst != nil {
		mapDstTxHash = response.Data.MapDst.TxHash
	}

	crossData := ormModels.CrossInfo{
		Project:      "", // todo
		Status:       int64(response.Data.Status),
		SrcChain:     response.Data.Src.Chain,
		SrcTxHash:    response.Data.Src.TxHash,
		RelayTxHash:  relayHash,
		DstChain:     dstChain,
		SrcInfo:      string(srcInfoBytes),
		RelayInfo:    string(relayInfoBytes),
		DstTxHash:    dstTxHash,
		DstInfo:      string(dstInfoBytes),
		MapDstTxHash: mapDstTxHash,
		MapDstInfo:   string(mapDstInfoBytes),
		CostTime:     usedTime,
		UpdatedAt:    time.Now(),
	}
	err = db.Table(crossData.TableName()).Where("id = ?", id).Updates(&crossData).Error
	if err != nil {
		return errors.Wrapf(err, "failed to update cross data with id %d", id)
	}

	return nil
}
