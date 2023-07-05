package config

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	minSubscriptionBufferSize     = 100
	defaultSubscriptionBufferSize = 200
)

// ApiConfig defines the configuration options for the API server
type ApiConfig struct {
	RootDir string `mapstructure:"home"`

	// TCP or UNIX socket address for the RPC server to listen on
	ListenAddress string `mapstructure:"laddr"`

	// A list of origins a cross-domain request can be executed from.
	// If the special '*' value is present in the list, all origins will be allowed.
	// An origin may contain a wildcard (*) to replace 0 or more characters (i.e.: http://*.domain.com).
	// Only one wildcard can be used per origin.
	CORSAllowedOrigins []string `mapstructure:"cors_allowed_origins"`

	// A list of methods the client is allowed to use with cross-domain requests.
	CORSAllowedMethods []string `mapstructure:"cors_allowed_methods"`

	// A list of non simple headers the client is allowed to use with cross-domain requests.
	CORSAllowedHeaders []string `mapstructure:"cors_allowed_headers"`

	CORSAllowedCredentials bool `mapstructure:"cors_allowed_credentials"`
	CORSAMaxAge            int  `mapstructure:"cors_max_age"`

	//// TCP or UNIX socket address for the gRPC server to listen on
	//// NOTE: This server only supports /broadcast_tx_commit
	//GRPCListenAddress string `mapstructure:"grpc_laddr"`
	//
	//// Maximum number of simultaneous connections.
	//// Does not include RPC (HTTP&WebSocket) connections. See max_open_connections
	//// If you want to accept a larger number than the default, make sure
	//// you increase your OS limits.
	//// 0 - unlimited.
	//GRPCMaxOpenConnections int `mapstructure:"grpc_max_open_connections"`

	// Activate unsafe RPC commands like /dial_persistent_peers and /unsafe_flush_mempool
	Unsafe bool `mapstructure:"unsafe"`

	// Maximum number of simultaneous connections (including WebSocket).
	// Does not include gRPC connections. See grpc_max_open_connections
	// If you want to accept a larger number than the default, make sure
	// you increase your OS limits.
	// 0 - unlimited.
	// Should be < {ulimit -Sn} - {MaxNumInboundPeers} - {MaxNumOutboundPeers} - {N of wal, db and other open files}
	// 1024 - 40 - 10 - 50 = 924 = ~900
	MaxOpenConnections int `mapstructure:"max_open_connections"`

	// Maximum number of unique clientIDs that can /subscribe
	// If you're using /broadcast_tx_commit, set to the estimated maximum number
	// of broadcast_tx_commit calls per block.
	MaxSubscriptionClients int `mapstructure:"max_subscription_clients"`

	// Maximum number of unique queries a given client can /subscribe to
	// If you're using GRPC (or Local RPC client) and /broadcast_tx_commit, set
	// to the estimated maximum number of broadcast_tx_commit calls per block.
	MaxSubscriptionsPerClient int `mapstructure:"max_subscriptions_per_client"`

	// The number of events that can be buffered per subscription before
	// returning `ErrOutOfCapacity`.
	SubscriptionBufferSize int `mapstructure:"experimental_subscription_buffer_size"`

	// The maximum number of responses that can be buffered per WebSocket
	// client. If clients cannot read from the WebSocket endpoint fast enough,
	// they will be disconnected, so increasing this parameter may reduce the
	// chances of them being disconnected (but will cause the node to use more
	// memory).
	//
	// Must be at least the same as `SubscriptionBufferSize`, otherwise
	// connections may be dropped unnecessarily.
	WebSocketWriteBufferSize int `mapstructure:"experimental_websocket_write_buffer_size"`

	// If a WebSocket client cannot read fast enough, at present we may
	// silently drop events instead of generating an error or disconnecting the
	// client.
	//
	// Enabling this parameter will cause the WebSocket connection to be closed
	// instead if it cannot read fast enough, allowing for greater
	// predictability in subscription behavior.
	CloseOnSlowClient bool `mapstructure:"experimental_close_on_slow_client"`

	// Maximum size of request body, in bytes
	MaxBodyBytes int64 `mapstructure:"max_body_bytes"`

	// Maximum size of request header, in bytes
	MaxHeaderBytes int `mapstructure:"max_header_bytes"`

	// The path to a file containing certificate that is used to create the HTTPS server.
	// Might be either absolute path or path related to AA-Scan config directory.
	//
	// If the certificate is signed by a certificate authority,
	// the certFile should be the concatenation of the server's certificate, any intermediates,
	// and the CA's certificate.
	//
	// NOTE: both tls_cert_file and tls_key_file must be present for AA-Scan to create HTTPS server.
	// Otherwise, HTTP server is run.
	TLSCertFile string `mapstructure:"tls_cert_file"`

	// The path to a file containing matching private key that is used to create the HTTPS server.
	// Might be either absolute path or path related to AA-Scan config directory.
	//
	// NOTE: both tls_cert_file and tls_key_file must be present for AA-Scan to create HTTPS server.
	// Otherwise, HTTP server is run.
	TLSKeyFile string `mapstructure:"tls_key_file"`

	// no set no enable.
	// /debug/pprof/
	PprofPrefix string `mapstructure:"pprof_prefix"`

	Prefork bool `mapstructure:"prefork"`
}

// DefaultApiConfig returns a default configuration for the RPC server
func DefaultApiConfig() *ApiConfig {
	return &ApiConfig{
		RootDir:                "",
		ListenAddress:          "127.0.0.1:9080",
		CORSAllowedOrigins:     []string{},
		CORSAllowedMethods:     []string{http.MethodHead, http.MethodGet, http.MethodPost},
		CORSAllowedHeaders:     []string{"Origin", "Accept", "Content-Type", "X-Requested-With", "X-Server-Time"},
		CORSAllowedCredentials: false,
		CORSAMaxAge:            86400,
		//GRPCListenAddress:      "",
		//GRPCMaxOpenConnections: 900,

		Unsafe:                    false,
		MaxOpenConnections:        900,
		MaxSubscriptionClients:    100,
		MaxSubscriptionsPerClient: 5,
		SubscriptionBufferSize:    defaultSubscriptionBufferSize,
		WebSocketWriteBufferSize:  defaultSubscriptionBufferSize,
		CloseOnSlowClient:         false,
		MaxBodyBytes:              int64(1000000), // 1MB
		MaxHeaderBytes:            1 << 20,        // same as the net/http default
		TLSCertFile:               "",
		TLSKeyFile:                "",
		PprofPrefix:               "",
		Prefork:                   false,
	}
}

// ValidateBasic performs basic validation (checking param bounds, etc.) and
// returns an error if any check fails.
func (cfg *ApiConfig) ValidateBasic() error {
	if cfg.MaxOpenConnections < 0 {
		return errors.New("max_open_connections can't be negative")
	}
	if cfg.MaxSubscriptionClients < 0 {
		return errors.New("max_subscription_clients can't be negative")
	}
	if cfg.MaxSubscriptionsPerClient < 0 {
		return errors.New("max_subscriptions_per_client can't be negative")
	}
	if cfg.SubscriptionBufferSize < minSubscriptionBufferSize {
		return fmt.Errorf(
			"experimental_subscription_buffer_size must be >= %d",
			minSubscriptionBufferSize,
		)
	}
	if cfg.WebSocketWriteBufferSize < cfg.SubscriptionBufferSize {
		return fmt.Errorf(
			"experimental_websocket_write_buffer_size must be >= experimental_subscription_buffer_size (%d)",
			cfg.SubscriptionBufferSize,
		)
	}
	if cfg.MaxBodyBytes < 0 {
		return errors.New("max_body_bytes can't be negative")
	}
	if cfg.MaxHeaderBytes < 0 {
		return errors.New("max_header_bytes can't be negative")
	}
	return nil
}
