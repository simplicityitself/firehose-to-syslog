package events_test

import (
	"github.com/Sirupsen/logrus"
	"github.com/cloudfoundry-community/firehose-to-syslog/events"
	. "github.com/cloudfoundry/noaa/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestEvents(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Events Suite")
}

var _ = Describe("Events", func() {
	Describe("GetSelectedEvents", func() {
		Context("called with a empty list", func() {
			It("should return a hash of only the default event", func() {
				expected := map[string]bool{"LogMessage": true}
				Expect(events.GetSelectedEvents("")).To(Equal(expected))
			})
		})

		Context("called with a list of bogus event names", func() {
			It("should return a hash of only the default event", func() {
				expected := map[string]bool{"LogMessage": true}
				Expect(events.GetSelectedEvents("bogus,bogus1")).To(Equal(expected))
			})
		})

		Context("called with a list of both real and bogus event names", func() {
			It("should return a hash of only the real events", func() {
				expected := map[string]bool{
					"HttpStartStop": true,
					"CounterEvent":  true,
				}
				Expect(events.GetSelectedEvents("bogus,HttpStartStop,bogus1,CounterEvent")).To(Equal(expected))
			})
		})
	})

	Describe("Constructing a Event from a LogMessage", func() {
		var eventType Envelope_EventType = 5
		var messageType LogMessage_MessageType = 1
		var posixStart int64 = 1
		origin := "yomomma__0"
		sourceType := "Kehe"
		logMsg := "Help, I'm a rock! Help, I'm a rock! Help, I'm a cop! Help, I'm a cop!"
		sourceInstance := ">9000"
		appID := "eea38ba5-53a5-4173-9617-b442d35ec2fd"

		logMessage := LogMessage{
			Message:        []byte(logMsg),
			AppId:          &appID,
			Timestamp:      &posixStart,
			SourceType:     &sourceType,
			MessageType:    &messageType,
			SourceInstance: &sourceInstance,
		}

		envelope := &Envelope{
			EventType:  &eventType,
			Origin:     &origin,
			LogMessage: &logMessage,
		}

		Context("given a envelope", func() {
			It("should give us what we want", func() {
				event := events.LogMessage(envelope)
				Expect(event.Fields["event_type"]).To(Equal("LogMessage"))
				Expect(event.Fields["origin"]).To(Equal(origin))
				Expect(event.Fields["cf_app_id"]).To(Equal(appID))
				Expect(event.Fields["timestamp"]).To(Equal(posixStart))
				Expect(event.Fields["source_type"]).To(Equal(sourceType))
				Expect(event.Fields["message_type"]).To(Equal("OUT"))
				Expect(event.Fields["source_instance"]).To(Equal(sourceInstance))
				Expect(event.Msg).To(Equal(logMsg))
			})
		})
	})

	Describe("AnnotateWithAppData", func() {
		Context("called with Fields set to empty map", func() {
			It("should do nothing", func() {
				event := events.Event{}
				wanted := events.Event{}
				event.AnnotateWithAppData()
				Expect(event).To(Equal(wanted))
			})
		})

		Context("called with Fields set to logrus.Fields", func() {
			It("should do nothing", func() {
				event := events.Event{logrus.Fields{}, ""}
				wanted := events.Event{logrus.Fields{}, ""}
				event.AnnotateWithAppData()
				Expect(event).To(Equal(wanted))
			})
		})

		Context("called with empty cf_app_id", func() {
			It("should do nothing", func() {
				event := events.Event{logrus.Fields{"cf_app_id": ""}, ""}
				wanted := events.Event{logrus.Fields{"cf_app_id": ""}, ""}
				event.AnnotateWithAppData()
				Expect(event).To(Equal(wanted))
			})
		})
	})
})
