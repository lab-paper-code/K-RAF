## Example workflow utilizing RPM 2. Application migration
**Requirement**
- You must have davfs2 installed on your physical device.
- Example 1. Silo - Utilize the PV data share and the daclab/inference_migration image.

**Edge server master node (the server running the RPM)
1. create a PVC with the Create Volume function.
2. Run the Mount Volume function with the created volume ID. 
3. register the application daclab/inference_migration:0.3 with the RPM (Register App)
4. Run the application Pod (Execute AppRun).
    - The AppID is the App you registered in step 3, the DeviceID is the DeviceID of the physical device you want to run the application on, and the VolumeID is the VolumeID of the Volume you created in step 1.
5. Connect to the application Pod.
    ```
    kubectl exec -it <inference_migration application-pod> -n ksv -- bash
    to the application pod.
    ```
6. Change the path to the image folder.
    ```
    sh setup.sh
    Move the images folder inside the /app directory to /uploads (the mounted path).
    ```
    
7. Modify /configs/config.py inside the /app directory.
    IMG_PATH: Enter the path where the images are stored mounted with WebDAV on a pod basis.
    OUTPUT_PATH: Enter the path to save inference results on a per-pod basis.
    HOST : Enter the internal IP (10.X.X.X) of the pod.
    PORT : <Port exposed when registering the app>.

8. (**Important**) Grant write permissions to the output, output2 folders in /app
```
    sudo chmod 777 output output2 
```
output is the folder where the results of processing in the Pod are stored, and output2 is the folder where the results of processing on the physical device are stored.
If you only grant write access to mounted folders on the physical device, the physical device results will not be saved properly.

For **SilverDevice**
1. Find the point where you want to mount the PV.
    ex. /mount
2. Mount it with the webdev server.
    1. 
    ```
        sudo mount -t davfs http://<master node address>:<port number>/uploads <point to mount>
    ```
    2. verify that it mounted.
        If it passes without mount failed, it's mounted successfully.
3. copy the repository with git clone.
```
    The relevant files are in the master branch.
    git clone -b master https://github.com/lab-paper-code/img_inference_migration.git
```
4. Modify device.config.py.
    In img_inference_migration/configs, modify device.config.py.
    1. to_send_addr : Enter the fastapi enabled pod address (the server address running the pod) and the service (NodePort) port. 
    2. OUTPUT_PATH: Enter the path to save inference results based on the actual device.
    3. IMG_PATH : Enter the path where the image is saved by mounting webdav on Raspberry Pi.

5. grant write permissions to the output, output2 folder on the mount point (**important**)
```
    sudo chmod 777 output output2
```

6. Run python3 server.py and sh run_device.sh on the edge server (/app in the application pod) and on the actual device img_inference_migration/ path, respectively.
- Server.py inside the application pod runs the server with FastAPI and sends requests to run_device.sh on the silicon device, so server.py must be run first.
