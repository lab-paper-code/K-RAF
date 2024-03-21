<script>
  import request from "../../lib/request";

  let authInfo = {
    id: "",
    password: "",
  };

  let volume_id = ""; 

  const unmountVolume = async () => {
    const url = `/mounts/${volume_id}`;

    try {
      const response = await request("DELETE", url, {}, {}, authInfo.id, authInfo.password);

      const parsedResponse = typeof response === 'string' ? JSON.parse(response) : response;
      console.log("parsed server response:", parsedResponse);

      alert(`Successfully unmounted volume with ID: ${volume_id}`);
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
      style="margin-right: 10px;"
      id="id"
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

    <label for="volume_id">VolumeID:</label>
    <input
      type="text"
      id="volume_id"
      style="margin-right: 10px;"
      bind:value={volume_id}
      placeholder="Enter Volume ID to Unmount"
    />
    <button on:click={unmountVolume}>Unmount Volume</button>
  </form>
</div>
