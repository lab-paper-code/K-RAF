<script>
    import request from "../../lib/request";

    let authInfo = {
        id: "",
        password: "",
    };

    let app = null;
    let apps = [];
    let appID = "";
    let isMulti = false;

    const clearOutput = () => {
        app = null;
        apps = [];
    };

    const formatDate = (isoDate) => {
        const date = new Date(isoDate);
        return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}:${String(date.getSeconds()).padStart(2, '0')}`;
    };
    
    const getApps = async () => {
        try {
            clearOutput();

            let url = isMulti ? "/apps" : `/apps/${appID}`;
            const response = await request("GET", url, {}, {}, authInfo.id, authInfo.password);
    
            const parsedResponse = typeof response === 'string' ? JSON.parse(response) : response;
            console.log("parsed server response:", parsedResponse);

            if (isMulti) {
                if (response && response.apps) {
                    apps = [...response.apps].map(dev => ({
                        ...dev,
                        created_at: formatDate(dev.created_at),
                        updated_at: formatDate(dev.updated_at)
                    }));
                    console.log("app list length:", apps.length);
                }
            } else {
                if (response) {
                    app = {
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
    <form>
        <label for="id">ID:</label>
        <input
            type="text"
            id="id"
            style="margin-bottom: 5pt; margin-right: 10px;"
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

        <!--App Search format-->
        {#if !isMulti}
            <label for="appID">AppID:</label>
            <input
                type="text"
                id="appID"
                style=" margin-right: 10px;"
                bind:value={appID}
                placeholder="Enter AppID"
            />
        {/if}

        <label>
            <input
                type="checkbox"
                
                bind:checked={isMulti}
                on:change={clearOutput}
            />
            Multi-app
        </label>
        <button style=" margin-left: 10px;" on:click={getApps}>Search App</button>
    </form>
</div>

<div class="display container mx-auto" style="margin-top: 20px; margin-left: 10px;">
    {#if isMulti}
        {#if apps.length > 0}
            <div class="app-container">
                <strong style="margin-left: 10px;" >App Info:</strong>
                {#each apps as app}
                    <div class="app">
                        <p>ID: {app.id}</p>
                        <p>Name: {app.name}</p>
                        <p>Require GPU: {app.require_gpu ? "Yes" : "No"}</p>
                        <p>Description: {app.description || "N/A"}</p>
                        <p>Docker Image: {app.docker_image}</p>
                        <p>Commands: {app.commands}</p>
                        <p>Arguments: {app.arguments}</p>
                        <p>Open Ports: {app.open_ports}</p>
                        <p>Created At: {app.created_at}</p>
                        <p>Updated At: {app.updated_at}</p>
                    </div>
                {/each}
            </div>
        {/if}
    {:else if app}
        <div class="app-container">
            <strong style="margin-left: 10px;">App Info:</strong>
            <div class="app">
                <p>ID: {app.id}</p>
                <p>Name: {app.name}</p>
                <p>Require GPU: {app.require_gpu ? "Yes" : "No"}</p>
                <p>Description: {app.description || "N/A"}</p>
                <p>Docker Image: {app.docker_image}</p>
                <p>Commands: {app.commands}</p>
                <p>Arguments: {app.arguments}</p>
                <p>Open Ports: {app.open_ports}</p>
                <p>Created At: {app.created_at}</p>
                <p>Updated At: {app.updated_at}</p>
            </div>
        </div>
    {/if}
</div>

<style>
    .app-container {
        border: 2px solid #323131;
    }

    .app {
        border-bottom: 1px solid #323131;
        padding: 10px;
    }

    .app:last-child {
        border-bottom: none;
    }
</style>
