<script>
  import request from "../../lib/request";

  let authInfo = {
    id: "",
    password: "",
  };

  let apprun_id = ""; // This will hold the ID of the AppRun you want to delete

  const terminateAppRun = async () => {
    const url = `/appruns/${apprun_id}`;

    try {
      const response = await request("DELETE", url, {}, {}, authInfo.id, authInfo.password);
      // console log to check server response
      const parsedResponse = typeof response === 'string' ? JSON.parse(response) : response;
      console.log("parsed server response:", parsedResponse);

      alert(`Successfully deleted AppRun with ID: ${apprun_id}`);
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

    <label for="apprun_id">AppRun ID:</label>
    <input
      type="text"
      id="apprun_id"
      style=" margin-right: 10px;"
      bind:value={apprun_id}
      placeholder="Enter AppRun ID"
    />
    <button on:click={terminateAppRun}>Terminate AppRun</button>
  </form>
</div>
