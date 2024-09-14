package safe

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

type PanicCallback func(err interface{})

func Go(f func(), panicCallback ...PanicCallback) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				const size = 64 << 10
				buf := make([]byte, size)
				buf = buf[:runtime.Stack(buf, false)]
				logrus.WithFields(logrus.Fields{"panic": r, "stack": string(buf)}).Error("Go catch panic")
				for _, cb := range panicCallback {
					cb(r)
				}
			}
		}()
		f()
	}()
}
