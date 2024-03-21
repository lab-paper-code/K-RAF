<script>
    import { link } from "svelte-spa-router";
</script>

<nav class="navbar fixed-top navbar-expand-sm navbar-light bg-light">
    <div class="container">
        <a use:link href="/" class="navbar-brand mb-0 h1">
            Edge Device Registration System(User)</a
        >
        <button
            type="button"
            data-bs-toggle="collapse"
            data-bs-target="#navbarNav"
            class="navbar-toggler"
            aria-controls="navbarNav"
            aria-expanded="false"
            aria-label="Toggle navigation"
        >
            <span class="navbar-toggler-icon" />
        </button>
        <div class="collapse navbar-collapse" id="navbarNav">
            <ul class="navbar-nav">
                <li class="nav-item dropdown">
                    <a
                        class="nav-link dropdown-toggle"
                        id="navbarDropdown"
                        role="button"
                        data-bs-toggle="dropdown"
                        aria-expanded="false"
                    >
                        Apps
                    </a>
                    <ul class="dropdown-menu" aria-labelledby="navbarDropdown">
                        <li>
                            <a
                                use:link
                                href="/apps/detail"
                                class="dropdown-item">Search App</a
                            >
                        </li>
                    </ul>
                </li>
                <li class="nav-item dropdown">
                    <a
                        class="nav-link dropdown-toggle"
                        id="navbarDropdown"
                        role="button"
                        data-bs-toggle="dropdown"
                        aria-expanded="false"
                    >
                        Appruns
                    </a>
                    <ul class="dropdown-menu" aria-labelledby="navbarDropdown">
                        <li>
                            <a use:link href="/appruns" class="dropdown-item"
                                >Execute AppRun</a
                            >
                        </li>
                        <li>
                            <a
                                use:link
                                href="/appruns/terminate"
                                class="dropdown-item">Terminate AppRun</a
                            >
                        </li>
                    </ul>
                </li>
                <li class="nav-item dropdown">
                    <a
                        class="nav-link dropdown-toggle"
                        id="navbarDropdown"
                        role="button"
                        data-bs-toggle="dropdown"
                        aria-expanded="false"
                    >
                        Devices
                    </a>
                    <ul class="dropdown-menu" aria-labelledby="navbarDropdown">
                        <li>
                            <a
                                use:link
                                href="/devices/update"
                                class="dropdown-item">Update Device</a
                            >
                        </li>
                    </ul>
                </li>
                <li class="nav-item dropdown">
                    <a
                        class="nav-link dropdown-toggle"
                        id="navbarDropdown"
                        role="button"
                        data-bs-toggle="dropdown"
                        aria-expanded="false"
                    >
                        Volumes
                    </a>
                    <ul class="dropdown-menu" aria-labelledby="navbarDropdown">
                        <li>
                            <a use:link href="/volumes" class="dropdown-item"
                                >Create Volume</a
                            >
                        </li>
                        <li>
                            <a
                                use:link
                                href="/volumes/detail"
                                class="dropdown-item">Search Volume</a
                            >
                        </li>
                        <li>
                            <a
                                use:link
                                href="/volumes/update"
                                class="dropdown-item">Update Volume</a
                            >
                        </li>
                        <li>
                            <a
                                use:link
                                href="/volumes/mount"
                                class="dropdown-item">Mount Volume</a
                            >
                        </li>
                        <li>
                            <a
                                use:link
                                href="/volumes/unmount"
                                class="dropdown-item">Unmount Volume</a
                            >
                        </li>
                    </ul>
                </li>
            </ul>
        </div>
    </div>
</nav>
