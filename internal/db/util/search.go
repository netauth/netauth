package util

import (
	"github.com/blevesearch/bleve"

	"github.com/NetAuth/NetAuth/internal/db"

	pb "github.com/NetAuth/Protocol"
)

// SearchIndex holds the methods to search entities and groups with
// blevesearch.  This is meant to be embedded into a db implementation
// to transparently give it the search functions.
type SearchIndex struct {
	eIndex bleve.Index
	gIndex bleve.Index
}

// NewIndex returns a new SearchIndex with the mappings configured and
// ready to use.  Mappings are statically defined for simplicity, and
// in general new mappings shouldn't be added without a very good
// reason.
func NewIndex() *SearchIndex {
	// Setup the mapping for entities and turn off certain sub
	// keys that shouldn't be indexed.
	eMapping := bleve.NewIndexMapping()
	eDocMap := bleve.NewDocumentMapping()
	eDocMap.AddSubDocumentMapping("secret", bleve.NewDocumentDisabledMapping())
	eDocMap.AddSubDocumentMapping("meta.Keys", bleve.NewDocumentDisabledMapping())
	eDocMap.AddSubDocumentMapping("meta.UntypedMeta", bleve.NewDocumentDisabledMapping())
	eMapping.AddDocumentMapping("_default", eDocMap)

	// The only real way to throw an error in here is if a mapping
	// is invalid, or if this were on disk if the backing boltdb
	// couldn't be allocated.  Since this is fully in memory and
	// uses a hard-coded mapping, there is no concievable way for
	// an error to be returned here.  The same is true of the
	// group mapping below.
	eIndex, _ := bleve.NewMemOnly(eMapping)
	eIndex.SetName("EntityIndex")

	// Setup the mapping for groups and turn off certain sub keys
	// that shouldn't be indexed.
	gMapping := bleve.NewIndexMapping()
	gDocMap := bleve.NewDocumentMapping()
	gDocMap.AddSubDocumentMapping("untypedmeta", bleve.NewDocumentDisabledMapping())
	gMapping.AddDocumentMapping("_default", gDocMap)
	gIndex, _ := bleve.NewMemOnly(gMapping)
	gIndex.SetName("GroupIndex")

	// Return the prepared struct
	return &SearchIndex{
		eIndex: eIndex,
		gIndex: gIndex,
	}
}

// SearchEntities searches the index for entities matching the
// qualities specified in the request.
func (s *SearchIndex) SearchEntities(r db.SearchRequest) ([]string, error) {
	if r.Expression == "" {
		return nil, db.ErrBadSearch
	}

	req := createSearchRequest(r)

	// This can only fail if the query is malformed, since the
	// worst that can happen is the query is empty, this can't
	// return an error.
	result, _ := s.eIndex.Search(req)
	slice := extractDocIDs(result)
	return slice, nil
}

// SearchGroups searches the index for groups matching the qualities
// specified in the request.
func (s *SearchIndex) SearchGroups(r db.SearchRequest) ([]string, error) {
	if r.Expression == "" {
		return nil, db.ErrBadSearch
	}

	req := createSearchRequest(r)

	// This can only fail if the query is malformed, since the
	// worst that can happen is the query is empty, this can't
	// return an error.
	result, _ := s.gIndex.Search(req)
	slice := extractDocIDs(result)
	return slice, nil
}

// IndexEntity adds or updates an entity in the index.
func (s *SearchIndex) IndexEntity(e *pb.Entity) error {
	return s.eIndex.Index(e.GetID(), e)
}

// IndexGroup adds or updates a group in the index.
func (s *SearchIndex) IndexGroup(g *pb.Group) error {
	return s.gIndex.Index(g.GetName(), g)
}

// createSearchRequest is a helper function which converts between a
// db.SearchRequest and a bleve.SearchRequest.
func createSearchRequest(r db.SearchRequest) *bleve.SearchRequest {
	q := bleve.NewQueryStringQuery(r.Expression)
	sr := bleve.NewSearchRequest(q)
	return sr
}

// extractDocIDs converts between a bleve.SearchResult and a []string
// which can be subsequently fetched by the storage layer.
func extractDocIDs(r *bleve.SearchResult) []string {
	slice := []string{}
	hits := r.Hits
	for i := range hits {
		slice = append(slice, hits[i].ID)
	}
	return slice
}
