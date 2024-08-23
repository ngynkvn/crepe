package index

import "net/http"

func (ix Indexer) HandleGetRepositories(response http.ResponseWriter, req *http.Request)  {}
func (ix Indexer) HandlePostRepositories(response http.ResponseWriter, req *http.Request) {}

func (ix Indexer) HandleGetRepositoriesRepoId(response http.ResponseWriter, req *http.Request)    {}
func (ix Indexer) HandleDeleteRepositoriesRepoId(response http.ResponseWriter, req *http.Request) {}

func (ix Indexer) HandleSearch(response http.ResponseWriter, req *http.Request) {}
