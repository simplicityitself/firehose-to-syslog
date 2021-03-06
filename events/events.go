package events

import (
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/cloudfoundry-community/firehose-to-syslog/caching"
	"github.com/cloudfoundry/noaa/events"
	"os"
	"strings"
)

type Event struct {
	Fields logrus.Fields
	Msg    string
}

func RouteEvents(in chan *events.Envelope, selectedEvents map[string]bool) {
	for msg := range in {
		eventType := msg.GetEventType()
		if selectedEvents[eventType.String()] {
			var event Event
			switch eventType {
			case events.Envelope_Heartbeat:
				event = Heartbeat(msg)
			case events.Envelope_HttpStart:
				event = HttpStart(msg)
			case events.Envelope_HttpStop:
				event = HttpStop(msg)
			case events.Envelope_HttpStartStop:
				event = HttpStartStop(msg)
			case events.Envelope_LogMessage:
				event = LogMessage(msg)
			case events.Envelope_ValueMetric:
				event = ValueMetric(msg)
			case events.Envelope_CounterEvent:
				event = CounterEvent(msg)
			case events.Envelope_Error:
				event = ErrorEvent(msg)
			case events.Envelope_ContainerMetric:
				event = ContainerMetric(msg)
			}

			event.AnnotateWithAppData()
			event.Log()
		}
	}
}

func GetSelectedEvents(wantedEvents string) map[string]bool {
	selectedEvents := make(map[string]bool)
	for _, event := range strings.Split(wantedEvents, ",") {
		if isAuthorizedEvent(event) {
			selectedEvents[event] = true
		} else {
			fmt.Fprintf(os.Stderr, "Rejected Event Name %s\n", event)
		}
	}
	// If any event is not authorize we fallback to the default one
	if len(selectedEvents) < 1 {
		selectedEvents["LogMessage"] = true
	}
	return selectedEvents
}

func isAuthorizedEvent(wantedEvent string) bool {
	for _, authorizeEvent := range events.Envelope_EventType_name {
		if wantedEvent == authorizeEvent {
			return true
		}
	}
	return false
}

func GetListAuthorizedEventEvents() (authorizedEvents string) {
	arrEvents := []string{}
	for _, listEvent := range events.Envelope_EventType_name {
		arrEvents = append(arrEvents, listEvent)
	}
	return strings.Join(arrEvents, ", ")

}

func getAppInfo(appGuid string) caching.App {
	if app := caching.GetAppInfo(appGuid); app.Name != "" {
		return app
	} else {
		caching.GetAppByGuid(appGuid)
	}
	return caching.GetAppInfo(appGuid)
}

func Heartbeat(msg *events.Envelope) Event {
	heartbeat := msg.GetHeartbeat()

	fields := logrus.Fields{
		"ctl_msg_id":     heartbeat.GetControlMessageIdentifier(),
		"error_count":    heartbeat.GetErrorCount(),
		"event_type":     msg.GetEventType().String(),
		"origin":         msg.GetOrigin(),
		"received_count": heartbeat.GetReceivedCount(),
		"sent_count":     heartbeat.GetSentCount(),
	}

	return Event{
		Fields: fields,
		Msg:    "",
	}
}

func HttpStart(msg *events.Envelope) Event {
	httpStart := msg.GetHttpStart()

	fields := logrus.Fields{
		"event_type":        msg.GetEventType().String(),
		"origin":            msg.GetOrigin(),
		"cf_app_id":         httpStart.GetApplicationId(),
		"instance_id":       httpStart.GetInstanceId(),
		"instance_index":    httpStart.GetInstanceIndex(),
		"method":            httpStart.GetMethod(),
		"parent_request_id": httpStart.GetParentRequestId(),
		"peer_type":         httpStart.GetPeerType(),
		"request_id":        httpStart.GetRequestId(),
		"remote_addr":       httpStart.GetRemoteAddress(),
		"timestamp":         httpStart.GetTimestamp(),
		"uri":               httpStart.GetUri(),
		"user_agent":        httpStart.GetUserAgent(),
	}

	return Event{
		Fields: fields,
		Msg:    "",
	}
}

func HttpStop(msg *events.Envelope) Event {
	httpStop := msg.GetHttpStop()

	fields := logrus.Fields{
		"event_type":     msg.GetEventType().String(),
		"origin":         msg.GetOrigin(),
		"cf_app_id":      httpStop.GetApplicationId(),
		"content_length": httpStop.GetContentLength(),
		"peer_type":      httpStop.GetPeerType(),
		"request_id":     httpStop.GetRequestId(),
		"status_code":    httpStop.GetStatusCode(),
		"timestamp":      httpStop.GetTimestamp(),
		"uri":            httpStop.GetUri(),
	}

	return Event{
		Fields: fields,
		Msg:    "",
	}
}

func HttpStartStop(msg *events.Envelope) Event {
	httpStartStop := msg.GetHttpStartStop()

	fields := logrus.Fields{
		"event_type":        msg.GetEventType().String(),
		"origin":            msg.GetOrigin(),
		"cf_app_id":         httpStartStop.GetApplicationId(),
		"content_length":    httpStartStop.GetContentLength(),
		"instance_id":       httpStartStop.GetInstanceId(),
		"instance_index":    httpStartStop.GetInstanceIndex(),
		"method":            httpStartStop.GetMethod(),
		"parent_request_id": httpStartStop.GetParentRequestId(),
		"peer_type":         httpStartStop.GetPeerType(),
		"remote_addr":       httpStartStop.GetRemoteAddress(),
		"request_id":        httpStartStop.GetRequestId(),
		"start_timestamp":   httpStartStop.GetStartTimestamp(),
		"status_code":       httpStartStop.GetStatusCode(),
		"stop_timestamp":    httpStartStop.GetStopTimestamp(),
		"uri":               httpStartStop.GetUri(),
		"user_agent":        httpStartStop.GetUserAgent(),
	}

	return Event{
		Fields: fields,
		Msg:    "",
	}
}

func LogMessage(msg *events.Envelope) Event {
	logMessage := msg.GetLogMessage()

	fields := logrus.Fields{
		"event_type":      msg.GetEventType().String(),
		"origin":          msg.GetOrigin(),
		"cf_app_id":       logMessage.GetAppId(),
		"timestamp":       logMessage.GetTimestamp(),
		"source_type":     logMessage.GetSourceType(),
		"message_type":    logMessage.GetMessageType().String(),
		"source_instance": logMessage.GetSourceInstance(),
	}

	return Event{
		Fields: fields,
		Msg:    string(logMessage.GetMessage()),
	}
}

func ValueMetric(msg *events.Envelope) Event {
	valMetric := msg.GetValueMetric()

	fields := logrus.Fields{
		"event_type": msg.GetEventType().String(),
		"origin":     msg.GetOrigin(),
		"name":       valMetric.GetName(),
		"unit":       valMetric.GetUnit(),
		"value":      valMetric.GetValue(),
	}

	return Event{
		Fields: fields,
		Msg:    "",
	}
}

func CounterEvent(msg *events.Envelope) Event {
	counterEvent := msg.GetCounterEvent()

	fields := logrus.Fields{
		"event_type": msg.GetEventType().String(),
		"origin":     msg.GetOrigin(),
		"name":       counterEvent.GetName(),
		"delta":      counterEvent.GetDelta(),
		"total":      counterEvent.GetTotal(),
	}

	return Event{
		Fields: fields,
		Msg:    "",
	}
}

func ErrorEvent(msg *events.Envelope) Event {
	errorEvent := msg.GetError()

	fields := logrus.Fields{
		"event_type": msg.GetEventType().String(),
		"origin":     msg.GetOrigin(),
		"code":       errorEvent.GetCode(),
		"delta":      errorEvent.GetSource(),
	}

	return Event{
		Fields: fields,
		Msg:    errorEvent.GetMessage(),
	}
}

func ContainerMetric(msg *events.Envelope) Event {
	containerMetric := msg.GetContainerMetric()

	fields := logrus.Fields{
		"event_type":     msg.GetEventType().String(),
		"origin":         msg.GetOrigin(),
		"cf_app_id":      containerMetric.GetApplicationId(),
		"cpu_percentage": containerMetric.GetCpuPercentage(),
		"disk_bytes":     containerMetric.GetDiskBytes(),
		"instance_index": containerMetric.GetInstanceIndex(),
		"memory_bytes":   containerMetric.GetMemoryBytes(),
	}

	return Event{
		Fields: fields,
		Msg:    "",
	}
}

func (e *Event) AnnotateWithAppData() {
	cf_app_id := e.Fields["cf_app_id"]
	if cf_app_id != nil && cf_app_id != "" {
		appInfo := getAppInfo(fmt.Sprintf("%s", cf_app_id))
		cf_app_name := appInfo.Name
		cf_space_id := appInfo.SpaceGuid
		cf_space_name := appInfo.SpaceName
		cf_org_id := appInfo.OrgGuid
		cf_org_name := appInfo.OrgName

		if cf_app_name != "" {
			e.Fields["cf_app_name"] = cf_app_name
		}

		if cf_space_id != "" {
			e.Fields["cf_space_id"] = cf_space_id
		}

		if cf_space_name != "" {
			e.Fields["cf_space_name"] = cf_space_name
		}

		if cf_org_id != "" {
			e.Fields["cf_org_id"] = cf_org_id
		}

		if cf_org_name != "" {
			e.Fields["cf_org_name"] = cf_org_name
		}
	}
}

func (e Event) Log() {
	logrus.WithFields(e.Fields).Info(e.Msg)
}
