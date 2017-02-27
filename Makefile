BINARY=ec2-manager

build:
	go build -o ${BINARY}

.PHONY: install
install:
	go install

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi