package main

type UserEntry struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Configuration struct {
	Datadir  string      `json:"datadir"`
	Port     int         `json:"port"`
	Hostname string      `json:"host"`
	Userdb   []UserEntry `json:"userdb"`
}
