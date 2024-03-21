<script>
  import request from "../../lib/request";

  let authInfo = {
    id: "",
    password: "",
  };

  let deviceRegister = {
    ip: "",
    password: "",
    description: "",
  };

  let device_id = null;
  let devices = [];

  const registerDevice = async () => {
    try {
      const url = "/devices";
      await request("POST", url, deviceRegister, {}, authInfo.id, authInfo.password);

      const response = await request("GET", url, {}, {}, authInfo.id, authInfo.password);
      console.log("Received from server:", response);

      const parsedResponse = typeof response === 'string' ? JSON.parse(response) : response;
      console.log("parsed server response:", parsedResponse);

      if (response && response.devices) {
        devices = [...response.devices];
        device_id = devices[devices.length - 1].id;
        console.log("device-id", device_id);
      }
    } catch (error) {
      alert(`Error: ${error.message}`);
    }
  };
</script>

<div class="input container mx-auto" style="margin-top: 60px">
  <form>
    <label for="id" style="margin-bottom: 5pt;">ID:</label>
    <input
      type="password"
      id="id"
      style="margin-right: 10pt;"
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

    <label for="device_id">IP:</label>
    <input
      type="text"
      id="ip"
      style="margin-right: 10px;"
      bind:value={deviceRegister.ip}
      placeholder="Enter IP"
    />

    <label for="volume_id">Device Password:</label>
    <input
      type="password"
      id="password"
      style="margin-right: 10px;"
      bind:value={deviceRegister.password}
      placeholder="Enter Device Password"
    />
    
    <label for="app_id">Description:</label>
    <input
      type="text"
      id="description"
      style="margin-right: 10px;"
      bind:value={deviceRegister.description}
      placeholder="Enter Description"
    />
    <button on:click={registerDevice}>Register Device</button>
  </form>
</div>

<div class="output container mx-auto" style="margin-top: 20px">
{#if device_id}
  <div>
    <strong>Your Device ID:</strong>
    {device_id}
  </div>
{/if}
</div>
