# MTWebViz

A real-time macOS multitouch trackpad visualizer using WebSockets.

## Overview

This application captures multitouch events from macOS trackpads and streams them over WebSocket for visualization and analysis. It uses the private `MultitouchSupport.framework` to access raw touch data.

## Architecture

- **[touch/touch.go](touch/touch.go)** - Multitouch device integration using CGo
- **[server/websocket.go](server/websocket.go)** - WebSocket server and client management
- **[server/frontend.go](server/frontend.go)** - HTML frontend
- **[main.go](main.go)** - Application entry point

## MultitouchSupport Framework Access
(Thanks to the excellent write-up by @rmhsilva [here](https://gist.github.com/rmhsilva/61cc45587ed34707da34818a76476e11))

The touch module interfaces with macOS's private `MultitouchSupport.framework` through CGo:

```c
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
```

### Touch Data Structure

Each touch point (`Finger`) provides:
- **Position** (`normalized.pos.x`, `normalized.pos.y`) - Normalized coordinates (0.0 to 1.0)
- **Velocity** (`normalized.vel.x`, `normalized.vel.y`) - Touch movement velocity
- **Identifier** - Unique ID for tracking individual touches
- **State** - Touch state (began, moved, ended, etc.)
- **Size** - Contact area size
- **Angle** - Touch rotation angle
- **Major/Minor Axis** - Ellipse dimensions of the touch area
- **Timestamp** - High-precision timestamp
- **Frame** - Frame number for sequencing

## Requirements

- macOS (uses private MultitouchSupport framework)
- Go 1.21+
- `github.com/gorilla/websocket`

## Installation

```bash
go get github.com/gorilla/websocket
go build
```

## Usage

```bash
./mtwebviz
```

Then open http://localhost:8080 in your browser.

The WebSocket endpoint is available at `ws://localhost:8080/ws`

## JSON Event Format

```json
{
  "timestamp": 1234567890.123,
  "frame": 12345,
  "fingers": [
    {
      "id": 1,
      "state": 4,
      "x": 0.5,
      "y": 0.5,
      "vx": 0.01,
      "vy": -0.02,
      "size": 0.3,
      "angle": 1.57,
      "majorAxis": 0.4,
      "minorAxis": 0.3
    }
  ]
}
```

## License

MIT
