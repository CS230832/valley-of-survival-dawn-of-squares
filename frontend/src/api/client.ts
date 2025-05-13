import { getSessionID, SESSION_KEY } from "../utils/session"

export async function apiFetch(url: string, options: RequestInit = {}): Promise<Response> {
    const headers = new Headers(options.headers)
    const session = getSessionID()

    if (session) {
        headers.set(SESSION_KEY, session)
    }

    return await fetch(url, { ...options, headers })
}