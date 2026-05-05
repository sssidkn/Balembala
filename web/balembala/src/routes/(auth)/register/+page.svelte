<script lang="ts">
    import Logo from "$lib/components/Logo.svelte"
    import { enhance } from "$app/forms"
    import type { PageProps } from './$types';

	let { form }: PageProps = $props();
</script>

<div class="vstack center">
    <Logo/>
    <div class=card>
        <h1 style:text-align=center>New account</h1>
        <form method="POST" use:enhance>
            <label class=text-form for="username">
                Username
            </label>
            <input
                class=text-form
                name="username"
                required
                autocomplete="off"
                value={form?.username ?? ""}
            >
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
                <small style:color=var(--error)>{form?.error}</small>
                <button type=submit>Create</button>
            </div>
        </form>
    </div>
</div>

<style>
    .center {
        position: absolute;
        top: 0;
        bottom: 0;
        left: 0;
        right: 0;
        transform: translateY(-5%);
    
        margin: auto;
    }
    .vstack {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        row-gap: 5pt;
    }
    button {
        margin-top: 10pt;
    }
    input {
        margin-bottom: 5pt;
    }
</style>
