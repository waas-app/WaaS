package config

const (
	Development string = "development"
	ServiceName string = "WaaS"
	CurrentUser string = "currentUser"
)

type Config struct {
	OTLPEndpoint  string `mapstructure:"otlp_endpoint"`
	RootURL       string `mapstructure:"root_url"`
	Environment   string `mapstructure:"environment"`
	AdminUserName string `mapstructure:"admin_username"`
	AdminPassword string `mapstructure:"admin_password"`
	ExternalHost  string `mapstructure:"externalHost"`
	Storage       string `mapstructure:"storage"`
	Port          int    `mapstructure:"port"`
	SessionSecret string `mapstructure:"session_secret"`
	CookieDomain  string `mapstructure:"cookie_domain"`
	WG            struct {
		// Set this to false to disable the embedded wireguard
		// server. This is useful for development environments
		// on mac and windows where we don't currently support
		// the OS's network stack.
		Enabled bool `mapstructure:"enabled"`
		// The network interface name of the WireGuard
		// network device.
		// Defaults to wg0
		Interface string `mapstructure:"interface"`
		// The WireGuard PrivateKey
		// If this value is lost then any existing
		// clients (WireGuard peers) will no longer
		// be able to connect.
		// Clients will either have to manually update
		// their connection configuration or setup
		// their VPN again using the web ui (easier for most people)
		PrivateKey string `mapstructure:"privateKey"`
		// The WireGuard ListenPort
		// Defaults to 51820
		Port int `mapstructure:"port"`
	} `mapstructure:"wg"`
	VPN struct {
		// CIDR configures a network address space
		// that client (WireGuard peers) will be allocated
		// an IP address from
		// defaults to 10.44.0.0/24
		CIDR string `mapstructure:"cidr"`
		// GatewayInterface will be used in iptable forwarding
		// rules that send VPN traffic from clients to this interface
		// Most use-cases will want this interface to have access
		// to the outside internet
		GatewayInterface string `mapstructure:"gatewayInterface"`
		// The "AllowedIPs" for VPN clients.
		// This value will be included in client config
		// files and in server-side iptable rules
		// to enforce network access.
		// defaults to ["0.0.0.0/0"]
		AllowedIPs []string `mapstructure:"allowedIPs"`
	} `mapstructure:"vpn"`
	// Configure the embeded DNS server
	DNS struct {
		// Enabled allows you to turn on/off
		// the VPN DNS proxy feature.
		// DNS Proxying is enabled by default.
		Enabled bool `mapstructure:"enabled"`
		// Upstream configures the addresses of upstream
		// DNS servers to which client DNS requests will be sent to.
		// Defaults the host's upstream DNS servers (via resolveconf)
		// or 1.1.1.1 if resolveconf cannot be used.
		// NOTE: currently wg-access-server will only use the first upstream.
		Upstream []string `mapstructure:"upstream"`
	} `mapstructure:"dns"`
}

var Spec Config
