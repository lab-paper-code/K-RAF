<script>
    import request from "../../lib/request";

    let authInfo = {
        id: "",
        password: "",
    };

    let appRun = null;
    let appRuns = [];
    let appRunID = "";
    let isMulti = false;

    const clearOutput = () => {
        appRun = null;
        appRuns = [];
    };

    const formatDate = (isoDate) => {
        const date = new Date(isoDate);
        return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}:${String(date.getSeconds()).padStart(2, '0')}`;
    };
    
    const getAppRuns = async () => {
        try {
            clearOutput();

            let url = isMulti ? "/appruns" : `/appruns/${appRunID}`;
            const response = await request("GET", url, {}, {}, authInfo.id, authInfo.password);

            const parsedResponse = typeof response === 'string' ? JSON.parse(response) : response;
            console.log("parsed server response:", parsedResponse);

            if (isMulti) {
                if (response && response.app_runs) {
                    appRuns = [...response.app_runs].map(dev => ({
                        ...dev,
                        created_at: formatDate(dev.created_at),
                        updated_at: formatDate(dev.updated_at)
                    }));
                }
            } else {
                if (response) {
                    appRun = {
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
    
    {#if !isMulti}
        <label for="appRunID">AppRun ID:</label>
        <input
            type="text"
            id="appRunID"
            style=" margin-right: 10px;"
            bind:value={appRunID}
            placeholder="Enter AppRun ID"
        />
    {/if}

    <label>
        <input type="checkbox" bind:checked={isMulti} on:change={clearOutput} />
        Multi-appRun
    </label>

    <button style=" margin-left: 10px;" on:click={getAppRuns}>Search AppRun</button>
</div>

<div class="output container mx-auto" style="margin-top: 20px; margin-left: 10px;">
{#if isMulti}
    {#if appRuns.length > 0}
        <div class="apprun-container">
            <strong style="margin-left: 10px;">AppRun Info:</strong>
            {#each appRuns as appRun}
                <div class="apprun">
                    <!-- Displaying fields according to AppRun format -->
                    {#each Object.keys(appRun) as key}
                        <p>{key}: {appRun[key]}</p>
                    {/each}
                </div>
            {/each}
        </div>
    {/if}
{:else if appRun}
    <div class="apprun-container">
        <strong style="margin-left: 10px;">AppRun Info:</strong>
        <div class="apprun">
            <!-- Displaying fields according to AppRun format -->
            {#each Object.keys(appRun) as key}
                <p>{key}: {appRun[key]}</p>
            {/each}
        </div>
    </div>
{/if}
</div>

<style>
    .apprun-container {
        border: 2px solid #323131;
    }

    .apprun {
        border-bottom: 1px solid #323131;
        padding: 10px;
    }

    .apprun:last-child {
        border-bottom: none;
    }
</style>
