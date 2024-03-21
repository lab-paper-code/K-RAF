<script>
    import request from "../../lib/request";

    let authInfo = {
        id: "",
        password: "",
    };

    let appRegister = {
        name: "",
        require_gpu: false,
        description: "",
        docker_image: "",
        arguments: "",
        open_ports: [],
    };

    let app_id = null;
    let apps = [];

    const registerApp = async () => {
        try {
            // Convert the comma-separated string to an array of numbers
            appRegister.open_ports = appRegister.open_ports.split(",").map(Number); 
        
            const url = "/apps";
            await request("POST", url, appRegister, {}, authInfo.id, authInfo.password);
            
            const response = await request("GET", url, {}, {}, authInfo.id, authInfo.password);
            console.log("Received from server:", response);

            const parsedResponse = typeof response === 'string' ? JSON.parse(response) : response;
            console.log("parsed server response:", parsedResponse);

            if (response && response.apps) {
                apps = [...response.apps];
                app_id = apps[apps.length - 1].id;
                console.log("app-id", app_id);
            }
        } catch (error) {
            alert(`Error: ${error.message}`);
        }
    };
</script>


<div class="input container" style="margin-top: 60px;">
    <form>
        <label for="id">ID:</label>
        <input
            type="password"
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

        <label for="name">Name:</label>
        <input
            type="text"
            id="name"
            style=" margin-right: 10px;"
            bind:value={appRegister.name}
            placeholder="Enter App Name"
        />
        </form>

        <form>
            <label for="description">Description:</label>
        <input
            type="text"
            id="description"
            style="  margin-top: 5px;margin-right: 10px;"
            bind:value={appRegister.description}
            placeholder="Enter Description"
        />

        <label for="docker_image">Docker Image:</label>
        <input
            type="text"
            id="docker_image"
            style="margin-right: 10px;"
            bind:value={appRegister.docker_image}
            placeholder="Enter Docker Image"
        />
    
        <label for="arguments">Arguments:</label>
        <input
            type="text"
            id="arguments"
            style=" margin-right: 10px;"
            bind:value={appRegister.arguments}
            placeholder="Enter Arguments"
        />
        </form>
        <form>
        <label for="open_ports">Open Ports:</label>
        <input
            type="text"
            id="open_ports"
            style=" margin-top: 5px; margin-right: 10px;"
            bind:value={appRegister.open_ports}
            placeholder="Enter Open Ports"
        />

        <label for="require_gpu">Require GPU:</label>
        <input
            type="checkbox"
            id="require_gpu"
            style=" margin-right: 10px;"
            bind:checked={appRegister.require_gpu}
        />
        <button on:click={registerApp}>Register App</button>
    </form>
</div>

<div class="output container mx-auto" style="margin-top: 20px">
    {#if app_id}
        <div>
            <strong>Your App ID:</strong>
            {app_id}
        </div>
    {/if}
</div>



