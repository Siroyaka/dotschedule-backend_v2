package abstruct

import (
	"context"

	"github.com/Siroyaka/dotschedule-backend_v2/utility"
)

type Firestore interface {
	Collection(string) FirestoreQuery
	Close() utility.IError
	GetContext() GCPContext
}

type GCPContext context.Context

type FirestoreQuery interface {
	Where(path, op string, value interface{}) FirestoreQuery
	OrderBy(path string, dir int32) FirestoreQuery
	Offset(n int) FirestoreQuery
	Limit(n int) FirestoreQuery
	LimitToLast(n int) FirestoreQuery
	StartAt(docSnapshotOrFieldValues ...interface{}) FirestoreQuery
	StartAfter(docSnapshotOrFieldValues ...interface{}) FirestoreQuery
	EndAt(docSnapshotOrFieldValues ...interface{}) FirestoreQuery
	EndBefore(docSnapshotOrFieldValues ...interface{}) FirestoreQuery

	Documents(ctx context.Context) DocumentIterator
}

type DocumentIterator interface {
	Next() (bool, FirestoreDocumentSnapshop, utility.IError)
	Stop()
}

type FirestoreDocumentSnapshop interface {
	Exists() bool
	Data() map[string]interface{}
}
