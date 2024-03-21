<script>
    import request from "../../lib/request";

    let authInfo = {
        id: "",
        password: "",
    };

    let volumeCreate = {
        device_id: "",
        volume_size: "",
    };

    let volume_id = null;

    const mountVolume = async () => {
        try {
            const url = `/mounts/${volume_id}`;
            const response = await request("POST", url, {}, {}, authInfo.id, authInfo.password);
            
            const parsedResponse = typeof response === 'string' ? JSON.parse(response) : response;
            console.log("parsed server response:", parsedResponse);

            alert(`"Successfully mounted Volume with ID: ${volume_id}"`);
        } catch (error) {
            alert(`Error: ${error.message}`);
        }
    };
</script>

<div class="container" style="margin-top: 60px">
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

        <!-- Mounting Volume form -->
        <label for="volume_id">VolumeID:</label>
        <input
            type="text"
            id="volume_id"
            style="margin-right: 10px;"
            bind:value={volume_id}
            placeholder="Enter VolumeID to Mount"
        />
        <button on:click={mountVolume}>Mount Volume</button>
    </form>
</div>
