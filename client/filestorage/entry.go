package filestorage

type EntryProvider interface {
}

type Entry struct {
	Path       string   `json:"path"`
	IPFSAddr   string   `json:"ipfs_addr"`
	Owner      string   `json:"owner"`
	SharedWith []string `json:"shared_with"`
	WAccess    []string `json:"w_access"`
}
