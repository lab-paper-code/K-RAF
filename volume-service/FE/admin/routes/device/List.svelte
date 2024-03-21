<script>
    import request from "../../lib/request";

    let authInfo = {
        id: "",
        password: "",
    };

    let device = null;
    let devices = [];
    let deviceID = "";
    let isMulti = false;

    const clearOutput = () => {
        device = null;
        devices = [];
    };

    const formatDate = (isoDate) => {
        const date = new Date(isoDate);
        return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}:${String(date.getSeconds()).padStart(2, '0')}`;
    };
    
    const getDevices = async () => {
        try {
            clearOutput();

            let url = isMulti ? "/devices" : `/devices/${deviceID}`;
            const response = await request("GET", url, {}, {}, authInfo.id, authInfo.password);

            const parsedResponse = typeof response === 'string' ? JSON.parse(response) : response;
            console.log("parsed server response:", parsedResponse);

            if (isMulti) {
                if (response && response.devices) {
                    devices = [...response.devices].map(dev => ({
                        ...dev,
                        created_at: formatDate(dev.created_at),
                        updated_at: formatDate(dev.updated_at)
                    }));
                }
            } else {
                if (response) {
                    device = {
                        ...response,
                        created_at: formatDate(response.created_at),
                        updated_at: formatDate(response.updated_at)
                    };
                }
            }
        } catch (error) {
            alert(`Error: ${error.message}`);
        }
    };
</script>
<div class="input container mx-auto" style="margin-top: 60px">
    <label for="id">ID:</label>
    <input type="text" id="id" style="margin-right: 10pt; margin-bottom: 5pt;" bind:value={authInfo.id} placeholder="Enter ID" />

    <label for="password">Password:</label>
    <input type="password" id="password" bind:value={authInfo.password} placeholder="Enter Password" />

<form>
    {#if !isMulti}
        <label for="deviceID">DeviceID:</label>
        <input type="text" id="deviceID" style="margin-right: 10pt;" bind:value={deviceID} placeholder="Enter DeviceID" />
    {/if}

    <label>
        <input type="checkbox" bind:checked={isMulti} on:change={clearOutput} /> Multi-device
    </label>

    <button style="margin-left: 10pt;" on:click={getDevices}>Search Device</button>
</form>

{#if isMulti}
    {#if devices.length > 0}
        <div class="device-container" style="margin-top: 20px">
            <strong style="margin-left: 10px;">Device Info:</strong>
            {#each devices as device}
                <div class="device">
                    <p>ID: {device.id}</p>
                    <p>IP: {device.ip}</p>
                    <p>Description: {device.description || "N/A"}</p>
                    <p>Password: {device.password}</p>
                    <p>Created At: {device.created_at}</p>
                    <p>Updated At: {device.updated_at}</p>
                </div>
            {/each}
        </div>
    {/if}
{:else if device}
    <div class="device-container" style="margin-top: 20px">
        <strong style="margin-left: 10px;">Device Info:</strong>
        <div class="device">
            <p>ID: {device.id}</p>
            <p>IP: {device.ip}</p>
            <p>Description: {device.description || "N/A"}</p>
            <p>Password: {device.password}</p>
            <p>Created At: {device.created_at}</p>
            <p>Updated At: {device.updated_at}</p>
        </div>
    </div>
{/if}
</div>

<style>
    .device-container {
        border: 2px solid #323131;
    }

    .device {
        border-bottom: 1px solid #323131;
        padding: 10px;
    }

    .device:last-child {
        border-bottom: none;
    }
</style>
