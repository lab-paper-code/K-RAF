<script>
  import request from "../../lib/request";

  let authInfo = {
    id: "",
    password: "",
  };

  let appRunExecute = {
    device_id: "",
    volume_id: "",
  };

  let app_id = "";
  let apprun_id = null;
  let app_runs = [];

  const executeAppRun = async () => {
    const url1 = `/appruns/${app_id}`;
    const url2 = "/appruns";

    try {
      await request("POST", url1, appRunExecute, {}, authInfo.id, authInfo.password);

      const response = await request("GET", url2, {}, {}, authInfo.id, authInfo.password);
      console.log("Received from server:", response);

      const parsedResponse = typeof response === 'string' ? JSON.parse(response) : response;
      console.log("parsed server response:", parsedResponse);

      if (response && response.app_runs) {
        app_runs = [...response.app_runs];
        apprun_id = app_runs[app_runs.length - 1].id;
        console.log("apprun_id", apprun_id);
      }
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
      style=" margin-right: 10px;"
      bind:value={authInfo.id}
      placeholder="Enter ID"
    />

    <label for="password">Password:</label>
    <input
      type="password"
      id="password"
      style=" margin-right: 10px;"
      bind:value={authInfo.password}
      placeholder="Enter Password"
    />

    <label for="device_id">Device ID:</label>
    <input
      type="text"
      id="device_id"
      style=" margin-right: 10px;"
      bind:value={appRunExecute.device_id}
      placeholder="Enter Device ID"
    />

    <label for="volume_id">Volume ID:</label>
    <input
      type="text"
      id="volume_id"
      style=" margin-right: 10px;"
      bind:value={appRunExecute.volume_id}
      placeholder="Enter Volume ID"
    />
</form>
<form> <!-- JS: may modify layout later -->
    <label for="app_id">App ID:</label>
    <input
      type="text"
      id="app_id"
      style="margin-top: 5px;margin-right: 10px;"
      bind:value={app_id}
      placeholder="Enter App ID"
    />
    <button on:click={executeAppRun}>Execute AppRun</button>
  </form>
</div>

<div class="output container mx-auto" style="margin-top: 20px">
{#if apprun_id}
  <div>
    <strong>Received AppRuns ID:</strong>
    {apprun_id}
  </div>
{/if}
</div>