<script lang="ts">
    import Logo from "$lib/components/Logo.svelte"
    import Login from "$lib/assets/icons/Login.svelte"
    
    import { resolve } from "$app/paths";
    import { enhance } from "$app/forms";
    
    import type { PageProps } from './$types';

	let { form }: PageProps = $props();
</script>

<div class=vstack style:margin-top=5em>
    <Logo/>
    <div class=card>
        <h1 style:text-align=center>Welcome</h1>
        <form method="POST" use:enhance>
            <label class=text-form for="email">
                Email
            </label>
            <input
                class=text-form
                type=email
                name="email"
                required
                autocomplete="off"
                value={form?.email ?? ""}
            >
            <label class=text-form for="password">
                Password
            </label>
            <input
                class=text-form
                type=password
                name="password"
                required
                autocomplete="off"
            >
            <div class=vstack>
                <small class=error>{form?.error ?? ""}</small>
                <button type=submit>
                    <div class=hstack>
                        <Login --size=16pt></Login>
                        Log in
                    </div>
                </button>
                <a href={resolve("/register")}>Create an account</a>
            </div>
        </form>
    </div>
</div>

<style>
    .vstack {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        row-gap: 5pt;
    }
    .hstack {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: center;
        column-gap: 5pt;
    }
    .error {
        color: var(--error);
        min-height: 1.2em;
    }
    input {
        margin-bottom: 5pt;
    }
</style>
