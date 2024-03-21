<script>
    import request from "../../lib/request";

    let authInfo = {
        id: "",
        password: "",
    };

    let volume = null;
    let volumes = [];
    let volumeID = "";
    let isMulti = false;

    const clearOutput = () => {
        volume = null;
        volumes = [];
    };

    const convertBytesToGB = (bytes) => {
        return Math.floor(bytes / 1_073_741_824);
    };

    const formatDate = (isoDate) => {
        const date = new Date(isoDate);
        return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}:${String(date.getSeconds()).padStart(2, '0')}`;
    };

    const getVolumes = async () => {
        try {
            clearOutput();

            const headersForGet = {
                dev_id: authInfo.id,
                dev_password: authInfo.password,
            };

            let url = isMulti ? "/volumes" : `/volumes/${volumeID}`;
            const response = await request("GET", url, {}, {}, authInfo.id, authInfo.password);
            
            const parsedResponse = typeof response === 'string' ? JSON.parse(response) : response;
            console.log("parsed server response:", parsedResponse);

            if (isMulti) {
                if (response && response.volumes) {
                    volumes = [...response.volumes].map(vol => ({
                        ...vol,
                        volume_size: convertBytesToGB(vol.volume_size),
                        created_at: formatDate(vol.created_at),
                        updated_at: formatDate(vol.updated_at)
                    }));
                }
            } else {
                if (response) {
                    volume = {
                        ...response,
                        volume_size: convertBytesToGB(response.volume_size),
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

<div class="input container mx-auto" style="margin-top: 60px; .align-center">
    <form>
        <label for="id">ID:</label>
        <input
            type="text"
            id="id"
            style="margin-right: 10px;"
            bind:value={authInfo.id}
            placeholder="Enter ID"
        />

        <label for="password">Password:</label>
        <input
            type="password"
            id="password"
            style="margin-right: 10px;"
            bind:value={authInfo.password}
            placeholder="Enter Password"
        />
        {#if !isMulti}
            <label for="volumeID">VolumeID:</label>
            <input
                type="text"
                id="volumeID"
                style="margin-right: 10px;"
                bind:value={volumeID}
                placeholder="Enter VolumeID"
            />
        {/if}

        <label>
            <input
                type="checkbox"
                bind:checked={isMulti}
                on:change={clearOutput}
            /> Multi-volume
        </label>
        <button style="margin-left: 10px;" on:click={getVolumes}
            >Search Volume</button
        >
    </form>
</div>

<div class="display container mx-auto" style="margin-top: 20px; margin-left: 10px;">
    {#if isMulti}
        {#if volumes.length > 0}
            <div class="volume-container">
                <strong style="margin-left: 10px">Volume Info:</strong>
                {#each volumes as volume}
                    <div class="volume">
                        <p>ID: {volume.id}</p>
                        <p>Device ID: {volume.device_id}</p>
                        <p>Volume Size: {volume.volume_size}GB</p>
                        <p>Mounted: {volume.Mounted ? "Yes" : "No"}</p>
                        <p>Created At: {volume.created_at}</p>
                        <p>Updated At: {volume.updated_at}</p>
                    </div>
                {/each}
            </div>
        {/if}
    {:else if volume}
        <div class="volume-container">
            <strong style="margin-left: 10px">Volume Info:</strong>
            <div class="volume">
                <p>ID: {volume.id}</p>
                <p>Device ID: {volume.device_id}</p>
                <p>Volume Size: {volume.volume_size}GB</p>
                <p>Mounted: {volume.Mounted ? "Yes" : "No"}</p>
                <p>Created At: {volume.created_at}</p>
                <p>Updated At: {volume.updated_at}</p>
            </div>
        </div>
    {/if}
</div>

<style>
    .volume-container {
        border: 2px solid #323131;
    }

    .volume {
        border-bottom: 1px solid #323131;
        padding: 10px;
    }

    .volume:last-child {
        border-bottom: none;
    }
</style>
