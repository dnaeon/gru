--
-- Gru module for installing and configuring memcached
--

-- Manage the memcached package
memcached_pkg = resource.package.new("memcached")
memcached_pkg.state = "present"

-- Path to the systemd drop-in unit directory
systemd_dir = "/etc/systemd/system/memcached.service.d/"

-- Manage the systemd drop-in unit directory
unit_dir = resource.file.new(systemd_dir)
unit_dir.state = "present"
unit_dir.filetype = "directory"
unit_dir.require = {
   memcached_pkg:ID(),
}

-- Manage the systemd drop-in unit
unit_file = resource.file.new(systemd_dir .. "override.conf")
unit_file.state = "present"
unit_file.mode = tonumber("0644", 8)
unit_file.source = "data/memcached/memcached-override.conf"
unit_file.require = {
   unit_dir:ID(),
}

-- Instruct systemd(1) to reload it's configuration
systemd_reload = resource.shell.new("systemctl daemon-reload")
systemd_reload.require = {
   unit_file:ID(),
}

-- Manage the memcached service
memcached_svc = resource.service.new("memcached")
memcached_svc.state = "running"
memcached_svc.enable = true
memcached_svc.require = {
   memcached_pkg:ID(),
   unit_file:ID(),
}

-- Finally, register the resources to the catalog
catalog:add(memcached_pkg, unit_dir, unit_file, systemd_reload, memcached_svc)
