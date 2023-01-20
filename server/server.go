package server

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"net/http"

	"github.com/notEpsilon/ss-cache/config"
	"github.com/notEpsilon/ss-cache/lru"
)

type Server struct {
	config *config.Config
	shard  *config.Shard
	cache  *lru.LRUCache[string, []byte]
}

type res map[string]any

func New(cacheCapacity int) (*Server, error) {
	cache, err := lru.New[string, []byte](cacheCapacity)
	if err != nil {
		return nil, err
	}

	return &Server{
		cache: cache,
	}, nil
}

func (s *Server) Start(listenAddr string) error {
	if listenAddr == "" {
		listenAddr = fmt.Sprintf("%s:%d", s.shard.ExposedAddress.Host, s.shard.ExposedAddress.Port)
	}

	http.HandleFunc("/get", s.getHandler)
	http.HandleFunc("/set", s.setHandler)

	return http.ListenAndServe(listenAddr, nil)
}

func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(res{
			"err": "Unable to parse query parameters",
		})
		return
	}

	key := r.Form.Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res{
			"err": "Please provide the `key` query parameter",
		})
		return
	}

	hash := fnv.New64()
	_, err := hash.Write([]byte(key))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res{
			"err": "Unable to hash the provided key",
		})
	}

	targetShardIndex := int(hash.Sum64() % uint64(len(s.config.Shards)))
	log.Printf("curr_shard=%d, target_shard=%d\n", s.shard.Idx, targetShardIndex)

	if targetShardIndex != s.shard.Idx {
		// not responsible for this shard, redirect to the appropriate server.
		url := fmt.Sprintf("http://%s:%d%s", s.config.Shards[targetShardIndex].ExposedAddress.Host, s.config.Shards[targetShardIndex].ExposedAddress.Port, r.RequestURI)
		log.Printf("redirected_request: %s\n", url)

		resp, err := http.Get(url)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(res{
				"err": fmt.Sprintf("%s, error: %s", "Unable to redirect request to the appropriate server", err.Error()),
			})
			resp.Body.Close()
			return
		}

		if resp.StatusCode != http.StatusOK {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(res{
				"err": "redirected request wasn't handled correctly by it's server",
			})
			resp.Body.Close()
			return
		}

		var jsonBody res
		json.NewDecoder(resp.Body).Decode(&jsonBody)
		json.NewEncoder(w).Encode(jsonBody)

		resp.Body.Close()
		return
	}

	rawData, _ := s.cache.Get(key) // no error check, if error occurred the data in the response is just null.
	json.NewEncoder(w).Encode(res{
		"data": rawData,
	})
}

func (s *Server) setHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(res{
			"err": "Unable to parse query parameters",
		})
		return
	}

	key := r.Form.Get("key")
	if key == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res{
			"err": "Please provide the `key` query parameter",
		})
		return
	}

	value := r.Form.Get("value")
	if value == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res{
			"err": "Please provide the `value` query parameter",
		})
		return
	}

	hash := fnv.New64()
	_, err := hash.Write([]byte(key))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(res{
			"err": "Unable to hash the provided key",
		})
	}

	targetShardIndex := int(hash.Sum64() % uint64(len(s.config.Shards)))
	log.Printf("curr_shard=%d, target_shard=%d\n", s.shard.Idx, targetShardIndex)

	if targetShardIndex != s.shard.Idx {
		// not responsible for this shard, redirect to the appropriate server.
		url := fmt.Sprintf("http://%s:%d%s", s.config.Shards[targetShardIndex].ExposedAddress.Host, s.config.Shards[targetShardIndex].ExposedAddress.Port, r.RequestURI)
		log.Printf("redirected_request: %s\n", url)

		resp, err := http.Get(url)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(res{
				"err": fmt.Sprintf("%s, error: %s", "Unable to redirect request to the appropriate server", err.Error()),
			})
			resp.Body.Close()
			return
		}

		if resp.StatusCode != http.StatusOK {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(res{
				"err": "redirected request wasn't handled correctly by it's server",
			})
			resp.Body.Close()
			return
		}

		var jsonBody res
		json.NewDecoder(resp.Body).Decode(&jsonBody)
		json.NewEncoder(w).Encode(jsonBody)

		resp.Body.Close()
		return
	}

	s.cache.Set(key, []byte(value))
}

func (s *Server) GetShard() *config.Shard {
	return s.shard
}

func (s *Server) SetShard(shard *config.Shard) {
	s.shard = shard
}

func (s *Server) GetConfig() *config.Config {
	return s.config
}

func (s *Server) SetConfig(config *config.Config) {
	s.config = config
}
