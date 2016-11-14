--
-- Example code for using triggers in resources
--

-- Manage the SNMP package
pkg = resource.package.new("net-snmp")
pkg.state = "present"

-- Manage the config file for SNMP daemon
config = resource.file.new("/etc/snmp/snmpd.conf")
config.state = "present"
config.content = "rocommunity public"
config.require = { pkg:ID() }

-- Manage the SNMP service
svc = resource.service.new("snmpd")
svc.state = "running"
svc.enable = true
svc.require = { pkg:ID(), config:ID() }

-- Subscribe for changes in the config file resource.
-- Reload the SNMP daemon service if the config file has changed.
svc.subscribe[config:ID()] = function()
   os.execute("systemctl reload snmpd")
end

-- Subscribe for changes in the package resource.
-- Restart the SNMP daemon service if the package has changed.
svc.subscribe[pkg:ID()] = function()
   os.execute("systemctl restart snmpd")
end

-- Add resources to the catalog
catalog:add(pkg, config, svc)
