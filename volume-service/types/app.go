package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rs/xid"
	"golang.org/x/xerrors"
)

const (
	appIDPrefix string = "app"
	runIDPrefix string = "run"
)

// App represents an app, holding all necessary info. about app
type App struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	RequireGPU  bool      `json:"require_gpu,omitempty"`
	Description string    `json:"description,omitempty"`
	DockerImage string    `json:"docker_image"`
	Commands    string    `json:"commands,omitempty"`   // a space-separated commands to run app, array/map not supported
	Arguments   string    `json:"arguments,omitempty"`  // a space-separated command-line arguments to run app, array/map not supported
	OpenPorts   []int     `json:"open_ports,omitempty"` // first element is the main service port to open with ingress setup
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

// App represents an app, holding all necessary info. about app
type AppSQLiteObj struct {
	ID          string    `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	RequireGPU  bool      `json:"require_gpu,omitempty"`
	Description string    `json:"description,omitempty"`
	DockerImage string    `json:"docker_image"`
	Commands    string    `json:"commands,omitempty"`   // command to set in container
	Arguments   string    `json:"arguments,omitempty"`  // arguments to use in container
	OpenPorts   string    `json:"open_ports,omitempty"` // store it as a comma-separated string
	CreatedAt   time.Time `json:"created_at,omitempty"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

func (app *App) ToAppSQLiteObj() AppSQLiteObj {
	openports := []string{}
	for _, port := range app.OpenPorts {
		portString := fmt.Sprintf("%d", port)
		openports = append(openports, portString)
	}

	openPortsCSV := strings.Join(openports, ",")

	return AppSQLiteObj{
		ID:          app.ID,
		Name:        app.Name,
		RequireGPU:  app.RequireGPU,
		Description: app.Description,
		DockerImage: app.DockerImage,
		Arguments:   app.Arguments,
		Commands:    app.Commands,
		OpenPorts:   openPortsCSV,
		CreatedAt:   app.CreatedAt,
		UpdatedAt:   app.UpdatedAt,
	}
}

func (app *AppSQLiteObj) ToAppObj() App {
	openports := []int{}
	openportsStringArr := strings.Split(app.OpenPorts, ",")

	for _, port := range openportsStringArr {
		portInt, _ := strconv.Atoi(port)
		openports = append(openports, portInt)
	}

	return App{
		ID:          app.ID,
		Name:        app.Name,
		RequireGPU:  app.RequireGPU,
		Description: app.Description,
		DockerImage: app.DockerImage,
		Commands:    app.Commands,
		Arguments:   app.Arguments,
		OpenPorts:   openports,
		CreatedAt:   app.CreatedAt,
		UpdatedAt:   app.UpdatedAt,
	}
}

func ValidateAppID(id string) error {
	if len(id) == 0 {
		return xerrors.Errorf("empty app id")
	}

	prefix := fmt.Sprintf("%s_", appIDPrefix)

	if !strings.HasPrefix(id, prefix) {
		return xerrors.Errorf("invalid app id - %s", id)
	}
	return nil
}

// NewAppID creates a new App ID
func NewAppID() string {
	return fmt.Sprintf("%s_%s", appIDPrefix, xid.New().String())
}

// AppRun represents an app execution (run), holding all necessary info. about app execution
type AppRun struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	DeviceID   string    `json:"device_id"`
	VolumeID   string    `json:"volume_id"`
	AppID      string    `json:"app_id"`
	Terminated bool      `json:"terminated"`
	CreatedAt  time.Time `json:"created_at,omitempty"`
	UpdatedAt  time.Time `json:"updated_at,omitempty"`
}

func ValidateAppRunID(id string) error {
	if len(id) == 0 {
		return xerrors.Errorf("empty app run id")
	}

	prefix := fmt.Sprintf("%s_", runIDPrefix)

	if !strings.HasPrefix(id, prefix) {
		return xerrors.Errorf("invalid app run id - %s", id)
	}
	return nil
}

// NewAppRunID creates a new AppRun ID
func NewAppRunID() string {
	return fmt.Sprintf("%s_%s", runIDPrefix, xid.New().String())
}
