import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
    const response = await fetch(import.meta.env.VITE_SOCKET_URL + '/api/v1/connections');
    const connections = await response.json();

    return {
        connections
    };
};
