package dcpredis

import (
	"fmt"
	"sync"
	"testing"

	"github.com/Trendyol/go-dcp-redis/config"
	"github.com/Trendyol/go-dcp-redis/couchbase"
)

func TestFindCollectionKeyMapping_ReturnsCorrectMapping(t *testing.T) {
	mappings := []config.CollectionKeyMapping{
		{Collection: "col1", KeyPrefix: "prefix1:", StorageType: "string"},
		{Collection: "col2", KeyPrefix: "prefix2:", StorageType: "hash"},
	}
	SetCollectionKeyMappings(&mappings)

	m := findCollectionKeyMapping("col1")
	if m.Collection != "col1" || m.KeyPrefix != "prefix1:" {
		t.Errorf("unexpected mapping: %+v", m)
	}

	m2 := findCollectionKeyMapping("col2")
	if m2.Collection != "col2" || m2.StorageType != "hash" {
		t.Errorf("unexpected mapping: %+v", m2)
	}
}

func TestFindCollectionKeyMapping_CachesResult(t *testing.T) {
	mappings := []config.CollectionKeyMapping{
		{Collection: "cached", KeyPrefix: "c:", StorageType: "string"},
	}
	SetCollectionKeyMappings(&mappings)

	// Call twice — second call should hit cache.
	m1 := findCollectionKeyMapping("cached")
	m2 := findCollectionKeyMapping("cached")

	if m1 != m2 {
		t.Errorf("expected same value from cache, got %+v and %+v", m1, m2)
	}
}

func TestFindCollectionKeyMapping_PanicsOnUnknown(t *testing.T) {
	mappings := []config.CollectionKeyMapping{
		{Collection: "known", StorageType: "string"},
	}
	SetCollectionKeyMappings(&mappings)

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for unknown collection, got none")
		}
	}()

	findCollectionKeyMapping("unknown-collection")
}

// TestFindCollectionKeyMapping_ConcurrentAccess verifies that concurrent calls
// to findCollectionKeyMapping do not trigger a data race on mappingCache.
// Run with -race to catch any synchronization issues.
func TestFindCollectionKeyMapping_ConcurrentAccess(t *testing.T) {
	const numCollections = 5
	const goroutines = 50

	mappings := make([]config.CollectionKeyMapping, numCollections)
	for i := 0; i < numCollections; i++ {
		mappings[i] = config.CollectionKeyMapping{
			Collection:  fmt.Sprintf("col%d", i),
			KeyPrefix:   fmt.Sprintf("prefix%d:", i),
			StorageType: "string",
		}
	}
	SetCollectionKeyMappings(&mappings)

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(idx int) {
			defer wg.Done()
			colName := fmt.Sprintf("col%d", idx%numCollections)
			_ = findCollectionKeyMapping(colName)
		}(i)
	}
	wg.Wait()
}

func TestDefaultMapper_MutationEvent(t *testing.T) {
	mappings := []config.CollectionKeyMapping{
		{Collection: "test", KeyPrefix: "t:", StorageType: "string"},
	}
	SetCollectionKeyMappings(&mappings)

	ctx := couchbase.Context{
		Event: couchbase.Event{
			CollectionName: "test",
			Key:            []byte("mykey"),
			Value:          []byte(`{"foo":"bar"}`),
			IsMutated:      true,
		},
	}

	actions := DefaultMapper(ctx)
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
}

func TestDefaultMapper_DeletionEvent(t *testing.T) {
	mappings := []config.CollectionKeyMapping{
		{Collection: "test", KeyPrefix: "t:", StorageType: "string"},
	}
	SetCollectionKeyMappings(&mappings)

	ctx := couchbase.Context{
		Event: couchbase.Event{
			CollectionName: "test",
			Key:            []byte("mykey"),
			IsDeleted:      true,
		},
	}

	actions := DefaultMapper(ctx)
	if len(actions) != 1 {
		t.Fatalf("expected 1 action, got %d", len(actions))
	}
}

func TestDefaultMapper_UnknownEventReturnsNil(t *testing.T) {
	mappings := []config.CollectionKeyMapping{
		{Collection: "test", StorageType: "string"},
	}
	SetCollectionKeyMappings(&mappings)

	ctx := couchbase.Context{
		Event: couchbase.Event{
			CollectionName: "test",
			// None of IsMutated/IsDeleted/IsExpired is set.
		},
	}

	actions := DefaultMapper(ctx)
	if actions != nil {
		t.Errorf("expected nil actions for unknown event, got %v", actions)
	}
}
