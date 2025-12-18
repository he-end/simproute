package goruntime

import (
	"bytes"
	"runtime"
	"strconv"
	"sync"

	"github.com/google/uuid"
)

var reqID sync.Map

func Goid() int64 {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)

	b := buf[:n]

	b = bytes.TrimPrefix(b, []byte(`goroutine `))
	idField := bytes.Fields(b)[0]

	id, _ := strconv.ParseInt(string(idField), 10, 64)
	return id
}

func GetCorelationID() uuid.UUID {
	gid := Goid()
	stockReqID, ok := reqID.Load(gid)
	if ok {
		return stockReqID.(uuid.UUID)
	}
	newID := uuid.New()
	reqID.Store(gid, newID)
	return newID
}

// Set Request ID to aware for global logger
func SetCorelationID(id uuid.UUID) {
	gid := Goid()
	reqID.Store(gid, id)
}

func ClearCorelationID() {
	gid := Goid()
	reqID.Delete(gid)
}
