package rest

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lab-paper-code/ksv/volume-service/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/xerrors"
)

// setupDeviceRouter setup http request router for device
func (adapter *RESTAdapter) setupDeviceRouter() {
	// any devices can call these APIs
	adapter.router.GET("/devices", adapter.basicAuthDeviceOrAdmin, adapter.handleListDevices)
	adapter.router.GET("/devices/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleGetDevice)
	adapter.router.PATCH("/devices/:id", adapter.basicAuthDeviceOrAdmin, adapter.handleUpdateDevice)

	// any devices can call these APIs
	adapter.router.POST("/devices", adapter.basicAuthAdmin, adapter.handleRegisterDevice)
}

func (adapter *RESTAdapter) handleListDevices(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleListDevices",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)

	type listOutput struct {
		Devices []types.Device `json:"devices"`
	}

	output := listOutput{}

	if adapter.isAdminUser(user) {
		// admin - returns all devices
		devices, err := adapter.logic.ListDevices()
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		redactedDevices := make([]types.Device, len(devices))
		for deviceIdx, device := range devices {
			redactedDevice := device.GetRedacted()
			redactedDevices[deviceIdx] = redactedDevice
		}

		output.Devices = redactedDevices
	} else {
		err := types.ValidateDeviceID(user)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// device - returns mine
		device, err := adapter.logic.GetDevice(user)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		output.Devices = []types.Device{device.GetRedacted()}
	}

	// success
	c.JSON(http.StatusOK, output)
}

func (adapter *RESTAdapter) handleGetDevice(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleGetDevice",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	deviceID := c.Param("id")

	err := types.ValidateDeviceID(deviceID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) && deviceID != user {
		// requesting other's device info
		err := xerrors.Errorf("failed to get device %s, you cannot access other device info", deviceID)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	device, err := adapter.logic.GetDevice(deviceID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// success
	c.JSON(http.StatusOK, device.GetRedacted())
}

func (adapter *RESTAdapter) handleRegisterDevice(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleRegisterDevice",
	})

	logger.Infof("access request to %s", c.Request.URL)

	type deviceRegistrationRequest struct {
		IP          string `json:"ip,omitempty"`
		Password    string `json:"password"`
		Description string `json:"description,omitempty"`
	}

	var input deviceRegistrationRequest

	err := c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.IP) == 0 {
		remoteAddrFields := strings.Split(c.Request.RemoteAddr, ":")
		if len(remoteAddrFields) > 0 {
			input.IP = remoteAddrFields[0]
		}
	}

	if len(input.Password) == 0 {
		// fail
		err = xerrors.Errorf("password is not given")
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device := types.Device{
		ID:          types.NewDeviceID(),
		IP:          input.IP,
		Password:    input.Password,
		Description: input.Description, // optional
	}

	err = adapter.logic.CreateDevice(&device)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, device.GetRedacted())
}

func (adapter *RESTAdapter) handleUpdateDevice(c *gin.Context) {
	logger := log.WithFields(log.Fields{
		"package":  "rest",
		"struct":   "RESTAdapter",
		"function": "handleUpdateDevice",
	})

	logger.Infof("access request to %s", c.Request.URL)

	user := c.GetString(gin.AuthUserKey)
	deviceID := c.Param("id")

	err := types.ValidateDeviceID(deviceID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !adapter.isAdminUser(user) && deviceID != user {
		// requesting other's device info
		err := xerrors.Errorf("failed to update device %s, you cannot access other device info", deviceID)
		logger.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	type deviceUpdateRequest struct {
		IP       string `json:"ip,omitempty"`
		Password string `json:"password,omitempty"`
	}

	var input deviceUpdateRequest

	err = c.BindJSON(&input)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(input.IP) > 0 {
		// update IP
		err = adapter.logic.UpdateDeviceIP(deviceID, input.IP)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if len(input.Password) > 0 {
		// update password
		err = adapter.logic.UpdateDevicePassword(deviceID, input.Password)
		if err != nil {
			// fail
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	device, err := adapter.logic.GetDevice(deviceID)
	if err != nil {
		// fail
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, device.GetRedacted())
}
