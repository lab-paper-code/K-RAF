<script>
    import request from "../../lib/request";

    let authInfo = {
        id: "",
        password: "",
    };

    let deviceUpdate = {
        mod_ip: "",
        mod_password: "",
    };

    

    const updateDevice = async () => {
        try {
            const url = "/devices/" + authInfo.id;
            const device_id = authInfo.id;
            const response = await request("PATCH", url, deviceUpdate, {}, authInfo.id, authInfo.password); // capture response

            alert(`Successfully updated Device with ID: ${device_id}`);
        } catch (error) {
            alert(`Error: ${error.message}`);
        }
    };
</script>

<div class="input container mx-auto" style="margin-top: 60px">
    <form>
        <label for="id">ID:</label>
        <input
            type="text"
            id="id"
            style="margin-bottom: 5px; margin-right: 10px;"
            bind:value={authInfo.id}
            placeholder="Enter ID"
        />

        <label for="password">Password:</label>
        <input
            type="password"
            id="password"
            bind:value={authInfo.password}
            placeholder="Enter Password"
        />

        <label for="mod_ip">IP:</label>
        <input
            type="text"
            id="mod_ip"
            style="margin-right: 10px;"
            bind:value={deviceUpdate.mod_ip}
            placeholder="Enter IP to Change"
        />

        <label for="mod_password">Device Password:</label>
        <input
            type="password"
            id="mod_password"
            style="margin-right: 10px;"
            bind:value={deviceUpdate.mod_password}
            placeholder="Enter Device Password to Change"
        />
        <button on:click={updateDevice}>Update Device</button>
    </form>
</div>
