## RPM 활용 워크플로 예시 1. 실디바이스 - PV 데이터 공유
**Requirement**
- 실디바이스에 davfs2가 설치되어 있어야 합니다.

**엣지 서버**
1. 볼륨 생성(Create Volume) 기능으로 PVC를 생성합니다.
2. 생성된 볼륨 ID로 볼륨 마운트(Mount Volume) 기능을 실행합니다. 

**실디바이스**
1. PV를 마운트할 포인트를 찾습니다.
    ex. /mount
2. 웹다브 서버와 마운트합니다.
    1. mount 명령어를 실행합니다.
    ```
        sudo mount -t davfs http://<master node 주소>:<포트번호>/uploads <마운트할 포인트>
    ```
    2. 마운트 되었는지 확인합니다.
        mount failed 없이 넘어간다면 정상적으로 마운트 된 것입니다.
