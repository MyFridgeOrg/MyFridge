# Deployment

This guide sets up the fridge unit so that MyFridge starts automatically when the Pi boots, and so that you can deploy a new version from your Mac with a single command:

```bash
make deploy
```

You create every file yourself. The complete contents are below.

## Overview

| What | Where | Purpose |
|---|---|---|
| `pi/Makefile` | Repo (Mac) | Build, deploy and manage the service from your Mac |
| `pi/deploy/myfridge.service` | Repo (Mac) | The systemd service definition |
| `pi/.gitignore` | Repo (Mac) | Keeps build output out of git |
| `/opt/myfridge/myfridge` | Pi | The installed binary |
| `/var/lib/myfridge/` | Pi | Where snapshots are saved |
| `/etc/systemd/system/myfridge.service` | Pi | Installed copy of the service definition |
| `/etc/sudoers.d/myfridge` | Pi | Lets `make deploy` install and restart without a password |

The service runs as **root**, so there are no users or groups to manage on the Pi.

The Pi is reached at `ghmeding@root.local`. If that changes, override it per command, for example `make deploy PI_HOST=ghmeding@192.168.1.16`.

## 1. Repo files

### `pi/deploy/myfridge.service`

```ini
[Unit]
Description=MyFridge motion-triggered camera

[Service]
ExecStart=/opt/myfridge/myfridge
StateDirectory=myfridge
WorkingDirectory=/var/lib/myfridge

[Install]
WantedBy=multi-user.target
```

Line by line:

- **`Description`**: the name shown in `systemctl status` and the logs.
- **`ExecStart`**: the program to run.
- **`StateDirectory=myfridge`**: systemd creates `/var/lib/myfridge` before starting the program, if it doesn't exist yet.
- **`WorkingDirectory=/var/lib/myfridge`**: the program runs from this folder. `cmd/main.go` saves snapshots to `"."`, so they end up here without any code change.
- There's no `Restart=` line, so if the program exits with an error (for example when the coordinator calls `log.Fatalf` after a failed snapshot), the service stays stopped. Check `make logs`, fix the problem, then run `make start`.
- **`WantedBy=multi-user.target`**: start the service during a normal boot. This is what `systemctl enable` hooks into.
- There's no `User=` line, so the service runs as root. That gives it access to the GPIO pins and the camera without any group setup.

### `pi/Makefile`

```make
PI_HOST ?= ghmeding@root.local

BINARY  := bin/myfridge
UPLOAD  := /home/ghmeding/myfridge-new
SERVICE := myfridge

.PHONY: build deploy install logs status start stop

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o $(BINARY) ./cmd

deploy: build
	scp $(BINARY) $(PI_HOST):$(UPLOAD)
	ssh $(PI_HOST) 'sudo -n /usr/bin/install -m 0755 $(UPLOAD) /opt/myfridge/myfridge && rm $(UPLOAD) && sudo -n /usr/bin/systemctl restart $(SERVICE)'

install:
	scp deploy/myfridge.service $(PI_HOST):myfridge.service
	ssh -t $(PI_HOST) 'sudo mkdir -p /opt/myfridge && sudo install -m 0644 myfridge.service /etc/systemd/system/myfridge.service && rm myfridge.service && sudo systemctl daemon-reload && sudo systemctl enable $(SERVICE)'

logs:
	ssh -t $(PI_HOST) 'journalctl -u $(SERVICE) -f'

status:
	ssh $(PI_HOST) 'systemctl status $(SERVICE) --no-pager'

start:
	ssh $(PI_HOST) 'sudo -n /usr/bin/systemctl start $(SERVICE)'

stop:
	ssh $(PI_HOST) 'sudo -n /usr/bin/systemctl stop $(SERVICE)'
```

> **Recipe lines must start with a tab, not spaces.** If `make` reports `missing separator`, your editor replaced the tab with spaces.

Variables:

- **`PI_HOST ?= ...`**: the `?=` means "use this unless set already", so `make deploy PI_HOST=...` overrides it.
- **`BINARY`**: where the build is written on your Mac.
- **`UPLOAD`**: where the binary is copied on the Pi before it's installed. It's in your home folder rather than `/tmp`, because anyone on the Pi can write to `/tmp`, and root is about to run whatever file is there. It's an absolute path because the sudoers rule (step 2) only matches exact commands.

Targets:

- **`build`**: cross-compiles for the Pi's 64-bit ARM CPU. `CGO_ENABLED=0` makes a fully static binary, so nothing extra is needed on the Pi.
- **`deploy`**: builds, copies the binary to the Pi, installs it into `/opt/myfridge`, deletes the uploaded copy and restarts the service. `restart` also starts the service if it wasn't running. `sudo -n` means "never ask for a password": if the sudoers rule is missing, the command fails straight away instead of hanging.
- **`install`**: one-time setup. It creates `/opt/myfridge`, installs the service file, reloads systemd and enables the service at boot. It uses `ssh -t` and plain `sudo`, so it asks for your Pi password. It doesn't start the service yet, because there's no binary until the first `make deploy`. Run it again whenever you change `myfridge.service`.
- **`logs`**: follows the service logs live. Press `Ctrl+C` to stop.
- **`status`**: shows whether the service is running and its last few log lines.
- **`start`** and **`stop`**: start or stop the service, for example around manual hardware tests.

### `pi/.gitignore`

```gitignore
bin/
```

Then delete the old binary that was built straight into `pi/`:

```bash
rm pi/myfridge
```

## 2. One-time setup on the Pi

### Allow deploys without a password

SSH into the Pi and open a new sudoers file. Always use `visudo`, which checks the syntax before saving. A broken sudoers file can lock you out of `sudo`.

```bash
sudo visudo -f /etc/sudoers.d/myfridge
```

Paste this as one line:

```
ghmeding ALL=(root) NOPASSWD: /usr/bin/install -m 0755 /home/ghmeding/myfridge-new /opt/myfridge/myfridge, /usr/bin/systemctl restart myfridge, /usr/bin/systemctl start myfridge, /usr/bin/systemctl stop myfridge
```

This allows exactly these four commands, with exactly these arguments, without a password. Every other `sudo` command still asks for one.

Check that it works:

```bash
sudo -l
```

The four commands should be listed under `NOPASSWD`.

### Remove the manual setup

If you're still running the old binary by hand, stop it with `Ctrl+C`. The old copy in your home folder is no longer needed:

```bash
rm ~/myfridge
```

## 3. First install

From the `pi/` folder on your Mac:

```bash
make install   # asks for your Pi password once
make deploy    # builds, uploads and starts the service
make status    # should show "active (running)"
```

Then check that it starts on boot:

```bash
ssh ghmeding@root.local 'sudo reboot'
```

Wait about a minute, then run `make status` again.

## 4. Daily workflow

All commands are run from `pi/` on your Mac.

| Task | Command |
|---|---|
| Deploy a new version | `make deploy` |
| Watch the logs (banner, snapshots, errors) | `make logs` |
| Check that it's running | `make status` |
| Stop it before manual hardware tests | `make stop` |
| Start it again | `make start` |

Copy the latest snapshot to your Mac:

```bash
scp ghmeding@root.local:/var/lib/myfridge/motion_.jpg ~/Downloads/
```

**Stop the service before manual hardware tests.** The service holds the GPIO pin and uses the camera. If it's running, `gpiomon`, `gpioget` or `rpicam-still` will fail or behave oddly. Run `make stop` first and `make start` afterwards.

## 5. Troubleshooting

**`make deploy` fails with `sudo: a password is required`**
The sudoers rule doesn't match the command. Check it with `sudo -l` on the Pi. The paths and arguments in the rule must be identical to those in the Makefile, including `-m 0755` and the full upload path.

**`make deploy` fails with `Failed to restart myfridge.service: Unit myfridge.service not found`**
Run `make install` first.

**`make status` shows `failed`**
Run `make logs` to see the error.
- `Failed to find pin 4`: the GPIO library didn't recognise the pin name. Try `"GPIO4"` in `cmd/main.go`.
- `CSI camera hardware failed execution`: check the camera with `make stop`, then `rpicam-still -o test.jpg -t 2000` on the Pi.

**`gpiomon` reports `Device or resource busy`, or `rpicam-still` reports that the camera is in use**
The service is still running. Run `make stop`.

**`make logs` shows `No journal files were found` or a permission error**
Your user can't read the system logs. Either change the `logs` target to use `sudo journalctl ...`, or add yourself to the `adm` group on the Pi:

```bash
sudo usermod -aG adm ghmeding
```

Then log out and back in.

**`make` reports `missing separator`**
A recipe line in the Makefile starts with spaces instead of a tab.
