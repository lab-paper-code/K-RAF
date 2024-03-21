<script>
  import Router from 'svelte-spa-router'
  import Home from "./routes/Home.svelte"
  import Navigation from "./components/Navigation.svelte"

  import registerApp from "./routes/app/Register.svelte"
  import getApps from "./routes/app/List.svelte"

  import executeAppRun from "./routes/apprun/Execute.svelte"
  import getAppRuns from "./routes/apprun/List.svelte"
  import terminateAppRuns from "./routes/apprun/Terminate.svelte"

  import createVolume from "./routes/volume/Create.svelte"
  import getVolumes from "./routes/volume/List.svelte"
  import updateVolume from "./routes/volume/Update.svelte"
  import deleteVolume from "./routes/volume/Delete.svelte"
  import mountVolume from "./routes/volume/Mount.svelte"
  import unmountVolume from "./routes/volume/Unmount.svelte"

  import registerDevice from "./routes/device/Register.svelte"
  import getDevices from "./routes/device/List.svelte"
  import updateDevice from "./routes/device/Update.svelte"
  
  const routes = {
    '/': Home,
    
    '/apps': registerApp,
    '/apps/detail' : getApps,

    '/appruns': executeAppRun,
    '/appruns/detail' : getAppRuns,
    '/appruns/terminate' : terminateAppRuns,

    '/volumes' : createVolume,
    '/volumes/detail' : getVolumes,
    '/volumes/update' : updateVolume,
    '/volumes/delete' : deleteVolume,
    '/volumes/mount' : mountVolume,
    '/volumes/unmount' : unmountVolume,

    '/devices': registerDevice,
    '/devices/detail' : getDevices,
    '/devices/update': updateDevice,
    
  }
</script>

<Navigation />
<Router {routes}/>

