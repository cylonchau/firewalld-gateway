package apis

const (
	INTERFACE      = "org.fedoraproject.FirewallD1"
	PATH           = "/org/fedoraproject/FirewallD1"
	DIRECT         = INTERFACE + ".direct"
	IPSET          = INTERFACE + ".ipset"
	POLICIES       = INTERFACE + ".policies"
	ZONE           = INTERFACE + ".zone"
	INTROSPECTABLE = "org.freedesktop.DBus.Introspectable"
	PROPERTIES     = "org.freedesktop.DBus.Properties"

	CONFIG_PATH               = PATH + "/config"
	CONFIG_INTERFACE          = INTERFACE + ".config"
	CONFIG_DIRECT_INTERFACE   = INTERFACE + ".config.direct"
	CONFIG_POLICIES_INTERFACE = INTERFACE + ".config.policies"

	ZONE_PATH      = CONFIG_PATH + "/zone"
	ZONE_INTERFACE = INTERFACE + ".config.zone"

	SERVICE_PATH      = PATH + "/config/service/i"
	SERVICE_INTERFACE = INTERFACE + ".config.service"

	IPSET_PATH      = PATH + "/config/ipset/i"
	IPSET_INTERFACE = INTERFACE + ".config.ipset"

	ICMP_PATH      = PATH + "/config/icmptype/i"
	ICMP_INTERFACE = INTERFACE + ".config.icmptype"

	// org.fedoraproject.FirewallD1
	INTERFACE_GETDEFAULTZONE     = INTERFACE + ".getDefaultZone"
	INTERFACE_SETDEFAULTZONE     = INTERFACE + ".setDefaultZone"
	INTERFACE_GETZONESETTINGS    = INTERFACE + ".getZoneSettings"
	INTERFACE_RELOAD             = INTERFACE + ".completeReload"
	INTERFACE_LISTSERVICES       = INTERFACE + ".listServices"
	INTERFACE_RUNTIMETOPERMANENT = INTERFACE + ".runtimeToPermanent"

	//config
	CONFIG_ADDSERVICE = CONFIG_INTERFACE + ".addService"

	// org.fedoraproject.FirewallD1.zone
	ZONE_ADDMASQUERADE     = ZONE + ".addMasquerade"
	ZONE_REMOVEMASQUERADE  = ZONE + ".removeMasquerade"
	ZONE_QUERYMASQUERADE   = ZONE + ".queryMasquerade"
	ZONE_ADDPORT           = ZONE + ".addPort"
	ZONE_REMOVEPORT        = ZONE + ".removePort"
	ZONE_ADDPROTOCOL       = ZONE + ".addProtocol"
	ZONE_ADDRICHRULE       = ZONE + ".addRichRule"
	ZONE_ADDSERVICE        = ZONE + ".addService"
	ZONE_GETSERVICES       = ZONE + ".getServices"
	ZONE_ADDSOURCE         = ZONE + ".addSource"
	ZONE_ADDINTERFACE      = ZONE + ".addInterface"
	ZONE_QUERYINTERFACE    = ZONE + ".queryInterface"
	ZONE_REMOVEINTERFACE   = ZONE + ".removeInterface"
	ZONE_REOMVERICHRULE    = ZONE + ".removeRichRule"
	ZONE_GETPORTS          = ZONE + ".getPorts"
	ZONE_ADDFORWARDPORT    = ZONE + ".addForwardPort"
	ZONE_GETFORWARDPORT    = ZONE + ".getForwardPorts"
	ZONE_REMOVEFORWARDPORT = ZONE + ".removeForwardPort"
	ZONE_QUERYFORWARDPORT  = ZONE + ".queryForwardPort"

	// get
	ZONE_GETZONES           = ZONE + ".getZones"
	ZONE_GETZONEOFINTERFACE = ZONE + ".getZoneOfInterface"
	ZONE_GETRICHRULES       = ZONE + ".getRichRules"
	ZONE_QUERYRICHRULE      = ZONE + ".queryRichRule"
	ZONE_QUERYSERVICE       = ZONE + ".queryService"
	ZONE_REMOVESERVICE      = ZONE + ".removeService"

	// org.fedoraproject.FirewallD1.config
	CONFIG_ADDZONE = CONFIG_INTERFACE + ".addZone"

	// org.fedoraproject.FirewallD1.config.zone
	CONFIG_ZONE                   = CONFIG_INTERFACE + ".zone"
	CONFIG_UPDATE                 = CONFIG_ZONE + ".update"
	CONFIG_ZONE_ADDRICHRULE       = CONFIG_ZONE + ".addRichRule"
	CONFIG_ZONE_REOMVERICHRULE    = CONFIG_ZONE + ".removeRichRule"
	CONFIG_ZONE_QUERYRICHRULE     = CONFIG_ZONE + ".queryRichRule"
	CONFIG_ZONE_ADDSERVICE        = CONFIG_ZONE + ".addService"
	CONFIG_ZONE_QUERYSERVICE      = CONFIG_ZONE + ".queryService"
	CONFIG_ZONE_REMOVESERVICE     = CONFIG_ZONE + ".removeService"
	CONFIG_ZONE_GETSERVICES       = CONFIG_ZONE + ".getServices"
	CONFIG_ZONE_ADDPORT           = CONFIG_ZONE + ".addPort"
	CONFIG_ZONE_GETPORTS          = CONFIG_ZONE + ".getPorts"
	CONFIG_ZONE_REMOVEPORT        = CONFIG_ZONE + ".removePort"
	CONFIG_ZONE_ADDMASQUERADE     = CONFIG_ZONE + ".addMasquerade"
	CONFIG_ZONE_REMOVEMASQUERADE  = CONFIG_ZONE + ".removeMasquerade"
	CONFIG_ZONE_QUERYMASQUERADE   = CONFIG_ZONE + ".queryMasquerade"
	CONFIG_ZONE_ADDINTERFACE      = CONFIG_ZONE + ".addInterface"
	CONFIG_ZONE_REMOVEINTERFACE   = CONFIG_ZONE + ".removeInterface"
	CONFIG_ZONE_ADDFORWARDPORT    = CONFIG_ZONE + ".addForwardPort"
	CONFIG_ZONE_REMOVEFORWARDPORT = CONFIG_ZONE + ".removeForwardPort"
	CONFIG_ZONE_QUERYFORWARDPORT  = CONFIG_ZONE + ".queryForwardPort"
	CONFIG_GETFORWARDPORT         = CONFIG_ZONE + ".getForwardPorts"
	CONFIG_REMOVEZONE             = CONFIG_ZONE + ".remove"
	CONFIG_DEFAULT_POLICY         = CONFIG_ZONE + ".getTarget"
)
