package collectiontypes

import (
	"log"
	"encoding/json"
	"bin-collections-api/internal/pkg/get-config-value"
)

// CollectionColourRegistryEntry f
type CollectionColourRegistryEntry struct {
	Colour string
	TypeName string
}

type indexedCollectionColourRegistryEntry struct {
	CollectionColourRegistryEntry
	Index int
}

// CollectionColourRegistry f
type CollectionColourRegistry map[string]indexedCollectionColourRegistryEntry

// NewCollectionColourRegistry f
func NewCollectionColourRegistry(colourTypesEntries []CollectionColourRegistryEntry) CollectionColourRegistry {
	var registry CollectionColourRegistry
	err := json.Unmarshal([]byte(getconfigvalue.ByKey("DEFAULT_COLLECTION_TYPES_BY_COLOUR")), &registry)

	if err != nil {
		log.Fatal(err)
	}

	for index, entry := range colourTypesEntries {
		registry[entry.Colour] = indexedCollectionColourRegistryEntry{
			entry,
			index,
		}
	}

	return registry
}
