## Enabling services during boot-time

In order to start your minions during boot-time you have two
options.

If your minions are running on a system that supports
[systemd](https://www.freedesktop.org/wiki/Software/systemd/),
you could use the provided systemd unit file for Gru.

Or you could use [supervisord](http://supervisord.org/) for process
control.

## systemd unit

Get the systemd unit file from the [contrib/systemd](../contrib)
directory and install it on your system.

Check [Unit File Load Path](https://www.freedesktop.org/software/systemd/man/systemd.unit.html#Unit%20File%20Load%20Path)
document for the location where you should install your unit file.

Once you've got the unit in place, execute these commands which will
enable and start your minion.

```bash
$ sudo systemctl daemon-reload
$ sudo systemctl enable gru-minion
$ sudo systemctl start gru-minion
```

## supervisord

Get the supervisord config file from [contrib/supervisord](../contrib/supervisord)
directory and place in under your supervisord `include` directory.

Once you've got the file in place, reload the supervisord
configuration, enable and start the service.

```bash
$ sudo supervisorctl reread
$ sudo supervisorctl reload
$ sudo supervisorctl start gru-minion
```
