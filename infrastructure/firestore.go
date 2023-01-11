package infrastructure

import (
	"cloud.google.com/go/firestore"
	"github.com/Siroyaka/dotschedule-backend_v2/adapter/abstruct"
	"github.com/Siroyaka/dotschedule-backend_v2/utility"
	gcp_context "golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type Firestore struct {
	ctx    abstruct.GCPContext
	client *firestore.Client
}

func NewFirestore(credentialsFilePath string, projectID string) abstruct.Firestore {
	ctx := gcp_context.Background()
	auth := option.WithCredentialsFile(credentialsFilePath)
	client, err := firestore.NewClient(ctx, projectID, auth)
	if err != nil {
		panic(err)
	}

	return Firestore{
		ctx:    ctx,
		client: client,
	}
}

func (fs Firestore) Close() utility.IError {
	err := fs.client.Close()
	if err != nil {
		return utility.NewError(err.Error(), "")
	}
	return nil
}

func (fs Firestore) GetContext() abstruct.GCPContext {
	return fs.ctx
}

func (fs Firestore) Collection(collenctionName string) abstruct.FirestoreQuery {
	return FirestoreQuery{
		query: fs.client.Collection(collenctionName).Query,
	}
}

type FirestoreQuery struct {
	query firestore.Query
}

func (fq FirestoreQuery) Where(path, op string, value interface{}) abstruct.FirestoreQuery {
	return FirestoreQuery{
		query: fq.query.Where(path, op, value),
	}
}

func (fq FirestoreQuery) OrderBy(path string, dir int32) abstruct.FirestoreQuery {
	direction := firestore.Direction(dir)
	return FirestoreQuery{
		query: fq.query.OrderBy(path, direction),
	}
}

func (fq FirestoreQuery) Offset(n int) abstruct.FirestoreQuery {
	return FirestoreQuery{
		query: fq.query.Offset(n),
	}
}

func (fq FirestoreQuery) Limit(n int) abstruct.FirestoreQuery {
	return FirestoreQuery{
		query: fq.query.Limit(n),
	}
}
func (fq FirestoreQuery) LimitToLast(n int) abstruct.FirestoreQuery {
	return FirestoreQuery{
		query: fq.query.LimitToLast(n),
	}
}
func (fq FirestoreQuery) StartAt(docSnapshotOrFieldValues ...interface{}) abstruct.FirestoreQuery {
	return FirestoreQuery{
		query: fq.query.StartAt(docSnapshotOrFieldValues...),
	}
}
func (fq FirestoreQuery) StartAfter(docSnapshotOrFieldValues ...interface{}) abstruct.FirestoreQuery {
	return FirestoreQuery{
		query: fq.query.StartAfter(docSnapshotOrFieldValues...),
	}
}
func (fq FirestoreQuery) EndAt(docSnapshotOrFieldValues ...interface{}) abstruct.FirestoreQuery {
	return FirestoreQuery{
		query: fq.query.EndAt(docSnapshotOrFieldValues...),
	}
}
func (fq FirestoreQuery) EndBefore(docSnapshotOrFieldValues ...interface{}) abstruct.FirestoreQuery {
	return FirestoreQuery{
		query: fq.query.EndBefore(docSnapshotOrFieldValues...),
	}
}

func (fq FirestoreQuery) Documents(ctx gcp_context.Context) abstruct.DocumentIterator {
	return &FirestoreDocumentIterator{
		iter: fq.query.Documents(ctx),
	}
}

type FirestoreDocumentIterator struct {
	iter *firestore.DocumentIterator
}

func (fdi *FirestoreDocumentIterator) Next() (bool, abstruct.FirestoreDocumentSnapshop, utility.IError) {
	snapshot, err := fdi.iter.Next()
	if err == iterator.Done {
		return false, snapshot, utility.NewError(err.Error(), "")
	}
	if err != nil {
		return false, snapshot, utility.NewError(err.Error(), "")
	}
	return true, snapshot, nil
}

func (fdi *FirestoreDocumentIterator) Stop() {
	fdi.Stop()
}

type FirestoreQueryValue struct {
	path  string
	op    string
	value interface{}
}

func NewFirestoreQueryValue(path string, op string, value interface{}) FirestoreQueryValue {
	return FirestoreQueryValue{
		path:  path,
		op:    op,
		value: value,
	}
}

func (fqv FirestoreQueryValue) getQueryValue() (string, string, interface{}) {
	return fqv.op, fqv.path, fqv.value
}
