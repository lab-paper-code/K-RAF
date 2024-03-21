package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
	"gopkg.in/resty.v1"
)

const (
	prometheusServiceIP       = "155.230.36.27"
	prometheusPort            = "30803"
	query                     = "sum(node_filesystem_avail_bytes) by (node)"
	volumeSizeMinimum   int64 = 1024 * 1024 * 1024 // 1GB
)

var (
	currentTotalVolumeSizeAvail, volumeErr = getAvailVolumeSize()
	D                                      = 100
	currentVolumeSizeAvail                 = int64(currentTotalVolumeSizeAvail / D)
)

func initializeVolumeInfo() {
	// Call getAvailVolumeSize() when setupVolumeRouter()
	fmt.Print(currentVolumeSizeAvail)
	// Handle the error if necessary
	if volumeErr != nil {
		// Handle the error, for example, print it or log it.
		fmt.Println("Error:", volumeErr)
	}
}

func getAvailVolumeSize() (int, error) {
	url := fmt.Sprintf("http://%s:%s/api/v1/query", prometheusServiceIP, prometheusPort)
	data := map[string]string{"query": query}

	client := resty.New()
	resp, err := client.R().
		SetFormData(data).
		Post(url)

	if err != nil {
		return 0, err
	}

	if resp.StatusCode() == http.StatusOK {
		var responseData map[string]interface{}
		if err := json.Unmarshal(resp.Body(), &responseData); err != nil {
			return 0, err
		}

		// Assuming the numbers are in a field called "values"
		result := responseData["data"].(map[string]interface{})["result"].([]interface{})
		sum := 0

		for _, value := range result {
			data := value.(map[string]interface{})["value"].([]interface{})

			// Adjust the index or key based on your JSON structure
			if len(data) >= 2 {
				// Check if the value at index 1 is a float64 or a string
				if num, ok := data[1].(float64); ok {
					sum += int(num)
				} else if numStr, ok := data[1].(string); ok {
					num, err := strconv.ParseFloat(numStr, 64)
					if err != nil {
						return 0, err
					}
					sum += int(num)
				}
			}
		}

		return sum, nil
	}

	return 0, fmt.Errorf("failed to execute the query. Status code: %d", resp.StatusCode())
}

// setupVolumeRouter setup http request router for volume
func (adapter *RESTAdapter) setupVolumeRouter() {
	// initialize currentTotalVolumeSizeAvail
	initializeVolumeInfo()
	// any devices can call these APIs
	adapter.router.GET("/volumes", adapter.basicAuthDeviceOrAdmin, adapter.handleListVolumes)
	adapter.router.GET("/volumes/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleGetVolume)
	adapter.router.POST("/volumes", adapter.basicAuthDeviceOrAdmin, adapter.handleCreateVolume)
	adapter.router.PATCH("/volumes/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleUpdateVolume)
	adapter.router.DELETE("/volumes/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleDeleteVolume)

	adapter.router.POST("/mounts/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleMountVolume)
	adapter.router.DELETE("/mounts/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleUnmountVolume)
}

func (adapter *RESTAdapter) handleListVolumes(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleListVolumes",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)

	type listOutput struct {
		Volumes []types.Volume `json:"volumes"`
	}

	output := listOutput{}

	if adapter.isAdminUser(user) {
		// admin - returns all volumes
		volumes, err := adapter.logic.ListAllVolumes()
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		output.Volumes = volumes
	} else {
		err := types.ValidateDeviceID(user)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// device - returns mine
		volumes, err := adapter.logic.ListVolumes(user)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		output.Volumes = volumes
	}

	// success
	c.JSON(http.StatusOK, output)
}

func (adapter *RESTAdapter) handleGetVolume(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleGetVolume",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	volumeID := c.Param("id")

	err := types.ValidateVolumeID(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	volume, err := adapter.logic.GetVolume(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) && volume.DeviceID != user {
		// requestiong other's volume info
		err := xerrors.Errorf("failed to get volume %s, you cannot access other devices' volume info", volumeID)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	// success
	c.JSON(http.StatusOK, volume)
}

func (adapter *RESTAdapter) handleCreateVolume(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleCreateVolume",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)

	type volumeCreationRequest struct {
		DeviceID   string `json:"device_id,omitempty"`
		VolumeSize string `json:"volume_size,omitempty"`
	}

	var input volumeCreationRequest

	err := c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	volumeSizeNum := types.SizeStringToNum(input.VolumeSize)
	if volumeSizeNum < volumeSizeMinimum {
		logger.Debugf("you cannot give volume size lesser than %d, set to %d", volumeSizeMinimum, volumeSizeMinimum)
		volumeSizeNum = volumeSizeMinimum
	}

	// check for available volumesize
	currentTotalVolumeSizeAvail, volumeErr = getAvailVolumeSize()
	if volumeErr != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// if user request exceeds available volumesize, set to currentVolumeSizeAvail
	if volumeSizeNum > currentVolumeSizeAvail {
		logger.Debugf("you cannot give volume size more than %d, set to %d", currentVolumeSizeAvail, currentVolumeSizeAvail)
		volumeSizeNum = currentVolumeSizeAvail
	}

	volume := types.Volume{
		ID:         types.NewVolumeID(),
		VolumeSize: volumeSizeNum,
	}

	if adapter.isAdminUser(user) {
		volume.DeviceID = input.DeviceID
	} else {
		volume.DeviceID = user
	}

	err = types.ValidateDeviceID(volume.DeviceID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = adapter.logic.CreateVolume(&volume)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, volume)
}

func (adapter *RESTAdapter) handleUpdateVolume(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleUpdateVolume",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	volumeID := c.Param("id")

	err := types.ValidateVolumeID(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type volumeUpdateRequest struct {
		VolumeSize string `json:"volume_size,omitempty"`
	}

	var input volumeUpdateRequest

	err = c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.VolumeSize) == 0 {
		// no change
		err := xerrors.Errorf("failed to update volume %s, no change", volumeID)
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// resize
	volumeSizeNum := types.SizeStringToNum(input.VolumeSize)
	if volumeSizeNum < volumeSizeMinimum {
		logger.Debugf("you cannot give volume size lesser than %d, set to %d", volumeSizeMinimum, volumeSizeMinimum)
		volumeSizeNum = volumeSizeMinimum
	}

	// check for available volumesize
	currentTotalVolumeSizeAvail, volumeErr = getAvailVolumeSize()
	if volumeErr != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// if user request exceeds available volumesize, set to currentVolumeSizeAvail
	if volumeSizeNum > currentVolumeSizeAvail {
		logger.Debugf("you cannot give volume size more than %d, set to %d", currentVolumeSizeAvail, currentVolumeSizeAvail)
		volumeSizeNum = currentVolumeSizeAvail
	}

	volume, err := adapter.logic.GetVolume(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) && volume.DeviceID != user {
		// requestiong other's volume info
		err := xerrors.Errorf("failed to get volume %s, you cannot access other devices' volume info", volumeID)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("ID: %s\tVolumeSize: %d", volumeID, volumeSizeNum)

	if volume.VolumeSize == volumeSizeNum {
		// no change
		err := xerrors.Errorf("failed to resize volume %s, no size change, current %d, new %d", volumeID, volume.VolumeSize, volumeSizeNum)
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = adapter.logic.ResizeVolume(volumeID, volumeSizeNum)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, volume)
}

func (adapter *RESTAdapter) handleDeleteVolume(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleDeleteVolume",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	volumeID := c.Param("id")

	err := types.ValidateVolumeID(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	volume, err := adapter.logic.GetVolume(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) && volume.DeviceID != user {
		// requestiong other's volume info
		err := xerrors.Errorf("failed to get volume %s, you cannot access other devices' volume info", volumeID)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	if volume.Mounted {
		// deleting mounted volume
		err := xerrors.Errorf("failed to delete volume %s, you must unmount before you delete volume", volumeID)
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Deleting Volume ID: %s", volumeID)

	err = adapter.logic.DeleteVolume(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, volume)
}

func (adapter *RESTAdapter) handleMountVolume(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleMountVolume",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	volumeID := c.Param("id")

	err := types.ValidateVolumeID(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	volume, err := adapter.logic.GetVolume(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) && volume.DeviceID != user {
		// requestiong other's volume info
		err := xerrors.Errorf("failed to get volume %s, you cannot access other devices' volume info", volumeID)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Mounting Volume ID: %s", volumeID)

	err = adapter.logic.MountVolume(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, volume)
}

/*
type volumeMountRequest struct {
}

var input volumeMountRequest

err = c.BindJSON(&input)

	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
*/
func (adapter *RESTAdapter) handleUnmountVolume(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleUnmountVolume",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	volumeID := c.Param("id")

	err := types.ValidateVolumeID(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	/*
		type volumeUnmountRequest struct {
			// define input required
		}

		var input volumeUnmountRequest

		err = c.BindJSON(&input)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	*/
	volume, err := adapter.logic.GetVolume(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) && volume.DeviceID != user {
		// requestiong other's volume info
		err := xerrors.Errorf("failed to get volume %s, you cannot access other devices' volume info", volumeID)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	logger.Debugf("Unmounting Volume ID: %s", volumeID)

	err = adapter.logic.UnmountVolume(volumeID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, volume)
}
