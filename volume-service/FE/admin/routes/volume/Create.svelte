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
  
  let volumes = [];
  let volume_id = null;

  const createVolume = async () => {
    try {
      const url = "/volumes";
      await request("POST", url, volumeCreate, {}, authInfo.id, authInfo.password);


      const response = await request("GET", url, {}, {}, authInfo.id, authInfo.password);
      console.log("Received from server:", response);

      const parsedResponse = typeof response === 'string' ? JSON.parse(response) : response;
      console.log("parsed server response:", parsedResponse);
      
      if (response && response.volumes) {
        volumes = [...response.volumes];
        volume_id = volumes[volumes.length - 1].id;
        console.log("volume-id", volume_id);
      }
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

    <label for="device_id">DeviceID:</label>
    <input
      type="text"
      id="device_id"
      style="margin-right: 10px;"
      bind:value={volumeCreate.device_id}
      placeholder="Enter DeviceID for Volume"
    />

    <label for="volume_size">Volume Size:</label>
    <input
      type="text"
      id="volume_size"
      style="margin-right: 10px;"
      bind:value={volumeCreate.volume_size}
      placeholder="Enter Volume Size"
    />
    <button on:click={createVolume}>Create Volume</button>
  </form>
</div>

<div class="output(volume_id) container" style="margin-top: 20px">
  {#if volume_id}
    <div>
      <strong>Your Volume ID:</strong>
      {volume_id}
    </div>
  {/if}
</div>
