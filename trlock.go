// Package trlock locks/unlocks an X11 screen while leaving windows visible.
//
// TODO:
//    * Hide/change cursor
//
// References:
//    http://www.x.org/archive/X11R7.7/doc/man/man3/xcb_grab_keyboard.3.xhtml
package trlock

import (
	"fmt"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

type grabStatus uint32

func (g grabStatus) String() string {
	switch g {
	case xproto.GrabStatusSuccess:
		return "GrabStatusSuccess"
	case xproto.GrabStatusAlreadyGrabbed:
		return "GrabStatusAlreadyGrabbed"
	case xproto.GrabStatusInvalidTime:
		return "GrabStatusInvalidTime"
	case xproto.GrabStatusNotViewable:
		return "GrabStatusNotViewable"
	case xproto.GrabStatusFrozen:
		return "GrabStatusFrozen"
	default:
		return fmt.Sprintf("Unknown grab status %d", g)
	}
}

// Lock grabs the keyboard and pointer locking the X11 display
func Lock(X *xgb.Conn) error {
	screen := xproto.Setup(X).DefaultScreen(X)

	passEventsToOwner := false
	kbCookie := xproto.GrabKeyboard(X, passEventsToOwner, screen.Root, xproto.TimeCurrentTime, xproto.GrabModeAsync, xproto.GrabModeAsync)
	kbReply, err := kbCookie.Reply()
	if err != nil {
		return err
	}
	if kbReply.Status != xproto.GrabStatusSuccess {
		return fmt.Errorf("GrabKeyboard status %v", grabStatus(kbReply.Status))
	}

	ptrCookie := xproto.GrabPointer(X, passEventsToOwner, screen.Root, 0, xproto.GrabModeAsync, xproto.GrabModeAsync, xproto.AtomNone, xproto.AtomNone, xproto.TimeCurrentTime)

	ptrReply, err := ptrCookie.Reply()
	if err != nil {
		xproto.UngrabKeyboard(X, xproto.TimeCurrentTime)
		return err
	}
	if ptrReply.Status != xproto.GrabStatusSuccess {
		return fmt.Errorf("GrabPointer status %v", grabStatus(kbReply.Status))
	}

	return nil
}

// Unlock restores user control of keyboard and pointer
func Unlock(X *xgb.Conn) {
	xproto.UngrabKeyboard(X, xproto.TimeCurrentTime)

	xproto.UngrabPointer(X, xproto.TimeCurrentTime)
}
