## Example workflow utilizing RPM 1. Real Device - PV Data Sharing
**Requirement**.
- DAVFS2 must be installed on the physical device.

**Edge Server
1. Create a PVC with the Create Volume function.
2. Run the Mount Volume function with the created volume ID. 

**On a physical device
1. Find the point to mount the PV.
    ex. /mount
2. Mount it with the webdev server.
    1. execute the mount command.
    ```
        sudo mount -t davfs http://<master node address>:<port number>/uploads <point to mount>
    ```
    2. Verify that it mounted.
        If it passes without mount failed, it is mounted successfully.