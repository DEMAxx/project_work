package app_errors

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/message"
)

type AlertType string

const (
	AlertCriticalType AlertType = "critical"
	AlertDangerType   AlertType = "danger"
	AlertWarningType  AlertType = "warning"
	AlertInfoType     AlertType = "info"
	AlertSuccessType  AlertType = "success"
)

type Alert struct {
	Type        AlertType
	Code        string
	Message     string
	Source      string
	Params      map[string]any
	Meta        map[string]any
	Recoverable bool
	Status      int
}

func (a *Alert) Error() string {
	return fmt.Sprintf("%s : %s", a.Code, a.Message)
}
func (a *Alert) IsRecoverable() bool {
	return a.Recoverable
}

func (a *Alert) FillStatus(status int) {
	if a.Status == 0 {
		a.Status = status
	}
}

func (a *Alert) FillMessage(p *message.Printer) {
	if a.Message == "" {
		a.Message = p.Sprintf(fmt.Sprintf("alerts.%s", a.Code))
	}
}

func (a *Alert) Alert(ctx context.Context) *Alert {
	return a
}

func (a *Alert) HasErrors() bool {
	return a.Type != AlertSuccessType
}

type Alerts []*Alert

func (a Alerts) Error() (s string) {
	for _, alert := range a {
		s += alert.Error() + "\n"
	}
	return
}

func (a Alerts) IsRecoverable() bool {
	for _, alert := range a {
		if !alert.IsRecoverable() {
			return false
		}
	}
	return true
}

func (a Alerts) FiberAnswer(p *message.Printer) (int, fiber.Map) {
	status := fiber.StatusInternalServerError

	for i := range a {
		a[i].FillMessage(p)
		if a[i].Status != 0 {
			status = a[i].Status
		}
	}

	return status, fiber.Map{"alerts": a}
}

func (a Alerts) FillStatus(status int) {
	for i := range a {
		a[i].FillStatus(status)
	}
}

func (a Alerts) FiberMap() fiber.Map {
	return fiber.Map{"alerts": a}
}

func (a Alerts) HasErrors() bool {
	for _, alert := range a {
		if alert.HasErrors() {
			return true
		}
	}

	return false
}

func (a Alerts) FillSource(val string) {
	for i := range a {
		if a[i].Source == "" {
			a[i].Source = val
		}
	}
}

type AlertConvertable interface {
	Alert(ctx context.Context) *Alert
	Error() string
}

type AlertRecoverable interface {
	IsRecoverable() bool
	Error() string
}

type AlertSuccess interface {
	HasErrors() bool
	Error() string
}
