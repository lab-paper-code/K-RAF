## RPM 활용 워크플로 예시 2. 애플리케이션 마이그레이션
**Requirement**
- 실디바이스에 davfs2가 설치되어 있어야 합니다.
- 예시 1. 실디바이스 - PV 데이터 공유 내용과 daclab/inference_migration 이미지를 활용합니다.

**엣지 서버 마스터노드(RPM을 실행하는 서버)**
1. 볼륨 생성(Create Volume) 기능으로 PVC를 생성합니다.
2. 생성된 볼륨 ID로 볼륨 마운트(Mount Volume) 기능을 실행합니다. 
3. RPM에 애플리케이션 daclab/inference_migration:0.3을 등록합니다.(Register App)
4. 애플리케이션 파드를 실행합니다.(Execute AppRun)
    - 이때 AppID는 3번에서 등록한 App, DeviceID는 애플리케이션을 실행하려는 실디바이스의 DeviceID, VolumeID는 1번에서 생성한 Volume의 VolumeID를 사용합니다.
5. 애플리케이션 파드에 접속합니다.
    ```
    kubectl exec -it <inference_migration 애플리케이션 파드> -n ksv -- bash
    로 접속합니다.
    ```
6. 이미지 폴더 경로를 변경합니다.
    ```
    sh setup.sh
    /app 디렉토리 내에 있는 이미지 폴더를 /uploads(마운트된 경로)로 이동시킵니다.
    ```
    
7. /app 디렉토리 내 /configs/config.py를 수정합니다.
    IMG_PATH : pod 기준으로 WebDAV와 마운트되어 이미지가 저장된 경로를 입력합니다.
    OUTPUT_PATH : pod 기준으로 inference 결과를 저장할 경로를 입력합니다.
    HOST : pod의 내부 IP(10.X.X.X)를 입력합니다.
    PORT : <App 등록 시 노출시킨 포트>

8. (**중요**) /app 내의 output, output2 폴더에 쓰기 권한을 부여합니다
```
    sudo chmod 777 output output2 
```
output은 파드에서 처리한 결과가 저장되는 폴더, output2는 실디바이스에서 처리한 결과가 저장되는 폴더입니다.
실디바이스에서 마운트된 폴더에만 쓰기 권한을 부여하면 실디바이스 결과가 제대로 저장되지 않습니다.

**실디바이스**
1. PV를 마운트할 포인트를 찾습니다.
    ex. /mount
2. 웹다브 서버와 마운트합니다.
    1. 
    ```
        sudo mount -t davfs http://<master node 주소>:<포트번호>/uploads <마운트할 포인트>
    ```
    2. 마운트 되었는지 확인합니다.
        mount failed 없이 넘어간다면 정상적으로 마운트 된 것입니다.
3. git clone으로 리포지토리를 복사해옵니다.
```
    master 브랜치 내에 관련 파일이 있습니다.
    git clone -b master https://github.com/lab-paper-code/img_inference_migration.git
```
4. device.config.py를 수정합니다.
    img_inference_migration/configs에서 device.config.py를 수정합니다.
    1. to_send_addr : fastapi 켜진 파드 주소(파드를 실행하는 서버 주소>와 서비스(NodePort) 포트를 입력합니다. 
    2. OUTPUT_PATH : 실디바이스 기준 inference 결과를 저장할 경로를 입력합니다.
    3. IMG_PATH : 라즈베리파이 기준 webdav와 마운트 되어서 이미지 저장된 경로를 입력합니다.

5. 마운트 포인트에서 output, output2 폴더에 쓰기 권한을 부여합니다.(**중요**)
```
    sudo chmod 777 output output2
```

6. 엣지 서버(애플리케이션 파드 내 /app), 실디바이스 img_inference_migration/ 경로 에서 각각 python3 server.py, sh run_device.sh 순으로 실행합니다.
- 애플리케이션 파드 내에서 server.py로 FastAPI로 서버를 실행하고, 실디바이스에서 run_device.sh로 요청을 보내므로, server.py를 먼저 실행해야 합니다.
