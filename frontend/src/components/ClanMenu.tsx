import { createSignal, Show } from "solid-js";
import { apiFetch } from "../api/client";
import type { Clan, User } from "../utils/types";

interface Props {
    onReturn: () => void;
    onLogout: () => void;
}

const ClanMenu = (props: Props) => {
    const [currentUser, setCurrentUser] = createSignal<User | null>(null);
    const [currentClan, setCurrentClan] = createSignal<Clan | null>(null);
    const [error, setError] = createSignal<string | null>(null);
    const [clanName, setClanName] = createSignal<string>("");
    const [clanPassword, setClanPassword] = createSignal<string>("");

    const fetchCurrentUser = async () => {
        try {
            const res = await apiFetch("/api/info/current_user", {
                method: "GET",
            });
            if (!res.ok) throw new Error("Failed to fetch current user");
            const data = await res.json();
            setCurrentUser(data);
        } catch (err) {
            setError("Failed to fetch current user info");
            console.error(err);
        }
    };

    const fetchCurrentClan = async () => {
        const res = await apiFetch("/api/info/current_clan", {
            method: "GET",
        });

        if (res.ok) {
            const data = await res.json();
            setCurrentClan(data);
            return;
        }

        setCurrentClan(null);

        if (res.status === 417) {
            setError("");
        } else {
            setError("Failed to fetch current clan info");
            console.error(res);
        }
    };

    fetchCurrentUser();
    fetchCurrentClan();

    const handleCreateClan = async () => {
        try {
            const res = await apiFetch("/api/clan/create", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    name: clanName(),
                    password: clanPassword(),
                }),
            });
            if (!res.ok) throw new Error("Failed to create clan");
            setClanName("");
            setClanPassword("");
            fetchCurrentClan();
        } catch (err) {
            setError("Failed to create clan");
            console.error(err);
        }
    };

    const handleDeleteClan = async () => {
        try {
            const res = await apiFetch("/api/clan/delete", {
                method: "DELETE",
            });
            if (!res.ok) throw new Error("Failed to delete clan");
            fetchCurrentUser();
            fetchCurrentClan();
        } catch (err) {
            setError("Failed to delete clan");
            console.error(err);
        }
    };

    const handleJoinClan = async () => {
        try {
            const res = await apiFetch("/api/clan/join", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    name: clanName(),
                    password: clanPassword(),
                }),
            });
            if (!res.ok) throw new Error("Failed to join clan");
            setClanName("");
            setClanPassword("");
            fetchCurrentUser();
            fetchCurrentClan();
        } catch (err) {
            setError("Failed to join clan");
            console.error(err);
        }
    };

    const handleLeaveClan = async () => {
        try {
            const res = await apiFetch("/api/clan/leave", {
                method: "POST",
            });
            if (!res.ok) throw new Error("Failed to leave clan");
            fetchCurrentUser();
            fetchCurrentClan();
        } catch (err) {
            setError("Failed to leave clan");
            console.error(err);
        }
    };

    return (
        <div class="clan-menu">
            <h2>Clan Menu</h2>
            {error() && <p class="error">{error()}</p>}

            {/* Display user and clan info */}
            <p>
                <strong>User:</strong>{" "}
                {currentUser() ? currentUser()!.username : "Loading..."}
            </p>
            <p>
                <strong>Clan:</strong>{" "}
                {currentClan() ? currentClan()!.name : "None"}
            </p>

            {/* If the user is not in a clan, show create/join options */}
            {!currentClan() ? (
                <>
                    <input
                        type="text"
                        placeholder="Clan Name"
                        value={clanName()}
                        onInput={(e) => setClanName(e.currentTarget.value)}
                    />
                    <input
                        type="password"
                        placeholder="Clan Password"
                        value={clanPassword()}
                        onInput={(e) => setClanPassword(e.currentTarget.value)}
                    />
                    <button onClick={handleCreateClan}>Create Clan</button>
                    <button onClick={handleJoinClan}>Join Clan</button>
                </>
            ) : (
                <>
                    {/* If the user is in a clan, show leave or delete options */}
                    <Show
                        when={currentClan()!.owner_id === currentUser()!.id}
                        fallback={
                            <button onClick={handleLeaveClan}>
                                Leave Clan
                            </button>
                        }
                    >
                        <button onClick={handleDeleteClan}>Delete Clan</button>
                    </Show>
                </>
            )}

            {/* Always show the return button */}
            <button onClick={props.onReturn}>Return</button>
            <button onClick={props.onLogout}>Logout</button>
        </div>
    );
};

export default ClanMenu;
