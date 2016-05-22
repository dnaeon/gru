package "memcached" {
  state = "present"
}

file "/etc/systemd/system/memcached.service.d" {
  state = "present"
  filetype = "directory"
  require = [
    "package[memcached]",
  ]
}

file "/etc/systemd/system/memcached.service.d/override.conf" {
  state = "present"
  mode = 0644
  source = "data/memcached/memcached-override.conf"
  require = [
    "file[/etc/systemd/system/memcached.service.d]",
  ]
}

shell "systemctl daemon-reload" {
  require = [
    "file[/etc/systemd/system/memcached.service.d/override.conf]",
  ]
}

service "memcached" {
  state = "running"
  enable = true
  require = [
    "package[memcached]",
    "file[/etc/systemd/system/memcached.service.d/override.conf]",
  ]
}
