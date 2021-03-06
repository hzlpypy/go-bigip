package bigip

import "encoding/json"

const (
	uriSys            = "sys"
	uriFolder         = "folder"
	uriSyslog         = "syslog"
	uriSoftware       = "software"
	uriVolume         = "volume"
	uriHardware       = "hardware"
	uriGlobalSettings = "global-settings"
	uriManagementIp   = "management-ip"
	uriCrypto         = "crypto"
	uriCert           = "cert"
	uriKey            = "key"
	//uriPlatform = "?$select=platform"
)

type Volumes struct {
	Volumes []Volume `json:"items,omitempty"`
}

type Volume struct {
	Name       string `json:"items,omitempty"`
	FullPath   string `json:"fullPath,omitempty"`
	Generation int    `json:"generation,omitempty"`
	SelfLink   string `json:"selfLink,omitempty"`
	Active     bool   `json:"active,omitempty"`
	BaseBuild  string `json:"basebuild,omitempty"`
	Build      string `json:"build,omitempty"`
	Product    string `json:"product,omitempty"`
	Status     string `json:"status,omitempty"`
	Version    string `json:"version,omitempty"`
}

// Volumes returns a list of Software Volumes.
func (b *BigIP) Volumes() (*Volumes, error) {
	var volumes Volumes
	err, _ := b.getForEntity(&volumes, uriSys, uriSoftware, uriVolume)
	if err != nil {
		return nil, err
	}

	return &volumes, nil
}

type ManagementIP struct {
	Addresses []ManagementIPAddress
}

type ManagementIPAddress struct {
	Name       string `json:"items,omitempty"`
	FullPath   string `json:"fullPath,omitempty"`
	Generation int    `json:"generation,omitempty"`
	SelfLink   string `json:"selfLink,omitempty"`
}

func (b *BigIP) ManagementIPs() (*ManagementIP, error) {
	var managementIP ManagementIP
	err, _ := b.getForEntity(&managementIP, uriSys, uriManagementIp)
	if err != nil {
		return nil, err
	}

	return &managementIP, nil
}

type SyslogRemoteServer struct {
	Name       string `json:"name,omitempty"`
	Host       string `json:"host,omitempty"`
	LocalIP    string `json:"localIp,omitempty"`
	RemotePort int    `json:"remotePort,omitempty"`
}

type Syslog struct {
	SelfLink      string               `json:"selfLink,omitempty"`
	RemoteServers []SyslogRemoteServer `json:"remoteServers,omitempty"`
}

func (b *BigIP) Syslog() (*Syslog, error) {
	var syslog Syslog

	err, _ := b.getForEntity(&syslog, uriSys, uriSyslog)
	if err != nil {
		return nil, err
	}

	return &syslog, nil
}

func (b *BigIP) SetSyslog(config Syslog) error {
	return b.put(config, uriSys, uriSyslog)
}

// Folders contains a list of every folder on the BIG-IP system.
type Folders struct {
	Folders []Folder `json:"items"`
}

type folderDTO struct {
	Name      string `json:"name,omitempty"`
	Partition string `json:"partition,omitempty"`
	SubPath   string `json:"subPath,omitempty"`
	FullPath  string `json:"fullPath,omitempty"`

	AppService  string `json:"appService,omitempty"`
	Description string `json:"description,omitempty"`
	// Set to "default" to inherit or a device group name to control. You can also set it to "non-default" to pin its device group to its current setting and turn off inheritance.
	DeviceGroup string `json:"deviceGroup,omitempty"`
	Hidden      string `json:"hidden,omitempty" bool:"true"`
	NoRefCheck  string `json:"noRefCheck,omitempty" bool:"true"`
	// Set to "default" to inherit or a traffic group name to control. You can also set it to "non-default" to pin its traffic group to its current setting and turn off inheritance.
	TrafficGroup string `json:"trafficGroup,omitempty"`

	// Read-only property. Set DeviceGroup to control.
	InheritedDeviceGroup string `json:"inheritedDevicegroup,omitempty" bool:"true"`

	// Read-only property. Set TrafficGroup to control.
	InheritedTrafficGroup string `json:"inheritedTrafficGroup,omitempty" bool:"true"`
}

type Folder struct {
	Name      string `json:"name,omitempty"`
	Partition string `json:"partition,omitempty"`
	SubPath   string `json:"subPath,omitempty"`
	FullPath  string `json:"fullPath,omitempty"`

	AppService   string `json:"appService,omitempty"`
	Description  string `json:"description,omitempty"`
	DeviceGroup  string `json:"deviceGroup,omitempty"`
	Hidden       *bool  `json:"hidden,omitempty"`
	NoRefCheck   *bool  `json:"noRefCheck,omitempty"`
	TrafficGroup string `json:"trafficGroup,omitempty"`

	// Read-only property. Set DeviceGroup to "default" or "non-default" to control.
	InheritedDeviceGroup *bool `json:"inheritedDevicegroup,omitempty"`

	// Read-only property. Set TrafficGroup to "default" or "non-default" to control.
	InheritedTrafficGroup *bool `json:"inheritedTrafficGroup,omitempty"`
}

func (f *Folder) MarshalJSON() ([]byte, error) {
	var dto folderDTO
	marshal(&dto, f)
	return json.Marshal(dto)
}

func (f *Folder) UnmarshalJSON(b []byte) error {
	var dto folderDTO
	err := json.Unmarshal(b, &dto)
	if err != nil {
		return err
	}
	return marshal(f, &dto)
}

// Folders returns a list of folders.
func (b *BigIP) Folders() (*Folders, error) {
	var folders Folders
	err, _ := b.getForEntity(&folders, uriSys, uriFolder)
	if err != nil {
		return nil, err
	}

	return &folders, nil
}

// CreateFolder adds a new folder to the BIG-IP system.
func (b *BigIP) CreateFolder(name string) error {
	config := &Folder{
		Name: name,
	}

	return b.post(config, uriSys, uriFolder)
}

// AddFolder adds a new folder by config to the BIG-IP system.
func (b *BigIP) AddFolder(config *Folder) error {

	return b.post(config, uriSys, uriFolder)
}

// GetFolder retrieves a Folder by name. Returns nil if the folder does not exist
func (b *BigIP) GetFolder(name string) (*Folder, error) {
	var folder Folder
	err, ok := b.getForEntity(&folder, uriSys, uriFolder, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &folder, nil
}

// DeleteFolder removes a folder.
func (b *BigIP) DeleteFolder(name string) error {
	return b.delete(uriSys, uriFolder, name)
}

// ModifyFolder allows you to change any attribute of a folder. Fields that can
// be modified are referenced in the Folder struct. This replaces the existing
// configuration, so use PatchFolder if you want to change only particular
// attributes.
func (b *BigIP) ModifyFolder(name string, config *Folder) error {
	return b.put(config, uriSys, uriFolder, name)
}

// PatchFolder allows you to change any attribute of a folder. Fields that can
// be modified are referenced in the Folder struct. This changes only the
// attributes provided, so use ModifyFolder if you want to replace the existing
// configuration.
func (b *BigIP) PatchFolder(name string, config *Folder) error {
	return b.patch(config, uriSys, uriFolder, name)
}

// Certificates represents a list of installed SSL certificates.
type Certificates struct {
	Certificates []Certificate `json:"items,omitempty"`
}

// Certificate represents an SSL Certificate.
type Certificate struct {
	APIRawValues *struct {
		CertificateKeySize string `json:"certificateKeySize,omitempty"`
		Expiration         string `json:"expiration,omitempty"`
		PublicKeyType      string `json:"publicKeyType,omitempty"`
	} `json:"apiRawValues,omitempty"`
	AppService             string `json:"appService,omitempty"`
	CertValidationOptions  string `json:"certValidationOptions,omitempty"`
	City                   string `json:"city,omitempty"`
	CommonName             string `json:"commonName,omitempty"`
	Command                string `json:"command,omitempty"`
	Consumer               string `json:"consumer,omitempty"`
	Country                string `json:"country,omitempty"`
	EmailAddress           string `json:"emailAddress,omitempty"`
	FromLocalFile          string `json:"from-local-file,omitempty"`
	FromURL                string `json:"from-url,omitempty"`
	FullPath               string `json:"fullPath,omitempty"`
	Generation             int    `json:"generation,omitempty"`
	IssuerCert             string `json:"issuerCert,omitempty"`
	Key                    string `json:"key,omitempty"`
	Lifetime               string `json:"lifetime,omitempty"`
	Name                   string `json:"name,omitempty"`
	Organization           string `json:"organization,omitempty"`
	Ou                     string `json:"ou,omitempty"`
	Partition              string `json:"partition,omitempty"`
	State                  string `json:"state,omitempty"`
	SubjectAlternativeName string `json:"subjectAlternativeName,omitempty"`
}

// Certificates returns a list of certificates.
func (b *BigIP) Certificates() (*Certificates, error) {
	var certs Certificates
	err, _ := b.getForEntity(&certs, uriSys, uriCrypto, uriCert)
	if err != nil {
		return nil, err
	}

	return &certs, nil
}

// AddCertificate installs a certificate.
func (b *BigIP) AddCertificate(cert *Certificate) error {
	return b.post(cert, uriSys, uriCrypto, uriCert)
}

// GetCertificate retrieves a Certificate by name. Returns nil if the certificate does not exist
func (b *BigIP) GetCertificate(name string) (*Certificate, error) {
	var cert Certificate
	err, ok := b.getForEntity(&cert, uriSys, uriCrypto, uriCert, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &cert, nil
}

// DeleteCertificate removes a certificate.
func (b *BigIP) DeleteCertificate(name string) error {
	return b.delete(uriSys, uriCrypto, uriCert, name)
}

// Keys represents a list of installed keys.
type Keys struct {
	Keys []Key `json:"items,omitempty"`
}

// Key represents a private key associated with a certificate.
type Key struct {
	AdminEmailAddress      string `json:"adminEmailAddress,omitempty"`
	AppService             string `json:"appService,omitempty"`
	ChallengePassword      string `json:"challengePassword,omitempty"`
	City                   string `json:"city,omitempty"`
	Command                string `json:"command,omitempty"`
	CommonName             string `json:"commonName,omitempty"`
	Consumer               string `json:"consumer,omitempty"`
	Country                string `json:"country,omitempty"`
	CurveName              string `json:"curveName,omitempty"`
	EmailAddress           string `json:"emailAddress,omitempty"`
	FromLocalFile          string `json:"from-local-file,omitempty"`
	FromURL                string `json:"from-url,omitempty"`
	FullPath               string `json:"fullPath,omitempty"`
	Generation             int    `json:"generation,omitempty"`
	KeySize                string `json:"keySize,omitempty"`
	KeyType                string `json:"keyType,omitempty"`
	Lifetime               string `json:"lifetime,omitempty"`
	Name                   string `json:"name,omitempty"`
	Organization           string `json:"organization,omitempty"`
	Ou                     string `json:"ou,omitempty"`
	Partition              string `json:"partition,omitempty"`
	Passphrase             string `json:"passphrase,omitempty"`
	SecurityType           string `json:"securityType,omitempty"`
	State                  string `json:"state,omitempty"`
	SubjectAlternativeName string `json:"subjectAlternativeName,omitempty"`
}

// Keys returns a list of keys.
func (b *BigIP) Keys() (*Keys, error) {
	var keys Keys
	err, _ := b.getForEntity(&keys, uriSys, uriCrypto, uriKey)
	if err != nil {
		return nil, err
	}

	return &keys, nil
}

// AddKey installs a key.
func (b *BigIP) AddKey(config *Key) error {
	return b.post(config, uriSys, uriCrypto, uriKey)
}

// GetKey retrieves a key by name. Returns nil if the key does not exist.
func (b *BigIP) GetKey(name string) (*Key, error) {
	var key Key
	err, ok := b.getForEntity(&key, uriSys, uriCrypto, uriKey, name)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, nil
	}

	return &key, nil
}

// DeleteKey removes a key.
func (b *BigIP) DeleteKey(name string) error {
	return b.delete(uriSys, uriCrypto, uriKey, name)
}
