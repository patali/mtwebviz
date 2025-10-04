package touch

/*
#cgo LDFLAGS: -framework CoreFoundation -F/System/Library/PrivateFrameworks -framework MultitouchSupport
#include <CoreFoundation/CoreFoundation.h>

typedef struct { float x, y; } mtPoint;
typedef struct { mtPoint pos, vel; } mtReadout;

typedef struct {
    int frame;
    double timestamp;
    int identifier, state, foo3, foo4;
    mtReadout normalized;
    float size;
    int zero1;
    float angle, majorAxis, minorAxis;
    mtReadout mm;
    int zero2[2];
    float unk2;
} Finger;

typedef void* MTDeviceRef;
typedef int (*MTContactCallbackFunction)(int, Finger*, int, double, int);

MTDeviceRef MTDeviceCreateDefault();
void MTRegisterContactFrameCallback(MTDeviceRef, MTContactCallbackFunction);
void MTDeviceStart(MTDeviceRef, int);
void MTDeviceStop(MTDeviceRef);

extern int goTouchCallback(int device, Finger* data, int nFingers, double timestamp, int frame);

static inline void registerCallback(MTDeviceRef device) {
    MTRegisterContactFrameCallback(device, goTouchCallback);
}
*/
import "C"

import (
	"log"
	"unsafe"
)

// TouchEvent represents a single touch point
type TouchEvent struct {
	ID        int     `json:"id"`
	State     int     `json:"state"`
	X         float32 `json:"x"`
	Y         float32 `json:"y"`
	VX        float32 `json:"vx"`
	VY        float32 `json:"vy"`
	Size      float32 `json:"size"`
	Angle     float32 `json:"angle"`
	MajorAxis float32 `json:"majorAxis"`
	MinorAxis float32 `json:"minorAxis"`
}

// FrameEvent represents a complete frame of touch events
type FrameEvent struct {
	Timestamp float64      `json:"timestamp"`
	Frame     int          `json:"frame"`
	Fingers   []TouchEvent `json:"fingers"`
}

var eventCallback func(FrameEvent)

// SetEventCallback sets the callback function for touch events
func SetEventCallback(callback func(FrameEvent)) {
	eventCallback = callback
}

// This is called from C
//export goTouchCallback
func goTouchCallback(device C.int, data *C.Finger, nFingers C.int, timestamp C.double, frame C.int) C.int {
	event := FrameEvent{
		Timestamp: float64(timestamp),
		Frame:     int(frame),
		Fingers:   make([]TouchEvent, 0, nFingers),
	}

	// If there are fingers, add them to the event
	if nFingers > 0 && data != nil {
		// Convert C array to Go slice
		fingers := unsafe.Slice(data, nFingers)

		for i := 0; i < int(nFingers); i++ {
			finger := fingers[i]
			event.Fingers = append(event.Fingers, TouchEvent{
				ID:        int(finger.identifier),
				State:     int(finger.state),
				X:         float32(finger.normalized.pos.x),
				Y:         float32(finger.normalized.pos.y),
				VX:        float32(finger.normalized.vel.x),
				VY:        float32(finger.normalized.vel.y),
				Size:      float32(finger.size),
				Angle:     float32(finger.angle),
				MajorAxis: float32(finger.majorAxis),
				MinorAxis: float32(finger.minorAxis),
			})
		}
	}

	// Send event to callback if set
	if eventCallback != nil {
		eventCallback(event)
	}

	return 0
}

// Start begins multitouch tracking
func Start() {
	device := C.MTDeviceCreateDefault()
	if device == nil {
		log.Fatal("Failed to create multitouch device")
	}

	C.registerCallback(device)
	C.MTDeviceStart(device, 0)

	log.Println("Multitouch tracking started")
}
