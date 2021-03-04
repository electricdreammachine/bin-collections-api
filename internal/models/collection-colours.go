package models

import (
	config "bin-collections-api/internal/services/config"
	"encoding/json"
	"log"
)

// CollectionColourRegistryEntry f
type CollectionColourRegistryEntry struct {
	Colour   string
	TypeName string
}

// IndexedCollectionColourRegistryEntry f
type IndexedCollectionColourRegistryEntry struct {
	CollectionColourRegistryEntry
	Index int
}

// CollectionColourRegistry f
type CollectionColourRegistry map[string]IndexedCollectionColourRegistryEntry

// NewCollectionColourRegistry f
func NewCollectionColourRegistry(colourTypesEntries []CollectionColourRegistryEntry) CollectionColourRegistry {
	var registry CollectionColourRegistry
	err := json.Unmarshal([]byte(config.ByKey("DEFAULT_COLLECTION_TYPES_BY_COLOUR")), &registry)

	if err != nil {
		log.Fatal(err)
	}

	for index, entry := range colourTypesEntries {
		registry[entry.Colour] = IndexedCollectionColourRegistryEntry{
			entry,
			index,
		}
	}

	return registry
}
