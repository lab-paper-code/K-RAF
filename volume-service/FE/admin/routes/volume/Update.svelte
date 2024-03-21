<script>
    import request from "../../lib/request";

    let authInfo = {
        id: "",
        password: "",
    };

    let volume_id = "";
    let volumeUpdate = {
        volume_size: "",
    } ;

    const updateVolume = async () => {
        try {
            const url = `/volumes/${volume_id}`;

            const response = await request("PATCH", url, volumeUpdate, {}, authInfo.id, authInfo.password); // capture response
            console.log("Received from server:", response); // Debugging line
            alert(`Successfully updated Volume with ID: ${volume_id}`);
        } catch (error) {
            alert(`Error: ${error.message}`);
        }
    };
</script>

<div class="container" style="margin-top: 60px">
    <!-- Authentication form -->
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
        
        <!-- Volume update form -->
        <label for="volume_id">VolumeID:</label>
        <input
            type="text"
            id="volume_id"
            style="margin-right: 10px;"
            bind:value={volume_id}
            placeholder="Enter Volume ID to Update"
        />

        <label for="mod_volume_size">Volume Size:</label>
        <input
            type="text"
            id="mod_volume_size"
            style="margin-right: 10px;"
            bind:value={volumeUpdate.volume_size}
            placeholder="Enter New Volume Size"
        />
        <button on:click={updateVolume}>Update Volume</button>
    </form>
</div>
