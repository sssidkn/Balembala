import { fail } from '@sveltejs/kit';
import type { Actions } from './$types';
export const actions = {
    default: async ({request}) => {
        const data = await request.formData();
        const username = data.get("username");
        const email = data.get("email");
        const password = data.get("password");
        if (username === "error") {
            return fail(400, {username, email, error: "Generic error"})
        }
        console.log({username, email, password})
        return { success: true }
    }
} satisfies Actions;