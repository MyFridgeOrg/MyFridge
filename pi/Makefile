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
	scp provision/myfridge.service $(PI_HOST):myfridge.service
	ssh -t $(PI_HOST) 'sudo mkdir -p /opt/myfridge && sudo install -m 0644 myfridge.service /etc/systemd/system/myfridge.service && rm myfridge.service && sudo systemctl daemon-reload && sudo systemctl enable $(SERVICE)'

logs:
	ssh -t $(PI_HOST) 'journalctl -u $(SERVICE) -f'

status:
	ssh $(PI_HOST) '\
	  systemctl status $(SERVICE) --no-pager -n 0; \
	  echo; echo "== Camera =="; rpicam-hello --list-cameras; \
	  echo; echo "== PIR sensor (GPIO4) =="; gpioinfo -c gpiochip0 4; pinctrl get 4; \
	  echo; echo "== Last snapshot =="; ls -l --time-style="+%F %T" /var/lib/myfridge/*.jpg 2>/dev/null || echo "none yet"'

start:
	ssh $(PI_HOST) 'sudo -n /usr/bin/systemctl start $(SERVICE)'

stop:
	ssh $(PI_HOST) 'sudo -n /usr/bin/systemctl stop $(SERVICE)'