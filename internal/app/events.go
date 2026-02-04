package app

import (
	"github.com/RayUI/RayUI/internal/model"
	"github.com/RayUI/RayUI/internal/service"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// Event name constants.
const (
	EventTraffic        = "stats:traffic"
	EventCoreStatus     = "core:status"
	EventCoreLog        = "core:log"
	EventNotification   = "notification"
	EventSpeedTest      = "speedtest:result"
	EventUpdateProgress = "update:progress"
)

func (a *App) emitTraffic(stats model.TrafficStats) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, EventTraffic, stats)
	}
}

func (a *App) emitCoreStatus(status model.CoreStatus) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, EventCoreStatus, status)
	}
}

func (a *App) emitCoreLog(line string) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, EventCoreLog, line)
	}
}

func (a *App) emitNotification(ntype, title, message string) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, EventNotification, map[string]string{
			"type":    ntype,
			"title":   title,
			"message": message,
		})
	}
}

func (a *App) emitUpdateProgress(p service.UpdateProgress) {
	if a.ctx != nil {
		runtime.EventsEmit(a.ctx, EventUpdateProgress, p)
	}
}
