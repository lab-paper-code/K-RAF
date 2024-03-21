## basic setting and docker install

### SSH 및 방화벽 해제
```
apt-get update
(sudo) apt-get update

ssh openssh-server를 설치합니다.
(sudo) apt-get install openssh-server

ssh 클라이언트와 서버를 동시에 설치합니다.
(sudo) apt-get install ssh

```

```
(sudo) ufw enable
(sudo) ufw allow 22
(sudo) ufw reload

ssh 서비스를 시작합니다.

(sudo) service ssh start

ssh daemon이 제대로 실행되는지 확인합니다.

(sudo) service ssh status
(sudo) ps -ef | grep sshd
(sudo) netstat -ntlp | grep sshd

```
### Docker

```
https를 사용해서 레포지토리를 사용할 수 있도록 필요한 패키지를 설치합니다.
$ sudo apt-get install -y  apt-transport-https ca-certificates curl software-properties-common

Docker 공식 리포지토리에서 패키지를 다운로드 받았을 때 위변조 확인을 위한 GPG 키를 추가합니다.
$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add

Docker.com 의 GPG 키가 등록됐는지 확인합니다.
$ apt-key fingerprint

Docker 공식 저장소를 리포지토리로 등록합니다.
$ add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

저장소 등록정보에 기록됐는지 확인합니다.
$ grep docker /etc/apt/sources.list
deb [arch=amd64] https://download.docker.com/linux/ubuntu focal stable

리포지토리 정보를 갱신합니다.
$ sudo-apt update

docker container engine 을 설치합니다.
$ apt-get install -y docker-ce

도커 서비스 상태 및 버전을 확인합니다.
$ ps -ef | grep docker

$ docker --version
```

[참고](https://kindloveit.tistory.com/18)