package main

const (
	APP_CONFIG        = "/run/secrets/config.json"
	APP_NAME 					= "mboard-go"
	APP_VERSION 			= "1.0.0"
	TEST_ADDRESS 	    = "127.0.0.1:8000"
	CLOUD_ADDRESS     = "madsportslab.com"
	MBOARD            = "mboard"
	MBOARD_WWW        = "www"
	QR_FILE						= "qr.png"
)

const (
	MODE_WIFI				= 0
	MODE_HOTSPOT   	= 1
	MODE_WIRED      = 2
	MODE_CLOUD      = 3
	MODE_TEST   		= 4
)

const (
	INTERFACE_WIFI 		= "en"
	INTERFACE_HOTSPOT	= "wlan"
	INTERFACE_WIRED   = "eth"
	INTERFACE_CLOUD   = "cloud"
	INTERFACE_TEST		= "lo"
	INTERFACE_ERROR   = ""
)

const (
	API_PARAM_GAME_CONFIG			= "gameConfig"
)

const (
	KEY_FIELD_CONFIG				= "config"
	KEY_FIELD_DATA         	= "data"
	KEY_FIELD_PLAYS       	= "plays"
)

const (
	HSET										= "HSET"
	HGET										= "HGET"
	HGETALL									= "HGETALL"
	LPUSH                   = "LPUSH"
	LRANGE                  = "LRANGE"
)
