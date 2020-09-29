install: swbka profile.ini services/swbka.timer services/swbka.service
	install -d /opt/swbka/bin/
	install -m 644 profile.ini /opt/swbka/profile.ini
	install -m 755 swbka /opt/swbka/bin/swbka
	install -d /usr/lib/systemd/system
	install -m 644 services/swbka.timer /usr/lib/systemd/system/ 
	install -m 644 services/swbka.service /usr/lib/systemd/system/
	systemctl daemon-reload
	systemctl enable --now swbka.timer swbka.service
build:
	go build -ldflags "-w -s"


.PHONY: install clean
